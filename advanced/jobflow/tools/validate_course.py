#!/usr/bin/env python3
"""Validate course structure, local Markdown links, fixtures and archived originals."""
import hashlib
import json
from pathlib import Path
import re
import sys
from urllib.parse import unquote
import render_course

ROOT=Path(__file__).resolve().parents[1]

def require(condition,message):
    if not condition:raise ValueError(message)

def main():
    rows=json.loads((ROOT/'course/lessons.json').read_text())
    require([r['n'] for r in rows]==list(range(1,25)),'lessons must be 01–24')
    for r in rows:
        n=r['n'];name=f'{n:02}_{r["slug"]}.md'
        require(r['duration_minutes']==50 and r['reserve_minutes']==10,f'{n}: timing')
        require(r['star_required'] is False and r['star'],f'{n}: optional star')
        require(r['prerequisite']==n-1,f'{n}: prerequisite')
        for folder in ['student','mentor','combat']:
            p=ROOT/folder/name
            require(p.is_file(),f'missing {p.relative_to(ROOT)}')
            text=p.read_text()
            require(len(text)>1600,f'{p}: insufficient content')
            require('★' in text,f'{p}: no star policy')
            require(r['invariant'] in text,f'{p}: invariant drift')
            renderer={'student':render_course.student,'mentor':render_course.mentor,'combat':render_course.homework}[folder]
            require(text==renderer(r).rstrip()+'\n',f'{p}: stale generated lesson; run make docs')
        require((ROOT/f'reference/{n:02}'/r['path']).is_file(),f'{n}: source anchor missing')
        module=(ROOT/f'reference/{n:02}/go.mod').read_text()
        require(re.search(r'^module jobflow$',module,re.M),f'{n}: not a standalone module')
        require((ROOT/f'checks/{n:02}_test.go').is_file(),f'{n}: check missing')
        if r['poster']:require((ROOT/'posters'/r['poster']).is_file(),f'{n}: poster missing')
        for migration,first in [('001_core.sql',7),('002_engine.sql',14),('003_receipts.sql',19)]:
            if n>=first:
                require((ROOT/f'reference/{n:02}/migrations'/migration).read_bytes()==(ROOT/'reference/24/migrations'/migration).read_bytes(),f'{n}: migration history drift')
    for folder in ['student','mentor','combat']:
        require(len(list((ROOT/folder).glob('[0-9][0-9]_*.md')))==24,folder+': expected 24 lessons')
    require((ROOT/'deploy/Dockerfile').read_bytes()==(ROOT/'reference/24/Dockerfile').read_bytes(),'Dockerfile drift')
    for p in ROOT.rglob('*.md'):
        rel=p.relative_to(ROOT)
        if rel.parts[0] in ['.git','tmp','.cache'] or (rel.parts[0]=='legacy' and len(rel.parts)>2):continue
        # Fenced code may contain illustrative Markdown. Only check real prose links.
        text=re.sub(r'```.*?```','',p.read_text(),flags=re.S)
        for m in re.finditer(r'(?<!!)\[[^\]\n]+\]\(([^)]+)\)',text):
            target=m.group(1).strip().split(' "',1)[0].strip('<>')
            if target.startswith(('#','http://','https://','mailto:','app:')):continue
            path=unquote(target.split('#',1)[0])
            require((p.parent/path).exists(),f'broken link in {rel}: {target}')
    archive=ROOT.parents[1]/'legacy/2026-10-19'
    manifest=json.loads((archive/'SHA256.json').read_text())
    require(len(manifest)==32,'archive expected 32 original files')
    for name,digest in manifest.items():
        require(hashlib.sha256((archive/name).read_bytes()).hexdigest()==digest,'archive changed: '+name)
    vectors=json.loads((ROOT/'checks/fixtures/frames.json').read_text())
    require(len(vectors)>=10,'not enough format vectors')
    empty=next(v for v in vectors if v['name']=='empty')
    raw=bytes.fromhex(empty['hex'])
    require(len(raw)==128 and raw[:4]==b'\x02\x00\x00\x00','empty header differs')
    require(raw[96:].hex()=='af1349b9f5f9a1a6a0404dea36dcc9499bcb25c9adc112b7cc9a93cae41f3262','BLAKE3 empty golden differs')
    print('PASS: 24 lessons, 72 documents, stage/source/check links, migration history, format fixtures, 32 archived originals.')
    return 0

if __name__=='__main__':
    try:sys.exit(main())
    except (ValueError,OSError,KeyError,StopIteration) as e:
        print('FAIL: '+str(e),file=sys.stderr);sys.exit(1)
