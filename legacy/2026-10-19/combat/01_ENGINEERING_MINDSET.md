<!-- artifact:practicum/01_ENGINEERING_MINDSET.md v1.0.0 2026-09-13T15:43:00Z Topic: Delivery Summary: Практикум первого занятия -->

# 01. ENGINEERING MINDSET — PRACTICUM

Цель практикума — натренировать привычку смотреть на backend-задачу через требования, инварианты, границы, механизмы и проверяемый результат. Здесь нет лекционной части: только боевые задачи. На каждой задаче нужно что-то написать, исправить, классифицировать, упорядочить или доказать. Если не получается решить самостоятельно, используй ментора или LLM для разбора и подсказки.

Все ответы выполняй в каталоге:

```text
practicum/01/
```

Главный результат практикума:

```text
artifacts/01_ENGINEERING_MINDSET.md
```

После каждой задачи добавляй результат в этот файл в соответствующую секцию `Battle NN`.

---

## BATTLE 01. Что здесь вообще требуется?

Бизнес-заказчик говорит:

> Пользователь должен иметь возможность создать фоновую задачу. Система должна быстро принять задачу и надёжно обработать её позже. Пользователь должен получить идентификатор созданной задачи.

Выпиши:

```text
REQUIREMENTS: 5
ASSUMPTIONS: 3
AMBIGUITIES: 2
```

Задание считается пройденным, если:

* каждый requirement описывает ожидаемое поведение;
* каждое assumption явно обозначено как допущение;
* каждая ambiguity показывает, какую информацию ещё нужно уточнить;
* requirements нельзя проверить только по вкусу или субъективному мнению.

---

## BATTLE 02. Требование или разговор?

Классифицируй следующие утверждения:

1. Система должна быть удобной.
2. API возвращает идентификатор созданной Job.
3. Job должна быть обработана позже.
4. API должно работать быстро.
5. API отвечает не более чем за 200 ms в заданном сценарии.
6. Система должна быть современной.
7. Ошибка создания Job не должна приводить к созданию Job.
8. Система должна быть надёжной.
9. Повторная отправка одной операции не должна создавать нежелательный повторный effect.
10. Код должен быть красивым.

Для каждого укажи одну категорию:

```text
REQUIREMENT
NEEDS CLARIFICATION
OPINION
```

Для всех `NEEDS CLARIFICATION` запиши улучшенную формулировку.

Результат:

```text
artifacts/requirements.txt
```

Задание считается пройденным, если каждый пункт имеет одну категорию, а расплывчатые утверждения превращены в проверяемые.

---

## BATTLE 03. Найди то, что нельзя нарушать

Для JobFlow определены следующие утверждения:

* Job имеет уникальный идентификатор.
* Job обычно быстро выполняется.
* Job всегда имеет тип.
* Job может находиться в состоянии `RUNNING`.
* `RUNNING` может перейти в `DONE`.
* `DONE` не должна самопроизвольно снова становиться `RUNNING`.
* Некорректная Job не должна становиться корректным сохранённым состоянием.
* Система должна быть удобной.
* Две worker-операции могут одновременно менять одну Job.
* Ошибка обработки не должна незаметно превращать `FAILED` в `DONE`.

Отметь:

```text
INVARIANT
PROPERTY
OPINION
INVALID DESIGN
```

Затем самостоятельно добавь ещё пять invariants для JobFlow.

Для каждого invariant запиши:

```text
INVARIANT:
HOW IT CAN BE VIOLATED:
```

Задание считается пройденным, если invariant можно однозначно проверить и представить сценарий его нарушения.

---

## BATTLE 04. Собери state machine

Для Job существуют состояния:

```text
NEW
RUNNING
DONE
FAILED
CANCELLED
```

Ниже перечислены возможные переходы:

```text
NEW → RUNNING
NEW → DONE
NEW → FAILED
RUNNING → DONE
RUNNING → FAILED
RUNNING → CANCELLED
DONE → RUNNING
DONE → FAILED
FAILED → RUNNING
CANCELLED → RUNNING
```

Раздели переходы на:

```text
ALLOWED
FORBIDDEN
REQUIRES REQUIREMENT
```

После этого выпиши итоговую state machine JobFlow.

Результат:

```text
artifacts/state_machine.txt
```

Задание считается пройденным, если для каждого состояния явно понятны допустимые переходы, а terminal states не имеют случайных обратных переходов.

---

## BATTLE 05. Кто владеет состоянием?

Есть компоненты:

```text
HTTP Handler
Job Service
Worker
Repository
Database
Cache
```

И данные:

```text
Job ID
Job Type
Job Status
Payload
CreatedAt
Retry Count
```

Для каждого поля определи:

```text
OWNER
WHO MAY CHANGE IT
WHERE THE AUTHORITATIVE VALUE LIVES
```

Заполни таблицу в `artifacts/01_ENGINEERING_MINDSET.md`.

Задание считается пройденным, если для каждого элемента существует один authoritative owner и понятно, каким путём происходит изменение.

---

## BATTLE 06. Найди границу

Дан проект, в котором бизнес-логика напрямую вызывает конкретный PostgreSQL client.

Завтра команда хочет заменить storage implementation.

Ответь письменно:

```text
1. Что сейчас связано с PostgreSQL?
2. Что из этого является business policy?
3. Что является infrastructure detail?
4. Какую часть системы нужно сделать стабильным contract?
5. Какой минимальный boundary нужен?
```

После этого нарисуй словами итоговую зависимость между:

```text
business policy
contract
storage implementation
```

Результат:

```text
artifacts/boundary.md
```

Задание считается пройденным, если замена storage implementation не требует переписывать business policy.

---

## BATTLE 07. Выбери механизм

Для каждой ситуации выбери минимальный механизм:

```text
interface
goroutine
channel
context.Context
transaction
```

Ситуации:

1. Потребитель должен зависеть от поведения, а не от конкретной реализации.
2. Нужно независимо выполнять concurrent work.
3. Несколько concurrent operations должны обмениваться значениями или сигналами.
4. Работа должна остановиться после отмены внешнего request.
5. Два изменения состояния должны выполниться атомарно.

Для каждого решения запиши:

```text
PROBLEM:
MECHANISM:
GUARANTEE:
WHAT THIS DOES NOT GUARANTEE:
```

Задание считается пройденным, если механизм объясняется через проблему и гарантию, а не через название технологии.

---

## BATTLE 08. Сломай универсальное решение

Разработчик предлагает правило:

> На каждый тип и каждую функцию нужно добавить interface. Тогда система всегда будет легко тестироваться и изменяться.

Разбери решение.

Нужно написать:

```text
2 BENEFITS
3 COSTS
2 FAILURE MODES
1 CASE WHERE INTERFACE IS JUSTIFIED
1 CASE WHERE INTERFACE IS NOT JUSTIFIED
```

После этого сформулируй собственное правило выбора abstraction.

Результат:

```text
artifacts/01_ENGINEERING_MINDSET.md
```

Задание считается пройденным, если аргументы опираются на coupling, cohesion, cost of change, testability или complexity, а не на фразу `best practice`.

---

## BATTLE 09. А если всё сломается?

Есть операция:

```text
POST /jobs
→ validate
→ create Job
→ save
→ return 201
```

Разбери следующие сценарии:

```text
invalid request
database failure
database timeout
client retry
concurrent request
process crash
partial failure
```

Для каждого сценария запиши:

```text
FAILURE
WHAT CAN GO WRONG
EXPECTED BEHAVIOR
HOW TO VERIFY
```

Результат:

```text
artifacts/failure_matrix.md
```

Задание считается пройденным, если каждый сценарий имеет определённое ожидаемое поведение, а не просто `return error`.

---

## BATTLE 10. Докажи, что это работает

Есть утверждения:

1. Validation Job корректна.
2. Job ID уникален.
3. Недопустимый state transition отклоняется.
4. Concurrent operations не повреждают state.
5. Cancellation доходит до нужной операции.
6. Ошибка storage не оставляет некорректное состояние.
7. Производительность соответствует требованию.

Для каждого выбери конкретное evidence и запиши способ проверки.

Разрешены:

```text
unit test
integration test
race test
benchmark
failure test
manual reproduction
runtime observation
```

Для каждой строки обязательно укажи:

```text
CLAIM
EVIDENCE
PROCEDURE
PASS
FAIL
```

Результат:

```text
artifacts/verification.md
```

Задание считается пройденным, если другой человек сможет воспроизвести проверку только по записанной процедуре.

---

## BATTLE 11. Исправь плохое решение

Перед тобой реализация JobService:

* получает конкретный PostgreSQL client;
* получает конкретный Redis client;
* возвращает `nil`, если Job невалидна;
* игнорирует ошибку записи в database;
* записывает Job в Redis до подтверждения основного состояния;
* запускает notification в отдельной goroutine без cancellation;
* не проверяет повторную обработку;
* не имеет явного owner состояния Job.

Сначала выпиши минимум восемь дефектов.

Для каждого:

```text
DEFECT
VIOLATED INVARIANT
FAILURE SCENARIO
MINIMAL FIX
HOW TO VERIFY
```

После анализа исправь решение.

Результаты:

```text
artifacts/defects.md
artifacts/repaired.go
```

Задание считается пройденным, если каждый исправленный дефект связан с конкретной причиной и имеет verification.

---

## BATTLE 12. Защити своё решение

Теперь тебе дают новую задачу.

> Нужно принять запрос на обработку изображения, вернуть идентификатор операции и выполнить обработку позже. Обработка может быть долгой. Клиент может повторить запрос. Процесс может завершиться в любой момент.

Не используй заранее заданные технологии.

Самостоятельно определи:

```text
REQUIREMENTS
INVARIANTS
STATE
OWNERSHIP
BOUNDARIES
MECHANISMS
FAILURES
EVIDENCE
VERIFICATION
NON-GUARANTEES
```

Главное условие: объяснение должно начинаться с задачи и требуемого поведения, а не с `Kafka`, `Redis`, `goroutine`, `interface` или другого механизма.

Результат:

```text
artifacts/transfer_test.md
```

Задание считается пройденным, если решение остаётся логичным при изменении конкретной технологии реализации.

---

# FINAL ARTIFACT

После завершения двенадцати боёв приведи:

```text
artifacts/01_ENGINEERING_MINDSET.md
```

к следующей структуре:

```text
# JobFlow Engineering Mindset

## Battle 01
Requirements
Assumptions
Ambiguities

## Battle 02
Classification

## Battle 03
Invariants

## Battle 04
State Machine

## Battle 05
Ownership

## Battle 06
Boundaries

## Battle 07
Mechanisms

## Battle 08
Trade-offs

## Battle 09
Failure Model

## Battle 10
Evidence

## Battle 11
Defects and Repairs

## Battle 12
Final Defense

## Transfer Test

## Final Principles
```

В конце документа сформулируй своими словами пять правил, которые ты собираешься применять при решении backend-задач.

---

# FINAL CHECK

Обязательно проверь себя в конце! Боевой практикум считается PASS, если ты способен взять незнакомую backend-задачу и без подсказки пройти последовательность:

* Что система должна делать?
* Что нельзя нарушить?
* Кто владеет состоянием?
* Где проходит граница?
* Какой минимальный механизм нужен?
* Что может сломаться?
* Как это доказать?
* Что считается PASS?
* Что система не гарантирует?

Если на вопрос хочется сразу ответить названием технологии, задача ещё не решена.
