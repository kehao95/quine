---
title: "How I Open-Sourced a Living Research Workbench"
paper_id: open-sourcing-a-research-workbench
claim_bundle: open-sourcing-a-research-workbench
venue: blog
format-envelope: blog-post
track: blogs
cluster: repository-self-model
status: maintained
revision: true
state_mode: full-manuscript-home
role: open-source-release-case-study
family: reflexive-repository
theme: repository-self-model
theme_home: ../../../theory/views/reflexive.md
publication_class: public-essay
maturity_target: synthesis
source_home: manuscript.md
plan_home: manuscript-plan.md
theory-objects: [carrier, trace-inheritance, delayed-addressability, protocol-emergence, structural-pull, narrative-claim]
gaps: [revisable-time-semantics]
evidence: [none-yet]
view-affinity: [reflexive]
---

# How I Open-Sourced a Living Research Workbench

**Working title:** How I Open-Sourced a Living Research Workbench

**Upstream release line:** [`../../../development/open-source/README.md`](../../../development/open-source/README.md)

## Current Question

How do we explain turning Quine from a live private research repository into a
public workbench without flattening the real work into "removed secrets and
clicked publish"?

## Thesis

Open-sourcing a live research workbench was not a one-shot hygiene task. It
required separating three boundaries that had been entangled inside one
repository: public source/history, retained private artifacts, and selectively
published public evidence. Once those boundaries were made explicit, the actual
release work became a series of bounded control-plane and history repairs
rather than a vague cleanup campaign.

## In Scope

- the release chronology from private `lab` history to a sanitized public
  branch
- the concrete pain points surfaced by retained artifacts, raw logs, provider
  profiles, drafts, and DVC topology
- why public DVC had to become manifest-scoped publication rather than default
  restore
- the owner surfaces, audits, and release controls needed before the repo
  could become public
- practical lessons for open-sourcing other research workbenches that contain
  real research residue

## Out Of Scope

- a full origin story of Quine as a research program
- an exhaustive commit-by-commit release diary
- generic open-source advice detached from research-artifact reality
- using one public essay as the canonical owner of the release protocol itself

## Current State

- status: revision — complete draft under editorial tightening
- live status: [`../../../../development/status/README.md`](../../../development/status/README.md)
- canonical plan: [`./manuscript-plan.md`](manuscript-plan.md)
- raw ideation archive:
  [`./analysis/ideation-capture.md`](analysis/ideation-capture.md)
- evidence collection: [`./analysis/pain-point-reference-index.md`](analysis/pain-point-reference-index.md)
- release owner surface:
  [`../../../development/open-source/README.md`](../../../development/open-source/README.md)
- theme owner: [`../../../theory/views/reflexive.md`](../../../theory/views/reflexive.md)
- essays registry: [`../../essays/PORTFOLIO.md`](../../essays/PORTFOLIO.md)

## Relationship To Theme

This dossier is not the canonical theory owner for repository self-maintenance.
It is a concrete Quine case study showing what repository-self-model work looks
like when the object is "make the living lab public without lying about what is
still private."

## Experiment Feeds

- `none-yet` - none - not-for-paper-yet - this essay is currently driven by the release control plane, history-rewrite record, and artifact-policy surfaces rather than a registered experiment feed

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - the current backbone is still a selected set of release-control and evolution surfaces rather than a frozen paper-style evidence spine

## Strongest Evidence

- [`../../../development/open-source/README.md`](../../../development/open-source/README.md)
  - canonical release register, scope, and blocker history
- [`../../../evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md`](../../../evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md)
  - main open-source evolution arc
- [`../../../development/open-source/public-branch-cut-2026-05-26.md`](../../../development/open-source/public-branch-cut-2026-05-26.md)
  - sanitized public-branch cut and rewrite notes
- [`../../../.dvc/public-manifest.txt`](../../../.dvc/public-manifest.txt)
  - manifest-scoped public restore boundary
- workspace-local companion surface `personal/quine-lab-public-artifacts/README.md`
  - public artifact-serving surface and operational boundary

## Blocking Gaps

- the release work is indexed, but the blog arc is not yet compressed into a
  clean set of pain/solution moves
- the strongest commits still need a shortlist grouped by narrative role rather
  than by raw chronology
- the draft needs a clean boundary between Quine-specific mechanics and the
  reusable method for other research repos

## Next Useful Work

- group the open-source rollout into five to eight pain/solution arcs with
  concrete owner surfaces and commits
- select the minimum commit set that makes each arc legible without turning the
  post into a release log
- decide how much of the essay should read like a first-person chronicle versus
  an operator's guide
- extract the release lessons that remain true outside Quine

## Unsafe Overclaims

- do not imply the release was only a secret-scrubbing exercise
- do not collapse private retention and public publication into one DVC story
- do not present the public branch as proof that every retained artifact is now
  public
- do not let this case study silently replace the actual release protocol
  surfaces it depends on
