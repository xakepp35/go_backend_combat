#!/usr/bin/env python3
"""Refresh marked navigation only. Handwritten prose is never generated."""
import json
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]
START='<!-- course-nav:start -->'
END='<!-- course-nav:end -->'
def table(rows):
    result=['| Неделя | Разговор | Ученику | Практика | Ментору |','|---|---|---|---|---|']
    for r in rows:
        name=f'{r["n"]:02}_{r["slug"]}.md'
        result.append(f'| {r["n"]:02} | {r["title"]} | [Урок](../student/{name}) | [ДЗ](../combat/{name}) | [План](../mentor/{name}) |')
    return '\n'.join(result)
def main():
    p=ROOT/'course/README.md'
    text=p.read_text()
    if text.count(START)!=1 or text.count(END)!=1:
        raise ValueError('course/README.md needs one explicit navigation block')
    before,tail=text.split(START,1)
    _,after=tail.split(END,1)
    rows=json.loads((ROOT/'course/lessons.json').read_text())
    p.write_text(before+START+'\n\n'+table(rows)+'\n\n'+END+after)
    print('Navigation refreshed; handwritten content preserved.')
if __name__=='__main__': main()
