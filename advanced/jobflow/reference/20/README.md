# JobFlow: эталон 20 — Время и планирование

Открытая самостоятельная версия **после** ДЗ 20. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Запуск невозможен до nbf или без зависимостей; на exp результат отклоняется.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=20 PROJECT=reference/20`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/20_SCHEDULING.md) · [ДЗ](../../combat/20_SCHEDULING.md) · [Открытые проверки](../../checks/20_test.go)
