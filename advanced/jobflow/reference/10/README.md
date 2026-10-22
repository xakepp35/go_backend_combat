# JobFlow: эталон 10 — Горутины и ограничение параллелизма

Открытая самостоятельная версия **после** ДЗ 10. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Активных обработчиков не больше заданного workers.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=10 PROJECT=reference/10`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/10_WORKER_POOL.md) · [ДЗ](../../combat/10_WORKER_POOL.md) · [Открытые проверки](../../checks/10_test.go)
