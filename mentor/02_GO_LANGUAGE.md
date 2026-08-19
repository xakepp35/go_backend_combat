# Go Backend Combat — 01. LANGUAGE

Плакат у нас перед глазами. Я не собираюсь сейчас читать вам справочник Go. Наша задача проще: за несколько минут собрать в голове модель языка, а потом сразу начать писать JobFlow.

Смотрите на плакат.

Здесь шесть блоков. Не нужно сейчас запоминать всё, что внутри них написано. Мне важно, чтобы вы увидели последовательность.

Первый вопрос: **как выразить состояние данных?**

Второй: **как выразить действие и его результат?**

Третий: **где живёт поведение типа?**

Четвёртый: **от чего должен зависеть потребитель?**

Пятый: **как выразить неуспешный результат?**

Шестой: **как доказать ожидаемое поведение?**

То есть мы движемся примерно так: данные, действие, поведение, контракт, ошибка, доказательство.

Я специально не начинаю с синтаксиса. Мы сначала задаём модель, а потом синтаксис будет просто способом её записать.

---

Теперь представьте, что вы до этого много писали на Python, JavaScript или PHP.

Я показываю вам вот это:

```go
type Job struct {
    ID      string
    Type    string
    Payload []byte
}
```

Что это такое?

<details><summary>Помощь</summary>Не думайте пока про Go. Что эта конструкция описывает с точки зрения предметной области?</details>

<details><summary>Ответ</summary>Это пользовательский тип, который объединяет связанные данные одной сущности Job. В Go для этого используется struct.</details>

Вот. Это наша первая конструкция.

Я хочу, чтобы вы перестали смотреть на `struct` как на очередное ключевое слово. Смысл выше синтаксиса: у нас есть сущность, у неё есть состояние, мы делаем это состояние явным типом.

---

Теперь маленький вопрос.

Почему здесь не написать просто набор отдельных переменных?

```go
jobID string
jobType string
payload []byte
```

<details><summary>Помощь</summary>Что произойдёт, если таких Job станет тысяча?</details>

<details><summary>Ответ</summary>Мы потеряем целостную модель сущности. `struct` позволяет собрать связанные данные в один тип и передавать Job как единое значение.</details>

И это одна из базовых идей Go: **делать модель состояния явной**.

---

Теперь я хочу показать ещё одну особенность Go.

```go
var count int = 10
```

А теперь:

```go
count := 10
```

Что изменилось с точки зрения типов?

<details><summary>Помощь</summary>Нас интересует не запись, а то, что делает compiler.</details>

<details><summary>Ответ</summary>Тип по-прежнему статический. Во втором случае compiler выводит `int` из значения `10`. Type inference не отключает static typing.</details>

И это хороший общий принцип Go:

> **Можем писать меньше текста, не теряя информацию, которую знает compiler.**

---

Теперь ещё один маленький эксперимент.

```go
var name string
var count int
var active bool
```

Каковы значения?

<details><summary>Помощь</summary>Мы ничего не присвоили явно.</details>

<details><summary>Ответ</summary>`name == ""`, `count == 0`, `active == false`.</details>

Да.

У каждого типа есть zero value.

Я хочу, чтобы эта идея осталась у вас в голове именно как инженерная конструкция:

> **Хороший тип имеет понятное состояние по умолчанию.**

Не просто «это такая фича Go», а часть проектирования API.

---

Теперь возвращаемся к JobFlow.

Наш Job должен иметь статус.

Можно сделать:

```go
Status string
```

Но я хочу спросить: а зачем мне тогда вообще отдельный тип?

```go
type JobStatus string
```

Какой смысл?

<details><summary>Помощь</summary>Что меняется в контракте функции `SetStatus`, если она получает `string` или `JobStatus`?</details>

<details><summary>Ответ</summary>С `JobStatus` тип прямо выражает смысл значения. API становится более точным: `SetStatus(status JobStatus)` ожидает именно статус Job, а не произвольную строку.</details>

И вот здесь появляется одна из моих любимых идей:

> **Хороший тип делает неправильное состояние труднее выразить.**

Добавляем:

```go
const (
    StatusPending   JobStatus = "pending"
    StatusRunning   JobStatus = "running"
    StatusCompleted JobStatus = "completed"
    StatusFailed    JobStatus = "failed"
)
```

Теперь у нашей доменной модели уже появляется vocabulary.

Не строки вообще.

А конкретные состояния.

---

Теперь давайте соберём Job.

Я бы предложил такую модель:

```go
type Job struct {
    ID        string
    Type      string
    Payload   []byte
    Status    JobStatus
    CreatedAt time.Time
}
```

Я специально не буду сразу объяснять каждое поле.

Попробуйте сами: почему здесь `[]byte`, а не `string`?

<details><summary>Помощь</summary>Подумайте, что нам может прийти в payload и хотим ли мы навязать ему конкретное текстовое представление.</details>

<details><summary>Ответ</summary>`[]byte` позволяет хранить произвольное бинарное содержимое и не привязывает доменную модель к конкретной кодировке или представлению текста.</details>

Хорошо.

Теперь вопрос важнее.

Какая часть этого `Job` — состояние, а какая часть должна стать поведением?

<details><summary>Помощь</summary>Поля сами по себе ничего не решают. Что Job должен уметь делать?</details>

<details><summary>Ответ</summary>Например, переходить в `running`, `completed`, `failed`, проверять, завершён ли он, валидироваться.</details>

Отлично. Значит, мы пришли к следующему блоку плаката.

**METHODS. Где живёт поведение типа?**

---

Допустим, хотим проверить, завершена ли задача.

Я могу написать обычную функцию:

```go
func IsFinished(job Job) bool
```

Но могу написать method:

```go
func (j Job) IsFinished() bool {
    return j.Status == StatusCompleted ||
        j.Status == StatusFailed
}
```

Какой вариант вам кажется естественнее для модели Job?

<details><summary>Помощь</summary>Спросите себя: кто является владельцем этого поведения?</details>

<details><summary>Ответ</summary>Method естественнее, потому что поведение принадлежит типу `Job` и вызывается как `job.IsFinished()`.</details>

Да.

И здесь важна не синтаксическая красота.

Мы говорим:

> **Job владеет своим поведением.**

Это уже маленькая архитектурная граница.

---

Теперь я сделаю с вами маленький эксперимент, который очень важен для Go.

Пишем:

```go
func (j Job) Complete() {
    j.Status = StatusCompleted
}
```

Потом:

```go
job.Complete()
fmt.Println(job.Status)
```

Что мы ожидаем увидеть?

<details><summary>Помощь</summary>Посмотрите на receiver: `Job` или `*Job`?</details>

<details><summary>Ответ</summary>Мы ожидаем, что `job.Status` останется прежним, например `pending`, потому что value receiver работает с копией значения.</details>

Запускаем.

Получаем то, что и ожидали.

Теперь меняем:

```go
func (j *Job) Complete() {
    j.Status = StatusCompleted
}
```

Что ожидаем теперь?

<details><summary>Ответ</summary>`job.Status` станет `completed`, потому что pointer receiver позволяет изменить исходный объект через указатель.</details>

Вот это я хочу, чтобы вы почувствовали руками.

Не надо учить это как фразу «звёздочка значит меняет».

Нужно понимать:

> **Value receiver работает со значением. Pointer receiver даёт доступ к исходному состоянию.**

И ещё одна важная вещь.

Pointer в Go — это не «ручное управление памятью», как в C.

Память управляется runtime.

Мы сейчас обсуждаем прежде всего **семантику значения и доступа к состоянию**.

---

Теперь двигаемся дальше.

У нас будет не один Job, а много.

Как представить несколько Job?

<details><summary>Помощь</summary>Нам нужен упорядоченный набор элементов одного типа.</details>

<details><summary>Ответ</summary>Slice: `[]Job`.</details>

Например:

```go
jobs := []Job{
    NewJob("job-1", "email", nil),
    NewJob("job-2", "sms", nil),
    NewJob("job-3", "email", nil),
}
```

Какой здесь интересный вопрос?

```go
len(jobs)
cap(jobs)
```

Что означают эти два значения?

<details><summary>Ответ</summary>`len` — текущее количество элементов. `cap` — ёмкость backing array от текущего начала slice.</details>

Хорошо.

Теперь:

```go
jobs = append(jobs, job)
```

Почему мы присваиваем результат обратно?

<details><summary>Помощь</summary>Что может произойти с backing array при нехватке capacity?</details>

<details><summary>Ответ</summary>`append` может выделить новый backing array и вернуть новый slice header, поэтому результат нужно присвоить обратно.</details>

Вот эта маленькая деталь потом очень хорошо всплывёт в разговорах про производительность и память.

---

Теперь задача.

Найдите Job по ID:

```go
func FindJob(jobs []Job, id string) *Job
```

Я не даю реализацию. Напишите сами.

[Жду, пока ученик пишет.]

Пока он пишет, я смотрю не только на результат. Мне интересно, как он мыслит.

Если он пишет `for _, job := range jobs`, я спрошу:

> «Если я хочу вернуть указатель на найденный элемент, что здесь важно?»

<details><summary>Помощь</summary>Нам нужен доступ к элементу самого slice, а не только к копии переменной цикла.</details>

<details><summary>Ответ</summary>Нужно идти по индексу: `for i := range jobs`, а затем вернуть `&jobs[i]`.</details>

Ожидаем:

```go
for i := range jobs {
    if jobs[i].ID == id {
        return &jobs[i]
    }
}

return nil
```

Это одновременно показывает range, индекс, pointer и nil.

---

Теперь нам нужен быстрый доступ по ID.

Какую структуру данных возьмём?

<details><summary>Помощь</summary>Нам нужно `JobID → Job`.</details>

<details><summary>Ответ</summary>Map: `map[string]Job`.</details>

Пишем:

```go
jobsByID := map[string]Job{}
```

Добавляем:

```go
jobsByID[job.ID] = job
```

Теперь самое интересное.

```go
job, ok := jobsByID["job-1"]
```

Зачем здесь второй результат?

<details><summary>Помощь</summary>Что вернёт map, если ключ отсутствует?</details>

<details><summary>Ответ</summary>Map вернёт zero value значения и `ok == false`. Второй результат позволяет отличить отсутствующий ключ от существующего значения, совпадающего с zero value.</details>

Это idiom, который я хочу, чтобы вы узнали сразу.

```go
value, ok := m[key]
```

Запомнить эту конструкцию недостаточно.

Нужно понимать, какую неопределённость она снимает.

---

И теперь я хочу сделать из этого маленькую часть настоящего backend.

У нас будет:

```go
type MemoryJobRepository struct {
    jobs map[string]Job
}
```

Почему здесь `map`?

<details><summary>Ответ</summary>Потому что нам нужен доступ к Job по идентификатору.</details>

Пишем constructor:

```go
func NewMemoryJobRepository() *MemoryJobRepository {
    return &MemoryJobRepository{
        jobs: make(map[string]Job),
    }
}
```

Зачем `make`?

<details><summary>Помощь</summary>Вспомните, что такое zero value map.</details>

<details><summary>Ответ</summary>Нулевой `map` равен `nil`. Читать из него можно, но запись в nil map приводит к panic. `make` создаёт готовую для записи map.</details>

Отлично.

Это очень хороший пример того, почему zero value — не просто абстрактная теория.

---

Теперь `Add`.

Какой контракт вы бы предложили?

```go
func (r *MemoryJobRepository) Add(...)
```

Что возвращать?

<details><summary>Помощь</summary>Представьте, что позже это будет PostgreSQL. Что может пойти не так при сохранении?</details>

<details><summary>Ответ</summary>Операция может завершиться ошибкой, поэтому лучше с самого начала заложить `error` в контракт.</details>

Тогда:

```go
func (r *MemoryJobRepository) Add(job Job) error {
    r.jobs[job.ID] = job
    return nil
}
```

Сейчас `nil`, потому что memory storage пока не имеет реальной причины для ошибки.

Но контракт уже допускает ошибку.

И теперь мы приходим к пятому блоку плаката.

**ERRORS. Как выразить неуспешный результат?**

---

Представьте:

```go
job, err := repo.Get(id)
```

Как мы сообщим вызывающему коду, что Job не существует?

<details><summary>Помощь</summary>Нам нужен machine-readable признак, который можно проверить программно.</details>

<details><summary>Ответ</summary>Вернуть специальную ошибку, например `ErrJobNotFound`.</details>

Пишем:

```go
var ErrJobNotFound = errors.New("job not found")
```

И:

```go
func (r *MemoryJobRepository) Get(id string) (Job, error) {
    job, ok := r.jobs[id]
    if !ok {
        return Job{}, ErrJobNotFound
    }

    return job, nil
}
```

Теперь следующий вопрос.

Зачем не сделать просто:

```go
return Job{}, errors.New("job not found")
```

каждый раз?

<details><summary>Помощь</summary>Нам нужно, чтобы программа могла распознать именно этот класс ошибки, даже если текст изменился.</details>

<details><summary>Ответ</summary>Нам нужна стабильная error identity, например `ErrJobNotFound`, которую можно проверять через `errors.Is`.</details>

Тогда можно добавить контекст:

```go
return Job{}, fmt.Errorf(
    "get job %s: %w",
    id,
    ErrJobNotFound,
)
```

И дальше:

```go
if errors.Is(err, ErrJobNotFound) {
    // not found
}
```

Вот это очень важный Go-паттерн:

> **Сообщение ошибки объясняет проблему. Identity ошибки позволяет программе принять решение.**

---

Теперь я хочу специально остановиться на `interface`.

Пока у нас есть `MemoryJobRepository`.

Но мы знаем, что дальше в курсе у нас будет PostgreSQL.

Как сделать так, чтобы `JobService` не зависел от конкретного `MemoryJobRepository`?

<details><summary>Помощь</summary>Нам нужно зависеть от того, что сервис умеет делать, а не от того, где физически лежат данные.</details>

<details><summary>Ответ</summary>Определить небольшой interface с необходимым поведением repository.</details>

Пишем:

```go
type JobRepository interface {
    Add(Job) error
    Get(string) (Job, error)
}
```

Теперь:

```go
type JobService struct {
    repo JobRepository
}
```

И:

```go
func NewJobService(repo JobRepository) *JobService {
    return &JobService{repo: repo}
}
```

И вот сейчас я задаю очень важный вопрос.

Где мы написали:

```text
implements JobRepository
```

<details><summary>Ответ</summary>Нигде. В Go interface реализуется неявно: тип удовлетворяет interface, если имеет необходимый method set.</details>

И это очень характерная вещь для Go.

Маленький контракт задаёт потребитель.

А конкретная реализация просто соответствует этому контракту.

---

Теперь провокационный вопрос.

Нужно ли нам добавлять в `JobRepository`:

```go
Delete(...)
List(...)
Update(...)
GetByType(...)
```

если `JobService` пока использует только `Add` и `Get`?

<details><summary>Помощь</summary>Спросите себя: какие именно возможности потребителю реально нужны?</details>

<details><summary>Ответ</summary>Нет. Большой интерфейс расширяет dependency surface без необходимости. Лучше минимальный контракт под конкретного потребителя.</details>

Вот это маленькая, но очень важная привычка.

> **Interface описывает необходимое поведение, а не все возможности реализации.**

Мы ещё отдельно разберём это на Design Principles, но уже сейчас вы должны это почувствовать на реальном коде.

---

Теперь ещё одна небольшая вещь.

У нас будут разные виды Job: email, webhook, возможно, HTTP callback.

Я не хочу делать огромную иерархию классов.

Покажу композицию:

```go
type EmailJob struct {
    Job
    To string
}
```

И:

```go
type WebhookJob struct {
    Job
    URL string
}
```

Что здесь происходит?

<details><summary>Помощь</summary>Что у нас общего и что специфично?</details>

<details><summary>Ответ</summary>Общее состояние и поведение представлены через `Job`, а конкретный вид задачи композиционно добавляет свои поля и поведение.</details>

Это одна из причин, почему Go очень хорошо подходит для backend:

> **Мы строим поведение из небольших типов и контрактов, а не из глубокой иерархии наследования.**

---

Теперь осталась последняя карточка: **TESTING**.

Я хочу, чтобы вы привыкли к следующему порядку.

Не:

> написал код → надеюсь.

А:

> написал поведение → сформулировал ожидание → проверил → получил evidence.

Давайте создадим тест для `Complete`.

Что здесь важно доказать?

<details><summary>Помощь</summary>Нас интересует не реализация метода, а наблюдаемое поведение после вызова.</details>

<details><summary>Ответ</summary>После `Complete()` статус Job должен стать `StatusCompleted`.</details>

Пишем:

```go
func TestJobComplete(t *testing.T) {
    job := NewJob("job-1", "email", nil)

    job.Complete()

    if job.Status != StatusCompleted {
        t.Fatalf(
            "expected %q, got %q",
            StatusCompleted,
            job.Status,
        )
    }
}
```

Запускаем:

```bash
go test ./...
```

Зелёный тест — это уже evidence.

Но я хочу, чтобы вы не превратили это в магическое:

> «Есть тест — значит всё правильно».

Что именно доказал этот тест?

<details><summary>Помощь</summary>Он проверяет только одно конкретное поведение.</details>

<details><summary>Ответ</summary>Он подтверждает конкретный контракт: вызов `Complete()` приводит Job к состоянию `StatusCompleted`. Он не доказывает правильность всей системы.</details>

Вот это очень важная граница.

---

Теперь сделаем несколько вариантов одного поведения.

Например validation.

У нас есть:

```go
func ValidateJob(job Job) error
```

Какие случаи нам нужно проверить?

<details><summary>Помощь</summary>Начните с одного валидного Job и нескольких очевидно невалидных состояний.</details>

<details><summary>Ответ</summary>Валидный Job, Job без ID, Job без Type, возможно другие нарушения доменного контракта.</details>

Здесь мы можем собрать table-driven test.

Например:

```go
tests := []struct {
    name    string
    job     Job
    wantErr bool
}{
    {
        name: "valid",
        job: Job{
            ID:   "job-1",
            Type: "email",
        },
        wantErr: false,
    },
    {
        name: "missing id",
        job: Job{
            Type: "email",
        },
        wantErr: true,
    },
}
```

Зачем вообще превращать cases в таблицу?

<details><summary>Помощь</summary>Спросите себя: что у тестовых случаев общего?</details>

<details><summary>Ответ</summary>У них одна и та же логика проверки, меняются только входные данные и ожидаемый результат. Поэтому cases удобно представить как данные.</details>

Отлично.

Это уже начинает выглядеть как idiomatic Go.

---

Теперь я хочу, чтобы вы сами открыли плакат и прошли его от начала до конца.

Не по терминам.

По вопросам.

Что находится на первом блоке?

<details><summary>Ответ</summary>Types — как выразить состояние данных.</details>

Что на втором?

<details><summary>Ответ</summary>Functions — как выразить действие и его результат.</details>

Третий?

<details><summary>Ответ</summary>Methods — где живёт поведение типа.</details>

Четвёртый?

<details><summary>Ответ</summary>Interfaces — от чего должен зависеть потребитель.</details>

Пятый?

<details><summary>Ответ</summary>Errors — как выразить неуспешный результат.</details>

Шестой?

<details><summary>Ответ</summary>Testing — как доказать ожидаемое поведение.</details>

Хорошо.

Теперь вы уже видите, что это не шесть несвязанных тем.

Мы начали с данных.

Из данных получили поведение.

Из поведения получили контракты.

Из контрактов получили ошибки как часть API.

И затем закрепили ожидаемое поведение тестом.

---

Теперь давайте не будем больше разговаривать о Go вообще.

Будем писать JobFlow.

Открываем пустой репозиторий.

Я хочу, чтобы вы сами создали:

```text
jobflow/
    go.mod
    main.go
```

И начали с:

```bash
go mod init example.com/jobflow
```

Не копируйте за мной всё подряд. На этом месте я хочу посмотреть, что вы помните.

Что нужно сделать первым?

<details><summary>Помощь</summary>Сначала создаём модуль. Потом минимальную программу, которую compiler и runtime смогут запустить.</details>

<details><summary>Ответ</summary>`go mod init example.com/jobflow`, затем package `main`, функция `main()`, после чего `go run .` и `go build .`.</details>

Создаём:

```go
package main

import "fmt"

func main() {
    fmt.Println("JobFlow")
}
```

Запускаем.

Получаем:

```text
JobFlow
```

Теперь я специально ломаю:

```go
fmt.println("JobFlow")
```

Что должен сделать compiler?

<details><summary>Помощь</summary>Посмотрите на регистр.</details>

<details><summary>Ответ</summary>Он сообщит, что `fmt.println` не является экспортируемым именем `Println`. В Go регистр имени участвует в visibility, а `Println` — exported identifier.</details>

И это первая маленькая проверка того, что compiler уже работает как наш помощник.

---

Теперь переходим к настоящей работе.

Я хочу, чтобы мы не писали весь JobFlow сверху вниз по моему шаблону.

Сейчас вы сами создаёте `Job`.

Подумайте: какие данные нам действительно нужны, чтобы описать задачу?

<details><summary>Помощь</summary>Минимум: идентификатор, тип, полезная нагрузка и текущее состояние. Временная метка понадобится для backend-логики.</details>

<details><summary>Ответ</summary>Например:
`ID string`,
`Type string`,
`Payload []byte`,
`Status JobStatus`,
`CreatedAt time.Time`.</details>

Пишите.

[Жду, пока ученик пишет.]

Теперь создаём `JobStatus`.

Зачем собственный тип?

[Жду ответ.]

Теперь методы.

Какое поведение нам нужно первым?

[Жду ответ.]

Если ученик предлагает `Complete`, задаю следующий вопрос:

> «Какой receiver здесь выберешь — value или pointer? Почему?»

[Жду ответ.]

Потом просим написать.

---

Когда Job уже существует, я говорю:

> «Теперь у нас есть объект. Но объект без проверки — просто набор полей. Как мы поймём, что Job вообще валиден?»

Жду.

Если появляется идея validation:

```go
func ValidateJob(job Job) error
```

продолжаем.

Если ученик начинает сразу писать десяток условий, останавливаю:

> «Стоп. Сначала назови контракт. Какие состояния мы вообще считаем допустимыми?»

<details><summary>Помощь</summary>Сначала правила, потом код.</details>

<details><summary>Ответ</summary>Например: ID не пустой, Type не пустой, начальный статус — `pending`, payload должен соответствовать требованиям конкретной задачи.</details>

---

Теперь говорю:

> «Отлично. У нас есть модель, поведение и validation. Теперь нам нужно где-то держать несколько Job. Какой минимальный storage мы возьмём прямо сейчас, не притаскивая PostgreSQL?»

Жду.

<details><summary>Ответ</summary>Простейший in-memory repository на `map[string]Job`.</details>

Пишем:

```go
type MemoryJobRepository struct {
    jobs map[string]Job
}
```

Я снова спрашиваю:

> «Почему здесь map, а не slice?»

<details><summary>Ответ</summary>Основной lookup у нас по ID, поэтому map естественно выражает `JobID → Job` и даёт эффективный доступ по ключу.</details>

Потом:

```go
type JobRepository interface {
    Add(Job) error
    Get(string) (Job, error)
}
```

И я сразу спрашиваю:

> «Зачем interface уже сейчас, если у нас всего одна реализация?»

<details><summary>Помощь</summary>Вспомните будущую PostgreSQL и consumer-side contract.</details>

<details><summary>Ответ</summary>Чтобы application logic зависела от необходимого поведения repository, а конкретное хранение можно было заменить без переписывания потребителя.</details>

---

Теперь я специально создаю маленький конфликт.

> «У нас сейчас одна implementation. Значит, можно же вообще без interface?»

Пусть ученик защищает своё решение.

Если он отвечает только «так принято», возвращаю:

> «Нет. Это не аргумент. Какую стоимость изменения или dependency surface мы сейчас уменьшаем?»

<details><summary>Ответ</summary>Interface имеет смысл, когда есть граница, по которой потребитель должен быть изолирован от конкретной реализации, или когда abstraction помогает тестированию и замене детали. Его не следует добавлять исключительно ради шаблона.</details>

Это очень важная точка. Пусть он научится **защищать решение причинностью, а не стилем**.

---

Теперь тест.

Я говорю:

> «У нас уже достаточно кода. Не продолжайте. Сначала докажем одну вещь».

Что тестируем?

<details><summary>Ответ</summary>Минимально полезное observable behavior: например, `Complete()` меняет статус на `completed`, `Get()` возвращает добавленный Job, а отсутствующий ID даёт `ErrJobNotFound`.</details>

Пишем первый тест.

Запускаем:

```bash
go test ./...
```

Если зелёный:

> «Хорошо. Мы только что впервые замкнули цикл плаката.»

И показываю глазами:

> **TYPES → METHODS → INTERFACE → ERRORS → TESTING**

Затем возвращаюсь к первому плакату:

> «А теперь вспомните вчерашний вопрос: где тут evidence?»

Ответ:

> **`go test` — наблюдаемое доказательство конкретного поведения.**

---

На этом месте я бы завершил первую короткую часть урока словами:

> «Сейчас у нас ещё почти нет backend. И это нормально. У нас есть самое важное — первая доменная модель и понимание, почему код выглядит именно так.
>
> Дальше мы начнём делать программу настоящей.
>
> Сначала дадим ей много задач.
>
> Потом заставим их выполняться одновременно.
>
> Потом увидим race.
>
> Потом будем учиться ограничивать concurrency.
>
> А потом весь этот маленький JobFlow превратится в настоящий backend.»
>
> «И каждый раз, когда мы добавляем новую конструкцию Go, я буду задавать один и тот же вопрос: **какую проблему она решает и какую гарантию даёт?**»

После этого можно переходить к практической реализации `JobFlow` и постепенно отдавать ученику управление клавиатурой.

## Продолжение после первого блока

Первый блок мы уже провели: показали `GO LANGUAGE`, быстро прошли модель языка, создали `go.mod`, начали `Job`, `JobStatus`, methods, repository, errors и первый test.

Теперь не надо резко переходить к следующей теме. Нам нужно **закончить JobFlow v0**, проверить, что ученик способен самостоятельно применять модель Go, и обязательно завершить занятие домашним заданием.

Главное правило второй части:

> **После объяснения ученик большую часть времени пишет сам.**

---

# 1. Возвращаемся к Job

Я говорю:

> «Сейчас у нас есть каркас. Но каркас ещё не является моделью. Давайте доведём `Job` до состояния, когда его можно использовать как нормальную доменную сущность.»

Показываю текущий тип:

```go
type JobStatus string

const (
    StatusPending   JobStatus = "pending"
    StatusRunning   JobStatus = "running"
    StatusCompleted JobStatus = "completed"
    StatusFailed    JobStatus = "failed"
)

type Job struct {
    ID        string
    Type      string
    Payload   []byte
    Status    JobStatus
    CreatedAt time.Time
}
```

Спрашиваю:

> «Какие действия над Job нам нужны прямо сейчас?»

<details><summary>Помощь</summary>Назовите изменения состояния и проверку состояния.</details>

<details><summary>Ответ</summary>Создание, переход в running, завершение, failure и проверка, завершена ли задача.</details>

Пусть ученик сам предлагает API.

Мы приходим примерно к:

```go
func NewJob(id, jobType string, payload []byte) Job
func (j *Job) Start()
func (j *Job) Complete()
func (j *Job) Fail()
func (j Job) IsFinished() bool
```

---

# 2. Состояния Job

Я спрашиваю:

> «Можно ли выполнить `Complete()` для `pending`?»

<details><summary>Помощь</summary>Нарисуйте допустимый lifecycle состояния.</details>

<details><summary>Ответ</summary>Нам нужен явный lifecycle. Например: `pending → running → completed` или `pending → running → failed`.</details>

Показываю:

```text
pending
   │
   ▼
running
  /   \
 ▼     ▼
done  failed
```

И спрашиваю:

> «Что здесь является invariant?»

<details><summary>Ответ</summary>Состояние Job должно находиться в допустимом наборе состояний, а переходы должны соответствовать определённому lifecycle.</details>

Здесь не надо ещё вводить сложную state machine. Это будет позже.

Наша задача — показать, что даже простой `struct` имеет доменную семантику.

---

# 3. Проверяем value / pointer ещё раз

Я предлагаю:

```go
func (j Job) Complete() {
    j.Status = StatusCompleted
}
```

И спрашиваю:

> «Оставляем так или меняем receiver?»

<details><summary>Ответ</summary>Для изменения состояния объекта нужен pointer receiver: `func (j *Job) Complete()`.</details>

Теперь ученик сам исправляет.

Дальше спрашиваю:

> «Почему `IsFinished` может использовать value receiver?»

<details><summary>Помощь</summary>Метод только читает state или изменяет его?</details>

<details><summary>Ответ</summary>Он только читает состояние, поэтому value receiver подходит.</details>

Это очень полезный момент: ученик начинает выбирать receiver **по семантике операции**, а не по шаблону.

---

# 4. Validation как часть контракта

Теперь:

```go
var ErrInvalidJob = errors.New("invalid job")
```

Спрашиваю:

> «Что делает Job невалидным?»

Пусть ученик назовёт правила.

Минимально:

```text
ID != ""
Type != ""
Status == StatusPending
```

Затем:

```go
func ValidateJob(job Job) error
```

Ученик пишет реализацию.

Если он начинает делать отдельную ошибку на каждое поле, спрашиваю:

> «Нам сейчас действительно нужна отдельная error identity для каждого нарушения?»

<details><summary>Ответ</summary>Нет. Для текущего контракта достаточно общей `ErrInvalidJob` с добавлением контекста о конкретном поле.</details>

Можно прийти к:

```go
if job.ID == "" {
    return fmt.Errorf("job id: %w", ErrInvalidJob)
}
```

и аналогично для `Type`.

---

# 5. Constructor

Теперь спрашиваю:

> «Если `Job` нельзя создавать без обязательных полей, где лучше сконцентрировать создание?»

<details><summary>Ответ</summary>В constructor `NewJob`, чтобы централизовать initial state и правила создания.</details>

Ученик пишет:

```go
func NewJob(
    id string,
    jobType string,
    payload []byte,
) Job {
    return Job{
        ID:        id,
        Type:      jobType,
        Payload:   payload,
        Status:    StatusPending,
        CreatedAt: time.Now(),
    }
}
```

Здесь я задаю вопрос:

> «Можно ли `NewJob` возвращать ошибку?»

<details><summary>Помощь</summary>А если входные аргументы нарушают обязательный контракт?</details>

<details><summary>Ответ</summary>Да, если создание должно гарантировать валидность Job, естественный контракт — `(Job, error)`. Для учебной первой версии можно оставить создание простым и отдельно использовать validation, но решение должно быть осознанным.</details>

Я здесь не заставляю ученика усложнять API.

Главная мысль:

> **Constructor — механизм централизации initial state, а не обязательный шаблон для каждого типа.**

---

# 6. Repository v0

Теперь создаём:

```go
type MemoryJobRepository struct {
    jobs map[string]Job
}
```

Спрашиваю:

> «Что будет, если забыть инициализировать map?»

<details><summary>Ответ</summary>Чтение из nil map допустимо, запись вызовет panic.</details>

Создаём constructor:

```go
func NewMemoryJobRepository() *MemoryJobRepository {
    return &MemoryJobRepository{
        jobs: make(map[string]Job),
    }
}
```

Теперь ученик реализует сам:

```go
func (r *MemoryJobRepository) Add(job Job) error
```

и:

```go
func (r *MemoryJobRepository) Get(id string) (Job, error)
```

---

# 7. Ошибка not found

Спрашиваю:

> «Что должен получить вызывающий код, если Job не существует?»

<details><summary>Ответ</summary>Нужно вернуть zero Job и `ErrJobNotFound` либо ошибку, обёрнутую вокруг неё.</details>

После реализации спрашиваю:

> «Почему мы не возвращаем просто `nil`?»

<details><summary>Помощь</summary>Возвращаемое значение у нас `Job`, а не `*Job`. Но даже если был бы pointer, что ещё должен узнать вызывающий код?</details>

<details><summary>Ответ</summary>Нужно явно различить успешный результат и failure. `error` является частью контракта.</details>

---

# 8. Repository interface

Теперь я говорю:

> «У нас есть память. Через несколько уроков будет PostgreSQL. Что должен видеть `JobService`?»

Пусть ученик сам напишет:

```go
type JobRepository interface {
    Add(Job) error
    Get(string) (Job, error)
}
```

Спрашиваю:

> «Почему interface находится здесь, а не рядом с MemoryRepository?»

<details><summary>Помощь</summary>Кому реально нужен этот контракт?</details>

<details><summary>Ответ</summary>Потребителю — `JobService`. Поэтому interface описывает необходимое поведение со стороны consumer.</details>

Это мост к `DESIGN PRINCIPLES`, но пока без лекции SOLID.

---

# 9. Dependency injection

Создаём:

```go
type JobService struct {
    repo JobRepository
}
```

и:

```go
func NewJobService(repo JobRepository) *JobService {
    return &JobService{
        repo: repo,
    }
}
```

Я спрашиваю:

> «Почему мы не делаем внутри `JobService` так?»

```go
func NewJobService() *JobService {
    repo := NewMemoryJobRepository()
    ...
}
```

<details><summary>Ответ</summary>Тогда `JobService` сам выбирает infrastructure detail и становится связан с конкретной реализацией. Передача зависимости через constructor делает dependency явной и заменяемой.</details>

И кратко:

> «Вот здесь уже начинает работать Dependency Inversion, но мы пока не изучаем его теоретически. Мы только используем его, потому что он решает реальную проблему.»

---

# 10. Первый самостоятельный кусок

Теперь отдаю клавиатуру ученику.

Я говорю:

> «Дальше пять минут я не пишу код. Вы сами реализуете этот кусок.»

Задача:

```text
1. NewJob
2. Start
3. Complete
4. Fail
5. IsFinished
6. ValidateJob
7. MemoryJobRepository
8. JobRepository
9. JobService
```

Я могу только задавать вопросы.

Если ученик застрял на `Start()`:

<details><summary>Помощь</summary>Какое состояние было до вызова и какое должно быть после?</details>

<details><summary>Ответ</summary>`pending → running`.</details>

Если застрял на `IsFinished()`:

<details><summary>Помощь</summary>Какие состояния больше не требуют дальнейшего выполнения?</details>

<details><summary>Ответ</summary>`completed` и `failed`.</details>

Если застрял на repository:

<details><summary>Помощь</summary>Как быстро найти Job по ID?</details>

<details><summary>Ответ</summary>`map[string]Job`.</details>

Если застрял на interface:

<details><summary>Помощь</summary>Что конкретно нужно сервису от storage?</details>

<details><summary>Ответ</summary>Только операции, которые реально вызывает сервис: сейчас `Add` и `Get`.</details>

---

# 11. Намеренно вводим ошибку

После того как рабочий вариант написан, я говорю:

> «Теперь я специально сломаю вашу программу.»

Показываю:

```go
func (j Job) Complete() {
    j.Status = StatusCompleted
}
```

Спрашиваю:

> «Что должен показать этот тест?»

```go
func TestJobComplete(t *testing.T) {
    job := NewJob("job-1", "email", nil)
    job.Complete()

    if job.Status != StatusCompleted {
        t.Fatal("job was not completed")
    }
}
```

<details><summary>Ответ</summary>Тест должен упасть, потому что value receiver изменяет копию Job.</details>

Запускаем:

```bash
go test ./...
```

После failure:

> «Вот так я хочу, чтобы вы учились в этом курсе. Не я рассказал вам про pointer receiver. Мы создали наблюдаемое противоречие между ожидаемым и фактическим поведением, нашли причину и исправили её.»

Потом меняем receiver на `*Job`.

Тест зелёный.

---

# 12. `defer`

Когда основной JobFlow уже работает, добавляем маленький эксперимент.

Я говорю:

> «Есть ещё одна конструкция, которая почти наверняка понадобится нам в следующем уроке.»

Пишем:

```go
func demo() {
    defer fmt.Println("third")
    defer fmt.Println("second")

    fmt.Println("first")
}
```

Спрашиваю:

> «Какой порядок вывода?»

<details><summary>Ответ</summary>Сначала `first`, потом `second`, потом `third`. Несколько defer выполняются в обратном порядке регистрации.</details>

И сразу связываем с практикой:

> «Чаще всего `defer` будет использоваться для cleanup: unlock, close, cancel, release.»

Не разворачиваем тему дальше.

---

# 13. Testing как завершение цикла

Теперь говорю:

> «До сих пор мы в основном проверяли поведение вручную. Давайте сделаем это воспроизводимым.»

Ученик пишет тест:

```go
func TestJobComplete(t *testing.T) {
    job := NewJob("job-1", "email", nil)

    job.Complete()

    if job.Status != StatusCompleted {
        t.Fatalf(
            "expected %q, got %q",
            StatusCompleted,
            job.Status,
        )
    }
}
```

Запускаем:

```bash
go test ./...
```

Теперь добавляем:

```go
func TestValidateJob(t *testing.T)
```

И несколько cases.

---

# 14. Table-driven test

Я спрашиваю:

> «У нас три похожих теста. Что в них одинаковое? Что меняется?»

<details><summary>Ответ</summary>Логика проверки одинаковая, меняются вход и ожидаемый результат.</details>

Отсюда:

```go
tests := []struct {
    name    string
    job     Job
    wantErr bool
}{
    {
        name: "valid",
        job: Job{
            ID:   "job-1",
            Type: "email",
        },
        wantErr: false,
    },
    {
        name: "missing id",
        job: Job{
            Type: "email",
        },
        wantErr: true,
    },
    {
        name: "missing type",
        job: Job{
            ID: "job-1",
        },
        wantErr: true,
    },
}
```

И:

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        err := ValidateJob(tt.job)

        if (err != nil) != tt.wantErr {
            t.Fatalf(
                "error = %v, wantErr = %v",
                err,
                tt.wantErr,
            )
        }
    })
}
```

Спрашиваю:

> «Почему это лучше, чем три почти одинаковые функции?»

<details><summary>Ответ</summary>Cases выражены как данные, а сама логика проверки находится в одном месте. Добавление нового случая не требует копирования тестовой логики.</details>

---

# 15. Финальный mini-review

Я говорю:

> «Теперь закройте код. Я задаю вопросы, а вы отвечаете без подсказок.»

### Что такое `struct`?

<details><summary>Ответ</summary>Тип, объединяющий связанные поля состояния.</details>

### Зачем собственный `JobStatus`?

<details><summary>Ответ</summary>Чтобы модель и API явно выражали допустимый смысл значения и делали неправильные состояния труднее выразить.</details>

### Когда pointer receiver?

<details><summary>Ответ</summary>Когда method должен работать с исходным объектом и изменять его состояние либо когда это необходимо по семантике/стоимости копирования.</details>

### Зачем slice?

<details><summary>Ответ</summary>Для последовательности элементов с длиной и capacity.</details>

### Зачем map?

<details><summary>Ответ</summary>Для доступа по ключу, здесь `JobID → Job`.</details>

### Что означает `value, ok := m[key]`?

<details><summary>Ответ</summary>Получаем значение и признак существования ключа.</details>

### Зачем `error`?

<details><summary>Ответ</summary>Чтобы явно передавать неуспешный результат как часть контракта функции.</details>

### Зачем `%w`?

<details><summary>Ответ</summary>Чтобы добавить контекст, сохранив исходную error identity для `errors.Is` / `errors.As`.</details>

### Зачем interface?

<details><summary>Ответ</summary>Чтобы потребитель зависел от необходимого поведения, а не от конкретной реализации.</details>

### Почему tests?

<details><summary>Ответ</summary>Чтобы получить воспроизводимое evidence конкретного поведения.</details>

---

# 16. Финальная задача без подсказок

Теперь даю ученику самостоятельную задачу.

> «Представим, что нам нужно добавить удаление Job.»

Нужно самому определить:

```go
Delete(string) error
```

Вопросы я задаю только если ученик застрял.

### Как должна вести себя операция?

<details><summary>Помощь</summary>Опишите отдельно существующий и отсутствующий Job.</details>

<details><summary>Ответ</summary>Существующий Job удаляется и возвращается `nil`; отсутствующий Job возвращает `ErrJobNotFound`.</details>

### Нужно ли менять interface?

<details><summary>Ответ</summary>Да, если `JobService` будет использовать Delete. Но изменение interface должно следовать реальной потребности consumer.</details>

### Какие tests добавить?

<details><summary>Ответ</summary>Удаление существующего Job, удаление отсутствующего Job, проверка что после удаления `Get` возвращает not found.</details>

Ученик реализует.

Запускаем:

```bash
go test ./...
```

Если всё зелёное:

> «Хорошо. Теперь вы не просто слушали про Go — вы добавили изменение в рабочую модель и доказали его тестами.»

---

# 17. Финальная сборка JobFlow v0

К концу занятия на руках должна быть такая модель:

```text
Job
├── JobStatus
├── NewJob
├── Start
├── Complete
├── Fail
└── IsFinished

Validation
└── ValidateJob

Repository
├── JobRepository
└── MemoryJobRepository

Service
└── JobService

Tests
├── Job behavior
├── validation
└── repository
```

Минимальный runtime flow:

```text
NewJob
   ↓
pending
   ↓
Start
   ↓
running
   ├── Complete → completed
   └── Fail     → failed
```

---

# 18. Что сознательно оставляем на следующий урок

Я не объясняю сейчас concurrency глубоко.

Но показываю один вопрос:

> «А что будет, если два worker одновременно вызовут `Complete()` или два потока одновременно начнут менять repository?»

Пауза.

> «Вот здесь наш сегодняшний `MemoryJobRepository` уже перестанет быть безопасной моделью.»

И всё.

Плакат `03. GO CONCURRENCY` пока не открываем.

Это создаёт естественный cliffhanger.

---

# 19. Домашнее задание — JobFlow v1

Вот здесь урок обязательно заканчивается **конкретным самостоятельным результатом**.

Я говорю:

> «Домашнее задание — не прочитать ещё двадцать страниц Go. Вы продолжаете тот же JobFlow.»

## Задача

Расширить repository:

```go
type JobRepository interface {
    Add(Job) error
    Get(string) (Job, error)
    Delete(string) error
    List() []Job
}
```

Добавить:

```text
Delete
List
FindByType
```

## Требуемое поведение

`Get`:

```text
existing → Job + nil
missing  → zero Job + ErrJobNotFound
```

`Delete`:

```text
existing → nil
missing  → ErrJobNotFound
```

`List`:

> возвращает все Job.

`FindByType`:

> возвращает все Job указанного типа.

## Domain

Добавить и протестировать:

```text
Start
Complete
Fail
IsFinished
```

Проверить недопустимые состояния.

## Tests

Обязательно написать tests для:

```text
create
get
not found
delete
delete missing
list
find by type
start
complete
fail
validation
```

## Toolchain

Перед завершением:

```bash
go fmt ./...
go test ./...
go vet ./...
```

Ожидаемый результат:

```text
PASS
go test ./...
PASS
go vet ./...
```

---

# 20. PASS / FAIL домашнего задания

## PASS

Ученик:

* самостоятельно расширил interface;
* самостоятельно расширил memory implementation;
* написал новые tests;
* не нарушил существующие tests;
* понимает, зачем добавляется каждый новый метод;
* может объяснить value/pointer receiver;
* умеет объяснить `slice`, `map`, `error`, interface;
* запускает `go test ./...` без подсказки.

## FAIL

Если решение работает только после копирования готового кода и ученик не может объяснить:

> «Почему это должно работать?»

Домашнее задание считается незавершённым.

---

# 21. Последние 2 минуты занятия

Я не заканчиваю урок перечислением изученных keywords.

Я говорю:

> «Сегодня мы сделали первую версию JobFlow на чистом Go.
>
> У нас есть состояние — `struct`.
>
> У нас есть собственные доменные типы.
>
> Есть поведение — methods.
>
> Есть коллекции — slice и map.
>
> Есть контракт — interface.
>
> Есть failure — error.
>
> Есть evidence — tests.
>
> И теперь у нас есть проблема.»

Показываю на JobFlow:

> «Он умеет обработать Job. Но он обрабатывает их последовательно.»

Пауза.

> «А теперь представьте десять тысяч Job.»

Пауза.

> «В следующем занятии мы заставим их выполняться одновременно.»

Показываю следующий плакат:

# 03. GO CONCURRENCY

И задаю последний вопрос:

> «Что может пойти не так, если тысяча goroutine одновременно начнёт менять одно состояние?»

Не отвечаю.

> «С этим начнём следующий бой.»

---

# 22. Артефакты после занятия

К завершению `01.LANGUAGE` должны существовать:

```text
01.LANGUAGE.md
01.LANGUAGE/
    job.go
    job_test.go
    repository.go
    repository_test.go
    service.go
    service_test.go
```

И состояние проекта должно быть воспроизводимым:

```bash
go fmt ./...
go test ./...
go vet ./...
```

### Следующая точка

`02.CONCURRENCY.md`

На вход ему уже есть:

```text
Job
MemoryJobRepository
JobService
tests
```

А главный новый вопрос:

> **Как сделать выполнение конкурентным, ограниченным и корректно завершаемым?**
