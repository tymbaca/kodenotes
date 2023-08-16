## Конфигурация

Для запуска необходимо создать `.env` файл со следующими необходимыми переменными:

```
TARGET_STAGE=run

POSTGRES_HOST=postgres
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword

SERVER_ADDRESS=:8080
POSTGRES_DB=kodenotes
```

### Обязательно

- `TARGET_STAGE`: `[ run | test ]` целевой этап сборки `Dockerfile`а. Используется в 
  `compose.yaml`. `run` полноценно запустит веб-сервер. `test` запустит тесты.
  Данная переменная нужна только для сборки с `docker compose`. При желании его
  можно убрать из `.env` файла, и пользоваться утилитой `make` + `Makefile`, где
  данная переменная уже задана.

- `POSTGRES_HOST`: адрес запущенного PostreSQL БЕЗ ПОРТА. При запуске с `docker compose`
  необходимо указать название сервиса с PostgreSQL (из `compose.yaml` файла).

- `POSTGRES_USER`: имя пользователя для соединения с PostgreSQL. 

- `POSTGRES_PASSWORD`: пароль пользователя для соединения с PostgreSQL. 


### Необязательно

- `POSTGRES_DB`: (опционально) название базы данных внутри PostgreSQL. Если не указать,
    то по умолчанию будет назначено название `postgres`

- `SERVER_ADDRESS`: (опционально) адрес, который будет слушать веб-сервер. Лучше указать
только порт в формате `:8080`. Так сервер будет слушать локальную
сеть (0.0.0.0) с указанным портом.

