#!/usr/bin/env python3
"""Verify two intentional mistakes are detected in temporary core submissions."""
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import coursecheck
ROOT=Path(__file__).resolve().parents[1]
def main():
    go=coursecheck.go_binary()
    with tempfile.TemporaryDirectory(prefix='combat-mutants-') as tmp:
        base=Path(tmp)
        validation=base/'validation';shutil.copytree(ROOT/'reference/07',validation)
        file=validation/'pkg/task.go';text=file.read_text();start=text.index('func (t Task) Validate() error {');end=text.index('\nfunc New(',start)
        text=text[:start]+'func (t Task) Validate() error { return nil }\n'+text[end:]
        # Remove the unused import as part of this deliberate behavioral mutation.
        text=text.replace('\n\t"strings"','').replace('\n\t"fmt"','')
        file.write_text(text)
        result=subprocess.run([sys.executable,str(ROOT/'tools/coursecheck.py'),'--lesson','7','--project',str(validation)],text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,env=os.environ.copy(),timeout=90)
        if result.returncode==0 or 'TestValidationAndCRUD' not in result.stdout:raise RuntimeError('Validation mutant was not rejected by a behavioral assertion:\n'+result.stdout[-3000:])
        print('PASS: validation bypass rejected by the expected behavioral test.')
        pool=base/'pool';shutil.copytree(ROOT/'reference/18',pool)
        file=pool/'internal/pool/pool.go';text=file.read_text();old='i < workers';assert old in text;text=text.replace(old,'i < workers+1',1);file.write_text(text)
        result=subprocess.run([go,'test','./internal/pool','-timeout=10s'],cwd=pool,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=30)
        if result.returncode==0 or 'третий исполнитель' not in result.stdout:raise RuntimeError('Worker mutant was not rejected by the expected behavioral test:\n'+result.stdout[-3000:])
        print('PASS: exceeding worker count rejected by the concurrency assertion.')
if __name__=='__main__':main()
