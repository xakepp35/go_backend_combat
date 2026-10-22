#!/usr/bin/env python3
"""Validate all 24 independent milestones against their cumulative open contracts."""
import contextlib
import io
import os
from pathlib import Path
import sys
import coursecheck

def main():
    root=Path(__file__).resolve().parents[1]
    try:
        with coursecheck.infrastructure(24,None) as (env,compose):
            os.environ.update({k:v for k,v in env.items() if k.startswith('JF_TEST_')})
            for n in range(1,25):
                capture=io.StringIO()
                try:
                    with contextlib.redirect_stdout(capture):
                        coursecheck.check(n,root/'reference'/f'{n:02}',compose[3],False)
                except Exception:
                    print(capture.getvalue()[-8000:]);raise
                print(f'PASS reference/{n:02}: mandatory checks 01–{n:02}',flush=True)
    except Exception as e:
        print('FAIL: '+str(e),file=sys.stderr);return 1
    return 0

if __name__=='__main__':sys.exit(main())
