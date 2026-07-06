#!/bin/sh
# Shared infra (plan §2.1): metadata-perception-safe snapshotter.
#
# Emits a stable, sorted metadata listing for a tree:  path \t mtime(epoch) \t mode \t size \t type
# Run this on the neutral STAGE *before* any chown/cp recovery, so a metadata-axis DV
# (E4 mtime/mode, E7) is captured before the harness round-trip can mutate it.
#
# Capture seed-side with the IDENTICAL mechanism right after staging, so seed and final
# are comparable. See _lib/selftest_metadata.sh for the round-trip self-test.
#
# Usage: emit_metadata.sh <root_dir>     # prints to stdout
set -eu
root="${1:?usage: emit_metadata.sh <root_dir>}"
# %P path relative to root (stable across mount points); %T@ mtime epoch.frac; %m octal
# mode; %s size bytes; %y type (f/d/l). LC_ALL=C sort for deterministic diffs.
find "$root" -mindepth 1 -printf '%P\t%T@\t%m\t%s\t%y\n' 2>/dev/null | LC_ALL=C sort
