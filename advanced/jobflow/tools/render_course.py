#!/usr/bin/env python3
"""Render synchronized Russian course navigation, lessons and mentor guides."""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
ROWS = json.loads((ROOT / 'course/lessons.json').read_text())
MODULES = [
    'Старт и язык Go', 'Границы приложения и данные',
    'Транзакции, конкурентность и заголовок', 'Бинарный поток и доставка',
    'Проекции, сигналы и исполнение', 'Восстановление и эксплуатация',
]
TOOLS = [
    'Go toolchain, терминал, Git, CI', 'Типы, функции, []byte, go.mod',
    'map, slice, методы, указатели', 'error, errors.Is, testing',
    'Пакеты, интерфейсы потребителя', 'net/http, encoding/json, httptest',
    'psql, SQL, миграции', 'pgxpool, context, SQL-адаптер',
    'Транзакция, rollback, local intent', 'goroutine, channel, sync.WaitGroup',
    'context, select, race detector', 'encoding/binary, LittleEndian',
    'io.Reader, BLAKE3, fuzz', 'DFS, UUIDv7, SQL-транзакция',
    'Redis, CAS, восстановление кэша', 'Kafka API, Redpanda, Outbox',
    'sub, Redis Hash, version, replay', 'Redis Streams, Consumer Groups, XAUTOCLAIM',
    'act, claim, lease, token, registry', 'nbf, exp, maintenance, clock',
    'Retry, receipt, effect ledger, quota', 'ret, retention, ссылки CAS',
    'Docker Compose, slog, metrics, E2E', 'Демонстрация, отказ, аргументация гарантий',
]


def put(path, text):
    p = ROOT / path
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(text.rstrip() + '\n')


def filename(row):
    return f'{row["n"]:02}_{row["slug"]}.md'


def link(row, folder):
    return f'[{row["n"]:02}. {row["title"]}](../{folder}/{filename(row)})'


def nav(row, folder):
    n = row['n']
    parts = ['[Программа](../course/README.md)']
    if n > 1:
        parts.append(f'← {link(ROWS[n-2], folder)}')
    if n < 24:
        parts.append(f'{link(ROWS[n], folder)} →')
    return ' · '.join(parts)


def resources(row):
    n = row['n']
    out = f'- [Исходник этапа {n:02}](../reference/{n:02}/{row["path"]})\n'
    out += f'- [Открытые проверки {n:02}](../checks/{n:02}_test.go)\n'
    out += f'- [Официальный материал]({row["url"]})\n'
    out += '- [Правила работы и сдачи](../docs/WORKFLOW.md)\n'
    if n <= 5:
        out += '- [Go: опорный материал первых недель](../docs/GO_PRIMER.md)\n'
    if row['poster']:
        out += f'- [Плакат для разбора](../posters/{row["poster"]})\n'
    if n >= 12:
        out += '- [Контракт бинарного формата](../docs/FRAMING.md)\n'
    if n >= 14:
        out += '- [Состояния, API и гарантии](../docs/PROJECT.md)\n'
    return out


def numbered(items):
    return '\n'.join(f'{i}. {text}' for i, text in enumerate(items, 1))


def student(row):
    n = row['n']
    before = 'Рабочий starter; знание переменных, условий, циклов и функций на другом языке.' if n == 1 else f'Обязательные проверки этапа {n-1:02} проходят в вашем репозитории. Повторите его API и ошибку, которую он уже обнаруживает.'
    lang = row.get('code_language', 'go')
    theory = '\n\n'.join(row['theory'])
    return f'''# {n:02}. {row['title']}

{nav(row, 'student')}

**Место в проекте:** месяц {(n-1)//4+1}, {MODULES[(n-1)//4]}. Занятие: 50 минут + 10 минут резерва. ДЗ: 8–12 часов, включая подготовку; ★ добровольна.

## Задача недели

{row['problem']}

**Инвариант:** {row['invariant']}

**Перед встречей:** {before}

Инструменты недели: {TOOLS[n-1]}.

## Разбор

{theory}

## Опорный фрагмент

Фрагмент показывает основную операцию. Имена из окружающего кода смотрите в исходнике этапа: это не отдельная готовая программа.

```{lang}
{row['code']}
```

## Практика на встрече

1. Назовите вход, результат и защищаемое правило для операции недели.
2. Реализуйте первый обязательный пункт [домашнего задания](../combat/{filename(row)}). До запуска предскажите результат.
3. Воспроизведите ошибку: {row['bug']}
4. Исправьте причину и объясните, какое наблюдение теперь подтверждает правило.

## Самопроверка

{row['question']}

<details>
<summary>Разбор ответа</summary>

{row['answer']}

</details>

Проверьте следующие ситуации до сдачи:

{numbered(row['cases'])}

## Самостоятельная работа

Продолжайте свой JobFlow по [ДЗ {n:02}](../combat/{filename(row)}). Оценка опирается на обязательный контракт; творческая часть не нужна следующему этапу. Открывайте эталон после собственной попытки и сравнивайте причины решений, а затем меняйте решение под новое условие.

## Материалы

{resources(row)}'''


def mentor(row):
    n = row['n']
    if n == 24:
        schedule = [
            ('0–10', 'Запуск проекта и объяснение DAG; ученик показывает итог, IDs и состояния.'),
            ('10–30', 'Выбрать отказ из ДЗ: остановка процесса, повтор или потеря проекции. Восстановить и доказать результат.'),
            ('30–45', 'Обсудить token, ledger, outbox и границы гарантий на исходном коде ученика.'),
            ('45–50', 'Зафиксировать обязательные критерии, ограничения и следующий самостоятельный эксперимент.'),
        ]
    else:
        schedule = [
            ('0–5', 'Восстановить входной этап: ученик запускает прошлую проверку и называет результат.'),
            ('5–15', f'Разобрать вопрос: {row["question"]}'),
            ('15–35', f'Написать законченный фрагмент: {row["homework"][0]}'),
            ('35–45', f'Вызвать ошибку: {row["bug"]} Исправить и повторить наблюдение.'),
            ('45–50', 'Ученик объясняет инвариант и выбирает следующий обязательный шаг ДЗ.'),
        ]
    timing = '\n'.join(f'| {time} | {action} |' for time, action in schedule)
    return f'''# Ментору: {n:02}. {row['title']}

{nav(row, 'mentor')}

## Цель и вход

Ученик должен уметь объяснить и защитить правило: **{row['invariant']}**

До встречи: прочитать [конспект](../student/{filename(row)}), проверить окружение и подготовить исходник [эталона {n:02}](../reference/{n:02}/README.md). {'Вход — starter и опыт другого языка.' if n == 1 else f'Вход — PASS {n-1:02}; не тратьте занятие на переписывание всего предыдущего этапа.'}

Инструмент: {TOOLS[n-1]}. Результат встречи — один работающий фрагмент в проекте ученика.

## План 50 + 10 минут

| Минуты | Действие и наблюдение |
|---|---|
{timing}
| 50–60 | Резерв на вопросы и восстановление отставшего шага; новая обязательная тема не вводится. |

## Контроль понимания

**Вопрос:** {row['question']}

**Ожидаемая мысль:** {row['answer']}

**Типичная ошибка:** {row['bug']}

Попросите ученика показать ошибку в данных или поведении, а затем место, где действует защита. Успешная сборка сама по себе не заменяет эту проверку.

## Приёмка недели

{numbered(row['cases'])}

Автоматическая часть: `make check LESSON={n:02} PROJECT=/absolute/path/to/jobflow` из корня курса. Проверки накопительные; отсутствие базы или брокера не считается успехом. Затем попросите объяснить одно решение и внести маленькое изменение входного условия. Не нужен построчный ручной review всего репозитория.

## ДЗ и темп

[Полное ДЗ](../combat/{filename(row)}) рассчитано на 8–12 часов с подготовкой. Если ученик не успевает, приоритет — обязательный контракт и понятный разбор первой ошибки. ★ не влияет на PASS и не блокирует следующий урок. Ускоряющемуся ученику предложите ★: {row['star']}

{resources(row)}'''


def homework(row):
    n = row['n']
    return f'''# ДЗ {n:02}. {row['title']}

{nav(row, 'combat')}

## Результат

{row['problem']} Продолжайте тот же репозиторий JobFlow.

**Защищаемое правило:** {row['invariant']}

Бюджет: 8–12 часов с подготовкой. Ориентир: 1–2 часа на чтение и план, 5–7 на реализацию, 1–2 на проверки и объяснение. ★ выбирается по оставшемуся времени; для обязательного PASS она не требуется.

## Обязательная часть

{numbered(row['homework'])}

Сохраняйте ранее введённые интерфейсы и обещания. Формальные сигнатуры и крайние случаи доступны в [проверке {n:02}](../checks/{n:02}_test.go) и [правилах проверки](../checks/README.md); зависимость инфраструктуры от домена направлена к домену.

## Критерии приёмки

{numbered(row['cases'])}

Из корня курса запустите:

```bash
make check LESSON={n:02} PROJECT=/absolute/path/to/jobflow
```

Команда должна завершиться кодом 0 и PASS обязательных проверок 01–{n:02}. Ошибка, пропуск инфраструктуры или отсутствие пакета дают FAIL. Собственные тесты добавляйте рядом с реализацией. Открытые тесты разрешено читать и запускать до сдачи.

## Что сдать

- PR или коммит со всеми обязательными изменениями и README запуска.
- Вывод автоматической проверки с номером этапа и ревизией кода.
- Краткое объяснение инварианта: где проверяется, какой отказ обнаружен, чем подтверждён результат.

## Необязательная ★

{row['star']}

Не меняйте обязательный внешний контракт ради ★. В PR опишите дополнительное условие и способ проверки. Результат ★ обсуждается отдельно и не меняет PASS основной части.

## Если застряли

Воспроизведите одну ошибку минимальным входом. Проверьте ожидаемый результат до правки, затем сравните свой участок с [эталоном {n:02}](../reference/{n:02}/README.md). Объясните отличие своими словами и повторите проверку без копирования всего проекта.

{resources(row)}'''


def main():
    for row in ROWS:
        name = filename(row)
        put('student/' + name, student(row))
        put('mentor/' + name, mentor(row))
        put('combat/' + name, homework(row))
        n = row['n']
        commands = 'go run ./cmd/hello' if n == 1 else 'go run ./cmd/pub'
        put(f'reference/{n:02}/README.md', f'''# JobFlow: эталон {n:02} — {row['title']}

Открытая самостоятельная версия **после** ДЗ {n:02}. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: {row['invariant']}

Из этой директории: `{commands}`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON={n:02} PROJECT=reference/{n:02}`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/{name}) · [ДЗ](../../combat/{name}) · [Открытые проверки](../../checks/{n:02}_test.go)
''')
    for folder, title in [('student','Конспекты ученика'), ('mentor','Руководства ментора'), ('combat','Домашние задания')]:
        listing = '\n'.join(f'- [{r["n"]:02}. {r["title"]}]({filename(r)})' for r in ROWS)
        put(folder + '/README.md', f'# {title}\n\n[Программа курса](../course/README.md) · [Правила работы](../docs/WORKFLOW.md)\n\n{listing}')
    plan = '# Программа: Go Backend Combat\n\n24 недели / 6 месяцев. На входе — опыт программирования на другом языке; Go изучается с нуля. Занятие: 50 минут + 10 минут резерва. ДЗ: 8–12 часов с подготовкой; ★ добровольна.\n\nСквозной результат — JobFlow: учебный движок DAG-задач из pub, sub, act, ret. PostgreSQL хранит истину, Kafka доставляет события, Redis ускоряет чтение и передаёт сигналы. Исполняются зарегистрированные кооперативные обработчики.\n\n[Как учиться и сдавать](../docs/WORKFLOW.md) · [Контракт проекта](../docs/PROJECT.md) · [Формат 128 байт](../docs/FRAMING.md)\n'
    for m, title in enumerate(MODULES):
        plan += f'\n## Месяц {m+1}. {title}\n\n| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |\n|---|---|---|---|\n'
        for r in ROWS[m*4:m*4+4]:
            n = r['n']
            plan += f'| {link(r,"student")} | {TOOLS[n-1]} | {r["invariant"]} | [{r["homework"][0]}](../combat/{filename(r)}) |\n'
    plan += '\n## Синхронизация материалов\n\nДля каждого этапа есть конспект, компактный план ментора, ДЗ с ★, самостоятельный эталон и открытые накопительные проверки. Исходные данные программы: [lessons.json](lessons.json). `make docs` пересобирает связанные документы; `make validate` проверяет ссылки и структуру. Исторические материалы сохранены в [legacy](../../../legacy/README.md).\n'
    put('course/README.md', plan)
    mapping = '# Плакаты и новая программа\n\nSVG сохранены как визуальные опоры: отдельный плакат может использоваться на нескольких встречах. Нумерация файлов историческая; последовательность обучения определяется [программой](../course/README.md). Термины, API и гарантии уточняются в действующих конспектах.\n\n| Плакат | Уроки |\n|---|---|\n'
    for path in sorted((ROOT/'posters').glob('*.svg')):
        uses = [f'[{r["n"]:02}](../student/{filename(r)})' for r in ROWS if r['poster'] == path.name]
        mapping += f'| [{path.stem}]({path.name}) | {", ".join(uses) if uses else "Дополнительная опора ментора"} |\n'
    mapping += '\nДля бинарного формата 12–13 используется нормативная [таблица полей](../docs/FRAMING.md), включая примеры padding и связи.\n'
    put('posters/README.md', mapping)
    put('README.md', '''# Go Backend Combat

**Go с нуля для разработчиков с опытом другого языка.** 24 недели, 6 месяцев; каждое занятие — 50 минут и 10 минут резерва. Домашняя работа — 8–12 часов с подготовкой, на каждом этапе есть добровольная ★.

Сквозной проект — **JobFlow**, учебный распределённый движок DAG-задач: `pub`, `sub`, `act`, `ret`. Курс ведёт от первой Go-программы до транзакций, собственного бинарного формата, Kafka, Redis Streams, безопасных повторов и очистки данных. Итог показывает устройство ограниченного движка задач; совместимость и полный набор возможностей Temporal не заявлены.

## Начать

1. Откройте [программу 24 уроков](course/README.md) и [правила обучения](docs/WORKFLOW.md).
2. Скопируйте [starter](starter/README.md) в собственный репозиторий.
3. Пройдите [урок 01](student/01_ENVIRONMENT.md) и [ДЗ 01](combat/01_ENVIRONMENT.md).
4. Из корня курса: `make check LESSON=01 PROJECT=/absolute/path/to/jobflow`.

## Материалы

- [Конспекты ученика](student/README.md), [руководства ментора](mentor/README.md), [домашние задания](combat/README.md).
- [Эталон итогового проекта](reference/24/README.md); `reference/01`–`reference/24` — открытые самостоятельные версии этапов.
- [Проверки и PASS](checks/README.md), [вопросы для защиты](interview/README.md), [плакаты](posters/README.md).
- [Нормативный формат 128 байт](docs/FRAMING.md), [состояния и API](docs/PROJECT.md), [запуск через Compose](deploy/README.md).
- [Отчёт о проверке выпуска](docs/VALIDATION.md), [архив исходного курса](../../legacy/README.md).

## Итоговый проект локально

Требуется Docker с Compose. Команды из корня репозитория:

```bash
make up
make demo
make down
```

Для локальной разработки и открытых проверок: Go 1.27.2, Python 3, make; с урока 11 — C-компилятор для race detector. `make validate` проверяет материалы; `make check-all` запускает все 24 накопительных контракта, включая настоящие PostgreSQL, Redis, Kafka и четыре процесса. `make smoke` проверяет полный Compose и восстановление после отказов, `make check-mutants` — обнаружение намеренных поломок. Проверки открыты, ★ на PASS не влияет.

## Обновление курса

`course/lessons.json` содержит содержание уроков; `make docs` обновляет связанные материалы. Каноническая реализация находится в `reference/24`; `python3 tools/materialize_stages.py` воспроизводит ранние этапы, затем выполните `make docs`. Исходные материалы и два локальных SQL-черновика сохранены в `legacy/2026-10-19` с SHA256-манифестом.
''')
    print('Rendered 72 lesson documents, 24 reference guides and course navigation.')


if __name__ == '__main__':
    main()
