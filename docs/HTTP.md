# Посмотреть HTTP-разговор, а не гадать по окну браузера

[Урок 10](../student/10_HTTP.md) · [API 11](../student/11_API.md) · [Договор проекта](PROJECT.md)

## 1. Запуск и адрес

В своём проекте с Compose: `docker compose up -d --build app`. До 11 задачи ещё не доступны по HTTP; на 10 есть /hello, /healthz и /about. Для хостового Go-запуска используй `go run ./cmd/jobflow`; не держи одновременно два сервера на одном опубликованном порту.

Внутри контейнера `HTTP_ADDR=:8080`, с хоста обращаемся к localhost:8080 (или своему HTTP_PORT). Ошибка соединения означает отсутствие HTTP-ответа. Статусы 404/405/500 — уже полученный ответ сервера; диагностика пойдёт по другому пути.

## 2. Первый GET

```sh
curl -i http://localhost:8080/hello
curl --get --data-urlencode "name=Аня" http://localhost:8080/hello
curl -i -X POST http://localhost:8080/hello
curl -I http://localhost:8080/hello
```

`-i` показывает заголовки вместе с телом, `-X` выбирает метод, `-I` делает HEAD. На Windows пиши `curl.exe`, если curl — псевдоним PowerShell. Персональное приветствие относится к примеру 10; на 11 /tasks заменяет вводную лабораторию.

В первой строке ответа ищи статус. После заголовков и пустой строки — тело. HEAD не передаёт тело ответа. Наличие текста в журнале приложения не означает, что клиент уже получил этот текст.

## 3. JSON в файле снимает войну с кавычками

С 11 создай request.json:

```json
{"id":"a","title":"Повторить Go","priority":2}
```

```sh
curl -i -H "Content-Type: application/json" --data-binary '@request.json' http://localhost:8080/tasks
curl -i http://localhost:8080/tasks/a
curl -i http://localhost:8080/tasks
curl -i -X PUT -H "Content-Type: application/json" --data-binary '@request.json' http://localhost:8080/tasks/a
curl -i -X POST http://localhost:8080/tasks/a/complete
curl -i -X DELETE http://localhost:8080/tasks/a
```

`--data-binary` отправляет содержимое файла; без явного метода curl здесь выберет POST. Для PUT метод указываем. В request.json можно не передавать ID для PUT: сервер возьмёт его из пути; другой непустой ID — 400. DELETE 204 без тела — успешный ответ, а не «сервер ничего не сделал».

Для повторных опытов выбирай понятные ID и учитывай прежнее состояние. Дубликат даёт 409, а не новый успех 201. На 11–13 перезапуск app очищает map; с 14 данные живут в PostgreSQL. Не удаляй том ради повторения примера.

## 4. Альтернатива через знакомый Python

```python
from urllib.error import HTTPError, URLError
from urllib.request import urlopen

try:
    with urlopen("http://localhost:8080/hello", timeout=5) as response:
        print(response.status)
        print(response.read().decode())
except HTTPError as error:
    print("Сервер ответил:", error.code)
    print(error.read().decode())
except URLError as error:
    print("HTTP-ответ не получен:", error.reason)
```

Сохрани client.py и запусти своей командой python/python3. HTTPError обрабатываем раньше URLError, от которого он наследуется. Таймаут клиента ограничивает его ожидание; он не обещает, что сервер ничего не успел сделать.

## 5. Небольшая процедура диагностики

1. Назови ожидаемые метод, путь, тело и статус.
2. Сохрани фактический ответ или ошибку соединения.
3. При отсутствии ответа проверь ps/logs, адрес и публикацию порта.
4. При ответе проверь маршрут, формат данных и причину отказа.
5. После операции прочитай состояние: один текст ошибки не доказывает сохранность данных.

Каркас JSON ограничивает тело 64 КиБ, неизвестные поля и второй документ. Он не отклоняет повторные ключи специально и не проверяет Content-Type; все ошибки Decode сейчас 400. Это точный предел примера, а не универсальная рекомендация для любого API.

Для существующих проверок HTTP в Go используется [httptest](https://pkg.go.dev/net/http/httptest). Такие проверки маршрутизатора не заменяют наблюдение настоящего сетевого сервера, контейнерных адресов или PostgreSQL.
