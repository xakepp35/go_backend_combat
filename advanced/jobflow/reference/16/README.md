# JobFlow: эталон 16 — Kafka и доставка Outbox

Открытая самостоятельная версия **после** ДЗ 16. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Недоставленное событие остаётся в Outbox до подтверждения публикации.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=16 PROJECT=reference/16`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/16_KAFKA_OUTBOX.md) · [ДЗ](../../combat/16_KAFKA_OUTBOX.md) · [Открытые проверки](../../checks/16_test.go)
