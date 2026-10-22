# Локальный JobFlow

Состав: PostgreSQL — подтверждённые данные; Redis — CAS-кэш, проекции и Streams; Redpanda — Kafka API; `pub`, `sub`, `act`, `ret` — четыре Go-процесса. Версии образов закреплены в [Compose](docker-compose.yml). Приложения работают от непривилегированного пользователя, порты инфраструктуры опубликованы только на localhost.

## Полная система, этапы 23–24

Команды из корня курса:

```bash
make up
make demo
curl -fsS http://localhost:8080/healthz
curl -fsS http://localhost:8080/metrics
docker compose -f deploy/docker-compose.yml --profile app logs --tail=50 pub sub act ret
make down
```

`make demo` отправляет diamond DAG из четырёх задач и ждёт SUCCEEDED у последней. Миграции применяет pub до открытия HTTP. Тема `job.events` создаётся отдельным init-контейнером. `healthz` показывает жизнь процесса; принятие и выполнение DAG проверяются отдельно.

`down` сохраняет тома. Удалить данные **этого проекта** можно командой `docker compose -f deploy/docker-compose.yml --profile app down -v`. Перед сбросом убедитесь, что это нужный локальный проект. Учебные имя/пароль БД записаны в Compose; развёртывание для внешних пользователей в программу не входит.

## Собственная реализация

В своём проекте добавьте [Dockerfile](Dockerfile) на этапе 23: он собирает команды pub, sub, act, ret и demo. Из корня курса:

```bash
JOBFLOW_PROJECT=/absolute/path/to/jobflow make up
make demo
make down
```

Default build context — `reference/24`. Для параллельного экземпляра задайте отдельные `COMPOSE_PROJECT_NAME`, `JOBFLOW_IMAGE_TAG` и свободные порты из [.env.example](.env.example). Переменные можно экспортировать в терминале или скопировать пример в `deploy/.env` при работе из папки deploy.

## Постепенная локальная разработка

С 07 поднимите инфраструктуру из корня курса:

```bash
docker compose -f deploy/docker-compose.yml up -d --wait postgres
```

На уроке 07 выполняйте ручные SQL-упражнения в отдельной схеме lesson07. Команды из корня курса; путь к файлам замените на свой:

```bash
docker compose -f deploy/docker-compose.yml exec -T postgres psql -U jobflow -d jobflow -c 'CREATE SCHEMA IF NOT EXISTS lesson07'
docker compose -f deploy/docker-compose.yml exec -T -e 'PGOPTIONS=-c search_path=lesson07' postgres psql -U jobflow -d jobflow -v ON_ERROR_STOP=1 < /absolute/path/to/jobflow/migrations/001_core.sql
docker compose -f deploy/docker-compose.yml exec -T -e 'PGOPTIONS=-c search_path=lesson07' postgres psql -U jobflow -d jobflow -v ON_ERROR_STOP=1 < /absolute/path/to/jobflow/sql/01_basics.sql
```

С 08 pub сам применяет миграции к пустой public и регистрирует checksum. Ручной SQL урока 07 в public не выполняйте: иначе таблицы уже будут существовать без записи в журнале миграций.

С 15 добавьте `redis`, с 16 — `redpanda redpanda-init`. Запускайте `go run ./cmd/pub` из **собственного проекта**, где находятся `go.mod` и `migrations`. С 17 отдельно запускается sub, с 19 act, с 22 ret.

```bash
export DATABASE_URL='postgres://jobflow:jobflow@localhost:15432/jobflow?sslmode=disable'
export REDIS_ADDR='localhost:16379'
export KAFKA_BROKERS='localhost:19092'
export KAFKA_TOPIC='job.events'
export KAFKA_GROUP='jobflow-sub'
export JF_NAMESPACE='jobflow'
export PUB_URL='http://localhost:8080'
export HTTP_ADDR=':8080'
export WORKER_COUNT='2'
go run ./cmd/pub
```

Эти переменные задаются в каждом терминале. В Docker используются внутренние имена postgres/redis/redpanda/pub. При нескольких реализациях нужны отдельные БД, topic/group и namespace. Времена nbf/exp — Unix seconds UTC.

## Восстановление

- Перезапуск act: pending сообщение забирается другим consumer после 6 секунд; SQL lease истекает через 5 секунд. Maintenance разрешает новую попытку с задержкой. Старый token не принимает результат.
- Потеря Redis: sub при старте читает историю и snapshot из pub. Фоновая сверка обнаруживает пустой индекс или недостающий hash, повторяет восстановление и сигналы. CAS-кэш заполняется из PostgreSQL при чтении.
- Kafka недоступна: принятие workflow оставляет событие в Outbox; после восстановления relay повторяет доставку. Число outbox-записей не должно уменьшаться до подтверждённой публикации.
- ret повторяет запрос очистки каждую минуту. Retention терминальных тел — час; живые ссылки и недоставленные события защищают CAS.

Участники выполняют только зарегистрированные кооперативные Go-обработчики. Система допускает повтор попытки; демонстрационный counter защищён ledger на стороне pub. Ограничения и API: [PROJECT](../docs/PROJECT.md).

## Автоматические проверки

`make check LESSON=NN PROJECT=...` создаёт отдельное окружение со случайными host-портами и удаляет его по завершении. Оно не использует пользовательские тома и не останавливает другие Docker-проекты. `--infra-project` — явный режим повторного использования существующего **учебного** окружения; требует JF_TEST_DATABASE_URL, JF_TEST_REDIS_ADDR и JF_TEST_KAFKA_BROKERS по этапу. Этот режим ничего не удаляет.
