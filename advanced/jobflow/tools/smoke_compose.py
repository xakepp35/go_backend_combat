#!/usr/bin/env python3
"""Build an isolated four-service system; verify DAG, Kafka/Redis recovery and ret."""
import json
import os
from pathlib import Path
import subprocess
import sys
import time
import urllib.request
import uuid
import coursecheck

ROOT=Path(__file__).resolve().parents[1]

def wait_for(predicate,description,seconds=35):
    deadline=time.monotonic()+seconds
    while time.monotonic()<deadline:
        if predicate():return
        time.sleep(.2)
    raise RuntimeError('timeout: '+description)

def main():
    name='jf-smoke-'+uuid.uuid4().hex[:10]
    env=os.environ.copy()
    env.update(COMPOSE_PROJECT_NAME=name,JOBFLOW_IMAGE_TAG=name)
    for key in ['POSTGRES_PORT','REDIS_PORT','KAFKA_PORT','HTTP_PORT']:
        env[key]=str(coursecheck.port())
    command=['docker','compose','-p',name,'-f',str(ROOT/'deploy/docker-compose.yml'),'--profile','app']
    active=[]
    def compose(*args,timeout=180):
        return coursecheck.run(command+list(args),env=env,timeout=timeout)
    def query(sql):
        return compose('exec','-T','postgres','psql','-U','jobflow','-d','jobflow','-Atc',sql).strip()
    def jobs():
        with urllib.request.urlopen('http://127.0.0.1:'+env['HTTP_PORT']+'/internal/tasks',timeout=3) as r:
            return json.load(r)
    def start_demo():
        p=subprocess.Popen(command+['exec','-T','pub','demo'],env=env,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
        active.append(p);return p
    def finish_demo(p):
        output,_=p.communicate(timeout=50)
        if p.returncode or 'SUCCEEDED node-3' not in output:
            raise RuntimeError('DAG failed: '+output[-3000:])
        print(output.strip(),flush=True)
    try:
        print('Building isolated '+name,flush=True)
        compose('up','-d','--build','--wait',timeout=600)
        finish_demo(start_demo())
        print('PASS: Compose diamond DAG',flush=True)
        compose('stop','-t','3','redpanda')
        p=start_demo()
        wait_for(lambda: len(jobs())>=8,'accept workflow while Kafka is down')
        pending=int(query('SELECT count(*) FROM events WHERE NOT delivered'))
        if pending<4:raise RuntimeError('Kafka outage lost durable events')
        compose('up','-d','--wait','redpanda')
        finish_demo(p)
        wait_for(lambda: int(query('SELECT count(*) FROM events WHERE NOT delivered'))==0,'relay backlog drain')
        print('PASS: Kafka outage preserves Outbox; restored delivery completes DAG',flush=True)
        compose('stop','-t','3','act')
        p=start_demo()
        wait_for(lambda: len(jobs())>=12,'third workflow accepted')
        compose('exec','-T','redis','redis-cli','FLUSHDB')
        compose('up','-d','act')
        finish_demo(p)
        wait_for(lambda: int(query('SELECT count(*) FROM events WHERE NOT delivered'))==0,'final outbox drain')
        print('PASS: complete Redis loss; sub rebuild and act execute accepted DAG',flush=True)
        # This schema belongs only to the unique smoke project. Artificial clock fixture avoids a 1h wait.
        query("UPDATE jobs SET doc=jsonb_set(doc,'{updated_at}',to_jsonb('2000-01-01T00:00:00Z'::text)) WHERE doc->>'state'='SUCCEEDED'")
        compose('restart','ret')
        wait_for(lambda: int(query('SELECT count(*) FROM jobs WHERE body_key IS NOT NULL OR result_key IS NOT NULL'))==0,'ret removes eligible payload references')
        if int(query('SELECT count(*) FROM payloads'))!=0:raise RuntimeError('ret left unreferenced CAS')
        states=jobs()
        if len(states)!=12 or any(t['state']!='SUCCEEDED' for t in states):raise RuntimeError('GC destroyed terminal identity')
        logs=compose('logs','--no-color','ret')
        if 'task reclaimed' not in logs or states[-1]['id'] not in logs:raise RuntimeError('ret correlation log missing')
        print('PASS: ret cleanup, correlation IDs and retained terminal identity',flush=True)
        print('PASS: isolated full Compose smoke',flush=True)
    except Exception as e:
        print('FAIL: '+str(e),file=sys.stderr)
        try:print(compose('logs','--tail','35','pub','sub','act','ret')[-10000:],file=sys.stderr)
        except Exception:pass
        return 1
    finally:
        for p in active:
            if p.poll() is None:p.terminate()
        compose('down','-v','--remove-orphans',timeout=90)
    return 0

if __name__=='__main__':sys.exit(main())
