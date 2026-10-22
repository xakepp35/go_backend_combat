# JobFlow: эталон 21 — Повторы и идемпотентный эффект

Открытая самостоятельная версия **после** ДЗ 21. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Повтор логической операции не увеличивает счётчик второй раз.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=21 PROJECT=reference/21`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/21_RECOVERY.md) · [ДЗ](../../combat/21_RECOVERY.md) · [Открытые проверки](../../checks/21_test.go)
