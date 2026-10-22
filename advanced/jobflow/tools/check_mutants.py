#!/usr/bin/env python3
"""Prove mandatory checks reject two deliberately broken temporary submissions."""
import argparse
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import coursecheck

ROOT=Path(__file__).resolve().parents[1]

def main():
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--infra-project')
    a=p.parse_args()
    try:
        with coursecheck.infrastructure(16,a.infra_project) as (env,compose):
            for name,lesson,stage,file,old,new,reason in [
                ('validation-bypass',4,4,'internal/memory/repository.go',
                 'if err := t.Validate(); err != nil {\n\t\treturn err\n\t}',
                 '// deliberately bypass validation', 'invalid not recognized'),
                ('outbox-confirmed-before-publish',16,16,'internal/bus/kafka.go',
                 'for _, e := range events {',
                 'for _, e := range events {\nif err = s.Delivered(ctx,e.ID); err != nil { return err }', 'intent lost'),
            ]:
                with tempfile.TemporaryDirectory(prefix='jobflow-mutant-') as temp:
                    project=Path(temp)/'project'
                    shutil.copytree(ROOT/'reference'/f'{stage:02}',project)
                    f=project/file;source=f.read_text()
                    if old not in source:raise RuntimeError('mutation location changed: '+name)
                    f.write_text(source.replace(old,new,1))
                    cmd=[sys.executable,str(ROOT/'tools/coursecheck.py'),'--lesson',str(lesson),'--project',str(project)]
                    if lesson>=7:cmd+=['--infra-project',compose[3]]
                    result=subprocess.run(cmd,env=env,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=240)
                    if result.returncode==0 or reason not in result.stdout:
                        print(result.stdout[-5000:]);raise RuntimeError('mutant not rejected by expected behavioral check: '+name)
                    print('PASS rejected '+name+' with nonzero exit and expected behavioral failure',flush=True)
    except Exception as e:
        print('FAIL: '+str(e),file=sys.stderr);return 1
    return 0

if __name__=='__main__':sys.exit(main())
