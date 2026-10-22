# JobFlow: эталон 15 — CAS, BLAKE3 и Redis-кэш

Открытая самостоятельная версия **после** ДЗ 15. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Потеря Redis не уничтожает подтверждённый body.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=15 PROJECT=reference/15`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/15_CAS_REDIS.md) · [ДЗ](../../combat/15_CAS_REDIS.md) · [Открытые проверки](../../checks/15_test.go)
