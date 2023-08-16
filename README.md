## Конфигурация

Для запуска необходимо создать `.env` файл со следующими необходимыми переменными:

```
POSTGRES_HOST=postgres
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword

SERVER_ADDRESS=:8080
POSTGRES_DB=kodenotes
```

### Обязательно

- `POSTGRES_HOST`: адрес запущенного PostreSQL БЕЗ ПОРТА. При запуске с `docker compose`
  необходимо указать название сервиса с PostgreSQL (из `compose.yaml` файла).

- `POSTGRES_USER`: 

- `POSTGRES_PASSWORD`: 

### Необязательно

- `POSTGRES_DB`: (опционально) название базы данных внутри PostgreSQL. Если не указать,
    то по умолчанию будет назначено название `postgres`

- `SERVER_ADDRESS`: (опционально) адрес, который будет слушать веб-сервер. Лучше указать
только порт в формате `:8080`. Так сервер будет слушать локальную
сеть (0.0.0.0) с указанным портом.

