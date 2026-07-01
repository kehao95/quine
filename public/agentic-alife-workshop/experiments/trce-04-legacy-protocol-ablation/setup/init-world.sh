#!/bin/bash

set -euo pipefail

WORKSPACE="${1:?workspace required}"
SPEC_PATH="${2:?spec path required}"
ITEM_COUNT="${3:-24}"

rm -rf "${WORKSPACE}"
mkdir -p "${WORKSPACE}/lineage_state"
mkdir -p "$(dirname "${SPEC_PATH}")"

python3 - "${SPEC_PATH}" "${ITEM_COUNT}" <<'PY'
from __future__ import annotations

import json
import sys

spec_path = sys.argv[1]
item_count = int(sys.argv[2])

def payload_for(g: int) -> str:
    family = ((g * 73) + 19) % 6
    part1 = (g * 2654435761 + 2246822519) & 0xffffffff
    part2 = (((g ^ 0x9E37) * 40503) + 0x1234) & 0xffff
    part3 = ((g * 31337) + 0x4567) & 0xffff
    part4 = ((g * 1103515245 + 12345) >> 8) & 0xffff
    part5 = (((g + 1) * 281474976710597) + 0x123456789AB) & 0xffffffffffff
    part6 = (g * 747796405 + 2891336453) & 0xffffffff
    part7 = (((g + 17) * 48271) + 0x5F5E) & 0xffffff
    part8 = ((((g + 101) * 214013) + 2531011) >> 4) & 0xffffffff

    if family == 0:
        return f"{part1:08x}-{part2:04x}-{part3:04x}-{part4:04x}-{part5:012x}"
    if family == 1:
        return f"{part1:08x}{part6:08x}{((part1 ^ part6) & 0xffffffff):08x}{((part6 + 0x9e3779b9) & 0xffffffff):08x}{((part1 + 0x7f4a7c15) & 0xffffffff):08x}"
    if family == 2:
        return f"{part7:06x}.{((part7 ^ 0x5a5a5a) & 0xffffff):06x}.{((part7 + 0x13579b) & 0xffffff):06x}.{((part7 + 0x2468ac) & 0xffffff):06x}"
    if family == 3:
        return f"{part2:04x}-{part3:04x}-{((part2 ^ part3) & 0xffff):04x}-{((part4 + 0x1111) & 0xffff):04x}-{((part2 + 0x2222) & 0xffff):04x}-{((part3 + 0x3333) & 0xffff):04x}"
    if family == 4:
        return f"{part8:08x}/{((part8 ^ 0xa5a5a5a5) & 0xffffffff):08x}/{((part8 + 0x01010101) & 0xffffffff):08x}"
    return f"{part6:08x}:{part2:04x}:{((part6 + part1) & 0xffffffff):08x}:{part3:04x}"

items = {str(i): payload_for(i - 1) for i in range(1, item_count + 1)}

with open(spec_path, "w", encoding="utf-8") as f:
    json.dump({"items": items}, f, indent=2, sort_keys=True)
    f.write("\n")
PY
