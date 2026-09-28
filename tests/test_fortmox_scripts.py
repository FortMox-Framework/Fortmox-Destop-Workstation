import copy
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "scripts"))

from fortmox_scripts import all_templates, enabled_vms, load_config, load_template, qm_create_args


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

    def test_input_config_is_not_mutated(self):
        config = copy.deepcopy(self.config)
        original = copy.deepcopy(config)
        enabled_vms(config)
        self.assertEqual(config, original)


if __name__ == "__main__":
    unittest.main()
