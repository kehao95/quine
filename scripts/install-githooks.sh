#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

git config core.hooksPath .githooks
printf 'Configured core.hooksPath=%s\n' ".githooks"
