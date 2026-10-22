#!/usr/bin/env python3
"""Compile complete package-main examples in isolated milestone copies."""
import json
from pathlib import Path
import re
import shutil
import tempfile
import coursecheck
ROOT=Path(__file__).resolve().parents[1]
def main():
    count=0;go=coursecheck.go_binary()
    for row in json.loads((ROOT/'course/lessons.json').read_text()):
        n=row['n'];text=(ROOT/'student'/f'{n:02}_{row["slug"]}.md').read_text()
        for code in re.findall(r'```go\n(.*?)```',text,re.S):
            if not code.lstrip().startswith('package main') or not re.search(r'func\s+main\s*\(',code):continue
            with tempfile.TemporaryDirectory(prefix='combat-example-') as tmp:
                project=Path(tmp)/'project';shutil.copytree(ROOT/'reference'/f'{n:02}',project)
                folder=project/'doc_example';folder.mkdir();(folder/'main.go').write_text(code)
                coursecheck.run([go,'build','-o',str(Path(tmp)/'example'),'./doc_example'],cwd=project)
                count+=1
    print(f'PASS: {count} самостоятельных Go-примеров из глав собираются.')
if __name__=='__main__':main()
