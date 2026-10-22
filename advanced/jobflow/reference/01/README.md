# JobFlow: эталон 01 — Инженерная задача и окружение

Открытая самостоятельная версия **после** ДЗ 01. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Другой разработчик воспроизводит запуск из README.

Из этой директории: `go run ./cmd/hello`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=01 PROJECT=reference/01`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/01_ENVIRONMENT.md) · [ДЗ](../../combat/01_ENVIRONMENT.md) · [Открытые проверки](../../checks/01_test.go)
