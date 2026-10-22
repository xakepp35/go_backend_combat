#!/usr/bin/env python3
"""Run the course's open contract checks on an isolated copy of a student project."""
import argparse
import contextlib
import os
from pathlib import Path
import re
import shutil
import socket
import subprocess
import sys
import tempfile
import uuid

ROOT = Path(__file__).resolve().parents[1]

def run(argv, *, cwd=None, env=None, data=None, timeout=300):
    p = subprocess.run(argv, cwd=cwd, env=env, input=data, text=True,
                       stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=timeout)
    if p.returncode:
        print(p.stdout[-16000:], file=sys.stderr)
        raise RuntimeError(f'{Path(argv[0]).name} failed (exit {p.returncode})')
    return p.stdout

def port():
    with socket.socket() as s:
        s.bind(('127.0.0.1', 0))
        return s.getsockname()[1]

@contextlib.contextmanager
def infrastructure(lesson, external):
    env = os.environ.copy()
    if lesson < 7:
        yield env, None
        return
    compose = ['docker', 'compose', '-p', external or ('jf-check-' + uuid.uuid4().hex[:12]),
               '-f', str(ROOT/'deploy/docker-compose.yml')]
    if external:
        for key in ['JF_TEST_DATABASE_URL'] + (['JF_TEST_REDIS_ADDR'] if lesson >= 15 else []) + (['JF_TEST_KAFKA_BROKERS'] if lesson >= 16 else []):
            if not env.get(key):
                raise RuntimeError(f'{key} is required with --infra-project')
    else:
        for key in ['POSTGRES_PORT', 'REDIS_PORT', 'KAFKA_PORT']:
            env[key] = str(port())
        env['JF_TEST_DATABASE_URL'] = 'postgres://jobflow:jobflow@127.0.0.1:' + env['POSTGRES_PORT'] + '/jobflow?sslmode=disable'
        env['JF_TEST_REDIS_ADDR'] = '127.0.0.1:' + env['REDIS_PORT']
        env['JF_TEST_KAFKA_BROKERS'] = '127.0.0.1:' + env['KAFKA_PORT']
    services = ['postgres'] + (['redis'] if lesson >= 15 else []) + (['redpanda'] if lesson >= 16 else [])
    try:
        if not external:
            print('Starting isolated infrastructure: ' + ', '.join(services), flush=True)
            run(compose + ['up', '-d', '--wait'] + services, env=env, timeout=300)
        yield env, compose
    finally:
        if not external:
            run(compose + ['down', '-v', '--remove-orphans'], env=env, timeout=90)

def check_sql(project, compose, env):
    sql = (project/'migrations/001_core.sql').read_text()
    schema = 'check_' + uuid.uuid4().hex
    prefix = f'BEGIN; CREATE SCHEMA {schema}; SET LOCAL search_path TO {schema};\n'
    probe = """
INSERT INTO drafts(id,doc) VALUES('a','{}');
DO $$ BEGIN
  BEGIN
    INSERT INTO drafts(id,doc) VALUES('a','{}');
    RAISE EXCEPTION 'duplicate primary key accepted';
  EXCEPTION WHEN unique_violation THEN NULL;
  END;
END $$;
ROLLBACK;
"""
    run(compose + ['exec', '-T', 'postgres', 'psql', '-U', 'jobflow', '-d', 'jobflow',
                   '-v', 'ON_ERROR_STOP=1'], env=env, data=prefix+sql+probe)

def check(lesson, source, external, fuzz):
    go = os.environ.get('GO') or shutil.which('go')
    if not go:
        raise RuntimeError('Go 1.27.2 is required. Install Go or set GO=/absolute/path/to/go.')
    source = source.resolve()
    if not source.is_dir() or not (source/'go.mod').is_file():
        raise RuntimeError('PROJECT must point to a standalone JobFlow module with go.mod')
    if not re.search(r'^module\s+jobflow\s*$', (source/'go.mod').read_text(), re.M):
        raise RuntimeError('The course module contract is: module jobflow')
    with tempfile.TemporaryDirectory(prefix='jobflow-check-') as temp:
        project = Path(temp)/'project'
        shutil.copytree(source, project, ignore=shutil.ignore_patterns('.git','.cache','bin','coverage.out','coursecheck'))
        target = project/'coursecheck'
        target.mkdir()
        for n in range(1,lesson+1):
            shutil.copy2(ROOT/'checks'/f'{n:02}_test.go',target)
        for threshold,name in [(8,'db'),(13,'fixtures'),(14,'engine'),(16,'kafka')]:
            if lesson >= threshold:
                shutil.copy2(ROOT/'checks/helpers'/f'{name}_test.go',target)
        if lesson >= 13:
            shutil.copytree(ROOT/'checks/fixtures',target/'fixtures')
        with infrastructure(lesson,external) as (env,compose):
            if lesson >= 7:
                check_sql(project,compose,env)
            run([go,'vet','./...'],cwd=project,env=env)
            args=[go,'test','-count=1','-timeout=150s']
            if lesson >= 11:
                args.append('-race')
            coverage=Path(temp)/'coverage.out'
            if lesson >= 13:
                args += ['-coverpkg=./internal/protocol','-coverprofile='+str(coverage)]
            output=run(args+['./...'],cwd=project,env=env,timeout=240)
            print(output.strip())
            if lesson >= 13:
                report=run([go,'tool','cover','-func='+str(coverage)],cwd=project,env=env)
                total=float(re.search(r'total:.*?([\d.]+)%',report).group(1))
                if total < 80:
                    raise RuntimeError(f'Protocol coverage {total:.1f}% is below 80%')
                print(f'Protocol coverage: {total:.1f}%')
                if fuzz:
                    print(run([go,'test','./coursecheck','-run=^$','-fuzz=FuzzFrame','-fuzztime=3s','-parallel=2'],cwd=project,env=env).strip())
    print(f'PASS: lessons 01–{lesson:02}; mandatory checks only. ★ is optional.')

def main():
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--lesson',type=int,required=True,choices=range(1,25))
    p.add_argument('--project',type=Path,required=True)
    p.add_argument('--infra-project',help='Explicitly reuse this course Compose project and JF_TEST_* environment; never remove it')
    p.add_argument('--fuzz',action='store_true')
    a=p.parse_args()
    try:
        check(a.lesson,a.project,a.infra_project,a.fuzz)
    except (RuntimeError,OSError,subprocess.TimeoutExpired) as exc:
        print('FAIL: '+str(exc),file=sys.stderr)
        return 1
    return 0

if __name__=='__main__':
    sys.exit(main())
