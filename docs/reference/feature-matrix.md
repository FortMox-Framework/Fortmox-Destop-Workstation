# Feature Matrix

| Capability | Status | Entry point |
|------------|--------|-------------|
| Strict system config parsing | Available | `fortmox config validate` |
| VM dependency plan | Available | `fortmox config plan` |
| Hardware report | Available on Linux | `fortmox detect` |
| Template and host checks | Available; Proxmox checks require Proxmox | `fortmox verify` |
| Host deployment | Plan only | `fortmox deploy --dry-run` |
| VM create command generation | Available; opt-in execution | `fortmox vm create --dry-run` |
| GPU method planning | Available; no driver/VM mutation | `fortmox gpu` |
| CPU governor profiles | Plan and explicit apply | `fortmox power balanced` |
| Web/API service | Not implemented | N/A |

See [features and implementation status](features.md) before using any command on a host. This project remains alpha and is not a production-ready hypervisor configuration tool.
