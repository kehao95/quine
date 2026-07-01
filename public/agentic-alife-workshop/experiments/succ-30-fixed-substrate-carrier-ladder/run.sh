#!/bin/bash
# Experiment 6M.02: Fixed-Substrate Carrier Ladder
#
# Thin wrapper around the shared p6m morphology ladder harness.

set -euo pipefail

MORPH_EXPERIMENT_DIR="$(cd "$(dirname "$0")" && pwd)"
MORPH_EXPERIMENT_ID="6M.02"
MORPH_EXPERIMENT_TITLE="Fixed-Substrate Carrier Ladder"
MORPH_DEFAULT_VARIANT="C04-public-ctl-raw-carrier"
MORPH_SELF_SOURCE_ENABLED="1"
MORPH_SELF_DESCRIPTION_LABEL="source-projection"
MORPH_STDIN_VARIANTS="C01-raw-stdin C04-public-ctl-raw-carrier"
MORPH_FD_VARIANTS="C02-inherited-fd"
MORPH_CONTROL_VARIANTS="C03-native-ctl-inject C03-native-ctl-wake C03b-native-ctl-no-preserved-body"
MORPH_OBSERVE_CTL_VARIANTS="C03-native-ctl-inject C03-native-ctl-wake C03b-native-ctl-no-preserved-body C04-public-ctl-raw-carrier"
MORPH_LAUNCH_ONLY_VARIANTS="C05-launch-only-control"
MORPH_NO_PRESERVED_BODY_VARIANTS="C03b-native-ctl-no-preserved-body"

source "${MORPH_EXPERIMENT_DIR}/../lib/morphology-ladder-harness.sh"
