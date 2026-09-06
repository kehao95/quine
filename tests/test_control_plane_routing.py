#!/usr/bin/env python3
"""Check the real hook domains and active frontmatter dispatch without model calls."""
import re
import shutil
import subprocess
import tempfile
import tomllib
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
config = tomllib.loads((ROOT / 'prek.toml').read_text())
hooks = {h['id']: h for repo in config['repos'] for h in repo['hooks']}


def reaches(hook, path):
    return bool(re.search(hook['files'], path)) and not re.search(hook.get('exclude', r'(?!)'), path)


aux = {'lib', '_lib', 'analysis', 'notes', 'profiles', 'records', 'assets', 'setup', 'seedsets'}
experiments = [p.relative_to(ROOT).as_posix() for p in ROOT.glob('experiments/active/*/*/README.md')
               if p.parent.name not in aux]
papers = [p.relative_to(ROOT).as_posix() for pattern in
          ('Paper/papers/*/README.md', 'Paper/essays/papers/*/README.md') for p in ROOT.glob(pattern)]
assert experiments and papers, 'current input domains must be nonempty'
exp_hooks = ('check-experiment-paper-feeds', 'check-paper-dossier-links-from-experiments',
             'check-paper-evidence-backbone-from-experiments')
paper_hooks = ('check-paper-dossier-links', 'check-paper-evidence-backbone')
for name in (*exp_hooks, *paper_hooks, 'check-control-plane-frontmatter'):
    current = experiments if name in exp_hooks else papers
    if name == 'check-control-plane-frontmatter':
        current = experiments + papers
    for path in current:
        assert reaches(hooks[name], path), (name, 'missed current owner', path)
    for path in ('experiments/archive/family/exp/README.md',
                 'experiments/active/family/exp/runs/run/README.md',
                 'Paper/archive/old/README.md', 'public/paper/README.md'):
        assert not reaches(hooks[name], path), (name, 'captured historical payload', path)
for name in exp_hooks:
    for folder in aux:
        assert not reaches(hooks[name], f'experiments/active/family/{folder}/README.md'), name
for name in exp_hooks[1:]:
    assert hooks[name].get('pass_filenames') is False, 'reverse links need the dossier scan'
    assert reaches(hooks[name], 'experiments/active/family/exp/DESIGN.md'), name

for folder in ('analysis', 'plans', 'assets', 'notes', 'workspace', 'data', 'figures', 'support', 'sections'):
    assert not reaches(hooks['check-control-plane-frontmatter'], f'Paper/papers/paper/{folder}/README.md')
assert reaches(hooks['check-control-plane-frontmatter'], 'Paper/papers/paper/venue/README.md')
assert reaches(hooks['check-experiment-subdir-index'], 'experiments/active/family/analysis/summary.md')
assert not reaches(hooks['check-experiment-subdir-index'], 'experiments/active/family/exp/runs/run/analysis/summary.md')

# Exercise the actual validators after the hook decision; green must mean scanned.
for script, path, counter in (
    ('check-experiment-paper-feeds.sh', experiments[0], 'experiments'),
    ('check-paper-dossier-links.sh', papers[0], 'dossiers'),
    ('check-paper-evidence-backbone.sh', papers[0], 'dossiers'),
):
    result = subprocess.run(['bash', str(ROOT / 'scripts' / script), '--strict', path],
                            cwd=ROOT, capture_output=True, text=True, check=True)
    assert re.search(rf'\b{counter}=1\b', result.stdout), result.stdout

# Temporary input exercises both discovery and dispatch, including rejection.
script = 'scripts/check-control-plane-frontmatter.sh'
with tempfile.TemporaryDirectory(prefix='quine-routing-') as scratch:
    root = Path(scratch)
    for rel in (script, 'scripts/lib/paper_dossier_labels.py'):
        (root / rel).parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(ROOT / rel, root / rel)
    rel = 'experiments/active/example/probe/README.md'
    path = root / rel
    path.parent.mkdir(parents=True)
    valid = '---\nsurface_kind: experiment\nfamily: example\nexperiment_id: active/example/probe\nstatus: design-only\n---\n'
    for content, expected in ((valid, 0), (valid.replace('family: example', 'family: wrong'), 1)):
        path.write_text(content)
        for args in ([], [rel]):
            result = subprocess.run(['bash', str(root / script), '--strict', *args],
                                    cwd=root, capture_output=True, text=True)
            assert result.returncode == expected, (args, result.stdout, result.stderr)
            assert 'docs_scanned=1' in result.stdout, result.stdout
            if expected:
                assert 'MISMATCH experiment family' in result.stdout, result.stdout
print(f'Control-plane routing: {len(experiments)} experiment owners, {len(papers)} paper homes; positive and negative checks passed')
