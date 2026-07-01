#!/usr/bin/env bash

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
manifest_rel=".dvc/public-manifest.txt"
remote_name="${QUINE_DVC_PUBLIC_PULL_REMOTE:-public}"

usage() {
  cat <<'EOF'
Usage:
  scripts/pull-public-dvc-manifest.sh [manifest-path]

Purpose:
  Pull only the `.dvc` pointers listed in the tracked public manifest.

Defaults:
  manifest-path: .dvc/public-manifest.txt
  remote:        public

Override remote with:
  QUINE_DVC_PUBLIC_PULL_REMOTE=<remote-name>

Notes:
  Plain `dvc pull` attempts every DVC pointer in the checkout. The public
  manifest is the release contract for selectively published public artifacts.
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ $# -gt 1 ]]; then
  usage >&2
  exit 2
fi

if [[ $# -eq 1 ]]; then
  manifest_rel="$1"
fi

manifest_path="${repo_root}/${manifest_rel}"

if [[ ! -f "${manifest_path}" ]]; then
  printf 'manifest not found: %s\n' "${manifest_path}" >&2
  exit 1
fi

mapfile -t pointers < <(
  sed 's/[[:space:]]*#.*$//' "${manifest_path}" \
    | sed '/^[[:space:]]*$/d' \
    | awk '!seen[$0]++'
)

if [[ "${#pointers[@]}" -eq 0 ]]; then
  printf 'public manifest is empty; nothing to pull\n'
  exit 0
fi

missing=0
for pointer in "${pointers[@]}"; do
  if [[ ! -f "${repo_root}/${pointer}" ]]; then
    printf 'manifest entry does not exist: %s\n' "${pointer}" >&2
    missing=1
  fi
done

if [[ "${missing}" -ne 0 ]]; then
  exit 1
fi

printf 'pulling %d public pointer(s) from %s via %s\n' "${#pointers[@]}" "${manifest_rel}" "${remote_name}"
dvc pull -r "${remote_name}" "${pointers[@]}"
