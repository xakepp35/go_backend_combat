# JobFlow: эталон 23 — Сборка и наблюдаемость

Открытая самостоятельная версия **после** ДЗ 23. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Чистый запуск приводит к проверяемому результату; логи связаны по workflow/task/attempt.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=23 PROJECT=reference/23`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/23_OPERATIONS.md) · [ДЗ](../../combat/23_OPERATIONS.md) · [Открытые проверки](../../checks/23_test.go)
