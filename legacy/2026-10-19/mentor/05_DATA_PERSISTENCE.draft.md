// reply:"resp.md" v1.0.0 | 2026-10-14T11:28:00Z | Topic: Go Backend Combat — PostgreSQL / SQL | Summary: Начинаем PostgreSQL с нуля: сначала модель данных и базовые CRUD-запросы, затем JOIN и агрегаты. Формат — живой урок ментора с вопросами, задачами и проверкой понимания.

# Занятие: PostgreSQL как persistence layer Go-сервиса

> **Ментор:**  
> Сегодня мы не будем начинать с `pgx`, репозиториев или каких-нибудь модных архитектурных слов.  
>
> Сначала нужно понять саму вещь, с которой мы работаем.
>
> **PostgreSQL — это не просто место, куда мы складываем структуры Go.**
>
> Это отдельная система со своими правилами хранения, поиска, изменения данных и конкурентного доступа.
>
> И если мы не понимаем эти правила, то любой repository поверх PostgreSQL превращается в красивую обёртку над неизвестностью.
>
> Поэтому сегодня идём **с самого нуля**.

---

# 0. Что мы вообще строим

Представим обычный Go-сервис.

Например, сервис управления устройствами:

```text
HTTP/gRPC
    ↓
Application Service
    ↓
Repository
    ↓
PostgreSQL
```

На предыдущих занятиях мы уже разобрали, **зачем нужны Service, Repository и Adapter**.

Сегодня нас интересует нижняя часть:

```text
Repository
    ↓
PostgreSQL
```

Но начинать будем ещё ниже:

```text
SQL
    ↓
PostgreSQL
```

> **Ментор:**
>
> Представьте, что PostgreSQL перед вами впервые.
>
> Что вы хотите у него спросить?
>
> Например:
>
> **«Дай мне все устройства».**
>
> Как бы вы сформулировали это на языке SQL?

### Интерактив №1

**Вопрос:** какой SQL-запрос, по-твоему, должен вернуть все строки таблицы `devices`?

```sql
???
```

Не подглядывай. Сначала попробуй написать сам.

---

# 1. Таблица — это не `struct`

Пусть у нас есть таблица:

```sql
CREATE TABLE devices (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    status     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Ментор:

> Здесь важно сразу разделить две вещи.
>
> В Go мы можем представить устройство так:

```go
type Device struct {
    ID        int64
    Name      string
    Status    string
    CreatedAt time.Time
}
```

> Но PostgreSQL ничего не знает про эту структуру.
>
> Для PostgreSQL существует **таблица**, состоящая из строк и колонок.
>
> Go-структура — это представление данных на стороне приложения.
>
> SQL-таблица — представление данных на стороне БД.

То есть:

```text
Go                       PostgreSQL

Device            ←→     devices row

ID                ←→     id
Name              ←→     name
Status            ←→     status
CreatedAt         ←→     created_at
```

Но это **не одно и то же**.

Это принципиально важно для дальнейшей работы.

---

# 2. SELECT — чтение данных

Теперь возвращаемся к нашему вопросу.

Нужно получить все устройства.

Ответ:

```sql
SELECT *
FROM devices;
```

> **Ментор:**
>
> Здесь есть две части:
>
> `SELECT` — какие данные мы хотим получить.
>
> `FROM` — откуда их получить.
>
> `*` означает «все колонки».

Но почти сразу возникает вопрос:

> А зачем нам вообще выбирать все колонки?

Например:

```sql
SELECT id, name, status
FROM devices;
```

Это уже более явный запрос.

И обычно это лучше, чем бездумное:

```sql
SELECT *
FROM devices;
```

Потому что запрос начинает явно описывать свой контракт.

---

# 3. WHERE — фильтрация

Теперь задача меняется.

> **Ментор:**
>
> У нас 10 000 устройств.
>
> Мне нужны только устройства со статусом `active`.
>
> Что добавим?

Ответ:

```sql
SELECT id, name, status
FROM devices
WHERE status = 'active';
```

Логика:

```text
FROM devices
       ↓
WHERE status = 'active'
       ↓
SELECT id, name, status
```

Важно понимать порядок **смысловых операций**, а не только порядок написания текста.

Мы говорим БД:

> «Работай с `devices`, оставь строки, удовлетворяющие условию, и верни указанные колонки».

---

# 4. ORDER BY

Теперь:

> **Ментор:**
>
> Хорошо. Но я хочу получить устройства от самых новых к самым старым.

```sql
SELECT id, name, status, created_at
FROM devices
ORDER BY created_at DESC;
```

`DESC`:

```text
9
8
7
...
1
```

`ASC`:

```text
1
2
3
...
9
```

По умолчанию `ASC`.

---

# 5. LIMIT

Теперь представим API:

```text
GET /devices
```

Мы не хотим случайно вернуть миллион строк.

```sql
SELECT id, name, status, created_at
FROM devices
ORDER BY created_at DESC
LIMIT 100;
```

> **Ментор:**
>
> Обратите внимание на важную инженерную мысль.
>
> PostgreSQL может содержать очень много данных.
>
> Поэтому запрос:
>
> ```sql
> SELECT * FROM devices;
> ```
>
> может быть совершенно нормальным учебным запросом и совершенно плохой основой production API.

---

# 6. INSERT — создание строки

Теперь нужно создать устройство.

```sql
INSERT INTO devices (
    name,
    status
)
VALUES (
    'sensor-01',
    'active'
);
```

> **Ментор:**
>
> Мы указали только `name` и `status`.
>
> А `id`?
>
> Его генерирует PostgreSQL благодаря `BIGSERIAL`.
>
> А `created_at`?
>
> Его генерирует `DEFAULT now()`.

То есть часть состояния создаётся самой БД.

---

# 7. RETURNING

И вот здесь появляется очень полезная возможность PostgreSQL.

После `INSERT` нам часто нужен созданный объект.

Не просто:

```text
INSERT успешно
```

а:

```text
вот тебе id новой записи
```

Для этого:

```sql
INSERT INTO devices (
    name,
    status
)
VALUES (
    'sensor-01',
    'active'
)
RETURNING id, name, status, created_at;
```

PostgreSQL вернёт созданную строку.

Например:

```text
id | name      | status | created_at
---+-----------+--------+-------------------
42 | sensor-01 | active | 2026-10-14 10:00:00
```

> **Ментор:**
>
> Запомните `RETURNING`.
>
> Для PostgreSQL это один из очень важных инструментов при работе из Go.
>
> Нам не обязательно делать:
>
> ```text
> INSERT
> ↓
> SELECT WHERE id = ...
> ```
>
> чтобы получить созданную запись.

---

# 8. UPDATE

Теперь устройство изменило состояние.

```sql
UPDATE devices
SET status = 'offline'
WHERE id = 42;
```

И снова можно использовать `RETURNING`:

```sql
UPDATE devices
SET status = 'offline'
WHERE id = 42
RETURNING id, name, status, created_at;
```

---

# 9. DELETE

Удаление:

```sql
DELETE FROM devices
WHERE id = 42;
```

И опять:

```sql
DELETE FROM devices
WHERE id = 42
RETURNING id, name, status;
```

может вернуть удалённую строку.

---

# 10. Очень важный вопрос: что будет без WHERE?

> **Ментор:**
>
> Сейчас проверим, насколько вы уже начали чувствовать SQL.
>
> Что произойдёт здесь?

```sql
UPDATE devices
SET status = 'offline';
```

Подумай.

Не отвечай автоматически.

---

### Ответ

Изменятся **все строки** таблицы `devices`.

То же самое:

```sql
DELETE FROM devices;
```

удалит все строки.

Поэтому одно из базовых правил работы с SQL:

> **UPDATE и DELETE без осознанного условия — потенциально разрушительная операция.**

В production-коде мы должны точно понимать, сколько строк ожидаем изменить.

---

# 11. А теперь первая инженерная проверка

Представим:

```sql
UPDATE devices
SET status = 'offline'
WHERE id = 42;
```

Мы ожидали:

```text
1 row affected
```

Но PostgreSQL сообщил:

```text
0 rows affected
```

Что это может означать?

Возможные причины:

```text
id = 42 не существует
```

или:

```text
мы ошиблись в условии
```

А теперь другой вариант:

```text
100 rows affected
```

при ожидании одной строки.

Это уже **серьёзный сигнал**.

> **Ментор:**
>
> Поэтому repository должен не просто «выполнить SQL».
>
> Он должен понимать семантику результата операции.
>
> Например:
>
> ```text
> UPDATE по ID
> ```
>
> обычно ожидает:
>
> ```text
> 0 → NotFound
> 1 → Success
> >1 → нарушение предположения/ошибка
> ```
>
> Это уже мост от SQL к нашему application layer.

---

# 12. DELETE и бизнес-смысл

Здесь остановимся.

Представим API:

```text
DeleteDevice(42)
```

Что означает:

```sql
DELETE FROM devices
WHERE id = 42;
```

?

Только:

> «Физически удалить строку».

Но бизнес-смысл может быть другим:

```text
device.status = 'deleted'
```

или:

```text
device.deleted_at = now()
```

То есть SQL-операция и бизнес-операция — не обязательно одно и то же.

Это очень важная граница:

```text
Application:
    DeleteDevice()

Persistence:
    DELETE FROM ...
```

Application решает **что означает удаление**.

Repository решает **как сохранить это состояние**.

---

# 13. Теперь несколько таблиц

До этого у нас была одна таблица.

Реальные системы быстро становятся такими:

```text
users
devices
projects
orders
payments
events
...
```

Например:

```sql
CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
```

и:

```sql
CREATE TABLE devices (
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    name    TEXT NOT NULL,
    status  TEXT NOT NULL
);
```

Теперь появляется связь:

```text
users
  │
  └── devices
        user_id → users.id
```

---

# 14. JOIN

> **Ментор:**
>
> Представьте, что мне нужно получить:
>
> ```text
> device.name
> user.name
> ```
>
> из двух таблиц.
>
> Как это сделать?

Через `JOIN`:

```sql
SELECT
    d.id,
    d.name,
    d.status,
    u.id AS user_id,
    u.name AS user_name
FROM devices AS d
JOIN users AS u
    ON u.id = d.user_id;
```

Здесь:

```text
devices AS d
users AS u
```

дают короткие имена:

```text
d.name
u.name
```

А:

```sql
ON u.id = d.user_id
```

описывает условие связи.

---

# 15. Что делает JOIN на самом деле?

Не воспринимай `JOIN` как магическую команду.

Ментор:

> Представьте две таблицы.
>
> `users`:
>
> ```text
> id | name
> ---+------
> 1  | Bob
> 2  | Alice
> ```
>
> `devices`:
>
> ```text
> id | user_id | name
> ---+---------+--------
> 10 | 1       | sensor
> 11 | 1       | camera
> 12 | 2       | phone
> ```
>
> После:
>
> ```sql
> JOIN users u ON u.id = d.user_id
> ```
>
> PostgreSQL сопоставляет строки по условию:
>
> ```text
> devices.user_id = users.id
> ```

Получаем логически:

```text
sensor  → Bob
camera  → Bob
phone   → Alice
```

---

# 16. INNER JOIN и LEFT JOIN

Теперь важное различие.

### INNER JOIN

```sql
FROM devices d
JOIN users u
    ON u.id = d.user_id
```

Вернутся только устройства, для которых найден соответствующий пользователь.

### LEFT JOIN

```sql
FROM users u
LEFT JOIN devices d
    ON d.user_id = u.id
```

Вернутся **все пользователи**, даже если у пользователя нет устройств.

Например:

```text
users:

1 Bob
2 Alice
3 John
```

а устройства есть только у Bob и Alice.

`LEFT JOIN` даст:

```text
Bob   → sensor
Alice → camera
John  → NULL
```

> **Ментор:**
>
> Вот здесь SQL начинает становиться действительно интересным.
>
> JOIN — это не просто «соединить две таблицы».
>
> Нужно понимать, **какую сторону результата мы обязаны сохранить**.

---

# 17. Первая задача на JOIN

Есть:

```text
users
    id
    name

devices
    id
    user_id
    name
```

Нужно получить:

```text
имя пользователя
имя устройства
```

Но **включая пользователей, у которых вообще нет устройств**.

### Интерактив №2

Какой вариант правильный?

**A**

```sql
SELECT u.name, d.name
FROM users u
JOIN devices d
    ON d.user_id = u.id;
```

**B**

```sql
SELECT u.name, d.name
FROM users u
LEFT JOIN devices d
    ON d.user_id = u.id;
```

**C**

```sql
SELECT u.name, d.name
FROM devices d
LEFT JOIN users u
    ON u.id = d.user_id;
```

**D**

```sql
SELECT u.name, d.name
FROM users u
CROSS JOIN devices d;
```

**Остановись здесь и выбери вариант.**

---

# 18. Следующий шаг — агрегаты

После `JOIN` естественный следующий вопрос:

> Сколько устройств у каждого пользователя?

Здесь появляются агрегатные функции:

```sql
COUNT()
SUM()
AVG()
MIN()
MAX()
```

Например:

```sql
SELECT
    user_id,
    COUNT(*) AS device_count
FROM devices
GROUP BY user_id;
```

Получаем:

```text
user_id | device_count
--------+-------------
1       | 2
2       | 1
```

---

# 19. GROUP BY

Ментор:

> Что делает `GROUP BY`?
>
> Он говорит:
>
> **«Сгруппируй строки по этому признаку, а затем вычисли агрегаты для каждой группы».**

До:

```text
device
user_id
------
1
1
2
2
2
```

После:

```text
user_id | rows
--------+-----
1       | 2
2       | 3
```

И затем:

```sql
COUNT(*)
```

даёт размер каждой группы.

---

# 20. JOIN + GROUP BY

Теперь можно получить уже человеческий результат:

```sql
SELECT
    u.id,
    u.name,
    COUNT(d.id) AS device_count
FROM users AS u
LEFT JOIN devices AS d
    ON d.user_id = u.id
GROUP BY
    u.id,
    u.name;
```

Результат:

```text
id | name  | device_count
---+-------+-------------
1  | Bob   | 2
2  | Alice | 1
3  | John  | 0
```

Обрати внимание на `LEFT JOIN`.

Если использовать обычный `JOIN`, John исчезнет из результата.

---

# 21. HAVING

Теперь задача:

> Показать только пользователей, у которых **два или более устройства**.

Нельзя просто написать:

```sql
WHERE COUNT(d.id) >= 2
```

Почему?

Потому что `WHERE` работает **до формирования групп**, а `COUNT()` появляется **после группировки**.

Для условий на агрегат используется:

```sql
HAVING
```

Поэтому:

```sql
SELECT
    u.id,
    u.name,
    COUNT(d.id) AS device_count
FROM users AS u
LEFT JOIN devices AS d
    ON d.user_id = u.id
GROUP BY
    u.id,
    u.name
HAVING COUNT(d.id) >= 2;
```

---

# 22. Ментальная модель SQL

К этому моменту нужно начать видеть запрос примерно так:

```text
FROM
  ↓
JOIN
  ↓
WHERE
  ↓
GROUP BY
  ↓
HAVING
  ↓
SELECT
  ↓
ORDER BY
  ↓
LIMIT
```

Это **концептуальный порядок обработки**, а не просто порядок написания SQL.

Например:

```sql
SELECT ...
FROM ...
JOIN ...
WHERE ...
GROUP BY ...
HAVING ...
ORDER BY ...
LIMIT ...
```

---

# 23. Контрольная задача

Есть:

```sql
users (
    id,
    name
)

devices (
    id,
    user_id,
    status
)
```

Нужно получить:

```text
user_id
user_name
active_device_count
```

для **всех пользователей**, включая тех, у кого `0` активных устройств.

И вывести только пользователей, у которых:

```text
active_device_count >= 2
```

### Интерактив №3

Напиши SQL самостоятельно.

Подсказки:

```text
1. Нам нужны все users.
2. Значит, какой JOIN?
3. Считать нужно только active devices.
4. Нужно сгруппировать по пользователю.
5. Фильтр по COUNT идёт куда?
```

Не смотри решение сразу.

---

# 24. Что мы должны понять после первой части

Не надо пока запоминать сотню SQL-конструкций.

Нужно понять несколько фундаментальных вещей:

```text
SELECT
    читаем данные

INSERT
    создаём данные

UPDATE
    изменяем данные

DELETE
    удаляем данные

WHERE
    фильтруем строки

JOIN
    связываем отношения

GROUP BY
    формируем группы

HAVING
    фильтруем группы

RETURNING
    получаем результат модификации

ORDER BY
    задаём порядок

LIMIT
    ограничиваем результат
```

И главное:

> **SQL — это не набор строк, которые надо механически выучить.**
>
> SQL описывает, **какой результат над данными мы хотим получить**.

---

# 25. Куда идём дальше

После того как эти основы закреплены, следующий слой будет уже гораздо интереснее:

```text
SQL
 │
 ├── CTE
 ├── subquery
 ├── EXISTS
 ├── ON CONFLICT
 ├── constraints
 │
 ↓
Индексы
 │
 ├── B-tree
 ├── составные индексы
 ├── unique indexes
 └── EXPLAIN
       ↓
    EXPLAIN ANALYZE
       ↓
    query planner
```

И только **после этого** мы перейдём к транзакциям:

```text
BEGIN
  ↓
SQL
  ↓
COMMIT / ROLLBACK
```

где уже начнутся:

```text
atomicity
isolation
concurrent transactions
lost update
race conditions
locks
isolation levels
```

А затем всё это будет соединено с Go:

```text
Go
 ↓
pgx/v5
 ↓
pgxpool
 ↓
Repository
 ↓
PostgreSQL
```

То есть мы не будем учить `pgx` как набор вызовов API. **Сначала понимаем, что именно должна делать БД, затем учимся корректно заставлять Go это делать.**
