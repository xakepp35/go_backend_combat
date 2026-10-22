# JobFlow: эталон 09 — Транзакция и намерение события

Открытая самостоятельная версия **после** ДЗ 09. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Черновик и его intent сохраняются вместе либо не сохраняется ни один.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=09 PROJECT=reference/09`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/09_TRANSACTIONS.md) · [ДЗ](../../combat/09_TRANSACTIONS.md) · [Открытые проверки](../../checks/09_test.go)
