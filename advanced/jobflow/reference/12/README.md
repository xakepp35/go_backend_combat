# JobFlow: эталон 12 — 128 байт: фиксированный заголовок

Открытая самостоятельная версия **после** ДЗ 12. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Проверки границ и размеров предшествуют доступу к срезам.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=12 PROJECT=reference/12`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/12_BINARY_HEADER.md) · [ДЗ](../../combat/12_BINARY_HEADER.md) · [Открытые проверки](../../checks/12_test.go)
