# Программа: Go Backend Combat

24 недели / 6 месяцев. На входе — опыт программирования на другом языке; Go изучается с нуля. Занятие: 50 минут + 10 минут резерва. ДЗ: 8–12 часов с подготовкой; ★ добровольна.

Сквозной результат — JobFlow: учебный движок DAG-задач из pub, sub, act, ret. PostgreSQL хранит истину, Kafka доставляет события, Redis ускоряет чтение и передаёт сигналы. Исполняются зарегистрированные кооперативные обработчики.

[Как учиться и сдавать](../docs/WORKFLOW.md) · [Контракт проекта](../docs/PROJECT.md) · [Формат 128 байт](../docs/FRAMING.md)

## Месяц 1. Старт и язык Go

| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |
|---|---|---|---|
| [01. Инженерная задача и окружение](../student/01_ENVIRONMENT.md) | Go toolchain, терминал, Git, CI | Другой разработчик воспроизводит запуск из README. | [Скопировать starter в собственный репозиторий с module jobflow.](../combat/01_ENVIRONMENT.md) |
| [02. Типы, функции и модель задачи](../student/02_GO_MODEL.md) | Типы, функции, []byte, go.mod | ID, тип работы и исходные байты тела сохраняют свой смысл. | [Создать internal/task и поля ID, Kind, Body.](../combat/02_GO_MODEL.md) |
| [03. Коллекции, указатели и CRUD](../student/03_MEMORY_CRUD.md) | map, slice, методы, указатели | Дубликат не перезаписывает запись; чтение не отдаёт изменяемое внутреннее состояние. | [Реализовать memory.New и Create/Get/List/Update/Delete с context.Context.](../combat/03_MEMORY_CRUD.md) |
| [04. Ошибки, валидация и тесты](../student/04_ERRORS_TESTS.md) | error, errors.Is, testing | Невалидный ввод отклоняется до изменения хранилища. | [Определить категории ошибок и метод Validate.](../combat/04_ERRORS_TESTS.md) |

## Месяц 2. Границы приложения и данные

| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |
|---|---|---|---|
| [05. Пакеты и границы зависимостей](../student/05_BOUNDARIES.md) | Пакеты, интерфейсы потребителя | Прикладной сервис зависит от поведения, требуемого от хранилища. | [Выделить app.Service и интерфейс Repository.](../combat/05_BOUNDARIES.md) |
| [06. HTTP как граница приложения](../student/06_HTTP.md) | net/http, encoding/json, httptest | Транспорт переводит представления и ошибки; бизнес-правило остаётся в приложении. | [Реализовать api.New(repo) и маршруты черновиков из контракта.](../combat/06_HTTP.md) |
| [07. SQL с нуля и первая миграция](../student/07_SQL.md) | psql, SQL, миграции | Схема создаётся из файла; БД сама защищает уникальность ID. | [Создать migrations/001_core.sql с drafts и draft_intents по контракту.](../combat/07_SQL.md) |
| [08. Go и PostgreSQL-адаптер](../student/08_POSTGRES.md) | pgxpool, context, SQL-адаптер | Данные читаются после открытия независимого подключения. | [Реализовать postgres.Open/Close/Migrate и Repository.](../combat/08_POSTGRES.md) |

## Месяц 3. Транзакции, конкурентность и заголовок

| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |
|---|---|---|---|
| [09. Транзакция и намерение события](../student/09_TRANSACTIONS.md) | Транзакция, rollback, local intent | Черновик и его intent сохраняются вместе либо не сохраняется ни один. | [Добавить CreateWithIntent.](../combat/09_TRANSACTIONS.md) |
| [10. Горутины и ограничение параллелизма](../student/10_WORKER_POOL.md) | goroutine, channel, sync.WaitGroup | Активных обработчиков не больше заданного workers. | [Реализовать pool.Run(ctx, workers, jobs).](../combat/10_WORKER_POOL.md) |
| [11. Отмена, backpressure и shutdown](../student/11_CANCELLATION.md) | context, select, race detector | Отмена освобождает ожидающих работников, а Run ждёт их завершения. | [Добавить отменяемое ожидание работы и очереди.](../combat/11_CANCELLATION.md) |
| [12. 128 байт: фиксированный заголовок](../student/12_BINARY_HEADER.md) | encoding/binary, LittleEndian | Проверки границ и размеров предшествуют доступу к срезам. | [Реализовать Header и ParseHeader/MarshalBinary.](../combat/12_BINARY_HEADER.md) |

## Месяц 4. Бинарный поток и доставка

| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |
|---|---|---|---|
| [13. Поток, связи, тело и padding](../student/13_FRAMING.md) | io.Reader, BLAKE3, fuzz | Разбор сохраняет исходное тело и не поглощает следующий фрейм. | [Реализовать Marshal/Unmarshal/Read/Skip/Prefix.](../combat/13_FRAMING.md) |
| [14. DAG и принятие workflow](../student/14_DAG.md) | DFS, UUIDv7, SQL-транзакция | Принятый граф не содержит циклов и не меняет связи после отправки. | [Реализовать ValidateDAG и разрешение ref/dep.](../combat/14_DAG.md) |
| [15. CAS, BLAKE3 и Redis-кэш](../student/15_CAS_REDIS.md) | Redis, CAS, восстановление кэша | Потеря Redis не уничтожает подтверждённый body. | [Сохранять payloads по BLAKE3.](../combat/15_CAS_REDIS.md) |
| [16. Kafka и доставка Outbox](../student/16_KAFKA_OUTBOX.md) | Kafka API, Redpanda, Outbox | Недоставленное событие остаётся в Outbox до подтверждения публикации. | [Добавить relay в pub.](../combat/16_KAFKA_OUTBOX.md) |

## Месяц 5. Проекции, сигналы и исполнение

| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |
|---|---|---|---|
| [17. sub и восстанавливаемые проекции](../student/17_PROJECTIONS.md) | sub, Redis Hash, version, replay | Проекция принимает только более новую версию задачи. | [Добавить процесс sub и проекцию по Task ID.](../combat/17_PROJECTIONS.md) |
| [18. Redis Streams и потерянный потребитель](../student/18_STREAMS.md) | Redis Streams, Consumer Groups, XAUTOCLAIM | Получение не означает завершение; pending-запись может получить другой потребитель. | [Добавить ready Stream и consumer group act.](../combat/18_STREAMS.md) |
| [19. act и право на попытку](../student/19_ACT.md) | act, claim, lease, token, registry | Claim атомарен; чужой token не позволяет записать результат. | [Добавить act с ограниченным числом потребителей.](../combat/19_ACT.md) |
| [20. Время и планирование](../student/20_SCHEDULING.md) | nbf, exp, maintenance, clock | Запуск невозможен до nbf или без зависимостей; на exp результат отклоняется. | [Добавить nbf/exp и available_at.](../combat/20_SCHEDULING.md) |

## Месяц 6. Восстановление и эксплуатация

| № / тема | Новый инструмент | Инвариант / граница | Результат обязательного ДЗ |
|---|---|---|---|
| [21. Повторы и идемпотентный эффект](../student/21_RECOVERY.md) | Retry, receipt, effect ledger, quota | Повтор логической операции не увеличивает счётчик второй раз. | [Добавить bounded retry и backoff.](../combat/21_RECOVERY.md) |
| [22. ret и безопасная очистка](../student/22_RETENTION.md) | ret, retention, ссылки CAS | GC сохраняет живые ссылки, недоставленные события и идентичность задачи. | [Добавить ret и команду GC у pub.](../combat/22_RETENTION.md) |
| [23. Сборка и наблюдаемость](../student/23_OPERATIONS.md) | Docker Compose, slog, metrics, E2E | Чистый запуск приводит к проверяемому результату; логи связаны по workflow/task/attempt. | [Собрать команды pub/sub/act/ret в один образ.](../combat/23_OPERATIONS.md) |
| [24. Защита проекта](../student/24_DEFENSE.md) | Демонстрация, отказ, аргументация гарантий | Обещания README подтверждены командами и сценарием отказа. | [Описать запуск, архитектуру, гарантии и ограничения в README.](../combat/24_DEFENSE.md) |

## Синхронизация материалов

Для каждого этапа есть конспект, компактный план ментора, ДЗ с ★, самостоятельный эталон и открытые накопительные проверки. Исходные данные программы: [lessons.json](lessons.json). `make docs` пересобирает связанные документы; `make validate` проверяет ссылки и структуру. Исторические материалы сохранены в [legacy](../../../legacy/README.md).
