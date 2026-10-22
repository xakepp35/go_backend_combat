# JobFlow: эталон 14 — DAG и принятие workflow

Открытая самостоятельная версия **после** ДЗ 14. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Принятый граф не содержит циклов и не меняет связи после отправки.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=14 PROJECT=reference/14`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/14_DAG.md) · [ДЗ](../../combat/14_DAG.md) · [Открытые проверки](../../checks/14_test.go)
