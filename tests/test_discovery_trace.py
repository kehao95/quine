#!/usr/bin/env python3
import json
import subprocess
import unittest
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent


class TestDiscoveryTrace(unittest.TestCase):
    def test_all_summary(self):
        cmd = ["python3", str(REPO_ROOT / "scripts" / "discovery-trace.py"), "--all", "--json"]
        res = subprocess.run(cmd, capture_output=True, text=True, cwd=REPO_ROOT)
        self.assertEqual(res.returncode, 0, res.stderr)
        data = json.loads(res.stdout)
        self.assertGreater(data.get("theory_objects", 0), 50)
        self.assertGreater(data.get("predictions", 0), 40)
        self.assertGreater(data.get("experiments", 0), 50)
        self.assertGreater(data.get("mechanisms", 0), 10)
        self.assertGreater(data.get("implementations", 0), 30)

    def test_trace_theory_object(self):
        cmd = ["python3", str(REPO_ROOT / "scripts" / "discovery-trace.py"), "carrier", "--json"]
        res = subprocess.run(cmd, capture_output=True, text=True, cwd=REPO_ROOT)
        self.assertEqual(res.returncode, 0, res.stderr)
        data = json.loads(res.stdout)
        self.assertEqual(data["target"]["id"], "carrier")
        self.assertTrue(len(data["predictions"]) > 0)
        self.assertTrue(len(data["mechanisms"]) > 0)
        self.assertTrue(len(data["implementations"]) > 0)

    def test_trace_prediction(self):
        cmd = ["python3", str(REPO_ROOT / "scripts" / "discovery-trace.py"), "delayed-addressability-P1", "--json"]
        res = subprocess.run(cmd, capture_output=True, text=True, cwd=REPO_ROOT)
        self.assertEqual(res.returncode, 0, res.stderr)
        data = json.loads(res.stdout)
        self.assertEqual(data["target"]["id"], "delayed-addressability-P1")
        self.assertEqual(data["target"]["status"], "confirmed")

    def test_trace_missing_entity(self):
        cmd = ["python3", str(REPO_ROOT / "scripts" / "discovery-trace.py"), "non-existent-entity-xyz", "--json"]
        res = subprocess.run(cmd, capture_output=True, text=True, cwd=REPO_ROOT)
        self.assertEqual(res.returncode, 0, res.stderr)
        data = json.loads(res.stdout)
        self.assertIn("error", data)


if __name__ == "__main__":
    unittest.main()
