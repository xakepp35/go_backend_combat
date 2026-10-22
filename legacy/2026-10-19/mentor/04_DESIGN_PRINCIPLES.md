# Go Backend Combat — 04. DESIGN

**Version:** 1.0.0
**Owner:** Mentor
**Duration:** 120 минут
**Format:** интерактивный mentor-led lesson
**Project:** JobFlow
**Prerequisite:** `01.LANGUAGE.md`, `02.CONCURRENCY.md`
**Posters:** `04A_SOLID.svg`, `04B_DESIGN_ENGINEERING.svg`
**Outcome:** ученик умеет объяснить SOLID на интервью и применять принципы как инструменты управления изменениями, контрактами и зависимостями

---

# 1. Purpose

После двух первых занятий ученик уже умеет писать:

```text
Job
Repository
Service
WorkerPool
goroutines
channels
context
synchronization
tests
```

Теперь возникает новая проблема.

Код работает, но его начинает становиться **дорого менять**.

Например:

```text
HTTP handler
 ├── validation
 ├── business rules
 ├── PostgreSQL
 ├── Kafka
 └── logging
```

На этом занятии мы не изучаем SOLID как набор определений.

Мы делаем два шага:

```text
04A — знать и объяснять SOLID
        ↓
04B — видеть стоимость изменения
        ↓
JobFlow refactoring
        ↓
новая boundary
        ↓
проверяемое изменение
```

Главная идея:

> **Не применяй принцип потому, что он существует. Сначала найди стоимость изменения, затем выбери минимальную границу, которая её уменьшает.**

---

# 2. Learning Contract

После занятия ученик умеет:

### Interview

* дать точное определение SOLID;
* объяснить каждый из пяти принципов;
* привести Go-specific пример;
* объяснить, почему маленькие интерфейсы особенно естественны для Go;
* объяснить DIP через constructor injection.

### Engineering

* определить ответственность компонента;
* определить, что должно изменяться вместе;
* находить coupling;
* оценивать cohesion;
* видеть dependency direction;
* выделять consumer-side interface;
* различать abstraction и premature abstraction;
* использовать KISS / DRY / YAGNI без догматизма;
* объяснить cost of change.

### Практика

Ученик должен самостоятельно:

```text
увидеть плохую границу
→ объяснить проблему
→ предложить boundary
→ изменить код
→ сохранить behavior
→ доказать tests
```

---

# 3. Главный Mentor Flow

Весь урок ведём одним циклом:

```text
OBSERVE
   ↓
ASK
   ↓
PREDICT
   ↓
IDENTIFY CHANGE
   ↓
DESIGN BOUNDARY
   ↓
CODE
   ↓
TEST
   ↓
ATTACK
   ↓
VERIFY
```

На этом занятии особенно важно:

> **Не спрашивать сразу «какой принцип SOLID здесь нарушен?»**

Сначала:

> «Что станет дорого менять?»

И только потом давать название принципу.

---

# 4. 00:00–00:10 — Вход в урок

Я открываю JobFlow.

Говорю:

> «На прошлом занятии мы сделали concurrency. У нас уже есть workers, queue, repository и service.
>
> Сегодня я хочу сделать с вами одну неприятную вещь: специально испортить архитектуру.
>
> А потом мы попробуем сделать её снова хорошей — не потому, что так написано в SOLID, а потому что изменение стало дорогим.»

Показываю:

```text
HTTP
  ↓
JobService
  ↓
Repository
```

Потом говорю:

> «Представим, что разработчик решил сэкономить время и написал всё в одном месте.»

Показываю:

```text
HTTP handler
 ├── decode JSON
 ├── validation
 ├── business rules
 ├── PostgreSQL
 ├── Kafka
 └── logging
```

Спрашиваю:

> «Что здесь тебя настораживает?»

<details style="display: inline;"><summary>Помощь</summary>Не ищите пока названия принципов. Подумайте, сколько разных причин может заставить этот код измениться.</details>

<details style="display: inline;"><summary>Ответ</summary>Здесь смешаны transport, validation, business rules, persistence, messaging и observability. Они могут изменяться независимо, поэтому одна часть системы затрагивает слишком много независимых изменений.</details>

---

# 5. 00:10–00:25 — Плакат 04A: SOLID

Открываем `04A`.

Я говорю:

> «Теперь познакомимся с языком, которым инженеры описывают такие проблемы.»

---

## S — SINGLE RESPONSIBILITY

> «Что такое S?»

<details style="display: inline;"><summary>Помощь</summary>Спросите: какая причина изменения должна существовать у одного модуля?</details>

<details style="display: inline;"><summary>Ответ</summary>Single Responsibility — одна ось ответственности, то есть одна основная причина для изменения.</details>

Я специально спрашиваю:

> «Одна функция?»

<details><summary>Ответ</summary>Нет. SRP не означает «ровно одна функция». Речь об оси ответственности и причине изменения.</details>

---

## O — OPEN / CLOSED

> «Что означает O?»

<details style="display: inline;"><summary>Ответ</summary>Стабильную политику следует расширять новым поведением так, чтобы не превращать каждое новое поведение в постоянное переписывание стабильного ядра.</details>

И сразу:

> «Это означает, что код вообще нельзя менять?»

<details><summary>Ответ</summary>Нет. OCP не запрещает изменения. Он управляет направлением и стоимостью изменений там, где граница действительно стабильна.</details>

---

## L — LISKOV SUBSTITUTION

> «Что здесь главное?»

<details><summary>Помощь</summary>Забудьте на минуту про inheritance. Представьте interface и две реализации.</details>

<details><summary>Ответ</summary>Одну реализацию можно заменить другой без нарушения ожиданий потребителя. Контракт должен сохранять совместимую семантику.</details>

Я показываю JobFlow:

```go
type JobRepository interface {
    Get(string) (Job, error)
}
```

и:

```text
MemoryRepository
PostgresRepository
```

Спрашиваю:

> «Что будет нарушением LSP?»

<details><summary>Ответ</summary>Например, если MemoryRepository на отсутствующий Job возвращает `ErrJobNotFound`, а PostgresRepository вместо этого делает panic или возвращает другое несовместимое состояние, потребитель не может безопасно заменить одну реализацию другой.</details>

---

## I — INTERFACE SEGREGATION

> «Почему большой interface может быть проблемой?»

<details><summary>Ответ</summary>Потребитель начинает зависеть от методов, которые ему не нужны. Это увеличивает dependency surface и стоимость изменений.</details>

Показываю:

```go
type Repository interface {
    Add(...)
    Get(...)
    Delete(...)
    List(...)
    Export(...)
    RebuildIndex(...)
}
```

И спрашиваю:

> «Что реально нужно `JobReader`?»

<details><summary>Ответ</summary>Только `Get`, если он действительно использует лишь чтение.</details>

---

## D — DEPENDENCY INVERSION

Показываю:

```text
Service → PostgreSQL
```

и спрашиваю:

> «Что здесь является деталью?»

<details><summary>Ответ</summary>PostgreSQL — инфраструктурная деталь.</details>

> «Что должно знать о ней business policy?»

<details><summary>Ответ</summary>Business policy не должна зависеть непосредственно от конкретного PostgreSQL adapter.</details>

Потом:

```text
Service → JobRepository ← PostgreSQL
```

И спрашиваю:

> «Где создадим зависимость?»

<details><summary>Ответ</summary>Через constructor injection, например `NewJobService(repo JobRepository)`.</details>

---

# 6. 00:25–00:35 — Interview Drill

Закрываем код и оставляем только плакат.

Я говорю:

> «Теперь у вас 30 секунд. Интервьюер спрашивает: что такое SOLID?»

[Жду ответ.]

Если ответ расплывается, помогаю:

> «Сначала одна фраза, потом пять букв. Не уходите в историю паттернов.»

Эталон:

<details><summary>Ответ</summary>«SOLID — это пять принципов проектирования, помогающих управлять ответственностями, изменениями, контрактами и зависимостями. S — ответственность, O — расширение, L — сохранение ожиданий при подстановке, I — маленькие контракты под потребителей, D — направление зависимостей к стабильной абстракции.»</details>

Потом по каждой букве даю быстрый вопрос.

### S

> «Какой главный вопрос?»

<details><summary>Ответ</summary>Какая причина изменения принадлежит этому модулю?</details>

### O

> «Что контролирует?»

<details><summary>Ответ</summary>Способ расширения стабильной политики.</details>

### L

> «Что нельзя нарушить?»

<details><summary>Ответ</summary>Ожидания и контракт потребителя.</details>

### I

> «Для кого проектируется interface?»

<details><summary>Ответ</summary>Для конкретного потребителя.</details>

### D

> «Куда направляется зависимость?»

<details><summary>Ответ</summary>К стабильному контракту, а не к инфраструктурной детали.</details>

---

# 7. 00:35–00:50 — Плакат 04B: DESIGN ENGINEERING

Открываем второй плакат.

Я говорю:

> «Теперь начинается самое важное. SOLID — это словарь. А этот плакат — способ принимать решения.»

Показываю заголовок.

> «Хороший дизайн делает изменение локальным, зависимость ясной, а сложность оправданной.»

---

# 8. 00:50–01:00 — CARD 01: RESPONSIBILITY

Вопрос:

> «Что должно изменяться вместе?»

<details><summary>Помощь</summary>Найдите причины изменений, а не функции.</details>

<details><summary>Ответ</summary>Связанное поведение с общей причиной изменения имеет смысл держать вместе; независимые оси изменения лучше разделять.</details>

Теперь показываю наш плохой handler:

```text
HTTP handler
 ├── JSON
 ├── validation
 ├── business rules
 ├── SQL
 ├── Kafka
 └── logs
```

И спрашиваю:

> «Назовите причины, по которым этот код может измениться.»

<details><summary>Ответ</summary>Изменение HTTP API/JSON, validation rules, business rules, PostgreSQL schema/queries, Kafka protocol или publishing policy, logging/observability.</details>

Я говорю:

> «Вот мы только что практически обнаружили SRP. Не по названию. По изменению.»

---

# 9. 01:00–01:10 — CARD 02: DEPENDENCY

Теперь:

> «От чего должна зависеть политика?»

Показываю:

```text
JobService
   ↓
PostgreSQL
```

Что здесь плохо?

<details><summary>Ответ</summary>Application/business policy знает конкретную инфраструктурную деталь. Это увеличивает coupling и затрудняет замену, тестирование и изменение storage.</details>

Как изменим?

<details><summary>Ответ</summary>Выделим consumer-side interface и внедрим зависимость через constructor.</details>

Ученик пишет:

```go
type JobRepository interface {
    Get(ctx context.Context, id string) (Job, error)
    Add(ctx context.Context, job Job) error
}
```

И:

```go
type JobService struct {
    repo JobRepository
}
```

---

# 10. 01:10–01:20 — CARD 03: COMPLEXITY

Теперь я даю провокацию.

> «Раз уж мы узнали SOLID, давайте сразу сделаем десять интерфейсов, factory, strategy, registry и generic abstraction.»

Пауза.

И спрашиваю:

> «Хорошая идея?»

<details><summary>Ответ</summary>Нет. Абстракция должна решать реальную проблему изменения, зависимости или тестирования. Иначе она увеличивает complexity без требования.</details>

Смотрим на KISS:

> **Минимальная конструкция, достаточная для требования.**

И спрашиваю:

> «Что нам сегодня реально требуется от repository?»

<details><summary>Ответ</summary>Только операции, которые используются текущей business policy.</details>

YAGNI:

> «Нужно ли нам уже сейчас `Export`, `BulkImport`, `Archive`, `Snapshot`?»

<details><summary>Ответ</summary>Нет, если текущего требования на них нет.</details>

---

# 11. 01:20–01:30 — CARD 04: COHESION / COUPLING

Теперь я показываю два фрагмента.

Вариант A:

```text
JobService
├── create
├── start
├── complete
└── fail
```

Вариант B:

```text
JobService
├── create
├── renderPrometheus
├── encodeJSON
└── reconnectKafka
```

Спрашиваю:

> «Какой вариант имеет выше cohesion?»

<details><summary>Ответ</summary>Первый. Его операции связаны одной доменной ответственностью.</details>

> «А coupling?»

<details><summary>Ответ</summary>Второй вариант сильнее связан с несколькими внешними деталями, поэтому его изменение затрагивает больше зависимостей.</details>

И формула:

> **Cohesion внутри — вверх. Coupling наружу — вниз.**

---

# 12. 01:30–01:40 — CARD 05: CHANGE

Теперь мы перестаём обсуждать термины.

Я говорю:

> «Представим, что завтра вместо PostgreSQL мы используем другой storage. Какие файлы должен потребоваться изменить?»

Пусть ученик отвечает.

Если он говорит:

> «Repository.»

Я уточняю:

> «Только repository? Почему?»

<details><summary>Ответ</summary>Если abstraction boundary корректна, business/application policy не должна знать concrete storage detail. Меняется adapter и wiring, а не доменная политика.</details>

Теперь:

> «А если мы меняем HTTP API?»

<details><summary>Ответ</summary>Основное изменение должно остаться в transport layer и mapping, если business contract не изменился.</details>

Теперь:

> «А если меняется business rule обработки Job?»

<details><summary>Ответ</summary>Изменение должно локализоваться в application/domain policy, а transport/storage adapters по возможности оставаться неизменными.</details>

Вот сейчас плакат начинает работать как инженерная карта.

---

# 13. 01:40–01:50 — CARD 06: CHANGE COST

Теперь:

> «От чего зависит стоимость изменения?»

<details><summary>Помощь</summary>Посмотрите последнюю карточку.</details>

<details><summary>Ответ</summary>От количества затрагиваемых ответственностей, направлений зависимостей и добавленной сложности. Чем шире blast radius изменения, тем выше cost.</details>

Я говорю:

> «И вот теперь у нас появляется настоящее назначение всех этих принципов.»

Не:

> «Нужно выполнить SOLID.»

А:

> **«Нужно уменьшить blast radius изменения с сохранением требуемого поведения».**

---

# 14. 01:50–02:10 — Combat Exercise: Refactor JobFlow

Теперь основная практика.

Я даю намеренно плохой код.

Примерно:

```go
func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateJobRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    if req.Type == "" {
        http.Error(w, "missing type", http.StatusBadRequest)
        return
    }

    db := postgres.Open(...)
    defer db.Close()

    job := Job{
        ID:      uuid.NewString(),
        Type:    req.Type,
        Payload: req.Payload,
        Status:  StatusPending,
    }

    if _, err := db.ExecContext(
        r.Context(),
        "INSERT INTO jobs ...",
    ); err != nil {
        http.Error(w, "database error", http.StatusInternalServerError)
        return
    }

    kafka.Publish("job.created", job)

    log.Printf("created job %s", job.ID)

    json.NewEncoder(w).Encode(job)
}
```

Я говорю:

> «Не переписывайте код. Сначала проведите diagnosis.»

### Вопрос 1

> «Какие независимые причины изменения здесь есть?»

<details><summary>Ответ</summary>HTTP contract, validation, domain creation/business rules, PostgreSQL persistence, Kafka publishing и logging.</details>

### Вопрос 2

> «Какая часть является policy?»

<details><summary>Ответ</summary>Правила создания Job и application decision.</details>

### Вопрос 3

> «Какие части являются details?»

<details><summary>Ответ</summary>HTTP, JSON decoding, PostgreSQL, Kafka и конкретный logging mechanism.</details>

### Вопрос 4

> «Что должно быть boundary?»

<details><summary>Ответ</summary>Transport, application/business logic и infrastructure adapters должны иметь отдельные ответственности и contracts.</details>

---

# 15. 02:10–02:25 — Строим новую границу

Теперь ученик должен написать:

```go
type JobRepository interface {
    Add(ctx context.Context, job Job) error
}
```

И:

```go
type EventPublisher interface {
    Publish(ctx context.Context, event JobCreated) error
}
```

А service:

```go
type JobService struct {
    repo      JobRepository
    publisher EventPublisher
}
```

Я спрашиваю:

> «Почему service не должен иметь `*sql.DB`?»

<details><summary>Ответ</summary>Потому что SQL database — infrastructure detail. Service должен зависеть от поведения, которое ему нужно, а конкретный adapter подключается снаружи.</details>

> «Почему publisher тоже отдельный interface?»

<details><summary>Ответ</summary>Kafka publishing имеет собственную инфраструктурную ответственность и собственную причину изменения. Service зависит только от нужного ему поведения публикации.</details>

---

# 16. 02:25–02:35 — Attack: Interface слишком большой

Теперь я специально предлагаю:

```go
type Infrastructure interface {
    SaveJob(...)
    PublishKafka(...)
    ReadConfig(...)
    WriteLog(...)
    GetRedis(...)
    SendHTTP(...)
}
```

Спрашиваю:

> «Принимаем?»

<details><summary>Ответ</summary>Нет. Это объединяет несвязанные зависимости и увеличивает coupling. Consumer должен зависеть от минимальных contracts, которые нужны именно ему.</details>

> «Какой принцип мы сейчас использовали?»

<details><summary>Ответ</summary>ISP. Но важнее то, что мы уменьшили dependency surface для consumer.</details>

---

# 17. 02:35–02:45 — Attack: «Сделаем всё через interfaces»

Я говорю:

> «А теперь я предлагаю обернуть вообще всё в interfaces: `Clock`, `UUID`, `JSONEncoder`, `Logger`, `JobFactory`, `Validator`, `Slice`, `Map`.»

Спрашиваю:

> «Что скажете?»

<details><summary>Ответ</summary>Это premature abstraction. Если зависимости не создают значимой границы изменения или тестирования, interfaces увеличивают complexity без полезной гарантии. KISS/YAGNI требуют минимальной конструкции, достаточной для реальной потребности.</details>

Это очень важный combat moment.

Я хочу, чтобы после курса человек не стал **SOLID-zealot**.

---

# 18. 02:45–02:55 — Tests как Evidence

Теперь мы должны убедиться, что refactoring не изменил behavior.

Я спрашиваю:

> «Что мы обязаны сохранить?»

<details><summary>Ответ</summary>Observable behavior и domain invariants: корректное создание Job, validation, сохранение, публикация соответствующего события и error semantics там, где они являются контрактом.</details>

Добавляем tests.

Минимум:

```text
Create valid job
Invalid job
Repository failure
Publisher failure
Existing error semantics
```

Запускаем:

```bash
go test ./...
```

Если есть concurrency code:

```bash
go test -race ./...
```

И говорим:

> «Архитектурный refactoring без сохранения behavior — это не доказанный refactoring.»

---

# 19. 02:55–03:05 — Финальный Architecture Attack

Я теперь меняю требования.

### Change 1

> «PostgreSQL меняется.»

<details><summary>Ответ</summary>Меняется adapter/wiring. Application policy по возможности остаётся без изменений.</details>

### Change 2

> «HTTP становится gRPC.»

<details><summary>Ответ</summary>Меняется transport adapter. Application contract сохраняется.</details>

### Change 3

> «Публикация становится не Kafka, а другой messaging system.»

<details><summary>Ответ</summary>Меняется EventPublisher adapter и wiring; business logic не должна знать конкретную messaging implementation.</details>

### Change 4

> «Меняется business rule перехода Job в `running`.»

<details><summary>Ответ</summary>Меняется application/domain policy, а transport и infrastructure остаются максимально неизменными.</details>

### Change 5

> «Добавили новую реализацию Processor.»

<details><summary>Ответ</summary>Если boundary спроектирована правильно, добавляется новая implementation без постоянного переписывания стабильной policy. Это практическое проявление OCP.</details>

---

# 20. 03:05–03:15 — Итоговая проверка

Теперь 10 вопросов.

### 1. Что такое SRP?

<details><summary>Ответ</summary>Одна ось ответственности и изменения для модуля.</details>

### 2. Что такое OCP?

<details><summary>Ответ</summary>Стабильную политику следует расширять новым поведением, минимизируя постоянное переписывание существующего ядра.</details>

### 3. Что такое LSP?

<details><summary>Ответ</summary>Подстановка реализации не нарушает ожиданий и контракта потребителя.</details>

### 4. Что такое ISP?

<details><summary>Ответ</summary>Потребитель зависит только от нужного ему поведения и не вынужден зависеть от ненужных методов.</details>

### 5. Что такое DIP?

<details><summary>Ответ</summary>Высокоуровневая policy не должна зависеть от infrastructure detail; зависимости направляются к стабильным контрактам.</details>

### 6. Где обычно определяем interface в Go?

<details><summary>Ответ</summary>Часто рядом с consumer — на стороне потребителя, исходя из требуемого ему поведения.</details>

### 7. Что такое cohesion?

<details><summary>Ответ</summary>Насколько тесно связаны между собой элементы внутри компонента.</details>

### 8. Что такое coupling?

<details><summary>Ответ</summary>Насколько изменение одного компонента затрагивает другие компоненты через зависимости.</details>

### 9. Что такое KISS?

<details><summary>Ответ</summary>Используй минимальную конструкцию, достаточную для требования.</details>

### 10. Что такое DRY?

<details><summary>Ответ</summary>Одно знание должно иметь один источник истины. DRY относится к дублированию знания, а не к любому совпадающему тексту.</details>

---

# 21. Финальный вопрос

Я закрываю оба плаката.

И спрашиваю:

> «Допустим, я вижу огромный interface. Что я должен спросить первым: “Какой здесь принцип SOLID нарушен?” или что-то другое?»

<details><summary>Помощь</summary>Вернитесь к первому плакату курса.</details>

<details><summary>Ответ</summary>Сначала нужно определить требование, invariant, boundary и стоимость изменения. Затем выбрать минимальный механизм и только потом назвать соответствующий принцип, если он помогает объяснить решение.</details>

И это главный результат урока.

---

# 22. Финальные 5 минут — Mentor Closing

Я говорю:

> «Сегодня мы научились двум разным вещам.
>
> Первое — говорить на языке SOLID.
>
> Если вас спросили на интервью, вы можете быстро и точно объяснить S, O, L, I и D.
>
> Но второе гораздо важнее.
>
> Вы увидели, что SOLID не является чеклистом качества.
>
> Мы не ставим галочку возле SRP.
>
> Мы сначала смотрим, что изменяется.
>
> Потом смотрим, где проходят зависимости.
>
> Потом спрашиваем, нужна ли вообще новая абстракция.
>
> И только после этого выбираем механизм.
>
> Хороший design — это не много интерфейсов и не много слоёв.
>
> Хороший design — это когда изменение остаётся настолько локальным, насколько позволяет требование.»

После паузы:

> «На следующем занятии мы наконец перестанем хранить Job только в памяти.
>
> Теперь наши invariants должны переживать процессы, worker crashes и конкурентные изменения состояния.
>
> Поэтому следующий бой — **DATA & TRANSACTIONS**.»

---

# 23. Домашнее задание — JobFlow v2

Здесь домашнее задание должно **закрепить оба плаката**, а не просто дать написать ещё код.

## Часть A — Refactoring

Перестроить JobFlow так, чтобы:

```text
HTTP / transport
        ↓
JobService
        ↓
small interfaces
        ↓
infrastructure
```

не смешивали свои ответственности.

Минимальные контракты:

```go
type JobRepository interface {
    Add(ctx context.Context, job Job) error
    Get(ctx context.Context, id string) (Job, error)
}

type EventPublisher interface {
    Publish(ctx context.Context, event JobCreated) error
}
```

Конкретные implementations должны подключаться через constructor injection.

---

## Часть B — Change Exercise

После refactoring реализовать **два изменения**:

### Change 1

Добавить второй `JobProcessor`.

Например:

```text
EmailProcessor
WebhookProcessor
```

Не переписывая существующую business policy без необходимости.

### Change 2

Добавить альтернативный `JobRepository`.

Например:

```text
MemoryJobRepository
FakeJobRepository
```

И проверить, что `JobService` не знает конкретную реализацию.

---

## Часть C — Tests

Тесты должны показать:

```text
service works with memory repository
service works with fake repository
publisher can be substituted
business rules do not depend on infrastructure
```

---

## Часть D — Design Note

Написать **не больше 10 строк**:

> Почему `JobRepository` определён на стороне consumer?

> Почему `EventPublisher` отдельный interface?

> Почему не сделали один `Infrastructure` interface?

> Какое изменение теперь локализовано лучше, чем было раньше?

Это проверяет **решение**, а не объём написанного кода.

---

# 24. Acceptance Criteria

## PASS

Ученик:

* объясняет пять букв SOLID без подсказки;
* приводит Go-specific пример каждого принципа;
* отличает principle от implementation pattern;
* может назвать причину изменения компонента;
* видит coupling и cohesion;
* определяет consumer-side interface;
* вводит abstraction только при наличии boundary/change case;
* умеет заменить infrastructure implementation без изменения policy;
* сохраняет tests после refactoring;
* может объяснить, какую стоимость изменения снизило каждое решение.

## FAIL

Если ученик говорит:

> «Нам нужен interface, потому что DIP.»

и не может ответить:

> **«Какое изменение или dependency мы этим изолируем?»**

то design decision пока не считается усвоенным.

---

# 25. Результат урока

После `04.DESIGN` JobFlow должен выглядеть концептуально так:

```text
                 TRANSPORT
                     │
                     ▼
               JOB SERVICE
              /           \
             ▼             ▼
      JOB REPOSITORY   EVENT PUBLISHER
             │             │
             ▼             ▼
          STORAGE        MESSAGING
```

Но ученик должен понимать:

> **это не «правильная архитектура потому, что SOLID».**

Это следствие конкретных изменений и boundaries:

```text
HTTP changes
     ↓
transport boundary

Storage changes
     ↓
repository boundary

Messaging changes
     ↓
publisher boundary

Business rules change
     ↓
policy remains local
```

Следующий урок получает уже **не монолитный JobFlow, а систему с явными boundaries и contracts**, что естественно подводит к `DATA & TRANSACTIONS`.
