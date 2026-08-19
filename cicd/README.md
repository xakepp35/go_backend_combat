<!-- artifact:cicd/README.md v1.0.0 2026-08-19T07:03:00Z Topic: Delivery Summary: Практическая инструкция по локальному CI/CD и запуску JobFlow через Docker Compose -->

# JobFlow — Local CI/CD

Эта папка содержит готовый минимальный набор для запуска JobFlow в одинаковом окружении у разработчика, на CI и в Docker.

Здесь нет Kubernetes и сложной DevOps-инфраструктуры. На этом курсе она не нужна.

Нам достаточно научиться правильно работать с:

* `Docker`
* `Docker Compose`
* `.env`
* `Makefile`
* `Dockerfile`

Это уже важный backend-навык.

Сеньор должен уметь не только написать `service.go`, но и ответить на вопросы:

> Как запустить проект с нуля?

> Какие зависимости ему нужны?

> Как проверить, что код действительно работает?

> Как воспроизвести окружение у другого разработчика?

> Как собрать production-like binary/image?

> Что запускает CI?

---

# 1. Как начать работу

Сначала создайте собственный репозиторий для проекта.

Например:

```bash
mkdir jobflow
cd jobflow

git init
```

Далее создайте структуру самого приложения.

В результате корнем проекта должен быть именно `jobflow`.

Файлы из этой папки `cicd/` нужно **скопировать в корень вашего JobFlow**.

То есть после копирования должно получиться примерно так:

```text
jobflow/
├── Dockerfile
├── Makefile
├── docker-compose.yml
├── .env
├── .env.example
├── .gitignore
├── README.md
├── cmd/
├── internal/
└── migrations/
```

`cicd/` после этого не является обязательной runtime-папкой приложения.

Она содержит шаблон developer workflow, а не исходный код JobFlow.

---

# 2. Что находится в этой папке

## `.env`

Ваши локальные настройки.

Например:

```text
PostgreSQL
Redis
Kafka / Redpanda
HTTP
worker count
connection pool
```

Это файл для **локальной машины**.

Не храните в нём реальные production passwords, API keys и tokens.

---

## `.env.example`

Шаблон конфигурации.

Его можно коммитить в Git.

Он показывает:

> «Какие настройки вообще нужны проекту?»

Обычно новый разработчик делает:

```bash
cp .env.example .env
```

После этого получает рабочую локальную конфигурацию.

---

## `.gitignore`

Говорит Git, какие файлы не нужно добавлять в repository.

Особенно важно:

```text
.env
```

Локальные credentials не должны случайно попасть в Git.

---

## `docker-compose.yml`

Поднимает локальную инфраструктуру проекта.

На нашем курсе это:

```text
PostgreSQL
Redis
Redpanda
```

То есть вам не нужно устанавливать все эти сервисы вручную на Windows/Linux/macOS.

Docker запускает их в containers.

---

## `Dockerfile`

Описывает, как собрать сам JobFlow в Docker image.

Используется multi-stage build:

сначала Go compiler собирает binary, затем binary переносится в маленький runtime image.

В результате production image не содержит:

```text
Go compiler
source code
Go modules cache
```

В нём остаётся только то, что нужно для запуска приложения.

---

## `Makefile`

Это удобная командная панель проекта.

Вместо длинных Docker/Go команд мы используем короткие:

```bash
make up
make test
make migrate
make verify
```

Это особенно полезно в команде: у проекта появляется единый developer interface.

---

# 3. Первый запуск

После копирования файлов в JobFlow проверьте, что у вас есть `.env`.

Если его ещё нет:

```bash
cp .env.example .env
```

На Windows PowerShell аналог:

```powershell
Copy-Item .env.example .env
```

После этого поднимите инфраструктуру:

```bash
make up
```

Docker скачает необходимые images и запустит контейнеры.

Проверить состояние:

```bash
make ps
```

Вы должны увидеть PostgreSQL, Redis и Redpanda.

Если хотите посмотреть логи:

```bash
make logs
```

---

# 4. Что делает `make up`

Команда:

```bash
make up
```

запускает:

```text
PostgreSQL
Redis
Redpanda
```

в фоне.

После этого ваш Go application можно запускать обычным способом с вашей машины:

```bash
make run
```

Это очень удобный режим разработки.

Go-код работает локально, а инфраструктура живёт в Docker.

---

# 5. База данных

Когда PostgreSQL запущен, нужно применить migrations.

```bash
make migrate
```

Это создаст schema JobFlow.

Посмотреть состояние migrations:

```bash
make migrate-status
```

Откатить последнюю migration:

```bash
make migrate-down
```

На практике сначала пользуйтесь:

```bash
make migrate
```

А `migrate-down` нужен прежде всего для обучения, экспериментов и локальной разработки.

---

# 6. Запуск самого JobFlow

Когда infrastructure поднята и migrations применены:

```bash
make run
```

Обычно приложение будет слушать:

```text
localhost:8080
```

После этого можно обращаться к API через browser, `curl`, Postman или любой другой HTTP client.

Например:

```bash
curl http://localhost:8080/healthz
```

---

# 7. Полностью остановить окружение

Остановить containers:

```bash
make down
```

Данные при этом обычно сохраняются в Docker volumes.

Это означает, что PostgreSQL data не исчезает просто потому, что вы остановили containers.

---

# 8. Полностью очистить окружение

Если хотите начать локальную инфраструктуру практически с нуля:

```bash
make clean
```

Эта команда удаляет containers и volumes.

Это означает, что локальные данные PostgreSQL, Redis и Redpanda будут потеряны.

Для курса это очень полезная команда.

Например, можно проверить:

> «Восстанавливается ли database schema только migrations?»

Делаем:

```bash
make clean
make up
make migrate
```

Если после этого JobFlow снова работает, значит schema действительно воспроизводима.

---

# 9. Как проверить код

Основные команды Go доступны через Makefile.

## Форматирование

```bash
make fmt
```

Запускает:

```bash
go fmt ./...
```

---

## Статическая проверка

```bash
make vet
```

Запускает:

```bash
go vet ./...
```

---

## Тесты

```bash
make test
```

Это:

```bash
go test ./...
```

---

## Race detector

Для нашего курса это особенно важно.

```bash
make race
```

Запускает:

```bash
go test -race ./...
```

После урока по concurrency эту команду нужно использовать регулярно.

---

## Полная локальная проверка

Самая важная команда:

```bash
make verify
```

Она выполняет:

```text
format
vet
test
race
```

То есть перед commit можно просто сделать:

```bash
make verify
```

и проверить основную инженерную поверхность проекта.

---

# 10. Сборка binary

Чтобы собрать приложение:

```bash
make build
```

Binary появится в:

```text
bin/jobflow
```

Это полезно для понимания разницы между:

```text
source code
→ build
→ executable
```

Go в production обычно запускается не через:

```bash
go run
```

а через заранее собранный binary.

---

# 11. Docker image

Когда приложение уже работает локально:

```bash
make image
```

Docker прочитает `Dockerfile` и соберёт image.

Проверить images:

```bash
docker images
```

---

# 12. Запуск приложения в Docker

После сборки:

```bash
make image-run
```

Приложение будет запущено уже внутри Docker container.

Это важный этап.

До этого было:

```text
Go application → local machine
```

Теперь:

```text
Docker image → container
```

То есть мы проверяем не только код, но и packaging приложения.

---

# 13. Docker Compose и Docker image — разные вещи

Это важно понимать.

`docker-compose.yml` нужен прежде всего для **окружения**:

```text
PostgreSQL
Redis
Redpanda
```

`Dockerfile` нужен для **нашего приложения**:

```text
JobFlow binary
```

То есть:

* Compose говорит, какие сервисы локально нужны.
* Dockerfile говорит, как упаковать JobFlow.

Это можно использовать вместе или отдельно.

---

# 14. Типичный рабочий день

Обычная последовательность разработки будет выглядеть примерно так:

```bash
make up
make migrate
make run
```

Пишете код.

После изменений:

```bash
make test
```

Если работали с concurrency:

```bash
make race
```

Перед commit:

```bash
make verify
```

Если хотите проверить Docker packaging:

```bash
make image
```

---

# 15. Если что-то перестало работать

Сначала посмотрите containers:

```bash
make ps
```

Потом:

```bash
make logs
```

Проверьте migration:

```bash
make migrate-status
```

Проверьте тесты:

```bash
make test
```

Если проблема похожа на race:

```bash
make race
```

Если environment начал вести себя странно:

```bash
make clean
make up
make migrate
```

Это хороший первый способ понять, действительно ли проблема в коде, или локальное состояние окружения стало неконсистентным.

---

# 16. Что нужно понимать про `.env`

`.env` нужен для локального запуска.

Например:

```dotenv
DATABASE_URL=postgres://jobflow:jobflow@localhost:5432/jobflow?sslmode=disable
REDIS_ADDR=localhost:6379
KAFKA_BROKERS=localhost:19092
```

Ваш код не должен содержать такие значения прямо в source code:

```go
const password = "jobflow"
```

Конфигурация должна приходить извне приложения.

Это будет особенно важно дальше, когда мы поговорим о production deployment.

---

# 17. Почему Compose полезен для senior backend

Docker Compose кажется простой штукой.

Но он учит важным инженерным навыкам.

Вы начинаете понимать:

* какие dependencies нужны приложению;
* какие ports они используют;
* какие credentials нужны;
* какие volumes сохраняют state;
* какие services должны быть готовы раньше других;
* какие health checks существуют;
* как application взаимодействует с infrastructure;
* что относится к application, а что к environment.

Это уже гораздо ближе к реальной backend-разработке, чем просто:

```bash
go run main.go
```

---

# 18. Почему Kubernetes мы здесь не используем

В этом курсе Kubernetes намеренно не является частью основной программы.

Причина простая.

Сначала нужно научиться понимать:

```text
application
configuration
process
container
network
database
queue
cache
health
shutdown
resource limits
```

Docker и Compose позволяют увидеть эти вещи без огромного слоя orchestration.

Когда человек действительно понимает lifecycle приложения, Kubernetes становится следующим инструментом, а не магическим набором YAML-файлов.

---

# 19. CI

После того как проект работает локально, те же проверки должны выполняться автоматически.

CI запускает:

```text
format
vet
test
race
docker build
```

Это означает:

> Если код работает у меня на ноутбуке, этого недостаточно.

Нужно, чтобы repository мог воспроизводимо проверить состояние проекта в чистом окружении.

Именно поэтому Makefile и CI должны использовать одну и ту же базовую модель команд.

---

# 20. Главное, что нужно запомнить

Для начинающего достаточно держать в голове несколько команд.

Запустить infrastructure:

```bash
make up
```

Посмотреть состояние:

```bash
make ps
```

Применить database schema:

```bash
make migrate
```

Запустить приложение:

```bash
make run
```

Запустить tests:

```bash
make test
```

Проверить concurrency:

```bash
make race
```

Проверить всё:

```bash
make verify
```

Собрать Docker image:

```bash
make image
```

Остановить environment:

```bash
make down
```

Полностью очистить local state:

```bash
make clean
```

---

# 21. Минимальный workflow ученика

При первом запуске:

```bash
cp .env.example .env
make up
make migrate
make run
```

После написания кода:

```bash
make test
```

После concurrency-изменений:

```bash
make race
```

Перед отправкой кода:

```bash
make verify
```

При необходимости проверить контейнер:

```bash
make image
```

---

# 22. Definition of Done для локальной среды

Работа считается корректно настроенной, когда новый разработчик может:

1. Скопировать CI/CD-файлы в корень нового `jobflow`.
2. Создать `.env` из `.env.example`.
3. Выполнить `make up`.
4. Выполнить `make migrate`.
5. Запустить `make run`.
6. Выполнить `make test`.
7. Выполнить `make race`.
8. Собрать `make image`.

И всё это без ручной установки PostgreSQL, Redis и Kafka/Redpanda.

Главная идея этой папки:

> **Окружение проекта тоже является частью инженерного продукта.**

Хороший backend начинается не с первого `struct`, а с возможности **воспроизводимо поднять, запустить, проверить, упаковать и остановить систему**.
