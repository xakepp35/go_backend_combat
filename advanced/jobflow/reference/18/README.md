# JobFlow: эталон 18 — Redis Streams и потерянный потребитель

Открытая самостоятельная версия **после** ДЗ 18. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Получение не означает завершение; pending-запись может получить другой потребитель.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=18 PROJECT=reference/18`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/18_STREAMS.md) · [ДЗ](../../combat/18_STREAMS.md) · [Открытые проверки](../../checks/18_test.go)
