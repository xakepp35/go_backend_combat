# JobFlow: эталон 05 — Пакеты и границы зависимостей

Открытая самостоятельная версия **после** ДЗ 05. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Прикладной сервис зависит от поведения, требуемого от хранилища.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=05 PROJECT=reference/05`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/05_BOUNDARIES.md) · [ДЗ](../../combat/05_BOUNDARIES.md) · [Открытые проверки](../../checks/05_test.go)
