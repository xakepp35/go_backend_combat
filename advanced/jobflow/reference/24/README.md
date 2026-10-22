# JobFlow: эталон 24 — Защита проекта

Открытая самостоятельная версия **после** ДЗ 24. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Обещания README подтверждены командами и сценарием отказа.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=24 PROJECT=reference/24`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/24_DEFENSE.md) · [ДЗ](../../combat/24_DEFENSE.md) · [Открытые проверки](../../checks/24_test.go)
