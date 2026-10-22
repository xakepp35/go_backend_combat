# JobFlow: эталон 13 — Поток, связи, тело и padding

Открытая самостоятельная версия **после** ДЗ 13. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Разбор сохраняет исходное тело и не поглощает следующий фрейм.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=13 PROJECT=reference/13`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/13_FRAMING.md) · [ДЗ](../../combat/13_FRAMING.md) · [Открытые проверки](../../checks/13_test.go)
