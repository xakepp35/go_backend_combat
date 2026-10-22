# JobFlow: первый SQL-урок — от миграции к истории событий

**Продолжительность:** 50 минут + 10 минут резерва на вопросы.  
**Цель:** создать схему JobFlow, записать события и прочитать историю задания.  
**Формат:** объяснить → спросить результат → написать SQL → выполнить → проверить.

Сегодня один проект, одна база PostgreSQL, один компонент bus. Следующие занятия развивают этот же проект: проекция состояния, JOIN/агрегаты, pgx, транзакции, индексы.

## 1. План встречи

| Время | Что делаем | Что должен сделать ученик |
|---|---|---|
| 00–05 | Требование и структура JobFlow | Объяснить различие события и состояния |
| 05–12 | Первая миграция | Найти файл схемы, применить миграцию, проверить таблицы |
| 12–18 | Разобрать поля и constraints | Объяснить id/key/sub, PRIMARY KEY и FOREIGN KEY |
| 18–30 | INSERT и RETURNING | Записать тело и два события одного задания |
| 30–40 | SELECT, WHERE, ORDER BY, LIMIT | Прочитать историю и выбрать события по теме |
| 40–46 | Самостоятельная задача и плохая вставка | Проверить результат и предсказать ошибку |
| 46–50 | Итог и ДЗ | Сохранить SQL и назвать следующий прирост проекта |

**До встречи:** PostgreSQL уже доступен, Goose установлен. Инфраструктуру не устанавливаем вместо SQL. Для PostgreSQL 18 используем volume `/var/lib/postgresql`, запускаем только postgres, без Redis и Kafka.

## 2. Начинаем с требования, затем создаём файлы

> «В JobFlow будут фоновые задания. Сегодня пишем историю: какое событие произошло, когда, кто его создал и к какому заданию оно относится. История должна пережить перезапуск».

```text
job.init → job.wait → job.exec → job.done
                             ↘ job.fail
```

bus сохраняет события; state позже вычисляет состояние из истории; job выполняет работу. Пока это роли одного приложения. context.Context и goroutines живут в процессе, параметры и результаты задания — в событиях.

Структура по PROJECT FLOW:

```text
jobflow/
├── go.mod
├── cmd/bus-demo/main.go              # Сборка приложения; подключим позднее
├── pkg/bus/                         # Публичные типы и контракт компонента
├── internal/bus/                    # Правила и сервис, порт Repository
├── internal/adapter/postgres/       # Будущий адаптер с SQL и pgx
├── migrations/00001_bus.sql         # Создаём сегодня
└── sql/01_basics.sql                # Сохраняем сегодняшние запросы
```

Сегодня обязательно создаём модуль, migrations/ и sql/. Go-пакеты добавляем вместе с реализацией, не пишем пустые интерфейсы ради дерева. Обсуждаем путь `main → service → repository → PostgreSQL`; сегодня выполняем SQL напрямую через psql, позднее те же операции выполнит адаптер.

**Вопрос:** «Где будет SQL, где проверка входных данных, кто создаст подключение?»

**Ожидаем:** SQL — адаптер; правила — сервис; подключение — main.

## 3. Миграция: схема как файл проекта

> «Миграция — версионированное изменение схемы. Новый разработчик должен получить те же таблицы из файлов, а не повторять нашу ручную настройку».

Ученик создаёт `migrations/00001_bus.sql`. Саму схему можно дать готовой, затем вместе объяснить ключевые строки.

Для выбранного CAS физически разделяем тела и события: одно тело может принадлежать нескольким событиям. Представление bus показывает семь привычных полей.

```sql
-- +goose Up
CREATE TABLE payloads (
    key  text PRIMARY KEY,
    raw  bytea NOT NULL,
    data jsonb NOT NULL,
    CHECK (key = encode(sha256(raw), 'hex')),
    CHECK (data = convert_from(raw, 'UTF8')::jsonb),
    CHECK (jsonb_typeof(data) = 'object')
);

CREATE TABLE events (
    id      uuid PRIMARY KEY DEFAULT uuidv7(),
    key     text NOT NULL REFERENCES payloads(key),
    iat     timestamptz NOT NULL DEFAULT clock_timestamp(),
    iss     text NOT NULL CHECK (btrim(iss) <> ''),
    sub     uuid NOT NULL,
    topic   text NOT NULL CHECK (btrim(topic) <> ''),
    version bigint NOT NULL CHECK (version > 0),
    UNIQUE (sub, version)
);

CREATE VIEW bus AS
SELECT e.id, e.key, e.iat, e.iss, e.sub, e.topic, p.data
FROM events AS e
JOIN payloads AS p ON p.key = e.key;
```

**Комментарии для объяснения:**

- id — ID события; sub — ID задания; key — SHA-256 тела. Это три разных назначения.
- raw — канонические JSON-байты; data — их JSONB-представление для запросов.
- Канонизация RFC 8785 выполняется в приложении. Для сегодняшних SQL-примеров даём уже канонические строки. База проверяет соответствие hash/bytes/JSONB, но не реализует канонизацию.
- payloads.key уникален; events.key не уникален: одинаковое тело не означает одинаковое событие.
- iat — момент времени, используем UTC. version — порядок внутри задания; UUIDv7 и время не заменяют версию.
- PRIMARY KEY защищает идентичность; FOREIGN KEY не допускает событие без тела; CHECK защищает указанное условие.
- VIEW уже содержит JOIN. Сегодня используем его как готовое представление; подробно JOIN разберём отдельно.

Применяем к **пустой учебной БД PostgreSQL 18**:

```bash
export DATABASE_URL='postgres://jobflow:jobflow@localhost:5432/jobflow?sslmode=disable'
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" status
docker compose exec postgres psql -U jobflow -d jobflow
```

В psql:

```sql
SET TIME ZONE 'UTC';
\dt
\d events
\d payloads
```

**Спросить:** «Если снова выполнить goose up, появится вторая таблица events?»  
**Ответ:** нет, Goose учитывает применённые версии. Здесь миграция forward-only: DROP TABLE с историей не называем безопасным откатом.

**Резерв, если CLI Goose недоступен:** выполнить файл на отдельной пустой учебной БД через `psql -f migrations/00001_bus.sql`. Это выполнение SQL без учёта миграций; не запускать затем Goose поверх этих таблиц без baseline. Если psql только в контейнере: `docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U jobflow -d jobflow < migrations/00001_bus.sql`.

## 4. Запись: INSERT и RETURNING

> «Сначала добавим тело, затем событие со ссылкой на него. RETURNING позволяет увидеть поля реально созданной строки».

```sql
INSERT INTO payloads (key, raw, data)
VALUES (
    encode(sha256(convert_to('{"priority":1,"title":"Send report"}', 'UTF8')), 'hex'),
    convert_to('{"priority":1,"title":"Send report"}', 'UTF8'),
    '{"priority":1,"title":"Send report"}'::jsonb
)
RETURNING key, data;

INSERT INTO events (key, iss, sub, topic, version)
VALUES (
    encode(sha256(convert_to('{"priority":1,"title":"Send report"}', 'UTF8')), 'hex'),
    'student',
    '01900000-0000-7000-8000-000000000001',
    'job.init',
    1
)
RETURNING id, iat, sub, topic, version;
```

**Перед запуском:** «Какие поля назначит БД? Что вернёт RETURNING?»

Теперь добавляем ожидание:

```sql
INSERT INTO payloads (key, raw, data)
VALUES (
    encode(sha256(convert_to('{"queue":"default"}', 'UTF8')), 'hex'),
    convert_to('{"queue":"default"}', 'UTF8'),
    '{"queue":"default"}'::jsonb
)
RETURNING key;

INSERT INTO events (key, iss, sub, topic, version)
VALUES (
    encode(sha256(convert_to('{"queue":"default"}', 'UTF8')), 'hex'),
    'student',
    '01900000-0000-7000-8000-000000000001',
    'job.wait',
    2
)
RETURNING id, topic, version;
```

Сегодня выполняем отдельные команды, чтобы изучить INSERT. В Go-адаптере запись тела и события будет одной транзакцией. SQL-примеры вставляем один раз; повторная вставка той же версии закономерно даст конфликт.

## 5. Чтение: вопросы к данным

**Сначала ученик формулирует вопрос словами, затем пишет SQL.**

```sql
-- Что сейчас лежит в шине?
SELECT * FROM bus;

-- Какие задания были созданы?
SELECT sub, iat, data
FROM bus
WHERE topic = 'job.init';

-- История одного задания в правильном порядке.
SELECT version, topic, iss, iat
FROM events
WHERE sub = '01900000-0000-7000-8000-000000000001'
ORDER BY version;

-- Последние записи по времени события: обзор, не алгоритм потребления bus.
SELECT id, sub, topic, iat
FROM bus
ORDER BY iat DESC, id DESC
LIMIT 5;

-- Достаём название из JSONB.
SELECT sub, data ->> 'title' AS title
FROM bus
WHERE topic = 'job.init';

-- Фильтр относится к известной схеме job.init, где priority — целое.
SELECT sub, data
FROM bus
WHERE topic = 'job.init'
  AND (data ->> 'priority')::integer >= 1;
```

**Вопросы:**

1. Что изменится, если убрать WHERE?
2. Гарантирован ли порядок без ORDER BY?
3. Почему историю сортируем по version?
4. Что возвращает `->>` и зачем `::integer`?
5. «Когда-либо было событие job.wait» означает «задание сейчас ждёт»?

**Ориентиры:** WHERE отбирает строки; без ORDER BY порядок не обещан; version задаёт порядок одного sub; `->>` возвращает текст; текущее состояние надо вычислять из всей истории.

## 6. Самостоятельная задача на 6 минут

1. Добавить init второго задания с `sub = '01900000-0000-7000-8000-000000000002'`, version = 1 и тем же key тела Send report. Новый payload не нужен.
2. Прочитать все init: должны получиться два события.
3. Прочитать только историю первого задания: init и wait.
4. Попробовать повторить wait первого задания с version = 2. Предсказать ошибку UNIQUE и проверить SELECT, что история осталась прежней.

**Ожидаемый результат:** два payload, три events. Одинаковое тело не склеило два события разных заданий.

Если есть резерв времени: вставка с пустым iss → CHECK; с несуществующим key → FOREIGN KEY. Ученик объясняет, какой invariant защищён.

## 7. ДЗ: развиваем этот JobFlow

**Бюджет: 3–4 часа.**

- Сохранить миграцию и успешные запросы в sql/01_basics.sql. Ожидаемые ошибки оставить закомментированными.
- Создать три истории: `init → wait → exec → done`, `init → wait`, `init → wait → exec → fail`. При каждом событии увеличивать version внутри sub.
- Для exec взять `{"worker":"worker-1"}`, done — `{"result":"ok"}`, fail — `{"error":"temporary failure"}`. Distinct тела сначала добавлять в payloads; одинаковые тела использовать повторно.
- Написать пять запросов: история по sub; все init; все fail; последние пять записей; init с priority >= 2.
- Записать в README команды запуска и применения миграций, объяснить id/key/sub и зачем разделены payloads/events.
- Проверить три плохие записи: пустой автор, неизвестный key, повторная версия. После отказа подтвердить состояние SELECT-запросом.

**Приёмка:** схема воспроизводится через Goose в новой учебной БД; история читается в нужном порядке; SQL отвечает на пять вопросов; ограничения действительно отвергают плохие записи; ученик объясняет результат.

## 8. Следующие занятия

1. **Проекция jobs_read:** INSERT/UPDATE/DELETE, ON CONFLICT, новая миграция. Историю bus не редактируем; CRUD тренируем на восстанавливаемом состоянии.
2. **Связи и отчёты:** JOIN между событиями, телами и проекцией; COUNT, GROUP BY, HAVING, CTE. Задачи — отчёты по своим заданиям.
3. **Go + PostgreSQL:** pgxpool, контекст, scanning, ошибки; тот же контракт service/repository, реальный PostgreSQL-адаптер.
4. **Транзакции и конкуренция:** атомарный append, проверка версии, конфликты и rollback.
5. **Проверка хранения:** интеграционные/контрактные тесты, индекс под запрос истории, EXPLAIN.

Не пытаемся пройти эти пункты в остаток сегодняшнего часа.

Технические основания: [PostgreSQL UUIDv7](https://www.postgresql.org/docs/18/functions-uuid.html), [JSONB](https://www.postgresql.org/docs/18/datatype-json.html), [SHA-256](https://www.postgresql.org/docs/18/functions-binarystring.html), [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785.html), [Goose](https://github.com/pressly/goose).
