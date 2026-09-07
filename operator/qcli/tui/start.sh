#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")"
exec bun run src/index.tsx
