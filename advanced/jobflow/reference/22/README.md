# JobFlow: эталон 22 — ret и безопасная очистка

Открытая самостоятельная версия **после** ДЗ 22. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: GC сохраняет живые ссылки, недоставленные события и идентичность задачи.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=22 PROJECT=reference/22`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/22_RETENTION.md) · [ДЗ](../../combat/22_RETENTION.md) · [Открытые проверки](../../checks/22_test.go)
