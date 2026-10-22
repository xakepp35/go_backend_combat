#!/usr/bin/env python3
"""Run the two-service core, check HTTP/jobs and restart persistence in isolation."""
import json
import os
from pathlib import Path
import socket
import time
import urllib.error
import urllib.request
import uuid
import coursecheck
ROOT=Path(__file__).resolve().parents[1]
def port():
    with socket.socket() as s:s.bind(('127.0.0.1',0));return s.getsockname()[1]
def wait(predicate,description,seconds=30):
    deadline=time.monotonic()+seconds
    while time.monotonic()<deadline:
        try:
            if predicate():return
        except (urllib.error.URLError,ConnectionError,TimeoutError):pass
        time.sleep(.1)
    raise RuntimeError('timeout: '+description)
def main():
    env=os.environ.copy();env.update(POSTGRES_PORT=str(port()),HTTP_PORT=str(port()))
    name='combat-core-smoke-'+uuid.uuid4().hex[:10]
    cmd=['docker','compose','-p',name,'-f',str(ROOT/'deploy/docker-compose.yml')]
    def compose(*args,timeout=120):return coursecheck.run(cmd+list(args),env=env,timeout=timeout)
    def request(method,path,value=None,want=200):
        data=None if value is None else json.dumps(value).encode()
        req=urllib.request.Request('http://127.0.0.1:'+env['HTTP_PORT']+path,data=data,method=method,headers={'Content-Type':'application/json'})
        try:r=urllib.request.urlopen(req,timeout=5)
        except urllib.error.HTTPError as err:r=err
        with r:
            body=r.read();status=r.status;rid=r.headers.get('X-Request-ID','')
        if status!=want:raise RuntimeError(f'{method} {path}: {status} != {want}: '+body.decode()[:500])
        if len(rid)!=32:raise RuntimeError('request ID missing')
        return json.loads(body) if body else None
    def finished(id):return request('GET','/jobs/'+id)['status']=='done'
    try:
        compose('up','-d','--build','--wait',timeout=600)
        request('GET','/healthz')
        value={'id':'a','title':'Купить чай','priority':1}
        request('POST','/tasks',value,201)
        request('POST','/tasks',dict(value,title='кофе'),409)
        if request('GET','/tasks/a')['title']!=value['title']:raise RuntimeError('duplicate changed task')
        request('POST','/tasks',{'id':'bad','title':' ','priority':1},400)
        job={'id':'j','text':'чай хлеб молоко'}
        request('POST','/jobs',job,202);wait(lambda:finished('j'),'word count')
        if request('GET','/jobs/j')['result']!='3':raise RuntimeError('word count result')
        request('POST','/jobs',job,200);request('POST','/jobs',dict(job,text='other'),409)
        compose('stop','-t','8','app')
        compose('exec','-T','postgres','psql','-U','jobflow','-d','jobflow','-v','ON_ERROR_STOP=1','-c',"INSERT INTO jobs(id,body,status,attempts) VALUES('interrupted','one two','running',1),('waiting','one','queued',0)")
        compose('up','-d','--wait','app')
        wait(lambda:finished('interrupted') and finished('waiting'),'restart recovery')
        if request('GET','/tasks/a')['title']!=value['title'] or request('GET','/jobs/j')['result']!='3':raise RuntimeError('restart lost data')
        logs=compose('logs','--no-color','app')
        if 'request_id=' not in logs or 'job_id=interrupted' not in logs:raise RuntimeError('correlation logs missing')
        print('PASS core Compose: CRUD/errors, request IDs, word count, idempotent acceptance, restart persistence/recovery, logs.')
    finally:compose('down','-v','--remove-orphans',timeout=90)
if __name__=='__main__':main()
