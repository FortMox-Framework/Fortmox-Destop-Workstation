import copy
import importlib.util
import sys
import unittest
from pathlib import Path
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "scripts"))

from fortmox_scripts import all_templates, detect_form_factor, enabled_vms, load_config, load_template, qm_create_args

POWER_SCRIPT = ROOT / "scripts/power-management.py"
POWER_SPEC = importlib.util.spec_from_file_location("power_management", POWER_SCRIPT)
power_management = importlib.util.module_from_spec(POWER_SPEC)
POWER_SPEC.loader.exec_module(power_management)


class FortmoxScriptTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config, cls.project_root = load_config(str(ROOT / "config/system.yaml"))

    def test_enabled_vms_follow_dependencies(self):
        names = [name for name, _ in enabled_vms(self.config)]
        self.assertEqual(names, ["opnsense", "clean", "gaming", "research"])

    def test_dependency_cycle_is_rejected(self):
        config = {
            "microvm": {"enabled": False},
            "vms": {
                "alpha": {"enabled": True, "depends_on": "beta"},
                "beta": {"enabled": True, "depends_on": "alpha"},
            },
        }
        with self.assertRaisesRegex(ValueError, "dependency cycle"):
            enabled_vms(config)

    def test_all_five_templates_are_present_and_parse(self):
        templates = all_templates(self.config)
        self.assertEqual(len(templates), 5)
        for _, template_path in templates:
            with self.subTest(template=template_path):
                self.assertIsInstance(load_template(self.project_root, template_path)["vm"], dict)

    def test_qm_args_use_template_resources(self):
        template = load_template(self.project_root, "config/vm-templates/clean-vm.yaml")
        args = qm_create_args("clean", template)
        self.assertEqual(args[0:4], ["create", "100", "--name", "clean"])
        self.assertIn("--cores", args)
        self.assertEqual(args[args.index("--cores") + 1], "4")
        self.assertEqual(args[args.index("--memory") + 1], "8192")
        self.assertEqual(args[args.index("--scsi0") + 1], "local-lvm:50G")

    def test_qm_args_fall_back_to_vm_key_when_template_name_is_empty(self):
        template = {
            "vm": {
                "name": "-vm",
                "vmid": 100,
                "specs": {"cpu": {"cores": 2}, "memory": {"max": 2048}, "disk": {"storage": "local-lvm", "size": "10G"}},
            }
        }
        args = qm_create_args("clean", template)
        self.assertEqual(args[args.index("--name") + 1], "clean")

    def test_form_factor_uses_known_chassis_and_battery_fallback(self):
        self.assertEqual(detect_form_factor("10"), "laptop")
        self.assertEqual(detect_form_factor("3"), "desktop")
        self.assertEqual(detect_form_factor("0", battery_present=True), "laptop")
        self.assertEqual(detect_form_factor("0"), "unknown")

    def test_input_config_is_not_mutated(self):
        config = copy.deepcopy(self.config)
        original = copy.deepcopy(config)
        enabled_vms(config)
        self.assertEqual(config, original)

    def test_power_management_without_governors_is_successful_dry_run(self):
        with patch.object(power_management.glob, "glob", return_value=[]), patch.object(sys, "argv", [str(POWER_SCRIPT), "--dry-run"]):
            self.assertEqual(power_management.main(), 0)

    def test_power_management_without_governors_fails_when_applying(self):
        with patch.object(power_management.glob, "glob", return_value=[]), patch.object(sys, "argv", [str(POWER_SCRIPT), "--apply"]):
            self.assertEqual(power_management.main(), 1)


if __name__ == "__main__":
    unittest.main()
