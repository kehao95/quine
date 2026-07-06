# check-control-plane — Makefile excerpt (curated)
#
# This is the relevant slice of the repository's root Makefile. `make
# check-control-plane` is the control-plane validation gate the paper's agents
# run: an aggregate target that fans out to ~50 sub-checks. The full dependency
# list is reproduced verbatim below so the reader can see the whole gate; only
# the sub-check RECIPES most directly relevant to *this* paper — the ones that
# police the maintenance contract and its inbound reference network — are
# expanded here and shipped as scripts alongside this file. The remaining
# sub-checks are omitted (they validate unrelated parts of the private repo's
# governance structure).

check-control-plane: check-githooks-active check-experiment-paper-feeds \
	check-experiment-prompt-posture check-paper-dossier-links \
	check-paper-evidence-backbone check-paper-metadata-vocabulary \
	check-configuration-shape check-venue-refs check-paper-readout-coverage \
	check-model-evaluations check-runtime-doc-sync check-validation-surface \
	check-active-doc-links check-stable-evidence-surfaces \
	check-experiment-note-index check-experiment-current-state \
	check-experiment-surface-map check-experiment-subdir-index \
	check-experiment-task-readmes check-experiment-incubator-readmes \
	check-experiment-archive-readmes check-experiment-family \
	check-experiment-theory-reciprocity check-paper-legacy-track-roots \
	check-paper-track-readmes check-paper-subdir-index \
	check-paper-track-meta-index check-paper-track-special-index \
	check-paper-template-surface-index check-root-entry-docs \
	check-tests-entry-docs check-tests-surface-ownership \
	check-obsidian-base-index check-obsidian-surface-ownership \
	check-scripts-readme check-scripts-surface-ownership check-githooks-readme \
	check-githooks-surface-ownership check-surface-inventory \
	check-surface-ownership check-paper-home-root-boundedness \
	check-paper-core-surface-index check-control-plane-frontmatter \
	check-control-plane-projection-maintenance check-control-plane-size \
	check-phase-interpretation-ownership check-validator-registration \
	check-evolution-hygiene check-citation-cff check-selfsource-freshness \
	check-structural-pressure

# --- sub-checks shipped with this artifact (the reference/contract police) ---

# Flags dangling references across active docs — the check whose failure a fresh
# agent repairs by restoring a deleted structure (the paper's recurrence leg).
check-active-doc-links:
	./scripts/check-active-doc-links.sh --strict

# Validate the control-plane structure itself: frontmatter, projection-maintenance
# blocks, and size bounds on the contract surface.
check-control-plane-frontmatter:
	./scripts/check-control-plane-frontmatter.sh --strict

check-control-plane-projection-maintenance:
	./scripts/check-control-plane-projection-maintenance.sh --strict

check-control-plane-size:
	./scripts/check-control-plane-size.sh
