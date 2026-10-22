# JobFlow: эталон 11 — Отмена, backpressure и shutdown

Открытая самостоятельная версия **после** ДЗ 11. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Отмена освобождает ожидающих работников, а Run ждёт их завершения.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=11 PROJECT=reference/11`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/11_CANCELLATION.md) · [ДЗ](../../combat/11_CANCELLATION.md) · [Открытые проверки](../../checks/11_test.go)
