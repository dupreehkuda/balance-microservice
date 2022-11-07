# Balance microservice
### Description
Микросервис для работы с балансом пользователей.

Сервис может вывести, зачислить, снять, зарезервировать и перевести средства счета по соответствующим хендлерам.

Реализована трехслойная архитектура `handler` -> `action` -> `storage`. 
Чтобы слои имели меньшее зацепление, все передается через интерфейс слоя.
Сервис также способен читать и писать запросы в gzip через `middleware`.

### Deploy
Поднять контейнеры с сервисом и с базой данных можно командой:
```sh 
$ make compose
```

Остановка контейнеров:
```sh
$ make compose-down
```

Пересобрать контейнеры можно командой 
```sh 
$ make rebuild
```

Для того чтобы запустить приложение не используя Docker, необходимо подтянуть зависимости и запустить с флагами:
```sh
$ go run cmd/main.go -a <port> -d <postgres URI>
```
Без флага `-a` запустится на порте `:8080`. Также эти параметры можно задать через переменные окружения: `RUN_ADDRESS` и `DATABASE_URI`.