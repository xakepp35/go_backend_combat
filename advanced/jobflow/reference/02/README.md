# JobFlow: эталон 02 — Типы, функции и модель задачи

Открытая самостоятельная версия **после** ДЗ 02. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: ID, тип работы и исходные байты тела сохраняют свой смысл.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=02 PROJECT=reference/02`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/02_GO_MODEL.md) · [ДЗ](../../combat/02_GO_MODEL.md) · [Открытые проверки](../../checks/02_test.go)
