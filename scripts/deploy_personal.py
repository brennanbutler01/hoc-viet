"""Deploy only tracked source to Brennan's personal Vercel project."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

root = Path(__file__).resolve().parent.parent
project_path = root / '.vercel/project.json'
project = json.loads(project_path.read_text())
if (project.get('orgId') != 'team_cwjFWlUVzkCIYgVepQ6Glc71'
        or project.get('projectName') != 'hoc-viet-demo'):
    raise SystemExit('Deployment requires the personal hoc-viet-demo Vercel project.')
files = subprocess.check_output(['git', 'ls-files', '-z'], cwd=root).decode().split('\0')
with tempfile.TemporaryDirectory(prefix='hoc-viet-personal-') as directory:
    target = Path(directory)
    for name in filter(None, files):
        source = root / name
        if not source.is_file():
            raise SystemExit('Stage deletions before deploying.')
        destination = target / name
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source, destination)
    (target / '.vercel').mkdir()
    shutil.copy2(project_path, target / '.vercel/project.json')
    subprocess.run([
        'npm', 'exec', '--yes', '--package=vercel@59.15.0', '--',
        'vercel', 'deploy', '--prod', '--yes', '--scope', 'brennanbutler01s-projects',
        '--env', 'PORT=8080',
    ], cwd=target, check=True)
