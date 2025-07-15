# Запуск проекта
Для запуска проекта выполните следующие шаги:
## Клонируйте репозиторий:
 ```sh
   git  clone https://github.com/nongrata2/wbtest1
   cd wbtest1
   ```
## Создайте файл .env в корневой папке со следующими параметрами:
```
HTTP_SERVER_ADDRESS=
HTTP_SERVER_TIMEOUT=
LOG_LEVEL=
POSTGRES_USER=
POSTGRES_PORT=
POSTGRES_PASSWORD=
POSTGRES_DB=
KAFKA_BROKERS=
KAFKA_TOPIC=
KAFKA_GROUP_ID=
```

Пример:
```
HTTP_SERVER_ADDRESS=:8080
HTTP_SERVER_TIMEOUT=5s
LOG_LEVEL=DEBUG
POSTGRES_USER=postgres
POSTGRES_PORT=5432
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=orders_topic
KAFKA_GROUP_ID=order_service_group
```
 
## Запустите проект с помощью Docker Compose:
```sh
docker compose up --build
```
После запуска API будет доступен по адресу http://localhost:8081. Для тестирования можно использовать curl или открыть страницу в браузере по адресу http://localhost:8081