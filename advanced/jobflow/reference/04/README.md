# JobFlow: эталон 04 — Ошибки, валидация и тесты

Открытая самостоятельная версия **после** ДЗ 04. Модуль `jobflow`; исходники не импортируют соседние этапы.

Инвариант: Невалидный ввод отклоняется до изменения хранилища.

Из этой директории: `go run ./cmd/pub`, `go build ./...`, `go vet ./...`, `go test ./...`.
Из корня курса: `make check LESSON=04 PROJECT=reference/04`.

С 06 работает HTTP для черновиков, с 08 нужен DATABASE_URL и PostgreSQL. С 14 pub принимает бинарный workflow; с 15 использует Redis, с 16 запускает Kafka relay; sub появляется в 17, Streams в 18, act в 19, ret в 22. Docker Compose всей системы вводится в 23. Переменные и запуск описаны в [deploy](../../deploy/README.md).

[Конспект](../../student/04_ERRORS_TESTS.md) · [ДЗ](../../combat/04_ERRORS_TESTS.md) · [Открытые проверки](../../checks/04_test.go)
