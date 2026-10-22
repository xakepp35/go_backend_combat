#!/usr/bin/env python3
"""Validate structure and links without enforcing a prose template."""
import hashlib
import json
from pathlib import Path
import re
import sys
from urllib.parse import unquote
from render_course import START, END, table
ROOT=Path(__file__).resolve().parents[1]
def require(ok,message):
    if not ok:raise ValueError(message)
def main():
    rows=json.loads((ROOT/'course/lessons.json').read_text())
    require([r['n'] for r in rows]==list(range(1,25)),'Нужны недели 01–24')
    nav=(ROOT/'course/README.md').read_text()
    require(nav.count(START)==1 and nav.count(END)==1,'Нужен один блок навигации')
    require(nav.split(START,1)[1].split(END,1)[0].strip()==table(rows),'Обнови навигацию: make docs')
    bank=(ROOT/'interview/README.md').read_text()
    numbers=[int(n) for n in re.findall(r'(?m)^(\d+)\. ',bank)]
    require(numbers==list(range(1,501)),'Банк должен сохранять все 500 номеров')
    original=(ROOT/'legacy/2026-10-19/interview/README.md').read_text().strip()
    require(original in bank,'Исходные вопросы/подсказки/ответы изменены или потеряны')
    require('LOGO.png' in (ROOT/'README.md').read_text() and (ROOT/'LOGO.png').is_file(),'Нужен исходный логотип в README')
    for r in rows:
        n=r['n'];name=f'{n:02}_{r["slug"]}.md'
        require(r['prerequisite']==n-1,f'{n}: неправильная предпосылка')
        require(re.fullmatch(r'\d+(?:,\d+)?–\d+(?:,\d+)?',r['homework_hours']) and r['star_required'] is False,f'{n}: нагрузка/★')
        for folder in ['student','mentor','combat']:
            p=ROOT/folder/name
            require(p.is_file(),f'Нет {p.relative_to(ROOT)}')
            require('TODO' not in p.read_text(),f'Незавершённый учебный текст: {p}')
        require('★' in (ROOT/'combat'/name).read_text(),f'{n}: нет добровольного расширения')
        require((ROOT/f'reference/{n:02}/go.mod').is_file(),f'{n}: нет самостоятельного эталона')
        require(any((ROOT/f'checks/{n:02}').rglob('*.go')),f'{n}: нет открытых проверок')
        require(all(q in numbers for q in r['questions']),f'{n}: неизвестный номер вопроса')
        require(f'week-{n:02}' in (ROOT/'interview/route.md').read_text(),f'{n}: нет маршрута интервью')
    for folder in ['student','mentor','combat']:
        require(len(list((ROOT/folder).glob('[0-9][0-9]_*.md')))==24,f'{folder}: ожидаются 24 документа')
    # Archived originals and external links are not rewritten or interpreted as course navigation.
    for p in ROOT.rglob('*.md'):
        rel=p.relative_to(ROOT)
        if rel.parts[0] in ['.git','legacy','tmp','.cache']:continue
        text=re.sub(r'```.*?```','',p.read_text(),flags=re.S)
        for m in re.finditer(r'!?\[[^\]\n]*\]\(([^)]+)\)',text):
            target=m.group(1).strip().split(' "',1)[0].strip('<>')
            if target.startswith(('#','https://','http://','mailto:','app:')):continue
            path=unquote(target.split('#',1)[0])
            require((p.parent/path).exists(),f'Сломанная ссылка {rel}: {target}')
    archive=ROOT/'legacy/2026-10-19';manifest=json.loads((archive/'SHA256.json').read_text())
    require(len(manifest)==32,'Ожидаются 32 архивных файла')
    for name,digest in manifest.items():require(hashlib.sha256((archive/name).read_bytes()).hexdigest()==digest,'Изменён архив: '+name)
    vectors=json.loads((ROOT/'advanced/jobflow/checks/fixtures/frames.json').read_text())
    empty=bytes.fromhex(next(v['hex'] for v in vectors if v['name']=='empty'))
    require(len(empty)==128 and empty[96:].hex()=='af1349b9f5f9a1a6a0404dea36dcc9499bcb25c9adc112b7cc9a93cae41f3262','Изменён бинарный вектор продолжения')
    print('PASS: 24 недели, 72 документа, 500 исходных вопросов, логотип, ссылки, архив и вектор протокола.')
if __name__=='__main__':
    try:main()
    except (ValueError,OSError,KeyError,StopIteration) as e:print('FAIL: '+str(e),file=sys.stderr);sys.exit(1)
