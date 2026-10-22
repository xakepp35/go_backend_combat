# JobFlow: эталон 03 — Коллекции, указатели и CRUD

Открытая самостоятельная версия **после** ДЗ 03. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Дубликат не перезаписывает запись; чтение не отдаёт изменяемое внутреннее состояние.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=03 PROJECT=reference/03`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/03_MEMORY_CRUD.md) · [ДЗ](../../combat/03_MEMORY_CRUD.md) · [Открытые проверки](../../checks/03_test.go)
