# Paper History DVC Policy, 2026-05-26

This records the owner-reviewed policy for the next Paper-specific public
history rewrite.

## Scope

- Branch scanned: `main`
- Branch tip scanned: `6e7e7c51f9a25a44b8c66605342bc1a2ba7c35d9`
- Target branch for rewrite work: disposable clone of `main`
- Retention target: preserve recoverability through private DVC where material
  is not suitable for public git history
- Publication target: do not add these private Paper pointers to
  [DVC Public Manifest](../../.dvc/public-manifest.txt) unless a later review
  explicitly promotes exact paths

The release rule remains default-open for owner-authored drafts, notes, and
unfinished research. This policy only moves sensitive, bulky, generated, or
third-party-heavy Paper artifacts out of public git history.

## Current Evidence

The current `Paper/` HEAD scan found:

| Surface | Current HEAD evidence | Release read |
|---------|-----------------------|--------------|
| `Paper/` total | 898 tracked files, 58,288,956 bytes | Small enough to review directly, but history has larger generated blobs. |
| `Paper/essays/_meta/external-sources/**/*.pdf` | 6 PDFs, 35,250,995 bytes | Third-party PDF payloads should be private DVC. |
| `Paper/archive/track-alife-presentations/materials/**` | 79 files, 4,829,828 bytes | Collaboration-adjacent presentation material should be private DVC. |
| `Paper/archive/**/*.pdf` | 2 PDFs, 3,295,963 bytes | Not secret, but PDFs should be private DVC to keep git lightweight. |
| `Paper/core/graphs/**` | 9 files, 2,593,488 bytes | Graph payloads should be private DVC. |
| `Paper/` history | 3,941 unique blobs, 154,453,297 summed blob bytes, 1,947 historical paths | History rewrite should handle historical generated outputs, not only HEAD. |

`gitleaks dir Paper` returned no findings for the current `Paper/` tree during
the initial audit. This policy is about sensitivity, size, recoverability, and
public artifact boundaries rather than a current secret-scanner hit.

## Owner Decisions

| Item | Path family | Decision |
|------|-------------|----------|
| A | `Paper/theory/views/ecology/plans/synthetic-file-system-autopoiesis.md` | Keep public. The internal marker is not treated as a secret for this release. |
| B | `Paper/**/_meta/external-sources/**/*.pdf` | Convert PDFs to private DVC pointers, including historical pre-rename locations such as `Paper/TRACK-blogs/_meta/external-sources/`. Keep source notes, indexes, and non-PDF metadata public unless a later scan finds a separate issue. |
| C | `Paper/archive/track-alife-presentations/materials/**` | Convert the full presentations directory to private DVC pointers. Do not promote to public DVC by default. |
| D | `Paper/archive/**/*.pdf` | Convert PDFs to private DVC pointers. Keep archive text and metadata public by default. |
| E5 | Historical `Paper/Quine:*` and `Paper/Quine_Autopoietic_Intelligence_via_Recursive_POSIX_Processes.pdf` output | Treat as legacy generated paper build/output material. Convert generated binaries and build outputs to private DVC or remove from public git history; keep only source-like markdown if later review decides it is worth preserving publicly. |
| E6 | Historical and current graph directories | Convert graph payload directories to private DVC pointers. |
| F | Private/lab refs and email-like prose references | Keep public for now. These are not release blockers by themselves. |
| G | Active Paper draft and theory directories | Keep public by default. |

## Rewrite Path Policy

The next Paper rewrite should use this path policy.

### Keep Public

- `Paper/theory/views/ecology/plans/synthetic-file-system-autopoiesis.md`
- active Paper draft, theory, registry, and source markdown directories
- non-PDF external-source notes and indexes under
  `Paper/essays/_meta/external-sources/`
- non-PDF archive notes and metadata under `Paper/archive/`
- private/lab refs and email-like prose references unless a
  later focused review identifies a concrete access or consent issue

### Convert To Private DVC

- `Paper/**/_meta/external-sources/**/*.pdf`
- `Paper/archive/track-alife-presentations/materials/**`
- `Paper/archive/**/*.pdf`
- graph payloads in current and historical locations, including:
  - `Paper/**/graphs/**`
  - historical `Paper/graphs/**`
  - historical `Paper/shared/graphs/**`
  - historical `Paper/sections/graphs/**`
  - historical `Paper/Quine:*/sections/graphs/**`
  - historical `Paper/Quine:*/arxiv_submission/graphs/**`
- generated legacy paper outputs, including:
  - `Paper/Quine_Autopoietic_Intelligence_via_Recursive_POSIX_Processes.pdf`
  - historical `Paper/Quine:*/output/**`
  - historical `Paper/Quine:*` binary outputs such as PDFs, DOCX, HTML, and
    generated TeX where they are build products rather than source edits

### Keep Out Of Public DVC

All private-DVC paths above stay out of
[DVC Public Manifest](../../.dvc/public-manifest.txt) until separately
reviewed and explicitly promoted. Public readers should not expect plain
`dvc pull` to restore these objects from `https://dvc-artifacts.kehao.me`.

## Rewrite Procedure

1. Create a disposable clone from the current `main` tip and record backup refs
   before rewriting.
2. Enumerate every commit reachable from `main` and classify matching Paper
   paths with the policy above.
3. For each affected commit, replace private-DVC path content with DVC pointer
   files while retaining the payload in the private DVC object store and backup
   remote.
4. Preserve public source-like markdown and metadata in git where the policy
   says keep-public.
5. Import the rewritten branch back as a local candidate, without pushing to
   public remotes.
6. Validate the candidate branch:
   - current-tree `gitleaks`
   - full-history `gitleaks git`
   - path scan proving no denied raw Paper payloads remain in git history
   - largest-blob scan proving Paper generated binaries no longer dominate git
     history
   - historical checkout plus private-DVC restore smoke for at least one commit
     that previously contained presentations, external PDFs, archive PDFs, and
     legacy `Paper/Quine:*` output
7. Keep the public DVC manifest unchanged unless exact public Paper pointers are
   separately approved.

## Open Decision

The only Paper-specific policy question still worth explicit owner review is
whether any legacy `Paper/Quine:*` markdown source history should be preserved
publicly. The conservative rewrite can keep source-like markdown public and
move only generated binary/build payloads to private DVC.
