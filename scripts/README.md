# Python 3 Utilities

The primary interface is the Go binary documented in [FortMox CLI](../docs/README.md). These Python 3 utilities are config-driven helpers and use PyYAML (`python3 -m pip install -r ../requirements.txt`). They do not replace the CLI.

| Script | Purpose | Changes system by default |
|--------|---------|---------------------------|
| `config-manager.py` | Load and summarize `config/system.yaml` | No |
| `detect-hardware.py` | Report CPU, GPU, IOMMU and form factor | No |
| `verify.py` | Check all configured templates, enabled VMs and firewall service state | No |
| `deploy.py` | Print the host deployment plan | No; apply operations are not implemented |
| `deploy-vms.py` | Build config-driven `qm create` commands | No; requires `--apply` |
| `gpu-setup.py` | Show configured and detected GPU devices | No |
| `power-management.py` | Plan CPU frequency governor changes | No; requires `--apply` |

Bash implementations have been removed. These Python 3 utilities and the Go CLI are the supported entrypoints.
