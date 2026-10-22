# 500 Interview Questions

**Цель** - на любой вопрос за 30–60 секунд без подсказки дать сильный ответ - уметь точно объяснить смысл, выбрать механизм и защитить trade-off.

- [500 Interview Questions](#500-interview-questions)
  - [GO LANGUAGE — FUNDAMENTALS](#go-language--fundamentals)
  - [TYPES, API и ERROR SEMANTICS](#types-api-и-error-semantics)
  - [INTERFACES / DESIGN BASICS](#interfaces--design-basics)
  - [TESTING / ENGINEERING EVIDENCE](#testing--engineering-evidence)
  - [CONCURRENCY — CORE](#concurrency--core)
  - [RESOURCE CONTROL / LIFECYCLE](#resource-control--lifecycle)
  - [GO CONCURRENCY — SYNCHRONIZATION / RACE](#go-concurrency--synchronization--race)
  - [GO CONCURRENCY — ADVANCED](#go-concurrency--advanced)
  - [WORKER POOL / RUNTIME](#worker-pool--runtime)
  - [GO RUNTIME / GMP / SCHEDULER](#go-runtime--gmp--scheduler)
  - [SQL — SCHEMA / QUERY FUNDAMENTALS](#sql--schema--query-fundamentals)
  - [POSTGRESQL — QUERY PLANNING / INDEXES](#postgresql--query-planning--indexes)
  - [TRANSACTIONS / CONCURRENCY](#transactions--concurrency)
  - [ACID / DATA INTEGRITY](#acid--data-integrity)
  - [Pgxpool / Go + PostgreSQL](#pgxpool--go--postgresql)
  - [MIGRATIONS / SCHEMA EVOLUTION](#migrations--schema-evolution)
  - [DATABASE DESIGN / JOBFLOW](#database-design--jobflow)
  - [DATA ACCESS / SQL SAFETY](#data-access--sql-safety)
  - [KAFKA / DELIVERY SEMANTICS](#kafka--delivery-semantics)
  - [IDEMPOTENCY](#idempotency)
  - [OUTBOX](#outbox)
  - [RETRY / BACKOFF / FAILURE POLICY](#retry--backoff--failure-policy)
  - [DISTRIBUTED FAILURE ANALYSIS](#distributed-failure-analysis)
  - [HTTP / API](#http--api)
  - [CONTEXT / HTTP LIFECYCLE](#context--http-lifecycle)
  - [APPLICATION ARCHITECTURE](#application-architecture)
  - [DESIGN / SOLID](#design--solid)
  - [REDIS / CACHE](#redis--cache)
  - [PRODUCTION / OBSERVABILITY](#production--observability)
  - [SYSTEM DESIGN / CROSS-SYSTEM](#system-design--cross-system)


---

## GO LANGUAGE — FUNDAMENTALS

**[К оглавлению](#500-interview-questions)**

1. Что такое `package` в Go?

<details style="display: inline;"><summary>Помощь</summary>Подумайте о namespace и единице компиляции.</details>  
<details style="display: inline;"><summary>Ответ</summary>`package` объединяет Go-файлы в логическую единицу компиляции с общим namespace и правилами visibility.</details>

2. Что такое Go module?

<details style="display: inline;"><summary>Помощь</summary>Вспомните `go.mod`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Module — единица управления версионируемым Go-кодом и его зависимостями; его идентичность и зависимости описываются в `go.mod`.</details>

3. Что делает `go mod init`?

<details style="display: inline;"><summary>Помощь</summary>Что появляется после команды?</details>  
<details style="display: inline;"><summary>Ответ</summary>Создаёт новый `go.mod` с указанным module path и базовой информацией о модуле.</details>

4. Чем `go run` отличается от `go build`?

<details style="display: inline;"><summary>Помощь</summary>Подумайте о конечном binary.</details>  
<details style="display: inline;"><summary>Ответ</summary>`go run` компилирует и запускает программу, а `go build` компилирует package/program и создаёт или проверяет возможность создания исполняемого результата.</details>

5. Что такое static typing в Go?

<details style="display: inline;"><summary>Помощь</summary>Когда compiler проверяет тип?</details>  
<details style="display: inline;"><summary>Ответ</summary>Типы являются частью программы и проверяются compiler-ом до запуска, поэтому многие ошибочные операции обнаруживаются на этапе компиляции.</details>

6. Что такое type inference?

<details style="display: inline;"><summary>Помощь</summary>Посмотрите на `count := 10`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Compiler самостоятельно выводит тип выражения из его значения или контекста, при этом типизация остаётся статической.</details>

7. Отключает ли `:=` static typing?

<details style="display: inline;"><summary>Помощь</summary>Тип известен compiler-у?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. `:=` только сокращает запись; переменная всё равно имеет статический тип, который compiler проверяет.</details>

8. Что такое zero value?

<details style="display: inline;"><summary>Помощь</summary>Что происходит при `var x T` без initializer?</details>  
<details style="display: inline;"><summary>Ответ</summary>Zero value — определённое начальное значение типа, которое получает переменная без явного initial value.</details>

9. Каковы zero values `string`, `int` и `bool`?

<details style="display: inline;"><summary>Помощь</summary>Базовые типы.</details>  
<details style="display: inline;"><summary>Ответ</summary>`string` → `""`, `int` → `0`, `bool` → `false`.</details>

10. Что такое `nil` в Go?

<details style="display: inline;"><summary>Помощь</summary>Какие типы могут иметь nil state?</details>  
<details style="display: inline;"><summary>Ответ</summary>`nil` — специальное нулевое состояние для pointer, slice, map, channel, function и некоторых interface values.</details>

11. Можно ли читать из `nil` map?

<details style="display: inline;"><summary>Помощь</summary>Разделите чтение и запись.</details>  
<details style="display: inline;"><summary>Ответ</summary>Да, чтение допустимо и возвращает zero value и, при двухзначном lookup, `ok == false`; запись в `nil` map вызывает panic.</details>

12. Как создать map, готовую к записи?

<details style="display: inline;"><summary>Помощь</summary>Вспомните `make`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать `make`, например `make(map[string]Job)`.</details>

13. Что такое `struct`?

<details style="display: inline;"><summary>Помощь</summary>Что им удобно моделировать?</details>  
<details style="display: inline;"><summary>Ответ</summary>`struct` — пользовательский составной тип, объединяющий связанные поля состояния одной сущности.</details>

14. Зачем использовать named fields при создании struct?

<details style="display: inline;"><summary>Помощь</summary>Что происходит при изменении порядка полей?</details>  
<details style="display: inline;"><summary>Ответ</summary>Named fields делают создание явным, уменьшают вероятность ошибки и не зависят от порядка полей в структуре.</details>

15. Что такое method в Go?

<details style="display: inline;"><summary>Помощь</summary>Посмотрите на receiver.</details>  
<details style="display: inline;"><summary>Ответ</summary>Method — функция с receiver, то есть поведение, привязанное к определённому типу.</details>

16. Что такое receiver?

<details style="display: inline;"><summary>Помощь</summary>Где в сигнатуре method находится объект?</details>  
<details style="display: inline;"><summary>Ответ</summary>Receiver — параметр method, через который method получает value или pointer на экземпляр типа.</details>

17. Чем value receiver отличается от pointer receiver?

<details style="display: inline;"><summary>Помощь</summary>Копия или исходное состояние?</details>  
<details style="display: inline;"><summary>Ответ</summary>Value receiver работает со значением-копией, а pointer receiver получает pointer на исходный объект и может изменять его состояние.</details>

18. Когда для method естественен pointer receiver?

<details style="display: inline;"><summary>Помощь</summary>Нужно ли изменять объект?</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда method должен изменять состояние объекта или когда pointer semantics лучше соответствуют размеру и модели использования типа.</details>

19. Означает ли pointer в Go ручное управление памятью?

<details style="display: inline;"><summary>Помощь</summary>Кто управляет allocation/lifetime?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Pointer определяет способ доступа к памяти, а allocation и lifetime управляются runtime и garbage collector.</details>

20. Что произойдёт здесь и почему?

```go
func (j Job) Complete() {
	j.Status = StatusCompleted
}
```

<details style="display: inline;"><summary>Помощь</summary>Что за receiver?</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменяется копия `Job`, поэтому исходный объект снаружи method не изменится.</details>

21. Что изменится после замены receiver на `*Job`?

<details style="display: inline;"><summary>Помощь</summary>Где теперь находится состояние?</details>  
<details style="display: inline;"><summary>Ответ</summary>Method получает pointer на исходный `Job`, поэтому изменение `Status` изменяет объект, с которым работает вызывающий код.</details>

22. Что такое slice?

<details style="display: inline;"><summary>Помощь</summary>Это самостоятельный массив или представление над ним?</details>  
<details style="display: inline;"><summary>Ответ</summary>Slice — дескриптор последовательности элементов поверх backing array; он содержит ссылку на storage, length и capacity.</details>

23. Что означает `len(slice)`?

<details style="display: inline;"><summary>Помощь</summary>Сколько элементов сейчас доступно?</details>  
<details style="display: inline;"><summary>Ответ</summary>Текущее количество элементов, доступных через slice.</details>

24. Что означает `cap(slice)`?

<details style="display: inline;"><summary>Помощь</summary>Сколько места доступно до возможного перевыделения?</details>  
<details style="display: inline;"><summary>Ответ</summary>Capacity — сколько элементов может быть доступно в backing array от текущего начала slice до необходимости нового backing array.</details>

25. Почему пишут `jobs = append(jobs, job)`?

<details style="display: inline;"><summary>Помощь</summary>Может ли `append` изменить backing array?</details>  
<details style="display: inline;"><summary>Ответ</summary>Да. `append` может вернуть slice, связанный с новым backing array, поэтому результат нужно использовать.</details>

26. Может ли `append` изменить исходный backing array?

<details style="display: inline;"><summary>Помощь</summary>Что если capacity ещё есть?</details>  
<details style="display: inline;"><summary>Ответ</summary>Да. Если capacity достаточна, новые элементы могут попасть в существующий backing array; если недостаточна, создаётся новый.</details>

27. Что такое map?

<details style="display: inline;"><summary>Помощь</summary>Какой access pattern он выражает?</details>  
<details style="display: inline;"><summary>Ответ</summary>Map выражает отображение `key → value` и предназначен для доступа по ключу.</details>

28. Что означает `value, ok := m[key]`?

<details style="display: inline;"><summary>Помощь</summary>Зачем второй результат?</details>  
<details style="display: inline;"><summary>Ответ</summary>`value` содержит найденное значение или zero value, а `ok` показывает, существует ли ключ.</details>

29. Почему нельзя определять наличие ключа только по `value`?

<details style="display: inline;"><summary>Помощь</summary>Что если реальное значение равно zero value?</details>  
<details style="display: inline;"><summary>Ответ</summary>Потому что отсутствующий ключ также возвращает zero value; нужен `ok` для различения этих случаев.</details>

30. Что такое function в Go?

<details style="display: inline;"><summary>Помощь</summary>Назовите её основные части.</details>  
<details style="display: inline;"><summary>Ответ</summary>Function задаёт именованную или анонимную операцию с параметрами и возвращаемыми значениями.</details>

31. Зачем Go позволяет возвращать несколько значений?

<details style="display: inline;"><summary>Помощь</summary>Вспомните `(value, error)`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Это позволяет явно возвращать основной результат вместе с дополнительным состоянием, например `error`, без скрытых исключений.</details>

32. Что такое `defer`?

<details style="display: inline;"><summary>Помощь</summary>Когда выполняется отложенная операция?</details>  
<details style="display: inline;"><summary>Ответ</summary>`defer` регистрирует вызов, который выполнится при выходе из текущей функции.</details>

33. В каком порядке выполняются несколько `defer`?

<details style="display: inline;"><summary>Помощь</summary>Представьте stack.</details>  
<details style="display: inline;"><summary>Ответ</summary>В обратном порядке регистрации, по принципу LIFO.</details>

34. Для чего обычно используют `defer` в backend-коде?

<details style="display: inline;"><summary>Помощь</summary>Подумайте о resource cleanup.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для `Unlock`, `Close`, `Rollback`, `cancel` и других действий, которые должны выполняться при выходе из функции.</details>

35. Что такое exported identifier в Go?

<details style="display: inline;"><summary>Помощь</summary>Регистр первой буквы.</details>  
<details style="display: inline;"><summary>Ответ</summary>Identifier с заглавной буквы доступен из другого package, если соответствующая область package/import visibility это допускает.</details>

36. Почему `fmt.Println` доступен, а `fmt.println` нет?

<details style="display: inline;"><summary>Помощь</summary>Это два разных identifier.</details>  
<details style="display: inline;"><summary>Ответ</summary>`Println` — exported identifier из package `fmt`, а `println` — другое имя и не является экспортируемым API `fmt`.</details>

37. Что такое zero-value design?

<details style="display: inline;"><summary>Помощь</summary>Подумайте, можно ли использовать тип без constructor.</details>  
<details style="display: inline;"><summary>Ответ</summary>Это проектирование типа так, чтобы его zero value имел корректную или полезную семантику и был безопасен для базового использования.</details>

38. Почему `nil` slice часто удобнее специального пустого состояния?

<details style="display: inline;"><summary>Помощь</summary>Как ведут себя `len` и `append`?</details>  
<details style="display: inline;"><summary>Ответ</summary>Nil slice имеет `len == 0`, допускает `append` и часто не требует отдельной инициализации для корректной работы.</details>

39. Можно ли делать `append` в nil slice?

<details style="display: inline;"><summary>Помощь</summary>Вспомните zero value slice.</details>  
<details style="display: inline;"><summary>Ответ</summary>Да. `append` корректно создаёт нужное backing storage.</details>

40. Чем `array` отличается от `slice`?

<details style="display: inline;"><summary>Помощь</summary>Размер фиксирован?</details>  
<details style="display: inline;"><summary>Ответ</summary>Array имеет фиксированный размер, являющийся частью его типа; slice — динамическое представление над backing array.</details>

---

## TYPES, API и ERROR SEMANTICS

**[К оглавлению](#500-interview-questions)**

41. Зачем `type JobStatus string`, если можно оставить `string`?

<details style="display: inline;"><summary>Помощь</summary>Сравните выразительность API.</details>  
<details style="display: inline;"><summary>Ответ</summary>Собственный тип выражает domain semantics и уменьшает вероятность передачи произвольной строки туда, где нужен конкретный статус.</details>

42. Что означает «тип делает неправильное состояние труднее выразить»?

<details style="display: inline;"><summary>Помощь</summary>Сравните `string` и `JobStatus`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чем точнее тип отражает domain constraints, тем меньше некорректных значений можно случайно передать через API.</details>

43. Почему `Status string` хуже выражает contract, чем `Status JobStatus`?

<details style="display: inline;"><summary>Помощь</summary>Что понимает consumer из сигнатуры?</details>  
<details style="display: inline;"><summary>Ответ</summary>`JobStatus` явно сообщает семантику параметра, тогда как `string` позволяет передать любой текст без дополнительного semantic signal.</details>

44. Что такое domain invariant?

<details style="display: inline;"><summary>Помощь</summary>Какое состояние нельзя нарушать?</details>  
<details style="display: inline;"><summary>Ответ</summary>Свойство состояния, которое должно оставаться истинным во всех допустимых состояниях и после допустимых операций.</details>

45. Чем validation отличается от invariant enforcement?

<details style="display: inline;"><summary>Помощь</summary>Где находится authoritative state?</details>  
<details style="display: inline;"><summary>Ответ</summary>Validation проверяет input на границе системы, а invariant enforcement гарантирует допустимость самого состояния там, где оно владеется и изменяется.</details>

46. Почему database constraint полезен даже при наличии Go validation?

<details style="display: inline;"><summary>Помощь</summary>Может ли в БД писать другой consumer?</details>  
<details style="display: inline;"><summary>Ответ</summary>Database constraint защищает persisted state независимо от того, какой application path выполняет запись.</details>

47. Что такое `error` как contract?

<details style="display: inline;"><summary>Помощь</summary>Что должен знать consumer функции?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer должен знать, какие failures возможны и как их интерпретировать через return value, error identity или type.

48. Почему error message и error identity — разные вещи?

<details style="display: inline;"><summary>Помощь</summary>Человеческое описание и программное решение.</details>  
<details style="display: inline;"><summary>Ответ</summary>Message предназначено для контекста и диагностики, а identity/type позволяет программе стабильно классифицировать ошибку.</details>

49. Что такое sentinel error?

<details style="display: inline;"><summary>Помощь</summary>Например `ErrJobNotFound`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Предопределённое error value, которое используется как стабильная identity определённого класса ошибки.</details>

50. Когда лучше typed error, а когда sentinel error?

<details style="display: inline;"><summary>Помощь</summary>Нужна только классификация или ещё структурированные данные?</details>  
<details style="display: inline;"><summary>Ответ</summary>Sentinel подходит для простой классификации, typed error — когда consumer должен получить структурированные сведения вместе с типом ошибки.

---

## INTERFACES / DESIGN BASICS

**[К оглавлению](#500-interview-questions)**

51. Что такое implicit interface implementation?

<details style="display: inline;"><summary>Помощь</summary>Где declaration `implements`?</details>  
<details style="display: inline;"><summary>Ответ</summary>Отдельного `implements` нет: type автоматически удовлетворяет interface, если его method set соответствует требованию.</details>

52. Почему consumer-side interface уменьшает coupling?

<details style="display: inline;"><summary>Помощь</summary>Сколько методов consumer реально видит?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer зависит от минимального поведения, а не от полного API реализации, поэтому уменьшается dependency surface.</details>

53. Что значит «interface принадлежит consumer»?

<details style="display: inline;"><summary>Помощь</summary>Кто определяет необходимые methods?</details>  
<details style="display: inline;"><summary>Ответ</summary>Потребитель определяет минимальный contract, который ему нужен для выполнения своей задачи.</details>

54. Почему не следует создавать interface на один в один с concrete type?

<details style="display: inline;"><summary>Помощь</summary>Что мы выиграли?</details>  
<details style="display: inline;"><summary>Ответ</summary>Если interface просто копирует весь API implementation и не создаёт полезной boundary, он добавляет abstraction cost без существенного снижения coupling.

55. Что такое dependency surface?

<details style="display: inline;"><summary>Помощь</summary>Сколько behavior должен знать consumer?</details>  
<details style="display: inline;"><summary>Ответ</summary>Набор external capabilities и контрактов, от которых зависит component.

56. Почему большой interface увеличивает blast radius?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт при изменении лишнего method?</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменения ненужных для конкретного consumer частей interface могут затронуть больше implementations и consumers.

57. Что такое constructor injection?

<details style="display: inline;"><summary>Помощь</summary>Сигнатура constructor.</details>  
<details style="display: inline;"><summary>Ответ</summary>Concrete dependency передаётся component при создании через constructor.

58. Как constructor injection помогает тестированию?

<details style="display: inline;"><summary>Помощь</summary>Можно ли подставить fake implementation?</details>  
<details style="display: inline;"><summary>Ответ</summary>Да. Test может передать controlled implementation контракта вместо реального infrastructure adapter.

59. Что такое composition root?

<details style="display: inline;"><summary>Помощь</summary>Где собираются concrete dependencies?</details>  
<details style="display: inline;"><summary>Ответ</summary>Место, где создаются concrete components и соединяются их зависимости в runtime graph приложения.

60. Почему service не должен сам открывать PostgreSQL?

<details style="display: inline;"><summary>Помощь</summary>Policy или infrastructure?</details>  
<details style="display: inline;"><summary>Ответ</summary>Это создаёт прямую зависимость application policy от infrastructure detail и усложняет замену, тестирование и управление lifecycle.

---

## TESTING / ENGINEERING EVIDENCE

**[К оглавлению](#500-interview-questions)**

61. Что такое observable behavior?

<details style="display: inline;"><summary>Помощь</summary>Что может наблюдать consumer?</details>  
<details style="display: inline;"><summary>Ответ</summary>Результат или эффект, который виден через публичный contract, состояние, error или внешний system behavior.

62. Почему тест лучше писать на behavior, а не на implementation detail?

<details style="display: inline;"><summary>Помощь</summary>Что должно пережить refactoring?</details>  
<details style="display: inline;"><summary>Ответ</summary>Tests должны продолжать защищать contract при изменении внутренней реализации.

63. Что такое regression?

<details style="display: inline;"><summary>Помощь</summary>Раньше работало, теперь нет.</details>  
<details style="display: inline;"><summary>Ответ</summary>Поломка ранее поддерживаемого поведения после изменения системы.

64. Что такое acceptance criterion?

<details style="display: inline;"><summary>Помощь</summary>Когда задача считается выполненной?</details>  
<details style="display: inline;"><summary>Ответ</summary>Проверяемое условие, по которому определяется PASS/FAIL результата.

65. Что значит reproducible evidence?

<details style="display: inline;"><summary>Помощь</summary>Другой человек должен получить тот же сигнал.</details>  
<details style="display: inline;"><summary>Ответ</summary>Проверяемый результат, который можно повторить при тех же условиях и получить сопоставимый вывод.

66. Может ли зелёный unit test скрывать production failure?

<details style="display: inline;"><summary>Помощь</summary>Что unit test обычно не моделирует?</details>  
<details style="display: inline;"><summary>Ответ</summary>Да. Он может не проверять concurrency, network failure, resource exhaustion, real database behavior, delivery semantics и другие system-level conditions.

67. Что такое integration test?

<details style="display: inline;"><summary>Помощь</summary>Boundary с реальной dependency.</details>  
<details style="display: inline;"><summary>Ответ</summary>Test, проверяющий взаимодействие нескольких компонентов или реальной external dependency, например PostgreSQL.

68. Что важно проверять после architecture refactoring?

<details style="display: inline;"><summary>Помощь</summary>Invariant или структура файлов?</details>  
<details style="display: inline;"><summary>Ответ</summary>Observable behavior и invariants должны сохраниться, если изменение не заявляет новый behavior.

69. Что означает `PASS` для engineering task?

<details style="display: inline;"><summary>Помощь</summary>Не только «код компилируется».</details>  
<details style="display: inline;"><summary>Ответ</summary>Требуемый behavior выполнен, invariants сохранены, failure paths приемлемы, а результат подтверждён воспроизводимыми tests или другими evidence.

70. Зачем `go vet` перед сдачей?

<details style="display: inline;"><summary>Помощь</summary>Static analysis.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы обнаружить набор потенциально ошибочных или подозрительных конструкций, которые не обязательно ловятся compiler-ом.

---

## CONCURRENCY — CORE

**[К оглавлению](#500-interview-questions)**

71. Что такое concurrency?

<details style="display: inline;"><summary>Помощь</summary>Как организована работа нескольких задач?</details>  
<details style="display: inline;"><summary>Ответ</summary>Concurrency — структура выполнения, при которой независимые задачи могут прогрессировать независимо и координироваться.

72. Что такое parallelism?

<details style="display: inline;"><summary>Помощь</summary>Фактически одновременно?</details>  
<details style="display: inline;"><summary>Ответ</summary>Parallelism — физическое одновременное выполнение нескольких частей работы на разных execution resources.

73. Почему concurrency и parallelism не одно и то же?

<details style="display: inline;"><summary>Помощь</summary>Много задач может быть при одном CPU.</details>  
<details style="display: inline;"><summary>Ответ</summary>Можно иметь много concurrent задач без одновременного физического выполнения всех них.

74. Что делает `go f()`?

<details style="display: inline;"><summary>Помощь</summary>Что появляется?</details>  
<details style="display: inline;"><summary>Ответ</summary>Запускает вызов функции как goroutine, которая может выполняться concurrently с вызывающей goroutine.

75. Почему `main` может завершиться раньше goroutine?

<details style="display: inline;"><summary>Помощь</summary>Что завершает process?</details>  
<details style="display: inline;"><summary>Ответ</summary>Завершение `main` завершается процессом; другие goroutines не получают гарантии завершить работу после этого.

76. Как дождаться группы goroutines?

<details style="display: inline;"><summary>Помощь</summary>Вспомните `sync`.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать `sync.WaitGroup` или другой явно определённый lifecycle coordination mechanism.

77. Почему `WaitGroup` не является queue?

<details style="display: inline;"><summary>Помощь</summary>Что именно он хранит?</details>  
<details style="display: inline;"><summary>Ответ</summary>Он учитывает количество незавершённых operations, но не передаёт сами данные работы.

78. Зачем channel в worker pool?

<details style="display: inline;"><summary>Помощь</summary>Как job попадает к worker?</details>  
<details style="display: inline;"><summary>Ответ</summary>Channel передаёт Job от producer к workers и одновременно координирует их доступ к очереди.

79. Почему unbuffered channel блокирует sender?

<details style="display: inline;"><summary>Помощь</summary>Кто должен принять значение?</details>  
<details style="display: inline;"><summary>Ответ</summary>Send требует соответствующего receive; sender ждёт, пока receiver сможет принять значение.

80. Для чего buffered channel?

<details style="display: inline;"><summary>Помощь</summary>Producer может кратковременно опередить consumer.</details>  
<details style="display: inline;"><summary>Ответ</summary>Buffer позволяет временно накапливать ограниченное количество сообщений и разделять producer/consumer по времени.

---

## RESOURCE CONTROL / LIFECYCLE

**[К оглавлению](#500-interview-questions)**

81. Что такое backpressure?

<details style="display: inline;"><summary>Помощь</summary>Producer быстрее consumer.</details>  
<details style="display: inline;"><summary>Ответ</summary>Backpressure — механизм, заставляющий producer учитывать ограниченную capacity consumer или queue вместо бесконечного накопления work.

82. Почему бесконечная queue опасна?

<details style="display: inline;"><summary>Помощь</summary>Какой ресурс растёт?</details>  
<details style="display: inline;"><summary>Ответ</summary>Backlog может неограниченно потреблять memory и увеличивать latency до resource exhaustion.

83. Что такое bounded concurrency?

<details style="display: inline;"><summary>Помощь</summary>Назовите upper bound.</details>  
<details style="display: inline;"><summary>Ответ</summary>Количество одновременно выполняемых operations ограничено заранее определённым resource bound.

84. Почему bounded worker pool — engineering guarantee?

<details style="display: inline;"><summary>Помощь</summary>Какой ресурс защищается?</details>  
<details style="display: inline;"><summary>Ответ</summary>Он ограничивает concurrent resource consumption и не позволяет input volume напрямую превратитьcя в неограниченное execution pressure.

85. Что такое worker ownership?

<details style="display: inline;"><summary>Помощь</summary>Кто отвечает за goroutine?</details>  
<details style="display: inline;"><summary>Ответ</summary>Явно определённый owner отвечает за запуск, resources, lifecycle и stop condition worker.

86. Что такое stop condition?

<details style="display: inline;"><summary>Помощь</summary>Какие события завершают работу?</details>  
<details style="display: inline;"><summary>Ответ</summary>Явное событие, означающее, что operation должна завершиться, например cancellation, timeout, закрытие input или успешное выполнение.

87. Что такое graceful shutdown для worker pool?

<details style="display: inline;"><summary>Помощь</summary>Stop → wait.</details>  
<details style="display: inline;"><summary>Ответ</summary>Прекратить приём новой работы, сигнализировать cancellation/close, дать допустимой работе завершиться и дождаться выхода workers.

88. Что может произойти, если worker не слушает context?

<details style="display: inline;"><summary>Помощь</summary>Что если внешний owner отменился?</details>  
<details style="display: inline;"><summary>Ответ</summary>Worker может продолжать ненужную работу, удерживать resources и не завершаться вовремя.

89. Что такое resource ownership?

<details style="display: inline;"><summary>Помощь</summary>Кто создаёт и закрывает?</details>  
<details style="display: inline;"><summary>Ответ</summary>Ясно определённый component отвечает за lifecycle ресурса и условия его освобождения.

90. Какие ресурсы нужно ограничивать в backend?

<details style="display: inline;"><summary>Помощь</summary>Не только goroutines.</details>  
<details style="display: inline;"><summary>Ответ</summary>Workers, DB connections, Redis connections, queue capacity, in-flight operations, memory, downstream concurrency и другие bounded resources.

---

## GO CONCURRENCY — SYNCHRONIZATION / RACE

**[К оглавлению](#500-interview-questions)**

91. Что такое data race?

<details style="display: inline;"><summary>Помощь</summary>Shared memory + concurrent access.</details>  
<details style="display: inline;"><summary>Ответ</summary>Data race возникает при конфликтующем concurrent доступе к одной памяти без необходимой synchronization, когда хотя бы одна операция изменяет данные.

92. Почему `counter++` не считается атомарной операцией?

<details style="display: inline;"><summary>Помощь</summary>Сколько шагов внутри?</details>  
<details style="display: inline;"><summary>Ответ</summary>Это read-modify-write sequence, состоящая из нескольких операций, между которыми может вмешаться другая goroutine.

93. Как проверить race в Go?

<details style="display: inline;"><summary>Помощь</summary>Команда test с флагом.</details>  
<details style="display: inline;"><summary>Ответ</summary>`go test -race ./...`.

94. Для чего `sync.Mutex`?

<details style="display: inline;"><summary>Помощь</summary>Что блокируется?</details>  
<details style="display: inline;"><summary>Ответ</summary>Для защиты critical section и shared mutable state от конфликтующего concurrent доступа.

95. Зачем делать `defer mu.Unlock()` после `Lock()`?

<details style="display: inline;"><summary>Помощь</summary>Что если будет early return?</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы unlock гарантированно выполнился при выходе из функции и не оставил mutex заблокированным на error path.

96. Что такое `atomic.Int64`?

<details style="display: inline;"><summary>Помощь</summary>Alternative to mutex for simple state.</details>  
<details style="display: inline;"><summary>Ответ</summary>Тип для выполнения атомарных операций над `int64` shared state без необходимости явно защищать каждую операцию mutex.

97. Когда atomic подходит лучше mutex?

<details style="display: inline;"><summary>Помощь</summary>Один независимый value.</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда нужен простой атомарный counter/flag или другая отдельная операция над одним shared value.

98. Когда mutex естественнее atomic?

<details style="display: inline;"><summary>Помощь</summary>Связанные значения.</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда нужно согласованно изменять несколько связанных values и сохранять composite invariant.

99. Почему обычная Go map опасна при concurrent writes?

<details style="display: inline;"><summary>Помощь</summary>Map не является универсальным concurrent container.</details>  
<details style="display: inline;"><summary>Ответ</summary>Concurrent writes к обычной map требуют внешней synchronization; иначе возможны race и некорректное выполнение.

100. Как проектировать concurrent component до написания кода?

<details style="display: inline;"><summary>Помощь</summary>Вернитесь к Engineering Mindset.</details>  
<details style="display: inline;"><summary>Ответ</summary>Сначала определить ownership, state и invariant, затем boundaries, concurrency/resource bounds, synchronization, cancellation/stop condition, failure policy и способ verification.

## GO CONCURRENCY — ADVANCED

**[К оглавлению](#500-interview-questions)**

101. Что такое `select` в Go?

<details style="display: inline;"><summary>Помощь</summary>Подумайте, если worker одновременно должен ждать Job и cancellation.</details>  
<details style="display: inline;"><summary>Ответ</summary>`select` позволяет одновременно ожидать несколько channel operations и выполнить готовую ветку.</details>

102. Что происходит, если несколько `case` в `select` готовы одновременно?

<details style="display: inline;"><summary>Помощь</summary>Должен ли порядок `case` задавать приоритет?</details>  
<details style="display: inline;"><summary>Ответ</summary>Go выбирает одну из готовых коммуникаций псевдослучайно; порядок `case` не задаёт обычного приоритета.</details>

103. Что делает `default` внутри `select`?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт, если ни один channel operation не готов?</details>  
<details style="display: inline;"><summary>Ответ</summary>`default` выполняется немедленно, поэтому `select` не блокируется.</details>

104. Почему `default` в циклическом `select` может быть опасен?

<details style="display: inline;"><summary>Помощь</summary>Что будет делать loop, если событий нет?</details>  
<details style="display: inline;"><summary>Ответ</summary>Он может превратиться в busy loop и расходовать CPU без полезной работы.

105. Как сделать worker, который ждёт либо Job, либо cancellation?

<details style="display: inline;"><summary>Помощь</summary>Нужны два готовых события.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать `select` с веткой чтения из `jobs` и веткой `<-ctx.Done()`.

106. Что произойдёт при чтении из закрытого channel?

<details style="display: inline;"><summary>Помощь</summary>Разделите поведение при remaining values и после их окончания.</details>  
<details style="display: inline;"><summary>Ответ</summary>Сначала можно прочитать уже отправленные значения; после исчерпания чтение немедленно возвращает zero value, а двухзначный receive даёт `ok == false`.

107. Что произойдёт при записи в закрытый channel?

<details style="display: inline;"><summary>Помощь</summary>Кто должен владеть close?</details>  
<details style="display: inline;"><summary>Ответ</summary>Send в закрытый channel вызывает panic.

108. Почему receiver обычно не должен закрывать jobs channel?

<details style="display: inline;"><summary>Помощь</summary>Кто знает, что новых сообщений больше не будет?</details>  
<details style="display: inline;"><summary>Ответ</summary>Обычно producer владеет sending side и знает, когда новые Job закончились; receiver не должен принимать решение за producer.

109. Что такое channel ownership?

<details style="display: inline;"><summary>Помощь</summary>Кто создаёт, пишет и закрывает?</details>  
<details style="display: inline;"><summary>Ответ</summary>Явная ответственность за lifecycle channel: кто владеет sending side, кто получает данные и кто имеет право закрыть channel.

110. Когда channel лучше mutex?

<details style="display: inline;"><summary>Помощь</summary>Данные нужно защищать или передавать?</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда задача естественно выражается через передачу work/events или coordination между goroutines, а не через защиту общего mutable state.

111. Когда mutex лучше channel?

<details style="display: inline;"><summary>Помощь</summary>Есть ли shared state, который нужно изменить на месте?</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда несколько goroutines должны безопасно обращаться к одному shared state и передача владения через channel не упрощает модель.

112. Правда ли, что channels всегда предпочтительнее mutex?

<details style="display: inline;"><summary>Помощь</summary>Это правило или design choice?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Channel и mutex решают разные задачи; выбор зависит от модели ownership, communication и shared state.

113. Что такое data ownership в concurrent design?

<details style="display: inline;"><summary>Помощь</summary>Кто имеет право изменять state?</details>  
<details style="display: inline;"><summary>Ответ</summary>Явное правило, какой component или goroutine отвечает за изменение конкретного состояния и как остальные получают доступ к нему.

114. Что такое transfer of ownership через channel?

<details style="display: inline;"><summary>Помощь</summary>Producer отправил объект. Кто теперь отвечает за него?</details>  
<details style="display: inline;"><summary>Ответ</summary>Модель, при которой sender передаёт данные consumer, а consumer становится ответственным за дальнейшее использование согласно contract.

115. Почему immutable data упрощает concurrency?

<details style="display: inline;"><summary>Помощь</summary>Можно ли одновременно менять одно и то же?</details>  
<details style="display: inline;"><summary>Ответ</summary>Если данные после публикации не изменяются, concurrent reads не требуют synchronization для самой mutation.

116. Что такое race-free и что такое logically correct?

<details style="display: inline;"><summary>Помощь</summary>Без race автоматически правильно?</details>  
<details style="display: inline;"><summary>Ответ</summary>Race-free означает отсутствие определённого класса некорректного concurrent access; логическая корректность требует ещё сохранения business invariants и правильной synchronization semantics.

117. Может ли программа быть race-free, но всё равно неправильной?

<details style="display: inline;"><summary>Помощь</summary>Представьте два последовательных, но неверных состояния.</details>  
<details style="display: inline;"><summary>Ответ</summary>Да. Например, mutex может убрать race, но неверная последовательность state transitions всё равно нарушит business invariant.

118. Что такое deadlock?

<details style="display: inline;"><summary>Помощь</summary>Два участника ждут друг друга.</details>  
<details style="display: inline;"><summary>Ответ</summary>Состояние, в котором goroutines не могут продолжить выполнение, потому что каждая ждёт ресурс или событие, которое другая не может предоставить.

119. Что такое livelock?

<details style="display: inline;"><summary>Помощь</summary>Система активна, но прогресса нет.</details>  
<details style="display: inline;"><summary>Ответ</summary>Goroutines продолжают выполняться и менять состояние, но полезного progress к цели не происходит.

120. Что такое starvation?

<details style="display: inline;"><summary>Помощь</summary>Операция долго не получает ресурс.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ситуация, когда часть работы систематически не получает необходимый execution/resource allocation и не делает progress.

121. Почему lock ordering важен?

<details style="display: inline;"><summary>Помощь</summary>Две goroutines берут два mutex.</details>  
<details style="display: inline;"><summary>Ответ</summary>Разный порядок захвата нескольких locks может создать circular wait и deadlock; единый lock ordering снижает этот риск.

122. Что такое critical section?

<details style="display: inline;"><summary>Помощь</summary>Какой код нельзя выполнять параллельно?</details>  
<details style="display: inline;"><summary>Ответ</summary>Участок работы, доступ к которому должен быть синхронизирован относительно shared state.

123. Что такое atomicity операции?

<details style="display: inline;"><summary>Помощь</summary>Может ли другой поток увидеть промежуточное состояние?</details>  
<details style="display: inline;"><summary>Ответ</summary>Операция воспринимается как неделимое изменение относительно других concurrent observers.

124. Чем atomic operation отличается от atomic transaction?

<details style="display: inline;"><summary>Помощь</summary>Один value или группа изменений?</details>  
<details style="display: inline;"><summary>Ответ</summary>Atomic operation — неделимая операция над отдельным состоянием; transaction — логическая атомарная граница группы изменений.

125. Что такое memory visibility между goroutines?

<details style="display: inline;"><summary>Помощь</summary>Если одна goroutine записала значение, когда другая обязана его увидеть?</details>  
<details style="display: inline;"><summary>Ответ</summary>Видимость определяется правилами memory model и synchronization operations; concurrent program должен использовать корректные synchronization primitives.

126. Почему нельзя просто надеяться, что scheduler «успеет»?

<details style="display: inline;"><summary>Помощь</summary>Timing ≠ synchronization.</details>  
<details style="display: inline;"><summary>Ответ</summary>Scheduling timing не является гарантией ordering или visibility; correctness должна опираться на явную synchronization semantics.

127. Что делает `sync.Once`?

<details style="display: inline;"><summary>Помощь</summary>Нужно выполнить инициализацию один раз.</details>  
<details style="display: inline;"><summary>Ответ</summary>Гарантирует однократное выполнение заданной function при concurrent вызовах.

128. Когда `sync.Once` полезен в backend?

<details style="display: inline;"><summary>Помощь</summary>One-time initialization.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для thread-safe ленивой инициализации, которая должна произойти ровно один раз в рамках lifetime process.

129. Что такое atomic counter?

<details style="display: inline;"><summary>Помощь</summary>Какая операция типично нужна?</details>  
<details style="display: inline;"><summary>Ответ</summary>Concurrent counter, который можно увеличивать, читать и иногда уменьшать атомарными операциями.

130. Почему atomic не заменяет transaction?

<details style="display: inline;"><summary>Помощь</summary>Сколько state transitions нужно объединить?</details>  
<details style="display: inline;"><summary>Ответ</summary>Atomic защищает отдельную memory operation, а transaction объединяет несколько изменений и гарантирует их совместную фиксацию.

---

## WORKER POOL / RUNTIME

**[К оглавлению](#500-interview-questions)**

131. Почему goroutine-per-request и worker pool решают разные проблемы?

<details style="display: inline;"><summary>Помощь</summary>Request concurrency и bounded processing.</details>  
<details style="display: inline;"><summary>Ответ</summary>Goroutine-per-request создаёт execution unit на входную операцию, а worker pool дополнительно ограничивает число одновременно выполняемых jobs.

132. Как выбрать размер worker pool?

<details style="display: inline;"><summary>Помощь</summary>Какие ресурсы ограничивают throughput?</details>  
<details style="display: inline;"><summary>Ответ</summary>По характеру workload, CPU/I/O ratio, downstream capacity, latency targets и допустимому resource consumption, затем проверять измерениями.

133. Почему worker count больше CPU cores иногда имеет смысл?

<details style="display: inline;"><summary>Помощь</summary>Что делают workers во время I/O?</details>  
<details style="display: inline;"><summary>Ответ</summary>Для I/O-bound workloads часть workers может ждать external I/O, поэтому больше concurrent tasks могут повышать utilization до момента насыщения downstream.

134. Почему worker count можно сделать слишком большим?

<details style="display: inline;"><summary>Помощь</summary>Что станет bottleneck?</details>  
<details style="display: inline;"><summary>Ответ</summary>Растут contention, scheduling overhead, memory, queue pressure и нагрузка на downstream resources.

135. Что такое backpressure в worker pool?

<details style="display: inline;"><summary>Помощь</summary>Producer быстрее fixed workers.</details>  
<details style="display: inline;"><summary>Ответ</summary>Механизм, при котором producer не может бесконтрольно увеличивать backlog; это может выражаться blocking, bounded queue, rejection или throttling.

136. Чем bounded queue лучше бесконечного backlog?

<details style="display: inline;"><summary>Помощь</summary>Какой failure раньше проявляется?</details>  
<details style="display: inline;"><summary>Ответ</summary>Она превращает неконтролируемый рост memory/latency в явный capacity limit и позволяет применить failure policy.

137. Что должен делать worker, если queue closed?

<details style="display: inline;"><summary>Помощь</summary>Есть ли ещё новая работа?</details>  
<details style="display: inline;"><summary>Ответ</summary>Обработать оставшиеся значения и завершиться.

138. Что должен делать worker, если context cancelled?

<details style="display: inline;"><summary>Помощь</summary>Работа ещё нужна?</details>  
<details style="display: inline;"><summary>Ответ</summary>Корректно остановить дальнейшую работу согласно cancellation policy и освободить resources.

139. Что происходит, если producer закрывает queue слишком рано?

<details style="display: inline;"><summary>Помощь</summary>Есть ли ещё send?</details>  
<details style="display: inline;"><summary>Ответ</summary>Последующие sends вызовут panic, а pending producers могут быть логически сломаны.

140. Что происходит, если producer никогда не закрывает queue, а workers ждут `range`?

<details style="display: inline;"><summary>Помощь</summary>Как worker узнаёт, что jobs больше нет?</details>  
<details style="display: inline;"><summary>Ответ</summary>Worker может ждать бесконечно, если нет другого stop condition, создавая потенциальный lifecycle leak.

141. Что такое worker lifecycle contract?

<details style="display: inline;"><summary>Помощь</summary>Start, work, stop.</details>  
<details style="display: inline;"><summary>Ответ</summary>Определённые правила: кто запускает worker, какую работу он получает, как его отменить и при каком событии он должен завершиться.

142. Почему `WaitGroup` не отвечает за cancellation?

<details style="display: inline;"><summary>Помощь</summary>Что он знает?</details>  
<details style="display: inline;"><summary>Ответ</summary>`WaitGroup` только ожидает завершения; он не сообщает worker-ам, что им нужно остановиться.

143. Почему `context.Context` не заменяет `WaitGroup`?

<details style="display: inline;"><summary>Помощь</summary>Signal vs synchronization.</details>  
<details style="display: inline;"><summary>Ответ</summary>Context сообщает cancellation, а `WaitGroup` позволяет дождаться фактического завершения операций.

144. Как обычно связать context и WaitGroup?

<details style="display: inline;"><summary>Помощь</summary>Один отвечает за stop, второй — за wait.</details>  
<details style="display: inline;"><summary>Ответ</summary>Context инициирует cancellation, worker наблюдает его и выходит, а `WaitGroup` дожидается выхода всех workers.

145. Что такое graceful drain?

<details style="display: inline;"><summary>Помощь</summary>Shutdown с допустимым завершением already accepted work.</details>  
<details style="display: inline;"><summary>Ответ</summary>Controlled shutdown, при котором новая работа прекращается, а уже принятая работа завершается согласно заданной policy до выхода process.

146. Почему shutdown policy должна быть явной?

<details style="display: inline;"><summary>Помощь</summary>Что делаем с in-flight jobs?</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы однозначно определить, что можно завершить, что отменить, что retry-ить и какие effects допустимы при остановке.

147. Как тестировать goroutine leak?

<details style="display: inline;"><summary>Помощь</summary>Нужен observable lifecycle.</details>  
<details style="display: inline;"><summary>Ответ</summary>Создать контролируемый запуск/stop path, дождаться completion и при необходимости проверять runtime state или отсутствие зависших goroutines через специализированные test techniques.

148. Можно ли доказать отсутствие всех goroutine leaks простым `WaitGroup`?

<details style="display: inline;"><summary>Помощь</summary>Что если leak goroutine вообще не участвует в WaitGroup?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. WaitGroup доказывает только lifecycle зарегистрированных операций.

149. Что такое cancellation-aware blocking operation?

<details style="display: inline;"><summary>Помощь</summary>Она должна реагировать на context.</details>  
<details style="display: inline;"><summary>Ответ</summary>Operation, которая может прекратить ожидание или работу при cancellation/deadline.

150. Почему timeout должен быть рядом с ownership policy?

<details style="display: inline;"><summary>Помощь</summary>Кто решает, сколько ждать?</details>  
<details style="display: inline;"><summary>Ответ</summary>Timeout является частью lifecycle и resource policy операции, поэтому должен определяться в соответствии с тем, кто владеет её deadline.

---

## GO RUNTIME / GMP / SCHEDULER

**[К оглавлению](#500-interview-questions)**

151. Что означает GMP в Go runtime?

<details style="display: inline;"><summary>Помощь</summary>Три сущности.</details>  
<details style="display: inline;"><summary>Ответ</summary>G — goroutine, M — OS thread, P — runtime scheduler resource для выполнения Go code.

152. Что такое G?

<details style="display: inline;"><summary>Помощь</summary>Work unit.</details>  
<details style="display: inline;"><summary>Ответ</summary>G — структура runtime, представляющая goroutine и её execution state.

153. Что такое M?

<details style="display: inline;"><summary>Помощь</summary>OS execution resource.</details>  
<details style="display: inline;"><summary>Ответ</summary>M представляет OS thread, используемый runtime для фактического выполнения Go code.

154. Что такое P?

<details style="display: inline;"><summary>Помощь</summary>Scheduler resource.</details>  
<details style="display: inline;"><summary>Ответ</summary>P — логический runtime resource, необходимый scheduler-у для выполнения goroutine на M.

155. Всегда ли один P означает один OS thread?

<details style="display: inline;"><summary>Помощь</summary>Сравните P и M.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. P — runtime resource, а M — OS thread; их количество и association не эквивалентны.

156. Что означает «goroutine runnable»?

<details style="display: inline;"><summary>Помощь</summary>Готова, но ещё не обязательно running.</details>  
<details style="display: inline;"><summary>Ответ</summary>Она готова к выполнению scheduler-ом, но может находиться в очереди ожидания CPU execution resource.

157. Что такое run queue?

<details style="display: inline;"><summary>Помощь</summary>Очередь runnable G.</details>  
<details style="display: inline;"><summary>Ответ</summary>Очередь goroutines, готовых к выполнению.

158. Что такое local run queue?

<details style="display: inline;"><summary>Помощь</summary>Связана с P.</details>  
<details style="display: inline;"><summary>Ответ</summary>Очередь runnable goroutines, локально связанная с P и используемая scheduler-ом для эффективного выполнения.

159. Зачем global run queue?

<details style="display: inline;"><summary>Помощь</summary>Shared runnable work.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для хранения runnable work, доступного scheduler-у глобально, и балансировки загрузки между P.

160. Что такое work stealing?

<details style="display: inline;"><summary>Помощь</summary>Idle scheduler resource.</details>  
<details style="display: inline;"><summary>Ответ</summary>Scheduler может брать runnable goroutines из очередей других P, чтобы находить работу для простаивающего execution resource.

161. Почему runtime нужен scheduler?

<details style="display: inline;"><summary>Помощь</summary>Много G, меньше execution resources.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы эффективно распределять множество goroutines между доступными OS threads и CPU execution resources.

162. Что такое preemption?

<details style="display: inline;"><summary>Помощь</summary>Нужно освободить execution resource.</details>  
<details style="display: inline;"><summary>Ответ</summary>Механизм, позволяющий остановить текущую goroutine и передать execution resource другой runnable работе.

163. Зачем preemption нужна backend-системе?

<details style="display: inline;"><summary>Помощь</summary>Не дать одной G монополизировать ресурс.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для поддержания прогресса других runnable goroutines и более предсказуемого распределения execution resources.

164. Что происходит с goroutine во время network wait?

<details style="display: inline;"><summary>Помощь</summary>Нужен ли CPU постоянно?</details>  
<details style="display: inline;"><summary>Ответ</summary>Она может перейти в waiting state, а runtime использует network polling для возобновления её выполнения, когда I/O готово.

165. Что такое netpoller?

<details style="display: inline;"><summary>Помощь</summary>Network readiness.</details>  
<details style="display: inline;"><summary>Ответ</summary>Runtime mechanism для ожидания готовности сетевых I/O и последующего пробуждения связанных goroutines.

166. Что означает wake-up после network I/O?

<details style="display: inline;"><summary>Помощь</summary>G снова может выполняться.</details>  
<details style="display: inline;"><summary>Ответ</summary>G становится runnable и возвращается в scheduler flow.

167. Что происходит в упрощённой модели после `go f()`?

<details style="display: inline;"><summary>Помощь</summary>Не путайте этапы.</details>  
<details style="display: inline;"><summary>Ответ</summary>`go f()` создаёт G → G становится runnable → scheduler планирует G → G выполняется на M с P → после wait может снова стать runnable.

168. Является ли `G → P → M` последовательным pipeline?

<details style="display: inline;"><summary>Помощь</summary>Три сущности или три шага?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. G, P и M — разные runtime entities, которые scheduler связывает для исполнения goroutines.

169. Почему модель GMP полезна при performance debugging?

<details style="display: inline;"><summary>Помощь</summary>Очереди и waiting.</details>  
<details style="display: inline;"><summary>Ответ</summary>Она помогает понимать runnable backlog, scheduling, blocked I/O, CPU saturation и причины задержек concurrent work.

170. Что такое scheduler contention?

<details style="display: inline;"><summary>Помощь</summary>Слишком много runnable work.</details>  
<details style="display: inline;"><summary>Ответ</summary>Состояние, когда большое количество concurrent work создаёт дополнительную scheduling overhead или competition за execution resources.

---

## SQL — SCHEMA / QUERY FUNDAMENTALS

**[К оглавлению](#500-interview-questions)**

171. Что такое relational table?

<details style="display: inline;"><summary>Помощь</summary>Rows + columns.</details>  
<details style="display: inline;"><summary>Ответ</summary>Relation, представленная набором rows с определёнными columns и типами.

172. Что такое primary key?

<details style="display: inline;"><summary>Помощь</summary>Identity.</details>  
<details style="display: inline;"><summary>Ответ</summary>Constraint, задающий уникальную identity row и запрещающий NULL.

173. Что такое foreign key?

<details style="display: inline;"><summary>Помощь</summary>Relationship между tables.</details>  
<details style="display: inline;"><summary>Ответ</summary>Constraint, требующий, чтобы referenced value соответствовала существующей row в другой table согласно foreign-key semantics.

174. Что делает `NOT NULL`?

<details style="display: inline;"><summary>Помощь</summary>Можно ли отсутствовать значению?</details>  
<details style="display: inline;"><summary>Ответ</summary>Запрещает сохранять SQL `NULL` в column.

175. Что делает `UNIQUE`?

<details style="display: inline;"><summary>Помощь</summary>Duplicate identity/value.</details>  
<details style="display: inline;"><summary>Ответ</summary>Не позволяет сохранять конфликтующие duplicate values в constrained columns согласно PostgreSQL NULL semantics.

176. Что делает `CHECK`?

<details style="display: inline;"><summary>Помощь</summary>Какое выражение должен выполнять row?</details>  
<details style="display: inline;"><summary>Ответ</summary>Задаёт predicate, которому должно соответствовать значение row, чтобы запись была допустимой.

177. Почему `CHECK` важен, если есть validation в Go?

<details style="display: inline;"><summary>Помощь</summary>Кто ещё может писать в DB?</details>  
<details style="display: inline;"><summary>Ответ</summary>Database constraint защищает persisted state независимо от конкретного application path.

178. Что такое `NULL` в SQL?

<details style="display: inline;"><summary>Помощь</summary>Не пустая строка.</details>  
<details style="display: inline;"><summary>Ответ</summary>Специальное отсутствие/неизвестность значения, участвующее в SQL three-valued logic.

179. Почему `WHERE value = NULL` не работает ожидаемым образом?

<details style="display: inline;"><summary>Помощь</summary>Как проверить NULL?</details>  
<details style="display: inline;"><summary>Ответ</summary>Сравнение с NULL даёт unknown; необходимо использовать `IS NULL` или `IS NOT NULL`.

180. Что такое normalization в relational database?

<details style="display: inline;"><summary>Помощь</summary>Как уменьшить redundant facts?</details>  
<details style="display: inline;"><summary>Ответ</summary>Набор принципов организации schema, уменьшающих избыточность и аномалии обновления через выделение независимых facts и relationships.

181. Что такое denormalization?

<details style="display: inline;"><summary>Помощь</summary>Когда мы намеренно дублируем data.</details>  
<details style="display: inline;"><summary>Ответ</summary>Намеренное добавление redundancy для конкретного performance/access pattern ценой усложнения consistency.

182. Что делает `INSERT`?

<details style="display: inline;"><summary>Помощь</summary>Создание data.</details>  
<details style="display: inline;"><summary>Ответ</summary>Добавляет rows в table.

183. Что делает `UPDATE`?

<details style="display: inline;"><summary>Помощь</summary>Изменение existing rows.</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменяет columns существующих rows, соответствующих условию.

184. Что делает `DELETE`?

<details style="display: inline;"><summary>Помощь</summary>Удаление rows.</details>  
<details style="display: inline;"><summary>Ответ</summary>Удаляет rows, соответствующие условию.

185. Что делает `SELECT`?

<details style="display: inline;"><summary>Помощь</summary>Read.</details>  
<details style="display: inline;"><summary>Ответ</summary>Извлекает и вычисляет данные по заданному query.

186. Почему `UPDATE` без `WHERE` опасен?

<details style="display: inline;"><summary>Помощь</summary>Какой набор rows станет target?</details>  
<details style="display: inline;"><summary>Ответ</summary>Все rows table.

187. Что такое parameterized query?

<details style="display: inline;"><summary>Помощь</summary>`$1`, `$2`.</details>  
<details style="display: inline;"><summary>Ответ</summary>SQL statement, в котором values передаются отдельно от query text через parameters.

188. Почему parameterized queries защищают от SQL injection?

<details style="display: inline;"><summary>Помощь</summary>SQL и data обрабатываются отдельно.</details>  
<details style="display: inline;"><summary>Ответ</summary>Значения не интерпретируются как часть SQL syntax, а передаются как parameters driver/database.

189. Что такое JOIN?

<details style="display: inline;"><summary>Помощь</summary>Данные двух relations.</details>  
<details style="display: inline;"><summary>Ответ</summary>Операция объединения rows разных relations по заданному condition.

190. Что такое `INNER JOIN`?

<details style="display: inline;"><summary>Помощь</summary>Matching rows only.</details>  
<details style="display: inline;"><summary>Ответ</summary>Возвращает только rows, для которых найдено соответствие по условию JOIN.

191. Что такое `LEFT JOIN`?

<details style="display: inline;"><summary>Помощь</summary>Сохраняем левую таблицу.</details>  
<details style="display: inline;"><summary>Ответ</summary>Возвращает все rows левой relation и matching rows правой, а при отсутствии match — NULL для правой стороны.

192. Когда нужен `LEFT JOIN` в JobFlow?

<details style="display: inline;"><summary>Помощь</summary>Job может существовать без attempts?</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда нужно получить Job даже если для неё ещё нет `job_attempts`.

193. Что такое `ORDER BY`?

<details style="display: inline;"><summary>Помощь</summary>Result ordering.</details>  
<details style="display: inline;"><summary>Ответ</summary>Задаёт порядок строк в результате query.

194. Что такое `LIMIT`?

<details style="display: inline;"><summary>Помощь</summary>Ограничение result size.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ограничивает количество rows, возвращаемых query.

195. Что такое `COUNT(*)`?

<details style="display: inline;"><summary>Помощь</summary>Aggregate.</details>  
<details style="display: inline;"><summary>Ответ</summary>Возвращает количество rows, попавших в соответствующий result set/group.

196. Что такое `GROUP BY`?

<details style="display: inline;"><summary>Помощь</summary>Несколько rows → groups.</details>  
<details style="display: inline;"><summary>Ответ</summary>Группирует rows по выражениям, чтобы применять aggregates к каждой группе.

197. Что такое index?

<details style="display: inline;"><summary>Помощь</summary>Отдельная структура для доступа.</details>  
<details style="display: inline;"><summary>Ответ</summary>Auxiliary data structure, позволяющая ускорить некоторые access patterns ценой storage и maintenance.

198. Почему index может замедлять INSERT/UPDATE?

<details style="display: inline;"><summary>Помощь</summary>Index тоже нужно поддерживать.</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменение данных может требовать обновления связанных index structures.

199. Гарантирует ли наличие index его использование planner-ом?

<details style="display: inline;"><summary>Помощь</summary>Кто выбирает execution plan?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. PostgreSQL planner выбирает plan на основе стоимости, статистики и характеристик query/data.

200. Что показывает `EXPLAIN`?

<details style="display: inline;"><summary>Помощь</summary>Planner до фактического выполнения.</details>  
<details style="display: inline;"><summary>Ответ</summary>`EXPLAIN` показывает выбранный PostgreSQL execution plan и оценочную стоимость операций; `EXPLAIN ANALYZE` дополнительно выполняет запрос и показывает фактические показатели.

## POSTGRESQL — QUERY PLANNING / INDEXES

**[К оглавлению](#500-interview-questions)**

201. Что такое `EXPLAIN ANALYZE`?

<details style="display: inline;"><summary>Помощь</summary>Чем он отличается от обычного `EXPLAIN`?</details>  
<details style="display: inline;"><summary>Ответ</summary>`EXPLAIN ANALYZE` реально выполняет запрос и показывает фактические execution statistics вместе с выбранным планом; обычный `EXPLAIN` показывает план без фактического выполнения.</details>

202. Почему `EXPLAIN ANALYZE` нужно осторожно использовать на `UPDATE` и `DELETE`?

<details style="display: inline;"><summary>Помощь</summary>Он действительно запускает запрос.</details>  
<details style="display: inline;"><summary>Ответ</summary>Потому что `EXPLAIN ANALYZE` выполняет statement, поэтому DML действительно может изменить данные.</details>

203. Что такое sequential scan?

<details style="display: inline;"><summary>Помощь</summary>Как таблица читается без index lookup?</details>  
<details style="display: inline;"><summary>Ответ</summary>План, при котором PostgreSQL последовательно читает relation и проверяет строки на соответствие условию.</details>

204. Когда sequential scan может быть лучше index scan?

<details style="display: inline;"><summary>Помощь</summary>Что если query нужен почти весь table?</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда нужно прочитать большую долю строк и последовательное чтение дешевле множества отдельных index/table accesses.</details>

205. Что такое selectivity запроса?

<details style="display: inline;"><summary>Помощь</summary>Какую долю данных выбирает predicate?</details>  
<details style="display: inline;"><summary>Ответ</summary>Характеристика того, какая доля rows удовлетворяет условию query; высокая selectivity обычно означает мало matching rows.</details>

206. Что такое cardinality?

<details style="display: inline;"><summary>Помощь</summary>Количество элементов.</details>  
<details style="display: inline;"><summary>Ответ</summary>Количество distinct или total элементов в зависимости от контекста; для query planning это одна из характеристик распределения данных и ожидаемого количества rows.

207. Почему index на `status` может оказаться бесполезным?

<details style="display: inline;"><summary>Помощь</summary>Представьте, что 90% строк имеют `status = 'running'`.</details>  
<details style="display: inline;"><summary>Ответ</summary>При низкой selectivity почти вся таблица всё равно нужна, поэтому sequential scan может оказаться дешевле index access.

208. Что такое composite index?

<details style="display: inline;"><summary>Помощь</summary>Один index на несколько columns.</details>  
<details style="display: inline;"><summary>Ответ</summary>Index, построенный по нескольким columns в заданном порядке.

209. Почему порядок columns в composite index важен?

<details style="display: inline;"><summary>Помощь</summary>Index не является просто множеством columns.</details>  
<details style="display: inline;"><summary>Ответ</summary>Порядок определяет, для каких prefix predicates, сортировок и условий index наиболее эффективно используется.

210. Для запроса `WHERE status = $1 ORDER BY created_at DESC` какой index потенциально полезен?

<details style="display: inline;"><summary>Помощь</summary>Подумайте о filter и ordering вместе.</details>  
<details style="display: inline;"><summary>Ответ</summary>Например, composite index `(status, created_at DESC)` может соответствовать такому access pattern, но конкретную эффективность нужно проверять через план и реальные данные.

211. Почему index проектируют от query patterns, а не от названий columns?

<details style="display: inline;"><summary>Помощь</summary>Index оптимизирует access pattern.</details>  
<details style="display: inline;"><summary>Ответ</summary>Индекс полезен не сам по себе, а для конкретных фильтров, ordering и способов доступа к данным.

212. Что такое covering index в общем смысле?

<details style="display: inline;"><summary>Помощь</summary>Может ли query получить всё из index?</details>  
<details style="display: inline;"><summary>Ответ</summary>Index, содержащий все данные, необходимые конкретному query, потенциально позволяющий избежать чтения основной table.

213. Почему нельзя добавлять indexes «на всякий случай»?

<details style="display: inline;"><summary>Помощь</summary>У index есть стоимость.</details>  
<details style="display: inline;"><summary>Ответ</summary>Indexes занимают storage, увеличивают стоимость writes и maintenance и могут не использоваться planner-ом.

214. Что такое query plan?

<details style="display: inline;"><summary>Помощь</summary>Как PostgreSQL собирается выполнить query?</details>  
<details style="display: inline;"><summary>Ответ</summary>План — выбранная planner-ом последовательность операций, например scans, joins, sorts и aggregates, для получения результата.

215. Почему estimated rows и actual rows могут сильно различаться?

<details style="display: inline;"><summary>Помощь</summary>Planner использует statistics.</details>  
<details style="display: inline;"><summary>Ответ</summary>Статистика может быть устаревшей или недостаточно точной для распределения данных, поэтому оценка cardinality отличается от фактической.

216. Почему stale statistics могут ухудшить performance?

<details style="display: inline;"><summary>Помощь</summary>Planner выбирает plan по оценкам.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ошибочные оценки приводят planner к неудачному выбору join strategy, scan method или других операций.

217. Что такое query latency?

<details style="display: inline;"><summary>Помощь</summary>Время запроса от начала до результата.</details>  
<details style="display: inline;"><summary>Ответ</summary>Время, необходимое query для завершения и возврата результата.

218. Почему среднее latency недостаточно для production?

<details style="display: inline;"><summary>Помощь</summary>Что могут скрывать spikes?</details>  
<details style="display: inline;"><summary>Ответ</summary>Среднее не показывает распределение и tail latency; обычно важны p95/p99 и worst-case behavior.

219. Что такое p95 latency?

<details style="display: inline;"><summary>Помощь</summary>95 процентов запросов где?</details>  
<details style="display: inline;"><summary>Ответ</summary>Значение, ниже которого находится примерно 95% измеренных latency observations.

220. Как связаны index, query plan и latency?

<details style="display: inline;"><summary>Помощь</summary>Index меняет доступные plan choices.</details>  
<details style="display: inline;"><summary>Ответ</summary>Index расширяет возможные access paths, planner выбирает стоимость между ними, и выбранный plan влияет на фактическую latency.

---

## TRANSACTIONS / CONCURRENCY

**[К оглавлению](#500-interview-questions)**

221. Почему несколько SQL statements могут требовать одной transaction?

<details style="display: inline;"><summary>Помощь</summary>Что если первый прошёл, второй упал?</details>  
<details style="display: inline;"><summary>Ответ</summary>Если statements образуют одну логическую state transition, они должны либо все зафиксироваться, либо все отмениться.

222. Что произойдёт без transaction в операции «Job + Attempt»?

<details style="display: inline;"><summary>Помощь</summary>Рассмотрите failure между двумя INSERT.</details>  
<details style="display: inline;"><summary>Ответ</summary>Можно получить частично записанное состояние: Job существует, а Attempt отсутствует, либо наоборот, в зависимости от порядка операций.

223. Что означает «transaction boundary соответствует business operation»?

<details style="display: inline;"><summary>Помощь</summary>Какое состояние должно меняться вместе?</details>  
<details style="display: inline;"><summary>Ответ</summary>Все database changes, необходимые для одной логической state transition, находятся внутри одной transaction.

224. Почему transaction не должна быть слишком длинной?

<details style="display: inline;"><summary>Помощь</summary>Locks, resources, conflicts.</details>  
<details style="display: inline;"><summary>Ответ</summary>Длинные transactions удерживают resources и locks дольше, увеличивают contention, latency и вероятность конфликтов.

225. Что такое row-level lock?

<details style="display: inline;"><summary>Помощь</summary>Блокируется вся table или отдельная row?</details>  
<details style="display: inline;"><summary>Ответ</summary>Lock, который ограничивает определённые конкурентные операции над конкретными rows.

226. Что делает `FOR UPDATE`?

<details style="display: inline;"><summary>Помощь</summary>Нам нужно безопасно изменить найденную row.</details>  
<details style="display: inline;"><summary>Ответ</summary>`SELECT ... FOR UPDATE` блокирует выбранные rows для конфликтующих изменений до завершения текущей transaction.

227. Когда `FOR UPDATE` особенно полезен?

<details style="display: inline;"><summary>Помощь</summary>Read → check → update.</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда сначала нужно прочитать состояние row, проверить invariant, а затем безопасно изменить эту же row.

228. Почему обычный `SELECT` перед `UPDATE` может быть недостаточен?

<details style="display: inline;"><summary>Помощь</summary>Между ними может вмешаться другая transaction.</details>  
<details style="display: inline;"><summary>Ответ</summary>Другой transaction может изменить state между чтением и записью, создавая race или lost update.

229. Что такое read-modify-write race?

<details style="display: inline;"><summary>Помощь</summary>Две transactions читают одинаковое значение.</details>  
<details style="display: inline;"><summary>Ответ</summary>Несколько transactions читают одно состояние, независимо рассчитывают новое и потом записывают его, потенциально теряя updates.

230. Как выглядит lost update?

<details style="display: inline;"><summary>Помощь</summary>Значение было 3.</details>  
<details style="display: inline;"><summary>Ответ</summary>A и B читают `3`, обе рассчитывают `4`, затем обе записывают `4`; один из двух logical increments теряется.

231. Как защитить increment от lost update?

<details style="display: inline;"><summary>Помощь</summary>Можно ли менять state прямо в SQL?</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать атомарный SQL update вроде `SET attempt = attempt + 1`, либо подходящую transaction/locking strategy.

232. Почему `SET attempt = attempt + 1` часто лучше read → calculate → write?

<details style="display: inline;"><summary>Помощь</summary>Где происходит изменение?</details>  
<details style="display: inline;"><summary>Ответ</summary>Read-modify-write выполняется внутри database operation, что уменьшает race window и позволяет СУБД координировать concurrent updates.

233. Достаточно ли одного atomic SQL update для любого invariant?

<details style="display: inline;"><summary>Помощь</summary>А если нужно проверить несколько условий и изменить две таблицы?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Для сложных state transitions могут понадобиться transaction, constraints и locking.

234. Что такое isolation level?

<details style="display: inline;"><summary>Помощь</summary>Как transactions взаимодействуют?</details>  
<details style="display: inline;"><summary>Ответ</summary>Набор гарантий относительно видимости и взаимодействия concurrent transactions.

235. Какой isolation level является default в PostgreSQL?

<details style="display: inline;"><summary>Помощь</summary>Базовая настройка PostgreSQL.</details>  
<details style="display: inline;"><summary>Ответ</summary>`READ COMMITTED`.

236. Что примерно означает `READ COMMITTED` в PostgreSQL?

<details style="display: inline;"><summary>Помощь</summary>Visibility statement-by-statement.</details>  
<details style="display: inline;"><summary>Ответ</summary>Каждый statement видит данные, committed к моменту его начала, с semantics PostgreSQL MVCC.

237. Что даёт `REPEATABLE READ`?

<details style="display: inline;"><summary>Помощь</summary>Более стабильный view.</details>  
<details style="display: inline;"><summary>Ответ</summary>Transaction получает более стабильное snapshot представление данных и более сильные repeatability guarantees.

238. Что делает `SERIALIZABLE`?

<details style="display: inline;"><summary>Помощь</summary>Как будто transactions выполнялись последовательно.</details>  
<details style="display: inline;"><summary>Ответ</summary>Пытается обеспечить эффект эквивалентный сериализуемому порядку выполнения и может отклонить конфликтующую transaction.

239. Почему `SERIALIZABLE` может вернуть ошибку вместо ожидания?

<details style="display: inline;"><summary>Помощь</summary>Conflict detection.</details>  
<details style="display: inline;"><summary>Ответ</summary>Если выбранный concurrency schedule нельзя безопасно считать serializable, PostgreSQL может завершить transaction serialization failure.

240. Как правильно retry-ить serialization failure?

<details style="display: inline;"><summary>Помощь</summary>Старая transaction уже закончилась.</details>  
<details style="display: inline;"><summary>Ответ</summary>Rollback текущую transaction, затем начать и выполнить всю transaction заново по bounded retry policy.

241. Почему нельзя продолжить использовать transaction после serialization failure?

<details style="display: inline;"><summary>Помощь</summary>Её outcome уже определён.</details>  
<details style="display: inline;"><summary>Ответ</summary>Transaction больше не является действующей успешной execution boundary; retry требует новой transaction.

242. Что такое MVCC?

<details style="display: inline;"><summary>Помощь</summary>Multi-version.</details>  
<details style="display: inline;"><summary>Ответ</summary>Multi-Version Concurrency Control — механизм, использующий версии данных и snapshots для координации concurrent access.

243. Почему MVCC уменьшает необходимость блокировать обычные чтения?

<details style="display: inline;"><summary>Помощь</summary>Читатель может видеть подходящую версию.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтение может работать с согласованной версией данных, не требуя блокировать writer во всех случаях.

244. Является ли MVCC заменой всех locks?

<details style="display: inline;"><summary>Помощь</summary>Что если нужно защитить изменение?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Для некоторых state transitions нужны row locks и другие synchronization mechanisms.

245. Что такое transaction isolation anomaly?

<details style="display: inline;"><summary>Помощь</summary>Unexpected interaction between transactions.</details>  
<details style="display: inline;"><summary>Ответ</summary>Наблюдение или effect concurrent execution, который допускается выбранным isolation level и не соответствовал бы более сильной модели.

246. Почему isolation level выбирают под invariant?

<details style="display: inline;"><summary>Помощь</summary>Не существует одной «самой правильной» настройки.</details>  
<details style="display: inline;"><summary>Ответ</summary>Разные invariants требуют разных гарантий, а более строгая isolation имеет resource/performance cost.

247. Что такое lock contention?

<details style="display: inline;"><summary>Помощь</summary>Много transactions ждут один ресурс.</details>  
<details style="display: inline;"><summary>Ответ</summary>Конкурентные operations конфликтуют за lock и вынуждены ждать друг друга.

248. Как уменьшить lock contention?

<details style="display: inline;"><summary>Помощь</summary>Уменьшите размер и длительность критической секции.</details>  
<details style="display: inline;"><summary>Ответ</summary>Делать transactions короче, уменьшать locked set, выбирать правильные indexes и избегать ненужной работы внутри transaction.

249. Почему порядок получения locks важен?

<details style="display: inline;"><summary>Помощь</summary>Две transactions берут A и B.</details>  
<details style="display: inline;"><summary>Ответ</summary>Разный порядок может создать circular wait и deadlock.

250. Что делать при deadlock detected?

<details style="display: inline;"><summary>Помощь</summary>Database уже выбрала victim.</details>  
<details style="display: inline;"><summary>Ответ</summary>Корректно завершить текущую transaction и при допустимой retry policy повторить операцию; дополнительно устранить причину lock cycle.

---

## ACID / DATA INTEGRITY

**[К оглавлению](#500-interview-questions)**

251. Что означает Atomicity в ACID?

<details style="display: inline;"><summary>Помощь</summary>Все изменения вместе.</details>  
<details style="display: inline;"><summary>Ответ</summary>Transaction фиксирует группу изменений целиком либо не оставляет её частично committed.

252. Что означает Consistency в ACID?

<details style="display: inline;"><summary>Помощь</summary>Что после commit должно быть истинно?</details>  
<details style="display: inline;"><summary>Ответ</summary>После успешного commit database state продолжает удовлетворять constraints и заданным invariants.

253. Что означает Isolation в ACID?

<details style="display: inline;"><summary>Помощь</summary>Concurrent transactions.</details>  
<details style="display: inline;"><summary>Ответ</summary>Определяет гарантии взаимодействия concurrent transactions и видимости их изменений.

254. Что означает Durability?

<details style="display: inline;"><summary>Помощь</summary>После commit.</details>  
<details style="display: inline;"><summary>Ответ</summary>Успешно committed changes сохраняются в пределах гарантий durability конкретной PostgreSQL configuration/environment.

255. Почему Consistency в ACID не означает «данные логически идеальны всегда»?

<details style="display: inline;"><summary>Помощь</summary>Кто задаёт invariants?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consistency относится к сохранению определённых constraints и invariants системы, а не к автоматическому пониманию бизнес-смысла всей предметной области.

256. Можно ли обеспечить Consistency только средствами PostgreSQL?

<details style="display: inline;"><summary>Помощь</summary>Какие invariants принадлежат database?</details>  
<details style="display: inline;"><summary>Ответ</summary>Часть invariants отлично выражается constraints и transactions, но business invariants могут требовать application logic и межсистемной координации.

257. Зачем нужен `CHECK` вместе с Go validation?

<details style="display: inline;"><summary>Помощь</summary>Defense in depth.</details>  
<details style="display: inline;"><summary>Ответ</summary>Go validation улучшает API feedback, а database constraint защищает persisted state от всех путей записи.

258. Что лучше: проверять `status` в Go или constraint в БД?

<details style="display: inline;"><summary>Помощь</summary>Они решают разные задачи.</details>  
<details style="display: inline;"><summary>Ответ</summary>Оба уровня полезны: Go даёт раннюю и понятную validation, а database constraint обеспечивает authoritative state invariant.

259. Что такое source of truth?

<details style="display: inline;"><summary>Помощь</summary>Откуда можно восстановить состояние?</details>  
<details style="display: inline;"><summary>Ответ</summary>Компонент, чьё состояние считается авторитетным и из которого система может восстановить корректное значение.

260. Почему Redis cache не должен автоматически становиться source of truth?

<details style="display: inline;"><summary>Помощь</summary>Cache может исчезнуть.</details>  
<details style="display: inline;"><summary>Ответ</summary>Cache может быть evicted, stale или недоступен; authoritative state должен иметь более подходящую persistence semantics.

---

## Pgxpool / Go + PostgreSQL

**[К оглавлению](#500-interview-questions)**

261. Что такое `pgxpool.Pool`?

<details style="display: inline;"><summary>Помощь</summary>Shared DB resource.</details>  
<details style="display: inline;"><summary>Ответ</summary>Connection pool, управляющий набором PostgreSQL connections и их lifecycle.

262. Почему backend использует pool, а не одно DB connection?

<details style="display: inline;"><summary>Помощь</summary>Много concurrent operations.</details>  
<details style="display: inline;"><summary>Ответ</summary>Pool позволяет обслуживать множество concurrent operations и ограничивать database concurrency.

263. Что означает `MaxConns`?

<details style="display: inline;"><summary>Помощь</summary>Верхняя граница.</details>  
<details style="display: inline;"><summary>Ответ</summary>Максимальное число database connections, одновременно выдаваемых pool.

264. Что произойдёт при 100 workers и `MaxConns = 20`?

<details style="display: inline;"><summary>Помощь</summary>Не каждый worker получает connection.</details>  
<details style="display: inline;"><summary>Ответ</summary>Не более 20 operations одновременно будут владеть pool connections; остальные могут ждать свободный connection.

265. Почему connection pool является backpressure mechanism?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт при saturation?</details>  
<details style="display: inline;"><summary>Ответ</summary>Pool ограничивает concurrent database access и заставляет лишнюю работу ждать или получать timeout.

266. Зачем передавать `context.Context` в `Query`?

<details style="display: inline;"><summary>Помощь</summary>Request cancellation.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы database operation могла остановиться при cancellation или deadline.

267. Что делать, если context уже cancelled перед DB call?

<details style="display: inline;"><summary>Помощь</summary>Нужна ли работа?</details>  
<details style="display: inline;"><summary>Ответ</summary>Операция должна быстро завершиться с cancellation/error semantics, не выполняя ненужный downstream work.

268. Как отобразить PostgreSQL `no rows` в application error?

<details style="display: inline;"><summary>Помощь</summary>Есть ли domain identity?</details>  
<details style="display: inline;"><summary>Ответ</summary>Adapter может преобразовать driver-specific no-row error в application/domain error, например `ErrJobNotFound`.

269. Почему service не должен сравнивать pgx-specific errors повсеместно?

<details style="display: inline;"><summary>Помощь</summary>Кому принадлежит infrastructure knowledge?</details>  
<details style="display: inline;"><summary>Ответ</summary>Это связывает application policy с infrastructure implementation; adapter должен переводить низкоуровневые errors в стабильный application contract.

270. Что такое repository adapter?

<details style="display: inline;"><summary>Помощь</summary>Port → concrete DB.</details>  
<details style="display: inline;"><summary>Ответ</summary>Concrete component, который реализует repository contract поверх PostgreSQL или другого storage.

---

## MIGRATIONS / SCHEMA EVOLUTION

**[К оглавлению](#500-interview-questions)**

271. Что такое migration?

<details style="display: inline;"><summary>Помощь</summary>Transition schema.</details>  
<details style="display: inline;"><summary>Ответ</summary>Versioned изменение database schema от одного известного состояния к следующему.

272. Почему schema changes нужно хранить в repository?

<details style="display: inline;"><summary>Помощь</summary>Application code тоже versioned.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы schema evolution была воспроизводимой, reviewable и синхронизированной с application code.

273. Что такое schema reproducibility?

<details style="display: inline;"><summary>Помощь</summary>Можно ли получить одинаковую schema с нуля?</details>  
<details style="display: inline;"><summary>Ответ</summary>Возможность построить ожидаемую schema на чистой database только последовательным применением versioned migrations.

274. Что такое schema drift?

<details style="display: inline;"><summary>Помощь</summary>Runtime schema отличается от migrations.</details>  
<details style="display: inline;"><summary>Ответ</summary>Расхождение фактической database schema с той, которую описывают versioned migrations.

275. Почему ручной `ALTER TABLE` в production без migration опасен?

<details style="display: inline;"><summary>Помощь</summary>Как повторить change?</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменение становится плохо воспроизводимым и может отсутствовать на других окружениях.

276. Что должен содержать migration `Up`?

<details style="display: inline;"><summary>Помощь</summary>Как движется schema?</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменения, переводящие schema к новому состоянию.

277. Что должен содержать migration `Down`?

<details style="display: inline;"><summary>Помощь</summary>Как вернуться назад?</details>  
<details style="display: inline;"><summary>Ответ</summary>Контролируемую обратную операцию, отменяющую изменения конкретной migration там, где rollback поддерживается policy проекта.

278. Почему rollback не всегда безопасен?

<details style="display: inline;"><summary>Помощь</summary>Data loss.</details>  
<details style="display: inline;"><summary>Ответ</summary>Некоторые schema changes могут уничтожать или необратимо трансформировать data, поэтому rollback должен учитывать фактические последствия.

279. Что такое backward-compatible schema change?

<details style="display: inline;"><summary>Помощь</summary>Старый code может продолжать работать.</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменение, при котором существующие consumers/application versions продолжают работать в пределах согласованного migration strategy.

280. Почему добавление nullable column часто безопаснее немедленного добавления `NOT NULL`?

<details style="display: inline;"><summary>Помощь</summary>Что делать с существующими rows?</details>  
<details style="display: inline;"><summary>Ответ</summary>Существующие rows могут не иметь значения; nullable column позволяет поэтапно заполнить data и затем усилить constraint.

---

## DATABASE DESIGN / JOBFLOW

**[К оглавлению](#500-interview-questions)**

281. Почему `job_attempts` лучше отдельной таблицы?

<details style="display: inline;"><summary>Помощь</summary>Одна Job может иметь много attempts.</details>  
<details style="display: inline;"><summary>Ответ</summary>Attempt имеет собственную identity, cardinality и lifecycle, поэтому отношение 1:N лучше выражается отдельной relation.

282. Как выразить связь Job → Attempts?

<details style="display: inline;"><summary>Помощь</summary>Foreign key.</details>  
<details style="display: inline;"><summary>Ответ</summary>Через `job_attempts.job_id REFERENCES jobs(id)`.

283. Как защитить attempt от несуществующей Job?

<details style="display: inline;"><summary>Помощь</summary>Database relationship.</details>  
<details style="display: inline;"><summary>Ответ</summary>Foreign key constraint.

284. Как запретить пустой `Job.type` на уровне DB?

<details style="display: inline;"><summary>Помощь</summary>Null ≠ empty string.</details>  
<details style="display: inline;"><summary>Ответ</summary>Например, `type text NOT NULL` защищает от NULL; если нужен запрет пустой строки, дополнительно требуется соответствующий `CHECK`.

285. Зачем `updated_at`?

<details style="display: inline;"><summary>Помощь</summary>State change tracking.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для фиксации времени последнего изменения состояния и диагностики/ordering/use cases.

286. Почему `created_at` и `updated_at` обычно `timestamptz`?

<details style="display: inline;"><summary>Помощь</summary>Timezone-aware timestamp.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для корректной модели абсолютного времени с timezone semantics; application может хранить и отображать его в UTC.

287. Что такое surrogate key?

<details style="display: inline;"><summary>Помощь</summary>Техническая identity.</details>  
<details style="display: inline;"><summary>Ответ</summary>Идентификатор строки, не имеющий непосредственного business meaning, например UUID.

288. Почему UUID может быть удобен для Job ID?

<details style="display: inline;"><summary>Помощь</summary>Генерация identity.</details>  
<details style="display: inline;"><summary>Ответ</summary>Он позволяет распределённо генерировать уникальные identifiers без центрального sequence coordinator.

289. Есть ли у UUID только преимущества?

<details style="display: inline;"><summary>Помощь</summary>Index locality и размер.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Он может иметь storage/index locality trade-offs и другой performance profile по сравнению с последовательными numeric keys.

290. Почему бизнес-состояние Job лучше хранить в явном `status`?

<details style="display: inline;"><summary>Помощь</summary>State machine.</details>  
<details style="display: inline;"><summary>Ответ</summary>Явный status делает lifecycle state machine наблюдаемой, проверяемой и удобной для conditional transitions.

---

## DATA ACCESS / SQL SAFETY

**[К оглавлению](#500-interview-questions)**

291. Почему query лучше параметризовать?

<details style="display: inline;"><summary>Помощь</summary>SQL structure vs values.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для безопасности, корректного escaping и явного разделения SQL syntax и user/data values.

292. Как защититься от SQL injection?

<details style="display: inline;"><summary>Помощь</summary>Не вставляйте user input в SQL text.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать parameterized queries и не строить SQL через ручную конкатенацию пользовательских значений.

293. Почему SQL injection — не только security problem?

<details style="display: inline;"><summary>Помощь</summary>Что ещё может ломаться?</details>  
<details style="display: inline;"><summary>Ответ</summary>Неправильное quoting и typing могут приводить к syntax/runtime errors и неожиданной интерпретации данных даже без злоумышленника.

294. Что должно произойти, если `Get(jobID)` не нашёл Job?

<details style="display: inline;"><summary>Помощь</summary>Application contract.</details>  
<details style="display: inline;"><summary>Ответ</summary>Adapter должен вернуть стабильную application error semantics, например `ErrJobNotFound`.

295. Почему repository должен возвращать domain/application error, а не HTTP 404?

<details style="display: inline;"><summary>Помощь</summary>Кто должен знать transport?</details>  
<details style="display: inline;"><summary>Ответ</summary>Repository не должен зависеть от HTTP transport; HTTP mapping выполняет transport adapter.

296. Что такое repository contract?

<details style="display: inline;"><summary>Помощь</summary>Какие capabilities нужны application?</details>  
<details style="display: inline;"><summary>Ответ</summary>Минимальный набор storage operations и их error semantics, необходимых consumer для выполнения use case.

297. Почему repository interface не должен отражать весь SQL client API?

<details style="display: inline;"><summary>Помощь</summary>Consumer требует меньше.</details>  
<details style="display: inline;"><summary>Ответ</summary>Application нужен business-oriented contract, а не огромная infrastructure abstraction.

298. Что такое N+1 query problem?

<details style="display: inline;"><summary>Помощь</summary>Один список и отдельный query для каждого элемента.</details>  
<details style="display: inline;"><summary>Ответ</summary>Сценарий, когда один query получает N entities, а затем выполняется ещё N отдельных queries для связанных данных, создавая лишние round trips.

299. Как уменьшить N+1?

<details style="display: inline;"><summary>Помощь</summary>Сделайте retrieval set-based.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать JOIN, batch query или заранее загружать нужный набор связанных данных вместо N отдельных round trips.

300. Как объяснить хороший database repository на интервью одной фразой?

<details style="display: inline;"><summary>Помощь</summary>Состояние, contract, SQL и errors.</details>  
<details style="display: inline;"><summary>Ответ</summary>«Repository должен скрывать конкретный storage mechanism за минимальным contract, выполнять parameterized SQL с корректным context/transaction handling и переводить infrastructure errors в понятную application semantics.»

## KAFKA / DELIVERY SEMANTICS

**[К оглавлению](#500-interview-questions)**

301. Что такое producer в Kafka?

<details style="display: inline;"><summary>Помощь</summary>Кто отправляет message в broker?</details>  
<details style="display: inline;"><summary>Ответ</summary>Producer — компонент, который публикует сообщения в Kafka topic.

302. Что такое consumer в Kafka?

<details style="display: inline;"><summary>Помощь</summary>Кто читает stream?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer — компонент, который читает сообщения из partitions и обрабатывает их.

303. Что такое broker?

<details style="display: inline;"><summary>Помощь</summary>Где durable message хранится?</details>  
<details style="display: inline;"><summary>Ответ</summary>Broker — серверный компонент messaging system, который принимает, хранит и доставляет сообщения producer/consumer.

304. Что такое topic?

<details style="display: inline;"><summary>Помощь</summary>Логический поток сообщений.</details>  
<details style="display: inline;"><summary>Ответ</summary>Topic — логический поток сообщений, разделённый на partitions.

305. Что такое partition?

<details style="display: inline;"><summary>Помощь</summary>Ordering + parallelism.</details>  
<details style="display: inline;"><summary>Ответ</summary>Partition — упорядоченный журнал сообщений и единица распределения consumer work внутри topic.

306. Где Kafka гарантирует ordering?

<details style="display: inline;"><summary>Помощь</summary>Вспомните partition.</details>  
<details style="display: inline;"><summary>Ответ</summary>Порядок гарантируется внутри отдельной partition; глобальный порядок между partitions отсутствует.

307. Почему partition key может быть `job.ID`?

<details style="display: inline;"><summary>Помощь</summary>Все события одной Job.</details>  
<details style="display: inline;"><summary>Ответ</summary>Одинаковый key позволяет направлять события одной Job в одну partition и сохранять порядок между ними.

308. Гарантирует ли partition key глобальный порядок всех Job?

<details style="display: inline;"><summary>Помощь</summary>Одна Job или весь topic?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Он может сохранить порядок для одной ключевой группы, но разные partitions не имеют общего порядка.

309. Что такое consumer group?

<details style="display: inline;"><summary>Помощь</summary>Несколько consumers работают вместе.</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer group — набор consumers, совместно распределяющих partitions одного topic.

310. Как consumer group масштабирует обработку?

<details style="display: inline;"><summary>Помощь</summary>Partitions распределяются между consumers.</details>  
<details style="display: inline;"><summary>Ответ</summary>Разные partitions могут обрабатываться разными consumers одновременно.

311. Почему 100 consumers не дают 100-кратный throughput при 3 partitions?

<details style="display: inline;"><summary>Помощь</summary>Сколько partitions можно одновременно назначить?</details>  
<details style="display: inline;"><summary>Ответ</summary>В одной consumer group одновременно полезно работать смогут не больше consumers, которым назначены partitions; при 3 partitions максимум одновременно активных partition owners — 3.

312. Что такое offset?

<details style="display: inline;"><summary>Помощь</summary>Позиция в partition.</details>  
<details style="display: inline;"><summary>Ответ</summary>Offset — позиция сообщения в partition и точка progress consumer.

313. Зачем consumer commit offset?

<details style="display: inline;"><summary>Помощь</summary>Что нужно знать после restart?</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы зафиксировать progress и определить, откуда продолжить обработку после restart/rebalance.

314. Что такое at-most-once delivery?

<details style="display: inline;"><summary>Помощь</summary>Ноль или один раз.</details>  
<details style="display: inline;"><summary>Ответ</summary>Сообщение обрабатывается не более одного раза, но может быть потеряно.

315. Что такое at-least-once delivery?

<details style="display: inline;"><summary>Помощь</summary>Что принимаем как нормальный failure?</details>  
<details style="display: inline;"><summary>Ответ</summary>Сообщение может быть доставлено повторно, поэтому consumer должен корректно переживать duplicates.

316. Почему at-least-once часто практичнее exactly-once?

<details style="display: inline;"><summary>Помощь</summary>Сравните сложность guarantee.</details>  
<details style="display: inline;"><summary>Ответ</summary>At-least-once проще обеспечить на transport boundary, а duplicate можно сделать безопасным через idempotent effect.

317. Что такое exactly-once semantics?

<details style="display: inline;"><summary>Помощь</summary>Уточните границу guarantee.</details>  
<details style="display: inline;"><summary>Ответ</summary>Гарантия, что определённый logical operation/effect учитывается ровно один раз в конкретной transaction boundary; она не означает автоматически exactly-once end-to-end external effect.

318. Почему нельзя просто сказать «Kafka даёт exactly-once»?

<details style="display: inline;"><summary>Помощь</summary>Transport ≠ external effect.</details>  
<details style="display: inline;"><summary>Ответ</summary>Exactly-once semantics ограничена конкретной областью broker/client transactional model и не распространяется автоматически на arbitrary database, HTTP или внешний side effect.

319. Что происходит при crash после effect, но до offset commit?

<details style="display: inline;"><summary>Помощь</summary>Broker считает progress зафиксированным?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Message может быть redelivered, поэтому effect должен быть idempotent.

320. Что происходит при offset commit до effect?

<details style="display: inline;"><summary>Помощь</summary>Что будет после crash?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer может потерять message без выполнения effect, потому что broker считает его обработанным.

321. Почему порядок `effect → offset commit` важен?

<details style="display: inline;"><summary>Помощь</summary>Сравните loss и duplicate.</details>  
<details style="display: inline;"><summary>Ответ</summary>Такой порядок позволяет при crash получить повторную доставку вместо потери effect; duplicate затем должен быть безопасно обработан.

322. Что такое consumer lag?

<details style="display: inline;"><summary>Помощь</summary>Producer position против consumer progress.</details>  
<details style="display: inline;"><summary>Ответ</summary>Разница между последним доступным положением в partition и текущим progress consumer.

323. Почему высокий lag не обязательно означает, что Kafka сломана?

<details style="display: inline;"><summary>Помощь</summary>Где может быть bottleneck?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer может быть медленным из-за CPU, database, downstream, недостатка partitions или собственной concurrency bound.

324. Что проверять при росте consumer lag?

<details style="display: inline;"><summary>Помощь</summary>Смотрите всю цепочку.</details>  
<details style="display: inline;"><summary>Ответ</summary>Partition count, consumer assignment, processing latency, DB/downstream latency, worker concurrency, resource saturation и retry backlog.

325. Что такое rebalance?

<details style="display: inline;"><summary>Помощь</summary>Изменился состав consumer group.</details>  
<details style="display: inline;"><summary>Ответ</summary>Процесс перераспределения partitions между consumers группы.

326. Что может вызвать rebalance?

<details style="display: inline;"><summary>Помощь</summary>Consumer joins/leaves/fails.</details>  
<details style="display: inline;"><summary>Ответ</summary>Изменение состава группы, failure consumer, restart и другие изменения topology/group state.

327. Как rebalance влияет на processing?

<details style="display: inline;"><summary>Помощь</summary>Ownership partitions.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ownership partitions временно меняется, что может вызвать задержки и требует корректной lifecycle/progress semantics.

328. Почему consumer должен быть idempotent при rebalance?

<details style="display: inline;"><summary>Помощь</summary>Progress и effect могут разделиться.</details>  
<details style="display: inline;"><summary>Ответ</summary>Во время restart/rebalance возможна повторная доставка сообщений, поэтому повторный effect должен быть безопасным.

329. Что такое retryable Kafka error?

<details style="display: inline;"><summary>Помощь</summary>Temporary или permanent?</details>  
<details style="display: inline;"><summary>Ответ</summary>Ошибка, для которой есть разумное ожидание, что повторная попытка позже может завершиться успешно.

330. Почему retry любого Kafka error опасен?

<details style="display: inline;"><summary>Помощь</summary>Что если payload навсегда invalid?</details>  
<details style="display: inline;"><summary>Ответ</summary>Permanent errors могут бесконечно повторяться, создавая retry storm и блокируя progress.

---

## IDEMPOTENCY

**[К оглавлению](#500-interview-questions)**

331. Что такое idempotency?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт при повторе одной логической операции?</details>  
<details style="display: inline;"><summary>Ответ</summary>Повторение одной logical operation не создаёт новый нежелательный effect.

332. Является ли idempotency тем же, что deterministic behavior?

<details style="display: inline;"><summary>Помощь</summary>Повторяемость и детерминизм — разные свойства.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Deterministic behavior зависит от одинаковых inputs/условий, а idempotency — от отсутствия дополнительного effect при повторе logical operation.

333. Что такое idempotency key?

<details style="display: inline;"><summary>Помощь</summary>Stable identity operation.</details>  
<details style="display: inline;"><summary>Ответ</summary>Уникальный стабильный идентификатор logical operation, позволяющий распознать повтор.

334. Где хранить idempotency state?

<details style="display: inline;"><summary>Помощь</summary>Где guarantee должна переживать restart?</details>  
<details style="display: inline;"><summary>Ответ</summary>В durable storage, если повтор должен быть распознан после process restart и между несколькими consumers.

335. Почему in-memory map недостаточен для глобальной idempotency?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт после restart или второго instance?</details>  
<details style="display: inline;"><summary>Ответ</summary>State исчезнет при restart и не будет shared между instances.

336. Как защитить `event_id` от duplicate на уровне PostgreSQL?

<details style="display: inline;"><summary>Помощь</summary>Identity constraint.</details>  
<details style="display: inline;"><summary>Ответ</summary>Создать `PRIMARY KEY` или `UNIQUE` constraint по `event_id`.

337. Что такое `processed_events`?

<details style="display: inline;"><summary>Помощь</summary>Persistent deduplication state.</details>  
<details style="display: inline;"><summary>Ответ</summary>Таблица, фиксирующая уже обработанные logical events для защиты от повторного effect.

338. Достаточно ли просто вставить `event_id` после обработки?

<details style="display: inline;"><summary>Помощь</summary>Что если crash между effect и INSERT?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Registration of processed event и business state change должны быть связаны в одной соответствующей transaction boundary, если это возможно.

339. Почему idempotency и transaction часто идут вместе?

<details style="display: inline;"><summary>Помощь</summary>Как атомарно отметить effect?</details>  
<details style="display: inline;"><summary>Ответ</summary>Transaction позволяет атомарно связать deduplication marker с authoritative state change.

340. Как сделать state transition idempotent без отдельной таблицы?

<details style="display: inline;"><summary>Помощь</summary>Условный UPDATE.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать transition, которая применяется только из ожидаемого предыдущего состояния, например `WHERE status = 'pending'`.

341. Что произойдёт при повторном `pending → running` через conditional UPDATE?

<details style="display: inline;"><summary>Помощь</summary>Как изменится WHERE после первого перехода?</details>  
<details style="display: inline;"><summary>Ответ</summary>Второй UPDATE затронет 0 rows, потому что status уже `running`.

342. Всегда ли 0 rows означает duplicate?

<details style="display: inline;"><summary>Помощь</summary>Что ещё могло случиться?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Это может означать unknown Job, invalid state, concurrent transition или другой бизнес-сценарий; semantics нужно определить явно.

343. Может ли HTTP `POST` быть idempotent?

<details style="display: inline;"><summary>Помощь</summary>Semantics важнее HTTP verb.</details>  
<details style="display: inline;"><summary>Ответ</summary>Да, если API использует idempotency key и повтор запроса не создаёт дополнительного logical effect.

344. Зачем idempotency особенно важна для payments?

<details style="display: inline;"><summary>Помощь</summary>Что страшнее duplicate?</details>  
<details style="display: inline;"><summary>Ответ</summary>Повторный financial effect может создать реальный двойной charge или другое критическое состояние.

345. Что такое exactly-once effect в практическом backend смысле?

<details style="display: inline;"><summary>Помощь</summary>Как достичь эффекта, если delivery at-least-once?</details>  
<details style="display: inline;"><summary>Ответ</summary>Можно получить exactly-once logical effect на конкретной authoritative boundary через transaction + idempotency, даже если transport доставляет message повторно.

346. Чем deduplication отличается от idempotency?

<details style="display: inline;"><summary>Помощь</summary>Discard duplicate или безопасно повторить.</details>  
<details style="display: inline;"><summary>Ответ</summary>Deduplication распознаёт и отбрасывает повтор; idempotency означает, что повторная операция безопасна относительно итогового effect. Они могут использоваться вместе.

347. Почему event ID должен быть stable?

<details style="display: inline;"><summary>Помощь</summary>Duplicate должен выглядеть как тот же event.</details>  
<details style="display: inline;"><summary>Ответ</summary>Повторное представление одного logical event должно иметь одинаковый identity, иначе deduplication не распознает повтор.

348. Что делать, если producer генерирует новый ID при каждом retry?

<details style="display: inline;"><summary>Помощь</summary>Для системы это один operation или несколько?</details>  
<details style="display: inline;"><summary>Ответ</summary>Если это одна logical operation, identity должна сохраняться между retries; иначе система увидит несколько разных operations.

349. Какой invariant можно сформулировать для `processed_events`?

<details style="display: inline;"><summary>Помощь</summary>Один event ID.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для каждого logical `event_id` должно быть не более одного зарегистрированного successful processing record.

350. Где проверять idempotency key?

<details style="display: inline;"><summary>Помощь</summary>До или после effect?</details>  
<details style="display: inline;"><summary>Ответ</summary>Проверка должна быть частью atomic processing boundary, чтобы concurrent duplicates не прошли одновременно до effect.

---

## OUTBOX

**[К оглавлению](#500-interview-questions)**

351. Что такое Outbox pattern?

<details style="display: inline;"><summary>Помощь</summary>DB state + event intent.</details>  
<details style="display: inline;"><summary>Ответ</summary>Паттерн, в котором изменение authoritative state и запись намерения опубликовать event фиксируются одной database transaction, после чего отдельный relay отправляет event в broker.

352. Какую проблему решает Outbox?

<details style="display: inline;"><summary>Помощь</summary>Gap между DB и Kafka.</details>  
<details style="display: inline;"><summary>Ответ</summary>Устраняет риск, что DB state committed, а message intent потерян из-за отказа между двумя независимыми системами.

353. Что было бы без Outbox?

<details style="display: inline;"><summary>Помощь</summary>`DB commit` и `Kafka publish` отдельно.</details>  
<details style="display: inline;"><summary>Ответ</summary>Мог возникнуть partial failure: Job committed, а Kafka event lost, либо event published для состояния, которое не было committed.

354. Что хранится в `outbox_events`?

<details style="display: inline;"><summary>Помощь</summary>Какое минимальное event state нужно relay?</details>  
<details style="display: inline;"><summary>Ответ</summary>Event identity, type, aggregate identity, payload, creation time и состояние доставки, например `published_at`.

355. Почему `event_id` в Outbox должен быть unique?

<details style="display: inline;"><summary>Помощь</summary>Identity logical event.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы один logical event имел стабильную identity и duplicate generation можно было контролировать.

356. Что происходит при `DB INSERT + outbox INSERT + COMMIT`?

<details style="display: inline;"><summary>Помощь</summary>Одна local transaction.</details>  
<details style="display: inline;"><summary>Ответ</summary>Job state и intent to publish event становятся атомарно зафиксированными в PostgreSQL.

357. Что происходит, если outbox INSERT падает?

<details style="display: inline;"><summary>Помощь</summary>Transaction boundary.</details>  
<details style="display: inline;"><summary>Ответ</summary>Transaction rollback отменяет и state change Job, если оба находятся в одной transaction.

358. Что происходит, если Kafka недоступна после DB commit?

<details style="display: inline;"><summary>Помощь</summary>Event уже durable где-то?</details>  
<details style="display: inline;"><summary>Ответ</summary>Outbox event остаётся в PostgreSQL и может быть повторно отправлен relay после восстановления broker.

359. Решает ли Outbox проблему duplicate publish?

<details style="display: inline;"><summary>Помощь</summary>Crash window relay.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Relay может опубликовать event и упасть до `published_at`, после чего event может быть опубликован повторно.

360. Какой delivery semantic естественно получается у простого Outbox relay?

<details style="display: inline;"><summary>Помощь</summary>Publish then mark.</details>  
<details style="display: inline;"><summary>Ответ</summary>Обычно at-least-once, потому что duplicate publish безопаснее, чем потеря event.

361. Почему нельзя сначала поставить `published_at`, а потом publish?

<details style="display: inline;"><summary>Помощь</summary>Что если process упадёт между шагами?</details>  
<details style="display: inline;"><summary>Ответ</summary>Event будет считаться опубликованным, хотя фактически мог не попасть в broker — возникает риск потери.

362. Что лучше: duplicate или loss для Outbox?

<details style="display: inline;"><summary>Помощь</summary>Какой failure consumer проще пережить?</details>  
<details style="display: inline;"><summary>Ответ</summary>Обычно выбирают duplicate с at-least-once и строят idempotent consumer, потому что повторный effect можно контролируемо обезвредить.

363. Как выбирать outbox events для relay?

<details style="display: inline;"><summary>Помощь</summary>Нужна bounded polling.</details>  
<details style="display: inline;"><summary>Ответ</summary>Например, выбирать ограниченными batch-ами непубликованные events по времени/ID, чтобы не загружать память и database.

364. Почему нельзя читать весь outbox в память?

<details style="display: inline;"><summary>Помощь</summary>Backlog может расти.</details>  
<details style="display: inline;"><summary>Ответ</summary>Большой backlog может привести к неограниченному memory usage и ухудшению latency.

365. Что такое outbox relay?

<details style="display: inline;"><summary>Помощь</summary>Кто доставляет event после commit?</details>  
<details style="display: inline;"><summary>Ответ</summary>Отдельный process/component, который читает unpublished events и публикует их в broker.

366. Может ли существовать несколько relay instances?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт с одним event?</details>  
<details style="display: inline;"><summary>Ответ</summary>Да, но нужно явно координировать ownership/claiming и всё равно проектировать downstream для duplicate.

367. Какой риск у нескольких relay без claiming?

<details style="display: inline;"><summary>Помощь</summary>Два relay видят одну row.</details>  
<details style="display: inline;"><summary>Ответ</summary>Они могут оба опубликовать один event, увеличив duplicate delivery.

368. Зачем batch processing в outbox relay?

<details style="display: inline;"><summary>Помощь</summary>Round trips.</details>  
<details style="display: inline;"><summary>Ответ</summary>Уменьшить database overhead и увеличить throughput, сохраняя bounded memory и predictable recovery.

369. Почему Outbox — local transaction pattern, а не distributed transaction?

<details style="display: inline;"><summary>Помощь</summary>Где commit?</details>  
<details style="display: inline;"><summary>Ответ</summary>Transaction охватывает только локальную authoritative database; broker delivery выполняется отдельным процессом.

370. Что проверять у Outbox на интервью?

<details style="display: inline;"><summary>Помощь</summary>С state, failure и duplicate.</details>  
<details style="display: inline;"><summary>Ответ</summary>Atomicity state+intent, delivery semantics, crash windows, duplicate handling, retry policy, cleanup и observability.

---

## RETRY / BACKOFF / FAILURE POLICY

**[К оглавлению](#500-interview-questions)**

371. Что такое retry policy?

<details style="display: inline;"><summary>Помощь</summary>Когда и сколько повторять?</details>  
<details style="display: inline;"><summary>Ответ</summary>Правила, определяющие retryable errors, количество попыток, интервалы, backoff, jitter и окончательное failure behavior.

372. Почему количество retries должно быть bounded?

<details style="display: inline;"><summary>Помощь</summary>Что если downstream не восстановится?</details>  
<details style="display: inline;"><summary>Ответ</summary>Бесконечные retries могут бесконечно удерживать resources и создавать дополнительную нагрузку.

373. Что такое exponential backoff?

<details style="display: inline;"><summary>Помощь</summary>100 → 200 → 400.</details>  
<details style="display: inline;"><summary>Ответ</summary>Стратегия увеличения времени между последовательными retry attempts, обычно экспоненциально.

374. Что такое jitter?

<details style="display: inline;"><summary>Помощь</summary>Случайная вариация delay.</details>  
<details style="display: inline;"><summary>Ответ</summary>Небольшая случайная составляющая retry delay, уменьшающая synchronized retries.

375. Что такое retry storm?

<details style="display: inline;"><summary>Помощь</summary>Много клиентов повторяют одновременно.</details>  
<details style="display: inline;"><summary>Ответ</summary>Массовое накопление повторных попыток, усиливающее нагрузку на уже неисправную систему.

376. Что такое transient error?

<details style="display: inline;"><summary>Помощь</summary>Может исчезнуть без изменения request.</details>  
<details style="display: inline;"><summary>Ответ</summary>Временная ошибка, при которой повтор позже потенциально может завершиться успешно.

377. Что такое permanent error?

<details style="display: inline;"><summary>Помощь</summary>Повтор не изменит причину.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ошибка, которая не исчезает обычным повтором того же operation, например invalid input.

378. Нужно ли retry `validation error`?

<details style="display: inline;"><summary>Помощь</summary>Изменится ли input сам?</details>  
<details style="display: inline;"><summary>Ответ</summary>Обычно нет, потому что повтор того же invalid input не меняет причину failure.

379. Нужно ли retry `timeout`?

<details style="display: inline;"><summary>Помощь</summary>Timeout доказывает operation failure?</details>  
<details style="display: inline;"><summary>Ответ</summary>Не автоматически. Нужно учитывать, является ли operation safely repeatable и что timeout может скрывать уже выполненный effect.

380. Почему timeout требует idempotency?

<details style="display: inline;"><summary>Помощь</summary>Unknown result.</details>  
<details style="display: inline;"><summary>Ответ</summary>После timeout операция могла выполниться; retry может создать duplicate effect без idempotency.

381. Что такое circuit breaker в общих чертах?

<details style="display: inline;"><summary>Помощь</summary>Downstream repeatedly fails.</details>  
<details style="display: inline;"><summary>Ответ</summary>Механизм, временно прекращающий новые calls к деградировавшему downstream после определённого порога ошибок и позволяющий ему восстановиться.

382. Зачем circuit breaker?

<details style="display: inline;"><summary>Помощь</summary>Не продолжать атаковать уже больной dependency.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы ограничить cascade failure и useless load на деградировавший downstream.

383. Чем timeout отличается от retry?

<details style="display: inline;"><summary>Помощь</summary>Wait limit vs повтор.</details>  
<details style="display: inline;"><summary>Ответ</summary>Timeout ограничивает время ожидания одной операции, retry повторяет операцию согласно policy.

384. Почему retry без timeout опасен?

<details style="display: inline;"><summary>Помощь</summary>Каждая попытка тоже может зависнуть.</details>  
<details style="display: inline;"><summary>Ответ</summary>Общее время и количество удерживаемых resources могут стать неограниченными.

385. Что такое retry budget?

<details style="display: inline;"><summary>Помощь</summary>Не бесконечные retries.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ограничение на retry work за определённый период/operation, предотвращающее чрезмерное amplification load.

386. Что делать с message, которая не проходит после max retries?

<details style="display: inline;"><summary>Помощь</summary>Нужна terminal policy.</details>  
<details style="display: inline;"><summary>Ответ</summary>Например, dead-letter queue, failed state или другой explicit recovery/manual intervention path.

387. Что такое dead-letter queue в общем смысле?

<details style="display: inline;"><summary>Помощь</summary>Message больше нельзя безопасно retry-ить автоматически.</details>  
<details style="display: inline;"><summary>Ответ</summary>Отдельное место для сообщений, которые исчерпали retry policy или требуют специального investigation.

388. Почему DLQ не должна быть «кладбищем без наблюдаемости»?

<details style="display: inline;"><summary>Помощь</summary>Кто узнает о failure?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нужно иметь metrics, logging, tracing/diagnostics и процесс повторного анализа или восстановления сообщений.

389. Как выбрать retry delay?

<details style="display: inline;"><summary>Помощь</summary>Ошибка, downstream и SLA.</details>  
<details style="display: inline;"><summary>Ответ</summary>С учётом ожидаемого recovery time, downstream capacity, request deadline и общей latency/reliability policy.

390. Почему retry должен учитывать context deadline?

<details style="display: inline;"><summary>Помощь</summary>Client уже ждёт ограниченное время.</details>  
<details style="display: inline;"><summary>Ответ</summary>Не имеет смысла начинать новую попытку, если общий operation deadline уже истёк или недостаточен для её завершения.

---

## DISTRIBUTED FAILURE ANALYSIS

**[К оглавлению](#500-interview-questions)**

391. Что такое partial failure?

<details style="display: inline;"><summary>Помощь</summary>Одна часть успела, другая нет.</details>  
<details style="display: inline;"><summary>Ответ</summary>Distributed operation частично завершилась, поэтому разные компоненты имеют разные сведения о её результате.

392. Почему distributed timeout создаёт unknown result?

<details style="display: inline;"><summary>Помощь</summary>Что если remote commit уже произошёл?</details>  
<details style="display: inline;"><summary>Ответ</summary>Инициатор потерял возможность узнать окончательный outcome, хотя remote operation могла уже завершиться.

393. Что делать при unknown result?

<details style="display: inline;"><summary>Помощь</summary>Не повторяйте слепо.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать idempotency, query-by-operation-ID или другую recovery semantics для определения/повторения результата безопасным способом.

394. Почему «timeout = failure» является опасным допущением?

<details style="display: inline;"><summary>Помощь</summary>Сеть сообщает состояние ожидания.</details>  
<details style="display: inline;"><summary>Ответ</summary>Timeout сообщает лишь, что ответ не получен вовремя; remote side могла успешно выполнить operation.

395. Что такое failure domain?

<details style="display: inline;"><summary>Помощь</summary>Что может отказать независимо?</details>  
<details style="display: inline;"><summary>Ответ</summary>Группа компонентов/resources, которая может испытывать общий failure относительно другой части системы.

396. Почему PostgreSQL и Kafka имеют разные failure domains?

<details style="display: inline;"><summary>Помощь</summary>У них разные lifecycle и network boundaries.</details>  
<details style="display: inline;"><summary>Ответ</summary>Это независимые системы с отдельным state, availability и commit semantics, поэтому их операции не являются одной локальной transaction.

397. Что происходит при `DB commit ✓, Kafka ✗` без Outbox?

<details style="display: inline;"><summary>Помощь</summary>Job state уже существует.</details>  
<details style="display: inline;"><summary>Ответ</summary>Authoritative state зафиксирован, но asynchronous event может быть потерян, из-за чего worker не узнает о Job.

398. Что происходит при `Kafka ✓, DB commit ✗`?

<details style="display: inline;"><summary>Помощь</summary>Consumer увидел effect, который не стал state.</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumer может получить событие о состоянии, которого фактически нет в authoritative database, что нарушает ожидаемую event/state consistency.

399. Как Outbox меняет failure model?

<details style="display: inline;"><summary>Помощь</summary>Что теперь commit-ится вместе?</details>  
<details style="display: inline;"><summary>Ответ</summary>DB state и event intent становятся одной local transaction; если broker недоступен, event остаётся durable в outbox до успешной доставки.

400. Какая главная инженерная мысль распределённого урока?

<details style="display: inline;"><summary>Помощь</summary>Не пытайтесь отменить сетевую реальность.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нужно сделать delivery semantics, retries, idempotent effects, timeouts и failure recovery явными и bounded, чтобы partial failure оставался управляемым.

## HTTP / API

**[К оглавлению](#500-interview-questions)**

401. Что такое HTTP handler в Go backend?

<details style="display: inline;"><summary>Помощь</summary>Что он получает и что должен вернуть?</details>  
<details style="display: inline;"><summary>Ответ</summary>Handler — transport adapter, который принимает HTTP request, преобразует его в application input, вызывает use case и формирует HTTP response.

402. Что не должен делать HTTP handler?

<details style="display: inline;"><summary>Помощь</summary>Вспомните separation of concerns.</details>  
<details style="display: inline;"><summary>Ответ</summary>Не должен содержать business policy, прямой SQL, Kafka/Redis infrastructure logic и другую логику, не относящуюся непосредственно к HTTP transport.

403. Почему HTTP transport должен быть отдельной boundary?

<details style="display: inline;"><summary>Помощь</summary>Что если HTTP заменится gRPC?</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы изменение внешнего protocol не заставляло переписывать application/domain policy.

404. Что такое request DTO?

<details style="display: inline;"><summary>Помощь</summary>Transport representation.</details>  
<details style="display: inline;"><summary>Ответ</summary>Структура, представляющая входные данные конкретного transport protocol и отделяющая их от domain model.

405. Почему не стоит напрямую использовать HTTP request struct как domain model?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт при смене transport?</details>  
<details style="display: inline;"><summary>Ответ</summary>Domain станет зависеть от transport-specific деталей, увеличивая coupling и затрудняя замену protocol.

406. Как должен проходить `POST /jobs` через систему?

<details style="display: inline;"><summary>Помощь</summary>Transport → policy → state/effect.</details>  
<details style="display: inline;"><summary>Ответ</summary>HTTP decode → transport validation → application command → business decision → state change/transaction → response.

407. Что такое HTTP status mapping?

<details style="display: inline;"><summary>Помощь</summary>Application error → HTTP code.</details>  
<details style="display: inline;"><summary>Ответ</summary>Преобразование application/domain outcomes в HTTP status codes и response representation на transport boundary.

408. Должен ли `ErrJobNotFound` знать о HTTP 404?

<details style="display: inline;"><summary>Помощь</summary>Кто владеет HTTP?</details>  
<details style="display: inline;"><summary>Ответ</summary>Нет. Error принадлежит application/domain semantics, а HTTP adapter переводит её в `404`.

409. Какой status code естественен для resource not found?

<details style="display: inline;"><summary>Помощь</summary>Клиент запросил отсутствующий resource.</details>  
<details style="display: inline;"><summary>Ответ</summary>`404 Not Found`.

410. Какой status code часто используется для validation failure?

<details style="display: inline;"><summary>Помощь</summary>Некорректный input.</details>  
<details style="display: inline;"><summary>Ответ</summary>Обычно `400 Bad Request`, если запрос не соответствует ожидаемому contract.

411. Когда `409 Conflict` уместен?

<details style="display: inline;"><summary>Помощь</summary>Request syntactically valid, но state конфликтует.</details>  
<details style="display: inline;"><summary>Ответ</summary>Когда request невозможно выполнить из-за конфликта с текущим resource state, например недопустимого state transition.

412. Когда использовать `500 Internal Server Error`?

<details style="display: inline;"><summary>Помощь</summary>Ошибка, которую клиент не может исправить request-ом.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для unexpected server-side failure, когда подходящий более конкретный status code не установлен.

413. Что такое API contract?

<details style="display: inline;"><summary>Помощь</summary>Что потребитель ожидает от endpoint?</details>  
<details style="display: inline;"><summary>Ответ</summary>Набор определённых request/response semantics, status codes, fields, validation и error behavior.

414. Почему API contract важнее внутренней реализации?

<details style="display: inline;"><summary>Помощь</summary>Кто зависит от endpoint?</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumers зависят от observable behavior API, а внутренняя реализация должна оставаться заменяемой.

415. Что такое idempotent HTTP operation?

<details style="display: inline;"><summary>Помощь</summary>Что происходит при повторном request?</details>  
<details style="display: inline;"><summary>Ответ</summary>Повтор одного логического request не создаёт дополнительного нежелательного effect.

416. Почему `GET` должен быть safe?

<details style="display: inline;"><summary>Помощь</summary>Что должен делать read endpoint?</details>  
<details style="display: inline;"><summary>Ответ</summary>GET по HTTP semantics предназначен для retrieval и не должен намеренно создавать изменения состояния.

417. Почему `GET /jobs/{id}` не должен создавать Job, если её нет?

<details style="display: inline;"><summary>Помощь</summary>Read vs command.</details>  
<details style="display: inline;"><summary>Ответ</summary>Это нарушило бы expected semantics read operation и сделало бы side effect скрытым.

418. Что такое health check?

<details style="display: inline;"><summary>Помощь</summary>Процесс жив?</details>  
<details style="display: inline;"><summary>Ответ</summary>Endpoint или механизм, показывающий состояние процесса согласно определённой health policy.

419. Что такое readiness?

<details style="display: inline;"><summary>Помощь</summary>Можно ли давать процессу production traffic?</details>  
<details style="display: inline;"><summary>Ответ</summary>Сигнал, что application готова принимать работу согласно требуемым dependencies и runtime state.

420. Чем liveness отличается от readiness?

<details style="display: inline;"><summary>Помощь</summary>Alive vs ready.</details>  
<details style="display: inline;"><summary>Ответ</summary>Liveness показывает, что process жив и не находится в критически зависшем состоянии; readiness — что instance может обслуживать traffic.

---

## CONTEXT / HTTP LIFECYCLE

**[К оглавлению](#500-interview-questions)**

421. Зачем передавать `context.Context` от HTTP handler вниз?

<details style="display: inline;"><summary>Помощь</summary>Request lifecycle.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы cancellation и deadline исходного request распространялись на application и downstream operations.

422. Что произойдёт, если service создаст `context.Background()` вместо request context?

<details style="display: inline;"><summary>Помощь</summary>Сигнал cancellation потеряется.</details>  
<details style="display: inline;"><summary>Ответ</summary>Downstream operation перестанет получать cancellation/deadline исходного request и может продолжить ненужную работу.

423. Для чего нужен context deadline?

<details style="display: inline;"><summary>Помощь</summary>Bounded waiting.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ограничить максимальное время выполнения operation и освободить resources после истечения срока.

424. Чем timeout отличается от cancellation?

<details style="display: inline;"><summary>Помощь</summary>Причина остановки.</details>  
<details style="display: inline;"><summary>Ответ</summary>Cancellation — общий сигнал прекращения работы, а timeout/deadline — cancellation, инициированная превышением допустимого времени.

425. Почему context не должен хранить произвольное application state?

<details style="display: inline;"><summary>Помощь</summary>Для чего он предназначен?</details>  
<details style="display: inline;"><summary>Ответ</summary>Context предназначен для cancellation, deadlines и request-scoped metadata, а не для передачи обязательных бизнес-зависимостей или domain state.

426. Что делать, если downstream operation не поддерживает cancellation?

<details style="display: inline;"><summary>Помощь</summary>Context signal должен доходить до actual operation.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нужно использовать API/driver с поддержкой context или ограничить внешний call другим механизмоm; нельзя считать cancellation гарантированной, если dependency её игнорирует.

427. Почему context важен при graceful shutdown?

<details style="display: inline;"><summary>Помощь</summary>Как сообщить всем workers stop?</details>  
<details style="display: inline;"><summary>Ответ</summary>Он позволяет распространить cancellation и завершить длительные operations согласованным способом.

428. Что такое request-scoped metadata?

<details style="display: inline;"><summary>Помощь</summary>Correlation/request ID.</details>  
<details style="display: inline;"><summary>Ответ</summary>Данные, относящиеся к lifecycle одного request, например request ID, которые должны распространяться по call chain.

429. Что нельзя делать с context cancellation?

<details style="display: inline;"><summary>Помощь</summary>Signal ≠ kill.</details>  
<details style="display: inline;"><summary>Ответ</summary>Нельзя считать, что cancellation принудительно уничтожает goroutine; код должен кооперативно наблюдать signal и завершаться.

430. Что должно произойти с DB query после HTTP disconnect?

<details style="display: inline;"><summary>Помощь</summary>Нужна ли ещё работа?</details>  
<details style="display: inline;"><summary>Ответ</summary>При переданном request context query должен получить cancellation/deadline и прекратить ненужное выполнение.

---

## APPLICATION ARCHITECTURE

**[К оглавлению](#500-interview-questions)**

431. Что такое application layer?

<details style="display: inline;"><summary>Помощь</summary>Use cases.</details>  
<details style="display: inline;"><summary>Ответ</summary>Слой, координирующий application use cases: validation, domain decisions, state changes и вызовы необходимых ports.

432. Что такое domain layer?

<details style="display: inline;"><summary>Помощь</summary>Business semantics.</details>  
<details style="display: inline;"><summary>Ответ</summary>Слой, содержащий domain state, business invariants и behavior, не зависящий от конкретного transport или infrastructure.

433. Что такое infrastructure layer?

<details style="display: inline;"><summary>Помощь</summary>Concrete dependencies.</details>  
<details style="display: inline;"><summary>Ответ</summary>Реализации работы с PostgreSQL, Kafka, Redis, HTTP clients и другими внешними технологиями.

434. Что такое adapter?

<details style="display: inline;"><summary>Помощь</summary>Boundary translation.</details>  
<details style="display: inline;"><summary>Ответ</summary>Компонент, преобразующий внешний/технологический API в contract, ожидаемый application layer, или обратно.

435. Что такое port?

<details style="display: inline;"><summary>Помощь</summary>Stable contract.</details>  
<details style="display: inline;"><summary>Ответ</summary>Application-facing contract, описывающий необходимое behavior без зависимости от конкретного adapter.

436. Чем port отличается от adapter?

<details style="display: inline;"><summary>Помощь</summary>Contract vs implementation.</details>  
<details style="display: inline;"><summary>Ответ</summary>Port задаёт контракт, adapter реализует или обслуживает этот контракт через конкретную technology boundary.

437. Что такое dependency direction?

<details style="display: inline;"><summary>Помощь</summary>Кто знает о деталях?</details>  
<details style="display: inline;"><summary>Ответ</summary>Направление зависимостей между слоями; application policy должна зависеть от stable contracts, а infrastructure detail — подключаться снаружи.

438. Почему dependency direction важнее структуры каталогов?

<details style="display: inline;"><summary>Помощь</summary>Можно иметь `internal/domain`, но wrong imports.</details>  
<details style="display: inline;"><summary>Ответ</summary>Реальная архитектура определяется отношениями между компонентами и их dependencies, а не именами директорий.

439. Что такое hexagonal architecture?

<details style="display: inline;"><summary>Помощь</summary>Ports and adapters.</details>  
<details style="display: inline;"><summary>Ответ</summary>Архитектурный стиль, в котором core application отделён от внешнего мира через ports, а concrete systems подключаются adapters.

440. Зачем hexagonal architecture в JobFlow?

<details style="display: inline;"><summary>Помощь</summary>PostgreSQL/Kafka/HTTP меняются независимо.</details>  
<details style="display: inline;"><summary>Ответ</summary>Она позволяет локализовать transport, storage и messaging details и сохранить application policy относительно независимой от технологий.

441. Что такое separation of concerns?

<details style="display: inline;"><summary>Помощь</summary>Разные причины изменения.</details>  
<details style="display: inline;"><summary>Ответ</summary>Разделение разных видов ответственности так, чтобы изменения в одной concern минимально затрагивали остальные.

442. Что такое cohesion?

<details style="display: inline;"><summary>Помощь</summary>Что находится вместе?</details>  
<details style="display: inline;"><summary>Ответ</summary>Степень того, насколько элементы внутри компонента связаны одной общей ответственностью или целью.

443. Что такое coupling?

<details style="display: inline;"><summary>Помощь</summary>Как сильно компоненты зависят друг от друга?</details>  
<details style="display: inline;"><summary>Ответ</summary>Степень зависимости компонентов друг от друга и влияния изменения одного на остальные.

444. Какой design target обычно желателен?

<details style="display: inline;"><summary>Помощь</summary>Что делать с cohesion/coupling?</details>  
<details style="display: inline;"><summary>Ответ</summary>Высокая cohesion внутри компонентов и минимально необходимый coupling между ними.

445. Что такое blast radius изменения?

<details style="display: inline;"><summary>Помощь</summary>Сколько компонентов затронуто?</details>  
<details style="display: inline;"><summary>Ответ</summary>Объём системы, который потенциально нужно изменить, проверить или затронуть при одном изменении.

446. Почему abstraction может увеличивать blast radius?

<details style="display: inline;"><summary>Помощь</summary>Слишком широкий interface.</details>  
<details style="display: inline;"><summary>Ответ</summary>Широкая abstraction surface может связать больше consumers и implementations, поэтому изменение contract затрагивает больше системы.

447. Что такое premature abstraction?

<details style="display: inline;"><summary>Помощь</summary>Abstraction до требования.</details>  
<details style="display: inline;"><summary>Ответ</summary>Создание abstraction до появления реальной boundary, variability или change case, из-за чего complexity растёт без достаточной выгоды.

448. Как понять, нужна ли abstraction?

<details style="display: inline;"><summary>Помощь</summary>Вернитесь к Engineering Mindset.</details>  
<details style="display: inline;"><summary>Ответ</summary>Определить requirement, invariant, вероятную ось изменения и dependency boundary; abstraction вводится только если она реально снижает стоимость изменения или изолирует dependency.

449. Что такое KISS?

<details style="display: inline;"><summary>Помощь</summary>Минимальная достаточная конструкция.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать минимальное решение, достаточное для требуемого behavior и гарантий.

450. Что такое YAGNI?

<details style="display: inline;"><summary>Помощь</summary>Requirement отсутствует.</details>  
<details style="display: inline;"><summary>Ответ</summary>Не создавать capability, пока нет реальной потребности и requirement.

---

## DESIGN / SOLID

**[К оглавлению](#500-interview-questions)**

451. Что такое Single Responsibility Principle?

<details style="display: inline;"><summary>Помощь</summary>Одна причина изменения.</details>  
<details style="display: inline;"><summary>Ответ</summary>Компонент должен иметь одну основную ось ответственности и изменения.

452. Что SRP не означает?

<details style="display: inline;"><summary>Помощь</summary>Это не «одна функция на файл».</details>  
<details style="display: inline;"><summary>Ответ</summary>SRP не требует ровно одной функции, struct или строки; он относится к причинам изменения и ответственности.

453. Что такое Open/Closed Principle?

<details style="display: inline;"><summary>Помощь</summary>Stable core + new behavior.</details>  
<details style="display: inline;"><summary>Ответ</summary>Стабильную policy желательно расширять новым behavior с минимальным переписыванием уже стабильного ядра в тех местах, где boundary действительно оправдана.

454. Что такое Liskov Substitution Principle?

<details style="display: inline;"><summary>Помощь</summary>Замена реализации.</details>  
<details style="display: inline;"><summary>Ответ</summary>Любая реализация, удовлетворяющая contract, должна быть заменяема другой без нарушения ожиданий consumer.

455. Что такое Interface Segregation Principle?

<details style="display: inline;"><summary>Помощь</summary>Не заставлять зависеть от лишнего.</details>  
<details style="display: inline;"><summary>Ответ</summary>Consumers не должны зависеть от методов и behavior, которые им не нужны.

456. Что такое Dependency Inversion Principle?

<details style="display: inline;"><summary>Помощь</summary>Policy vs detail.</details>  
<details style="display: inline;"><summary>Ответ</summary>Высокоуровневая policy не должна напрямую зависеть от низкоуровневых details; обе стороны должны быть связаны через stable abstraction.

457. Как реализовать DIP в Go?

<details style="display: inline;"><summary>Помощь</summary>Consumer interface + constructor.</details>  
<details style="display: inline;"><summary>Ответ</summary>Определить минимальный interface на consumer side и передать concrete implementation через constructor.

458. Почему Go особенно хорошо подходит для ISP?

<details style="display: inline;"><summary>Помощь</summary>Implicit interfaces.</details>  
<details style="display: inline;"><summary>Ответ</summary>Implicit implementation и лёгкость создания маленьких interfaces делают consumer-oriented contracts естественными.

459. Почему SOLID нельзя применять как checklist?

<details style="display: inline;"><summary>Помощь</summary>Сначала requirement, потом principle.</details>  
<details style="display: inline;"><summary>Ответ</summary>Принципы являются инструментами design decisions; применение без реальной проблемы может добавить complexity и ухудшить систему.

460. Как объяснить SOLID на интервью кратко?

<details style="display: inline;"><summary>Помощь</summary>Одна фраза + пять букв.</details>  
<details style="display: inline;"><summary>Ответ</summary>SOLID — пять принципов управления ответственностями, расширением, контрактами и зависимостями: SRP, OCP, LSP, ISP, DIP.

---

## REDIS / CACHE

**[К оглавлению](#500-interview-questions)**

461. Зачем Redis в JobFlow?

<details style="display: inline;"><summary>Помощь</summary>PostgreSQL остаётся source of truth.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для снижения latency и нагрузки на PostgreSQL при горячих reads через derived cache state.

462. Что такое cache-aside?

<details style="display: inline;"><summary>Помощь</summary>Кто читает cache первым?</details>  
<details style="display: inline;"><summary>Ответ</summary>Application сначала читает cache; при miss читает source of truth и заполняет cache результатом.

463. Что такое cache hit?

<details style="display: inline;"><summary>Помощь</summary>Данные уже в cache.</details>  
<details style="display: inline;"><summary>Ответ</summary>Запрошенное значение успешно найдено в cache и может быть возвращено без обращения к source of truth.

464. Что такое cache miss?

<details style="display: inline;"><summary>Помощь</summary>Нет cached value.</details>  
<details style="display: inline;"><summary>Ответ</summary>Запрошенное значение отсутствует или недоступно в cache, поэтому требуется fallback к authoritative source.

465. Что такое TTL?

<details style="display: inline;"><summary>Помощь</summary>Время жизни entry.</details>  
<details style="display: inline;"><summary>Ответ</summary>Time-to-live — ограничение времени, после которого cache entry считается истёкшим.

466. Зачем TTL?

<details style="display: inline;"><summary>Помощь</summary>Stale data.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы ограничивать срок жизни stale cached data и автоматически обновлять derived state.

467. Что такое stale cache?

<details style="display: inline;"><summary>Помощь</summary>Source of truth уже изменился.</details>  
<details style="display: inline;"><summary>Ответ</summary>Cache содержит значение, которое больше не соответствует актуальному authoritative state.

468. Что делать при Redis outage в cache-aside architecture?

<details style="display: inline;"><summary>Помощь</summary>Redis не source of truth.</details>  
<details style="display: inline;"><summary>Ответ</summary>По возможности деградировать к PostgreSQL, сохранив correctness и ограничивая дополнительную нагрузку.

469. Почему нельзя бесконечно fallback-ить на PostgreSQL при Redis outage?

<details style="display: inline;"><summary>Помощь</summary>Что будет с DB?</details>  
<details style="display: inline;"><summary>Ответ</summary>Массовый cache miss может создать DB overload; необходимо иметь resource limits, rate control и подходящую degradation policy.

470. Что такое cache stampede?

<details style="display: inline;"><summary>Помощь</summary>Много misses одновременно.</details>  
<details style="display: inline;"><summary>Ответ</summary>Ситуация, когда большое количество requests одновременно обращается к source of truth из-за отсутствия/истечения одного cache value.

471. Как уменьшить cache stampede?

<details style="display: inline;"><summary>Помощь</summary>Скоординируйте refresh.</details>  
<details style="display: inline;"><summary>Ответ</summary>Использовать request coalescing, distributed lock, jittered expiry, prewarming или другие bounded caching strategies.

472. Когда cache лучше обновлять, а когда удалять?

<details style="display: inline;"><summary>Помощь</summary>Consistency trade-off.</details>  
<details style="display: inline;"><summary>Ответ</summary>Зависит от access pattern; для простой cache-aside часто проще invalidate entry после authoritative update и восстановить её на следующем read.

473. Что такое cache invalidation?

<details style="display: inline;"><summary>Помощь</summary>Old value больше не должен использоваться.</details>  
<details style="display: inline;"><summary>Ответ</summary>Удаление или маркировка cached value как непригодного после изменения source of truth.

474. Может ли Redis хранить authoritative state?

<details style="display: inline;"><summary>Помощь</summary>Технически может, но вопрос в contract.</details>  
<details style="display: inline;"><summary>Ответ</summary>Да, если architecture специально задаёт Redis как durable/authoritative store; сам факт использования Redis не делает его cache автоматически.

475. Как измерять эффективность cache?

<details style="display: inline;"><summary>Помощь</summary>Что важно кроме latency?</details>  
<details style="display: inline;"><summary>Ответ</summary>Cache hit rate, miss rate, latency, DB load, eviction/TTL behavior и error rate.

---

## PRODUCTION / OBSERVABILITY

**[К оглавлению](#500-interview-questions)**

476. Что такое observability?

<details style="display: inline;"><summary>Помощь</summary>Как узнать внутреннее состояние по external signals?</details>  
<details style="display: inline;"><summary>Ответ</summary>Способность понимать состояние и поведение системы по её внешне наблюдаемым signals, прежде всего logs, metrics и traces.

477. Для чего нужны logs?

<details style="display: inline;"><summary>Помощь</summary>Discrete events.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для записи контекста и событий выполнения, необходимых для диагностики конкретных операций и failures.

478. Для чего нужны metrics?

<details style="display: inline;"><summary>Помощь</summary>Aggregate behavior over time.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для количественного наблюдения за throughput, latency, errors, resource usage и состоянием системы во времени.

479. Для чего нужны traces?

<details style="display: inline;"><summary>Помощь</summary>Одна операция через несколько components.</details>  
<details style="display: inline;"><summary>Ответ</summary>Для анализа distributed execution path одного request/operation через несколько сервисов и boundaries.

480. Зачем correlation/request ID?

<details style="display: inline;"><summary>Помощь</summary>Один request, много logs.</details>  
<details style="display: inline;"><summary>Ответ</summary>Чтобы связывать logs и события одного logical request между несколькими components.

481. Что логировать в backend?

<details style="display: inline;"><summary>Помощь</summary>Не всё подряд.</details>  
<details style="display: inline;"><summary>Ответ</summary>Значимые lifecycle/failure events с достаточным context: operation identity, outcome, latency, dependency и correlation information, без утечки sensitive data.

482. Что такое structured logging?

<details style="display: inline;"><summary>Помощь</summary>Fields вместо длинной строки.</details>  
<details style="display: inline;"><summary>Ответ</summary>Логи с машинно читаемыми структурированными fields, а не только свободным текстом.

483. Какие metrics полезны для JobFlow?

<details style="display: inline;"><summary>Помощь</summary>Состояние jobs и pipeline.</details>  
<details style="display: inline;"><summary>Ответ</summary>Jobs created/processed/failed, processing latency, worker utilization, queue/backlog, Kafka lag, DB pool saturation и cache hit rate.

484. Почему latency нужно измерять отдельно для dependencies?

<details style="display: inline;"><summary>Помощь</summary>Где bottleneck?</details>  
<details style="display: inline;"><summary>Ответ</summary>Общая latency может быть высокой из-за конкретной dependency; per-dependency measurements помогают локализовать bottleneck.

485. Что такое SLI?

<details style="display: inline;"><summary>Помощь</summary>Измеряемый reliability signal.</details>  
<details style="display: inline;"><summary>Ответ</summary>Service Level Indicator — измеряемая характеристика service behavior, например availability или request latency.

486. Что такое SLO?

<details style="display: inline;"><summary>Помощь</summary>Target для SLI.</details>  
<details style="display: inline;"><summary>Ответ</summary>Service Level Objective — целевое значение или диапазон для определённого SLI.

487. Что такое graceful degradation?

<details style="display: inline;"><summary>Помощь</summary>Не всё обязательно должно падать.</details>  
<details style="display: inline;"><summary>Ответ</summary>Способность сохранить core functionality при отказе необязательной dependency ценой ограниченного качества или дополнительной latency.

488. Как JobFlow должен деградировать при Redis failure?

<details style="display: inline;"><summary>Помощь</summary>Redis — derived state.</details>  
<details style="display: inline;"><summary>Ответ</summary>GET операции должны по возможности fallback к PostgreSQL, сохраняя correctness.

489. Должна ли система скрывать Kafka outage от пользователя любой ценой?

<details style="display: inline;"><summary>Помощь</summary>Состояние уже persisted?</details>  
<details style="display: inline;"><summary>Ответ</summary>Если Job и outbox intent успешно committed, можно принять request и асинхронно доставить event позже; если authoritative state не удалось зафиксировать, нельзя выдавать ложный success.

490. Что такое failure mode analysis?

<details style="display: inline;"><summary>Помощь</summary>Что сломается и что увидит система?</details>  
<details style="display: inline;"><summary>Ответ</summary>Систематический анализ возможных отказов, их эффектов, observable signals и recovery policy.

---

## SYSTEM DESIGN / CROSS-SYSTEM

**[К оглавлению](#500-interview-questions)**

491. Спроектируйте путь `POST /jobs` в JobFlow.

<details style="display: inline;"><summary>Помощь</summary>Пройдите transport → state → effect.</details>  
<details style="display: inline;"><summary>Ответ</summary>HTTP handler декодирует request → application validates/decides → PostgreSQL transaction сохраняет Job + Outbox event → commit → response; relay позже публикует event в Kafka.

492. Что является source of truth в JobFlow?

<details style="display: inline;"><summary>Помощь</summary>Где durable state?</details>  
<details style="display: inline;"><summary>Ответ</summary>PostgreSQL для authoritative Job state; Redis является derived cache, Kafka — transport/effect boundary.

493. Где находится asynchronous boundary?

<details style="display: inline;"><summary>Помощь</summary>State commit и worker execution.</details>  
<details style="display: inline;"><summary>Ответ</summary>Между persistent outbox/event и Kafka consumer/worker processing.

494. Где возникает duplicate в JobFlow?

<details style="display: inline;"><summary>Помощь</summary>Crash windows.</details>  
<details style="display: inline;"><summary>Ответ</summary>Например, после Kafka publish до отметки outbox published или после database effect до offset commit.

495. Где обеспечивается idempotency?

<details style="display: inline;"><summary>Помощь</summary>Посмотрите на consumer boundary.</details>  
<details style="display: inline;"><summary>Ответ</summary>В consumer processing через deduplication state и/или atomic conditional state transitions в PostgreSQL.

496. Что произойдёт, если PostgreSQL недоступен при создании Job?

<details style="display: inline;"><summary>Помощь</summary>Authoritative state не committed.</details>  
<details style="display: inline;"><summary>Ответ</summary>Request не должен получить ложный успешный результат; operation завершается контролируемой ошибкой согласно retry/HTTP policy.

497. Что произойдёт, если Kafka недоступна после создания Job?

<details style="display: inline;"><summary>Помощь</summary>Есть ли Outbox?</details>  
<details style="display: inline;"><summary>Ответ</summary>Job и outbox event остаются committed; relay повторяет publication после восстановления Kafka.

498. Что произойдёт, если Redis недоступен при `GET /jobs/{id}`?

<details style="display: inline;"><summary>Помощь</summary>Что является authoritative source?</details>  
<details style="display: inline;"><summary>Ответ</summary>Application может обратиться непосредственно к PostgreSQL и вернуть корректное состояние, если cache предусмотрен как optional derived layer.

499. Как объяснить архитектуру JobFlow на senior-интервью за одну минуту?

<details style="display: inline;"><summary>Помощь</summary>Не перечисляйте технологии — объясните границы.</details>  
<details style="display: inline;"><summary>Ответ</summary>«HTTP — transport boundary; application/domain содержит policy и invariants; PostgreSQL — authoritative state; outbox атомарно фиксирует message intent; Kafka — asynchronous delivery boundary с at-least-once semantics; consumer обеспечивает idempotent effect; Redis — derived cache; worker pool ограничивает concurrency; context и graceful shutdown управляют lifecycle».

500. Какой главный вопрос нужно задавать себе на любом backend system-design интервью?

<details style="display: inline;"><summary>Помощь</summary>Вернитесь к первому плакату курса.</details>  
<details style="display: inline;"><summary>Ответ</summary>«Какое требование мы должны выполнить, какой invariant сохранить, где проходит boundary, какой минимальный механизм это гарантирует, что произойдёт при отказе и чем мы докажем корректность?»

---

**[К оглавлению](#500-interview-questions)**
