# FortMoxDesktop Workstation

A Proxmox workstation project for security research, malware analysis and gaming on laptops and desktops. Its intended design uses declarative YAML and an OPNsense firewall VM; the current CLI does not configure the host networking or firewall path.

> **Status: Alpha.** The Go CLI supports config validation, hardware detection, read-only verification, deployment planning, VM command planning/creation, GPU strategy planning and CPU power profiles. Host deployment is not production-ready; review every proposed system change and test on disposable hardware first.

## What you get

| VM | ID | Purpose | Intended network |
|----|----|---------|---------|
| OPNsense (MicroVM) | 50 | Intended firewall/router VM | Uplink (not configured by the CLI) |
| Clean | 100 | Daily work, browsing | Intended via OPNsense (not configured) |
| Gaming | 101 | GPU passthrough + Looking Glass | Intended via OPNsense (not configured) |
| Research | 102 | Malware analysis | No network (intended only; not enforced) |
| Tools | 103 | Network analysis and testing | Intended isolated network (not configured) |

These are intended roles and VMIDs, not a deployed topology. The current CLI does not configure bridges, OPNsense routing, or VM network isolation.

Planned configuration areas include IOMMU and GPU passthrough (virtIO, SR-IOV, vGPU, full), encrypted storage vaults, kernel hardening, laptop/desktop power profiles, and a Wayland compositor on the host. Their presence in YAML or a guide does not mean the CLI applies them.

## Project website

The static project overview and local documentation viewer are in [site](site/index.html). Build the site and render every guide, reference, and planning page with:

```bash
python -m pip install -r site/requirements.txt
python site/build.py
```

Open `_site/index.html` to preview the complete site. GitHub Actions publishes the generated site to [GitHub Pages](https://fortmox-framework.github.io/Fortmox-Destop-Workstation/) on pushes to `main` and manual workflow runs; see [the Pages workflow](.github/workflows/pages.yml). If Pages is not enabled for this repository, set **Settings → Pages → Build and deployment → Source** to **GitHub Actions**.

## Quick start

Start with the [getting-started guide](docs/guides/getting-started.md). Build the CLI from the project root:

```bash
go build -o fortmox ./cmd/fortmox
./fortmox config validate
./fortmox detect
./fortmox deploy --dry-run
./fortmox vm create --dry-run
```

`deploy` currently plans host changes only. VM creation is dry-run by default; use `--dry-run=false` only on a Proxmox host after reviewing the generated `qm` commands. `verify` is read-only and reports Proxmox-only checks as skipped when run elsewhere.

## Project layout

```
.
├── README.md
├── config/
│   ├── system.yaml              # main configuration
│   └── vm-templates/            # opnsense, clean, gaming, research, tools
├── cmd/fortmox/                 # Cobra CLI entry point
├── internal/                    # config, hardware, verify, deploy plan, vm, gpu, power
├── scripts/                     # Python 3 utilities
├── requirements.txt             # PyYAML for Python utilities
└── docs/
    ├── README.md                # documentation index, start here
    ├── guides/                  # how-to: install, deploy, GPU, firewall, hardening
    ├── reference/               # architecture, features, quick reference
    ├── planning/                # roadmaps and designs (Go CLI, GPU tool)
```

## Go CLI

The CLI is the primary interface. Python 3 scripts are the only script entrypoints; Bash implementations have been removed.
See [TODO.md](TODO.md) for the remaining implementation and release work.

```bash
go mod tidy                                # first time only (resolve Go dependencies)
go test ./...
go run ./cmd/fortmox config validate
go run ./cmd/fortmox config plan
go run ./cmd/fortmox verify                # read-only host and VM checks
go run ./cmd/fortmox gpu --method auto     # strategy plan only
go run ./cmd/fortmox power balanced        # dry run; --apply requires root
```

## Security notes

This project is intended for security research but does not currently enforce VM isolation. Do not treat the Research VM as air-gapped or run untrusted malware until its network isolation is independently configured and tested. Keep logs, protect credentials, and verify firewall rules and backups.

## Requirements

- CPU with IOMMU and VT-x/AMD-V (VT-d/AMD-Vi for passthrough)
- 16 GB RAM minimum (32 GB+ recommended for gaming)
- 256 GB storage minimum (NVMe SSD 500 GB+ recommended)
- A GPU (NVIDIA, AMD or Intel); a dedicated one for passthrough
- Go 1.22+ to build the CLI
