"""Shared, config-driven helpers for the FortMox Python 3 command scripts."""

from __future__ import annotations

import os
import platform
import shutil
import subprocess
from pathlib import Path
from typing import Any

import yaml

def load_config(path: str) -> tuple[dict[str, Any], Path]:
    """Load YAML and return it with the project root used for template paths."""
    config_path = Path(path)
    if not config_path.is_absolute():
        config_path = Path.cwd() / config_path
    config_path = config_path.resolve()
    with config_path.open(encoding="utf-8") as stream:
        config = yaml.safe_load(stream)
    if not isinstance(config, dict):
        raise ValueError(f"{config_path} must contain a YAML mapping")
    return config, config_path.parent.parent


def enabled_vms(config: dict[str, Any]) -> list[tuple[str, dict[str, Any]]]:
    """Return enabled VMs in dependency order, including OPNsense; reject cycles."""
    entries = config.get("vms", {})
    result = [(name, vm) for name, vm in entries.items() if vm.get("enabled", False)]
    microvm = config.get("microvm", {})
    firewall = microvm.get("opnsense_firewall", {})
    if microvm.get("enabled") and firewall.get("enabled"):
        opnsense = entries.get("opnsense", {})
        if not any(name == "opnsense" for name, _ in result):
            result.append(("opnsense", {**opnsense, "template_file": opnsense.get("template_file") or firewall.get("template_file"), "enabled": True}))
    enabled = dict(result)
    ordered: list[tuple[str, dict[str, Any]]] = []
    visiting: set[str] = set()
    visited: set[str] = set()

    def visit(name: str) -> None:
        if name in visited:
            return
        if name in visiting:
            raise ValueError(f"VM dependency cycle involving {name}")
        visiting.add(name)
        dependencies = config.get("vms", {}).get(name, {}).get("depends_on", [])
        if isinstance(dependencies, str):
            dependencies = [dependencies]
        # Visit prerequisites first so command order follows the configured dependency graph.
        for dependency in sorted(dependencies):
            if dependency in enabled:
                visit(dependency)
        visiting.remove(name)
        visited.add(name)
        ordered.append((name, enabled[name]))

    for name in sorted(enabled):
        visit(name)
    return ordered


def template_file(config: dict[str, Any], name: str, entry: dict[str, Any]) -> str:
    """Resolve a template path, falling back to the microvm firewall settings."""
    file_name = entry.get("template_file")
    if not file_name and name == "opnsense":
        file_name = config.get("microvm", {}).get("opnsense_firewall", {}).get("template_file")
    if not file_name:
        raise ValueError(f"vms.{name}.template_file is not set")
    return str(file_name)


def load_template(root: Path, file_name: str) -> dict[str, Any]:
    """Load a template and require the VM mapping consumed by command builders."""
    path = (root / file_name).resolve()
    with path.open(encoding="utf-8") as stream:
        template = yaml.safe_load(stream)
    if not isinstance(template, dict) or not isinstance(template.get("vm"), dict):
        raise ValueError(f"{path} must define a vm mapping")
    return template


def all_templates(config: dict[str, Any]) -> list[tuple[str, str]]:
    """List configured templates, including those belonging to disabled VMs."""
    values = []
    for name, vm in config.get("vms", {}).items():
        file_name = template_file(config, name, vm)
        values.append((name, file_name))
    return values


def qm_create_args(name: str, template: dict[str, Any]) -> list[str]:
    """Build the supported subset of Proxmox qm create arguments for a template."""
    vm = template["vm"]
    specs = vm.get("specs", {})
    cpu = specs.get("cpu", {})
    memory = specs.get("memory", {})
    disk = specs.get("disk", {})
    vmid = int(vm.get("vmid", 0))
    cores = int(cpu.get("cores", 0))
    max_memory = int(memory.get("max", 0))
    storage, size = disk.get("storage"), disk.get("size")
    if vmid < 1 or cores < 1 or max_memory < 1 or not storage or not size:
        raise ValueError(f"template for {name} needs a positive vmid, cores, memory and disk")
    vm_name = str(vm.get("name", name)).removesuffix("-vm")
    return ["create", str(vmid), "--name", vm_name, "--cores", str(cores), "--memory", str(max_memory), "--scsi0", f"{storage}:{size}", "--net0", "virtio,bridge=vmbr0,firewall=1", "--ostype", "l26", "--agent", "1"]


def hardware_report() -> dict[str, Any]:
    """Collect best-effort CPU, virtualization, IOMMU, GPU, and form-factor data."""
    cpu = "unknown"
    flags = ""
    cpuinfo = Path("/proc/cpuinfo")
    if cpuinfo.exists():
        for line in cpuinfo.read_text(errors="replace").splitlines():
            if line.startswith(("model name", "Hardware")) and cpu == "unknown":
                cpu = line.partition(":")[2].strip()
            if line.startswith(("flags", "Features")) and not flags:
                flags = line.partition(":")[2]
    # Linux exposes Intel and AMD virtualization support as vmx and svm CPU flags.
    groups = list(Path("/sys/kernel/iommu_groups").glob("[0-9]*"))
    chassis = Path("/sys/class/dmi/id/chassis_type")
    kind = "unknown"
    if chassis.exists():
        value = chassis.read_text().strip()
        # DMI chassis codes 8-14 generally denote portable systems.
        kind = "laptop" if value in {"8", "9", "10", "11", "12", "14"} else "desktop"
    elif Path("/sys/class/power_supply/BAT0").exists():
        kind = "laptop"
    gpus = []
    if shutil.which("lspci"):
        output = subprocess.run(["lspci", "-nn"], text=True, capture_output=True, check=False).stdout
        gpus = [line.strip() for line in output.splitlines() if any(label in line.lower() for label in ("vga compatible controller", "3d controller", "display controller"))]
    return {"cpu": cpu, "logical_cpus": os.cpu_count() or 0, "virtualization": any(flag in f" {flags} " for flag in ("vmx", "svm")), "iommu": bool(groups), "form_factor": kind, "gpus": gpus, "platform": platform.platform()}
