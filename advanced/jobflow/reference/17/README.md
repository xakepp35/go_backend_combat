# JobFlow: эталон 17 — sub и восстанавливаемые проекции

Открытая самостоятельная версия **после** ДЗ 17. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Проекция принимает только более новую версию задачи.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=17 PROJECT=reference/17`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/17_PROJECTIONS.md) · [ДЗ](../../combat/17_PROJECTIONS.md) · [Открытые проверки](../../checks/17_test.go)
