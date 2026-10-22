#!/usr/bin/env python3
"""Check an isolated submission; core and advanced have independent lesson numbers."""
import argparse
import contextlib
import os
from pathlib import Path
import shutil
import socket
import subprocess
import sys
import tempfile
import uuid
ROOT=Path(__file__).resolve().parents[1]

def run(argv, *, cwd=None, env=None, data=None, timeout=240):
    p=subprocess.run(argv,cwd=cwd,env=env,input=data,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=timeout)
    if p.returncode:
        raise RuntimeError(f'{Path(argv[0]).name} failed ({p.returncode}):\n'+p.stdout[-12000:])
    return p.stdout

def go_binary():
    go=os.environ.get('GO') or shutil.which('go')
    if not go: raise RuntimeError('Нужен Go. См. docs/START_HERE.md или задай GO=/полный/путь/к/go.')
    return go

@contextlib.contextmanager
def infrastructure(lesson, database_url=None):
    env=os.environ.copy()
    if lesson<12:
        yield env; return
    provided=database_url or env.get('JF_TEST_DATABASE_URL')
    if provided:
        env['JF_TEST_DATABASE_URL']=provided
        yield env; return
    if not shutil.which('docker'):
        raise RuntimeError('SQL-этап требует PostgreSQL: установи Docker или задай JF_TEST_DATABASE_URL; пропуск не является PASS.')
    with socket.socket() as sock:
        sock.bind(('127.0.0.1',0)); port=sock.getsockname()[1]
    env['POSTGRES_PORT']=str(port)
    env['JF_TEST_DATABASE_URL']=f'postgres://jobflow:jobflow@127.0.0.1:{port}/jobflow?sslmode=disable'
    compose=['docker','compose','-p','combat-core-check-'+uuid.uuid4().hex[:10],'-f',str(ROOT/'deploy/docker-compose.yml')]
    try:
        run(compose+['up','-d','--wait','postgres'],env=env,timeout=180)
        yield env
    finally:
        run(compose+['down','-v','--remove-orphans'],env=env,timeout=90)

def sql_checks(lesson,project,env,go):
    parts=[]
    for p in sorted((project/'migrations').glob('*.sql')):
        sql=p.read_text()
        if lesson>=14:
            # Only the Up DDL belongs in this rolled-back SQL assertion helper.
            # Actual Goose execution/repeated startup is covered by existing Go checks.
            # Core fixtures use ordinary transactional SQL, not Goose Go migrations.
            if '-- +goose Up' not in sql:
                raise RuntimeError('Нужна секция Goose Up: '+p.name)
            sql=sql.split('-- +goose Up',1)[1].split('-- +goose Down',1)[0]
        parts.append(sql)
    if not parts: raise RuntimeError('Нет миграций: SQL-этап не может пройти без схемы.')
    for p in sorted((ROOT/'checks'/f'{lesson:02}').glob('*.sql')): parts.append(p.read_text())
    # Each invocation creates a schema inside a rolled-back transaction.
    run([go,'run','.'],cwd=ROOT/'tools/sqlcheck',env=env,data='\n'.join(parts))

def check(lesson,source,env):
    go=go_binary();source=Path(source).resolve()
    if not (source/'go.mod').is_file(): raise RuntimeError('PROJECT должен указывать на папку с go.mod.')
    overlay=ROOT/'checks'/f'{lesson:02}'
    if not overlay.is_dir(): raise RuntimeError('Нет открытых проверок для этой недели.')
    with tempfile.TemporaryDirectory(prefix='combat-core-') as temp:
        project=Path(temp)/'project'
        shutil.copytree(source,project,ignore=shutil.ignore_patterns('.git','.cache','bin','coverage.out','__pycache__'))
        for p in overlay.rglob('*.go'):
            target=project/p.relative_to(overlay);target.parent.mkdir(parents=True,exist_ok=True)
            if target.exists():raise RuntimeError('Проверка не может перезаписать файл проекта: '+str(p.relative_to(overlay)))
            shutil.copy2(p,target)
        testenv=env.copy()
        # Nested go commands in early exercise checks use the same chosen toolchain.
        testenv['PATH']=str(Path(go).resolve().parent)+os.pathsep+testenv.get('PATH','')
        run([go,'vet','./...'],cwd=project,env=testenv)
        args=[go,'test','-count=1','-timeout=90s']
        if lesson>=17:args+=['-race']
        output=run(args+['./...'],cwd=project,env=testenv)
        if lesson>=12:sql_checks(lesson,project,testenv,go)
        for name,first in [('PORTFOLIO.md',23),('CAREER.md',23),('NEXT_STEPS.md',24)]:
            if lesson>=first and not (project/name).is_file():raise RuntimeError('Нужен '+name+'; содержание обсуждается на ревью.')
        if lesson>=21:
            for name in ['Dockerfile','compose.yaml','.env.example']:
                if not (project/name).is_file():raise RuntimeError('Нужен '+name)
    print(f'PASS core {lesson:02}: перечисленные открытые сценарии. Объяснение и ★ обсуждаются отдельно.')
    return output

def main():
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--track',choices=['core','advanced'],default='core')
    p.add_argument('--lesson',type=int,required=True,choices=range(1,25))
    p.add_argument('--project',type=Path,required=True)
    p.add_argument('--database-url',help='Dedicated test PostgreSQL; schemas are created and removed')
    p.add_argument('--infra-project',help='Advanced track only')
    p.add_argument('--fuzz',action='store_true',help='Advanced track only')
    a=p.parse_args()
    if a.track=='advanced':
        args=[sys.executable,str(ROOT/'advanced/jobflow/tools/coursecheck.py'),'--lesson',str(a.lesson),'--project',str(a.project.resolve())]
        if a.infra_project:args+=['--infra-project',a.infra_project]
        if a.fuzz:args+=['--fuzz']
        return subprocess.call(args)
    if a.infra_project or a.fuzz: p.error('--infra-project/--fuzz относятся к advanced')
    try:
        with infrastructure(a.lesson,a.database_url) as env:check(a.lesson,a.project,env)
        return 0
    except (OSError,RuntimeError,subprocess.TimeoutExpired) as err:
        print('FAIL: '+str(err),file=sys.stderr);return 1
if __name__=='__main__':sys.exit(main())
