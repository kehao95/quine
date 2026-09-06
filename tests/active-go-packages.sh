#!/usr/bin/env bash
#
# Emit the active Go package set used by substrate validation.
# Historical evidence trees under */runs/* are repo-retained artifacts, not
# active substrate packages.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

go list ./... | awk '$0 !~ /\/runs\//'
