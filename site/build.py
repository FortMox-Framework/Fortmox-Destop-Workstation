#!/usr/bin/env python3
"""Build the static FortMox site and render repository Markdown as local pages."""

from __future__ import annotations

import argparse
from dataclasses import dataclass
from html import escape
from html.parser import HTMLParser
from pathlib import Path
import posixpath
import re
import shutil
from urllib.parse import unquote, urlsplit, urlunsplit

import markdown

ROOT = Path(__file__).resolve().parents[1]
SITE = ROOT / "site"
DOCS = ROOT / "docs"
REPO_URL = "https://github.com/FortMox-Framework/Fortmox-Destop-Workstation"
GROUPS = ("guides", "reference", "planning", "project")
GROUP_TITLES = {
    "guides": "Guides",
    "reference": "Reference",
    "planning": "Planning",
    "project": "Project",
}


@dataclass(frozen=True)
class Page:
    source: Path
    output: Path
    section: str
    title: str


class LinkRewriter(HTMLParser):
    """Rewrite source Markdown links to generated local pages when available."""

    def __init__(self, source: Path, current_output: Path, output_root: Path):
        super().__init__(convert_charrefs=False)
        self.source = source
        self.current_output = current_output
        self.output_root = output_root
        self.parts: list[str] = []

    def rewritten_href(self, href: str) -> str:
        parsed = urlsplit(href)
        if not parsed.path or parsed.scheme or parsed.netloc:
            return href
        target = (self.source.parent / unquote(parsed.path)).resolve()
        if not target.is_file():
            return href

        if target == ROOT / "docs/README.md":
            destination = self.output_root / "docs/index.html"
        elif target == ROOT / "README.md":
            destination = self.output_root / "docs/project.html"
        elif target == ROOT / "TODO.md":
            destination = self.output_root / "docs/todo.html"
        elif target.is_relative_to(DOCS) and target.suffix.lower() == ".md":
            destination = self.output_root / "docs" / target.relative_to(DOCS).with_suffix(".html")
        elif target.is_relative_to(ROOT):
            repository_path = target.relative_to(ROOT).as_posix()
            return f"{REPO_URL}/blob/main/{repository_path}" + (f"#{parsed.fragment}" if parsed.fragment else "")
        else:
            return href

        relative = posixpath.relpath(destination, self.current_output.parent).replace("\\", "/")
        return urlunsplit(("", "", relative, parsed.query, parsed.fragment))

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        values = []
        for name, value in attrs:
            if value is None:
                values.append(name)
                continue
            if name == "href":
                value = self.rewritten_href(value)
            values.append(f'{name}="{escape(value, quote=True)}"')
        suffix = " />" if tag in {"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr"} else ">"
        self.parts.append(f"<{tag}{(' ' + ' '.join(values)) if values else ''}{suffix}")

    def handle_endtag(self, tag: str) -> None:
        if tag not in {"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr"}:
            self.parts.append(f"</{tag}>")

    def handle_data(self, data: str) -> None:
        self.parts.append(data)

    def handle_entityref(self, name: str) -> None:
        self.parts.append(f"&{name};")

    def handle_charref(self, name: str) -> None:
        self.parts.append(f"&#{name};")

    def handle_comment(self, data: str) -> None:
        self.parts.append(f"<!--{data}-->")

    def handle_decl(self, decl: str) -> None:
        self.parts.append(f"<!{decl}>")


class LocalLinkCollector(HTMLParser):
    def __init__(self):
        super().__init__()
        self.links: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        for name, value in attrs:
            if name in {"href", "src"} and value:
                self.links.append(value)


def page_title(source: Path, text: str) -> str:
    heading = re.search(r"(?m)^#\s+(.+?)\s*$", text)
    if heading:
        return re.sub(r"[`*_]", "", heading.group(1)).strip()
    return source.stem.replace("-", " ").title()


def collect_pages() -> list[Page]:
    pages = []
    for source in sorted(DOCS.rglob("*.md")):
        if source == DOCS / "README.md":
            continue
        relative = source.relative_to(DOCS)
        section = relative.parts[0] if len(relative.parts) > 1 else "reference"
        if section not in GROUPS:
            section = "reference"
        output = Path("docs") / relative.with_suffix(".html")
        pages.append(Page(source, output, section, page_title(source, source.read_text(encoding="utf-8"))))
    for source, output, title in (
        (ROOT / "README.md", Path("docs/project.html"), "Project Overview"),
        (ROOT / "TODO.md", Path("docs/todo.html"), "Project TODO"),
    ):
        pages.append(Page(source, output, "project", title))
    return sorted(pages, key=lambda page: (GROUPS.index(page.section), page.title.casefold()))


def relative_href(current: Path, target: Path) -> str:
    return posixpath.relpath(target, current.parent).replace("\\", "/")


def navigation(pages: list[Page], current: Path) -> str:
    chunks = []
    for group in GROUPS:
        members = [page for page in pages if page.section == group]
        if not members:
            continue
        links = []
        for page in members:
            href = relative_href(current, page.output)
            active = ' aria-current="page"' if page.output == current else ""
            links.append(f'<li><a class="doc-link" href="{escape(href, quote=True)}"{active}>{escape(page.title)}</a></li>')
        chunks.append(
            f'<section class="doc-nav-group" data-group="{group}">'
            f'<h2>{escape(GROUP_TITLES[group])}</h2><ul>{"".join(links)}</ul></section>'
        )
    overview = relative_href(current, Path("docs/index.html"))
    return f'<a class="overview-link" href="{escape(overview, quote=True)}">All documentation <span aria-hidden="true">↗</span></a>{"".join(chunks)}'


def docs_shell(output: Path, title: str, section: str, body: str, pages: list[Page], is_index: bool = False) -> str:
    css = relative_href(output, Path("docs.css"))
    js = relative_href(output, Path("docs.js"))
    logo = relative_href(output, Path("assets/fortmox-logo.png"))
    home = relative_href(output, Path("index.html"))
    index = relative_href(output, Path("docs/index.html"))
    content = f"""<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="theme-color" content="#101511">
  <meta name="description" content="FortMox project documentation: guides, references, and implementation plans.">
  <title>{escape(title)} | FortMox Docs</title>
  <link rel="icon" type="image/png" href="{escape(logo, quote=True)}">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600&family=Space+Grotesk:wght@400;500;600;700&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="{escape(css, quote=True)}">
  <script src="{escape(js, quote=True)}" defer></script>
</head>
<body class="docs-page">
  <a class="docs-skip" href="#doc-main">Skip to documentation</a>
  <header class="docs-topbar">
    <a class="docs-brand" href="{escape(home, quote=True)}"><img src="{escape(logo, quote=True)}" alt="" width="38" height="38"><span>FORTMOX<span>.</span></span></a>
    <div class="docs-top-title">DOCUMENTATION <span>/</span> {escape(GROUP_TITLES.get(section, "Library").upper())}</div>
    <a class="docs-github" href="{escape(home, quote=True)}">Project home <span aria-hidden="true">↗</span></a>
  </header>
  <button class="docs-menu-button" type="button" aria-expanded="false" aria-controls="docs-sidebar"><span aria-hidden="true">☰</span> Browse docs</button>
  <div class="docs-layout">
    <aside class="docs-sidebar" id="docs-sidebar" aria-label="Documentation navigation">
      <a class="docs-home-link" href="{escape(home, quote=True)}">← Project home</a>
      <label class="docs-search-label" for="docs-search">FIND A PAGE</label>
      <input class="docs-search" id="docs-search" type="search" placeholder="Search documentation" autocomplete="off">
      <nav class="docs-nav">{navigation(pages, output)}</nav>
      <div class="docs-sidebar-note"><span>ALPHA</span><p>Instructions describe manual procedures unless marked implemented.</p></div>
    </aside>
    <main class="docs-main" id="doc-main">
      <div class="docs-crumb"><a href="{escape(index, quote=True)}">Docs</a><span>/</span>{escape(GROUP_TITLES.get(section, "Library"))}</div>
      {f'<section class="docs-welcome"><p class="docs-kicker">FORTMOX / DOCS</p><h1>{escape(title)}</h1><p>Guides, implementation references, and plans for the FortMox Proxmox workstation project.</p><div class="docs-status"><span>ALPHA</span><span>Some architecture pages are proposals, not active features.</span></div></section>' if is_index else f'<header class="docs-article-head"><p class="docs-kicker">{escape(GROUP_TITLES.get(section, "Project").upper())} / FORTMOX</p><h1>{escape(title)}</h1><p class="docs-article-state">Project documentation <span>·</span> Check the current CLI boundary before applying manual steps.</p></header>'}
      {body}
    <footer class="docs-page-footer"><a href="{escape(index, quote=True)}">← Documentation index</a><a href="{escape(home, quote=True)}">FortMox project home ↗</a></footer>
    </main>
  </div>
</body>
</html>
"""
    return content


def home_cards(pages: list[Page], output: Path) -> str:
    summaries = {
        "guides": "Installation, CLI usage, and manual system-topic procedures.",
        "reference": "Current implementation status, architecture boundaries, and command reference.",
        "planning": "Design notes and proposals; not current implementation commitments.",
        "project": "Project overview and outstanding implementation work.",
    }
    groups = []
    for group in GROUPS:
        members = [page for page in pages if page.section == group]
        if not members:
            continue
        links = "".join(
            f'<li><a href="{escape(relative_href(output, page.output), quote=True)}">{escape(page.title)}<span aria-hidden="true">↗</span></a></li>'
            for page in members
        )
        groups.append(
            f'<section class="docs-category"><p class="docs-kicker">{escape(GROUP_TITLES[group].upper())} / {len(members):02d}</p>'
            f'<h2>{escape(GROUP_TITLES[group])}</h2><p>{escape(summaries[group])}</p><ul>{links}</ul></section>'
        )
    return f'<div class="docs-category-grid">{"".join(groups)}</div>'


def rewrite_links(rendered: str, source: Path, output: Path, output_root: Path) -> str:
    parser = LinkRewriter(source, output, output_root)
    parser.feed(rendered)
    parser.close()
    return "".join(parser.parts)


def validate_local_links(output_root: Path) -> int:
    broken = []
    checked = 0
    for page in output_root.rglob("*.html"):
        collector = LocalLinkCollector()
        collector.feed(page.read_text(encoding="utf-8"))
        for link in collector.links:
            parsed = urlsplit(link)
            if parsed.scheme or parsed.netloc or not parsed.path:
                continue
            checked += 1
            destination = (page.parent / unquote(parsed.path)).resolve()
            if not destination.exists():
                broken.append(f"{page.relative_to(output_root)} -> {link}")
    if broken:
        raise ValueError("broken local links:\n" + "\n".join(broken))
    return checked


def build(output: Path) -> None:
    pages = collect_pages()
    output.mkdir(parents=True, exist_ok=True)
    shutil.copy2(SITE / "index.html", output / "index.html")
    shutil.copy2(SITE / "styles.css", output / "styles.css")
    shutil.copy2(SITE / "script.js", output / "script.js")
    shutil.copy2(SITE / "docs.css", output / "docs.css")
    shutil.copy2(SITE / "docs.js", output / "docs.js")
    shutil.copytree(SITE / "assets", output / "assets", dirs_exist_ok=True)

    docs_index = Path("docs/index.html")
    docs_index_path = output / docs_index
    docs_index_path.parent.mkdir(parents=True, exist_ok=True)
    docs_index_path.write_text(docs_shell(docs_index, "Documentation", "guides", home_cards(pages, docs_index), pages, is_index=True), encoding="utf-8")

    for page in pages:
        page_output = output / page.output
        page_output.parent.mkdir(parents=True, exist_ok=True)
        source_text = page.source.read_text(encoding="utf-8")
        source_text = re.sub(r"(?m)^#\s+.+?\s*\n", "", source_text, count=1)
        rendered = markdown.markdown(source_text, extensions=["extra", "toc"], extension_configs={"toc": {"permalink": True}})
        rendered = rewrite_links(rendered, page.source, page_output, output)
        article = f'<article class="markdown-body">{rendered}</article>'
        page_output.write_text(docs_shell(page.output, page.title, page.section, article, pages), encoding="utf-8")

    links = validate_local_links(output)
    print(f"Built site and {len(pages)} documentation pages in {output}; validated {links} local links")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, default=ROOT / "_site", help="directory for the generated static site")
    args = parser.parse_args()
    output = args.output if args.output.is_absolute() else ROOT / args.output
    build(output.resolve())
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
