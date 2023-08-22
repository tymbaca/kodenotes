# KODENotes
## Описание
Для общения с PostgreSQL используется стандартный пакет `database/sql` с драйвером
`github.com/lib/pq`.

Аутентификация базовая (Basic Auth). Пароли хранятся в базе данных в зашифрованном 
виде (1 итерация, SHA-256), без соли. Есть регистрация, делается путем POST запроса на `/register` с хеддером Basic Auth. При попытке зарегистрировать пользователя с уже существующим username'ом будет возвращена ошибка `422 Unprocessable Content`. Присутствует ограничение на длину логина/пароля — 250 символов на каждый.

В Яндекс.Спеллер запросы отправляются методом POST. По документации серсива, 
максимальный размер текста в этом случае составляет 10 000 символов. API будет 
возвращать ошибку `413 Payload Too Large`, если текст превышает это ограничение. Также через переменную окружения `YANDEX_SPELLER_TIMEOUT` можно задать время, по истечению которого, если Яндекс.Спеллер не отвечает, будет возвращена ошибка `502 Bad Gateway`.

При наличии орфографической ошибки в тексте создаваемой заметки, будет возвращена ошибка `400 Bad Request` вместе с `json` ответом, который возвращает Яндекс.Спеллер (см. `spellcheck/yandex.md`). При этом сама заметка создана не будет.

## Запуск

Для запуска в корне проекта необходимо создать `.env` файл со следующими переменными:

```
POSTGRES_HOST=postgres
POSTGRES_PASSWORD=mypassword
YANDEX_SPELLER_URL=https://speller.yandex.net/services/spellservice.json/checkText

SERVER_PORT=8080
YANDEX_SPELLER_TIMEOUT=10
```

После чего запустить следующей командой:
```bash
docker compose up
```

> Docker скачает [образ](https://hub.docker.com/r/tymbaca/kodenotes/tags) из docker.hub. Если что-то пойдет не так — замените в `compose.yaml` строку `images: tymbaca/kodenotes:latest` на `build: .`.

Альтернативно, сервер можно запустить локально (вне контейнера) следующим образом:
```bash
make pg-up
make run
```

Тесты можно запустить следующим образом:
```bash
make pg-up
make test
```

В `Makefile` созданы команды для удобного запуска и использования контейнера с PostgreSQL:
- `make pg-up`: поднимает контейнер с PostgreSQL (с именем `kodenotes-postgres`)
- `make pg-down`: останавливает контейнер с PostgreSQL
- `make pg-restart`: перезапускает контейнер с PostgreSQL (комбинация `pg-down` и `pg-up`)
- `make pg-shell`: входит в командную оболочку PostgreSQL (с опциями `-U postgres -d postgres`). Эта команда может также использоваться при поднятии сервиса через `docker compose` (в `compose.yaml` у PostgreSQL задано такое же название контейнера)

> **ВАЖНО** Перед использованием `docker compose up` убедитесь, что контейнер, оставшийся от запуска PostgreSQL через `make` удален. Сделать это можно командой `make pg-down`.

## Тестирование через Postman
Файл с коллекцией запросов лежит в `docs/KODENotes.postman_collection.json`. После его импорта 
рекомендуется запускать запросы в том порядке, в котором они представлены в коллекции, для большей 
наглядности.

## Конфигурация
### Обязательно

- `POSTGRES_HOST`: адрес запущенного PostreSQL БЕЗ ПОРТА. При запуске с `docker compose`
  необходимо указать название сервиса с PostgreSQL (из `compose.yaml` файла).

- `POSTGRES_PASSWORD`: пароль пользователя для соединения с PostgreSQL. 

- `YANDEX_SPELLER_URL`: URL API Яндекс.Спеллер, который принимает POST запросы. 
  На момент создания программы корректный URL - `https://speller.yandex.net/services/spellservice.json/checkText`.
  Подробную документацию по API можно найти по [ссылке](https://yandex.ru/dev/speller/doc/ru/reference/checkText).

- `SERVER_PORT`: порт, который будет слушать веб-сервер (по умолчанию будет задан порт `8080`). 

### Необязательно

- `YANDEX_SPELLER_TIMEOUT`: время в секундах, после которого запрос на сервер Яндекс.Спеллер
  отменяется и возвращается ошибка `502 Bad Gateway`. По умолчанию будет установлено 10 секунд.

## Формат запросов

### Регистрация
Для регистрации должны соблюдаться следующие условия:
- POST запрос на `/register`
- Заголовок Basic Auth с логином и паролем создаваемого пользователя
- Логин должен быть уникальным. Если в базе данных уже есть пользователь с тем же логином, то будет возвращена ошибка `422 Unprocessable Content`

### Создание заметки
Для создании заметки должны соблюдаться следующие условия:
- POST запрос на `/notes`
- Заголовок Basic Auth с логином и паролем пользователя (например: заголовок `Authorization` со значением `Basic dXNlcm5hbWU6cGFzc3dvcmQ=`, где `dXNlcm5hbWU6cGFzc3dvcmQ=` == `username:password`, который закодирован в `base64`)
- Заголовок `Content-Type` должен создержать значение `application/json`
- Тело запроса представляет `json` объект с обязательным полем `text` типа `string`, которое хранит текст заметки. При необходимости использования специальных символом, кавычек и пр., текст вожно отправлять в виде URL Encoded строки.

Пример тела запроса:

```json
{
  "text": "Your note text here"
}
```

Пример тела запроса с URL Encoded текстом (содержит закодированные переносы строк, кавычки разных типов и `\`'ы):

```json
{
  "text": "%D0%9F%D1%80%D0%B8%D0%B2%D0%B5%D1%82%20%D1%8D%D1%82%D0%BE%20%D0%BC%D0%BE%D0%B9%20%D1%82%D0%B5%D0%BA%D1%81%D1%82%20%D1%8F%20%D1%82%D1%83%D1%82%20%D0%BF%D0%B8%D1%88%D1%83%0A%0A%D0%9C%D0%BD%D0%BE%D0%B3%D0%BE%20%D0%B0%D0%B1%D0%B7%D0%B0%D1%86%D0%B5%D0%B2%20%D0%BA%D0%B0%D0%BA%20%D0%B2%D0%B8%D0%B4%D0%B8%D1%88%D1%8C.%20%22%D0%98%20%D0%BA%D0%B0%D0%B2%D1%8B%D1%87%D0%B5%D0%BA%22%20%27%D1%80%D0%B0%D0%B7%D0%BD%D1%8B%D1%85%27%0A%0A%D0%98%20%D0%B1%D0%B5%D0%BA%D1%81%D0%BB%D0%B5%D1%88%D0%B5%D0%B9%20%D1%82%D0%BE%D0%B6%D0%B5%20%D0%BC%D0%BD%D0%BE%D0%B3%D0%BE%20%5C%5C%27%5C%22%5C%22%22%20%5C%5C%5C%5C%5C%20%D1%81%20%D0%BA%D0%B0%D0%B2%D1%8B%D1%87%D0%BA%D0%B0%D0%BC%D0%B8%20%5Cn%20%5Cr%20%0A7%20%D1%8D%D1%82%D0%BE%20%D0%BA%D1%83%D1%81%D0%BE%D0%BA%20%D1%81%D1%8B%D1%80%D0%BE%D0%B3%D0%BE%20%D1%83%D1%80%D0%BB%20%D0%B5%D0%BD%D0%BA%D0%BE%D0%B4%D0%B5%D0%B4%20%D1%82%D0%B5%D0%BA%D1%81%D1%82%D0%B0"
}
```

### Получение заметок
Для получения заметок текущего пользователя должны соблюдаться следующие условия:
- GET запрос на `notes`
- Заголовок Basic Auth с логином и паролем пользователя (например: заголовок `Authorization` со значением `Basic dXNlcm5hbWU6cGFzc3dvcmQ=`, где `dXNlcm5hbWU6cGFzc3dvcmQ=` == `username:password`, который закодирован в `base64`)



## Трудности
При интеграции с Яндекс.Спеллер возникла проблема: в документации нигде не указано 
в каком именно формате должно быть тело запроса при использовании через POST запрос. 
Альтернативный вариант с GET запросом требует отправки текста в виде path-параметра, 
а это имеет серьезный недостаток — максимальный размер URL — 2083 символов для 
Google Chrome, что довольно мало для приложения заметок.

В итоге, с помощью Chrome DevTools получилось узнать, что тело запроса отправляется 
в виде `form-data`.

Были также трудности с системой аутентификации. Сначала было задумно сделать более сложную 
версию с созданием временных сессий, куками и прочим удовольствием, но понял, что:

- это займет много времени,
- я не смог заставить работать автоудаление устаревших сессий. Была проблема с "прокидыванием"
  необходимых значений времени в SQL формулу ([ссылка](https://github.com/tymbaca/kodenotes/blob/835253728ded9a932784f90cab7c3edf5d20cbfa/database/postgres.go#L137C1-L151C2)).

После этого, очередная трудность настигла меня в самом неприятном месте - в тестах. После того, как я "накатил" необходимые тесты, я заметил одну странную вещь: все тесты проходили успешно, когда я запускал их в VSCode, однако, если запускал через `go test ./...`, то они падали с ошибкой от PostgreSQL "pq: function uuid_generate_v4() does not exist". В итоге ошибка оказалась в том, что, видимо, драйвер базы данных не особо умеет в concurency, а команда `go test` запускает тесты *параллельно*. Решением стал запуск тестов с флагом `-p 1`, который ограничивает параллельный запуск тестов до 1 шт. ([см. коммит](https://github.com/tymbaca/kodenotes/commit/ba34c1baf8112ecef7bc06b38d0e2b3cc25854e6)).

## Что не нравится
- Организация хелпер-функций в тестах
- Не сделал нормальную сервисную оболочку между базой и хендлерами
- Не поделил логи на разные файлы, в зависимости от вызывающего (база, хендлер, тест (?) и т.д.)
- Не mock'ал неинтеграционные тесты, которые обращаются к Яндекс.Спеллер (хотя основа для этого есть - интерфейс `SpellChecker`)
- Отсутствие middleware

## Какой-то чеклист

- [x] Интеграция с PostgreSQL
- [x] Интеграция с Яндекс.Спеллер
- [x] Регистрация
- [x] Крутой логгер
- [x] ~~Тесты~~


