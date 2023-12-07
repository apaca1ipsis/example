# Конфигурация "Две партиции - четыре обработчика - два consumer'а"

Запустите контейнер с kafka

``` sh
docker-compose up -d
```

Создайте топик numbers с двумя партициями

``` sh
docker exec -it kafka kafka-topics --create --topic numbers --bootstrap-server localhost:9092 --partitions 2
```

Запустите модели баз данных (каждую в отдельном терминале)

``` sh
# A.db
PORT=10000 go run cmd/db/main.go
```

``` sh
# B.db
PORT=10001 go run cmd/db/main.go
```

``` sh
# C.db
PORT=10002 go run cmd/db/main.go
```

Запустите один экземпляр приложения `service` в отдельном терминале. Укажите количество каналов-очередей равным 2 и количество consumer'ов равным 2.

``` sh
CONSUMERS=2 QUEUES=2 go run cmd/service/main.go
```

Запустите приложение `gen` в отдельном терминале

``` sh
go run cmd/gen/main.go
```

Приложение `gen` напишет в терминал о том, что сообщения отправлены в kafka

``` sh
YYYY/MM/DD 12:00:00 INFO messages sent
```

В каждом из терминалов моделей базы данных даты записи одинаковых чисел совпадают. Это означает, что, например, три сообщения A:1, B:1 и C:1 располагались в параллельных очередях и начали обработку почти одновременно.


``` sh
# A.db
YYYY/MM/DD 12:00:01 INFO written number=A:1
YYYY/MM/DD 12:00:02 INFO written number=A:2
YYYY/MM/DD 12:00:03 INFO written number=A:3
```

``` sh
# B.db
YYYY/MM/DD 12:00:01 INFO written number=B:1
YYYY/MM/DD 12:00:02 INFO written number=B:2
YYYY/MM/DD 12:00:03 INFO written number=B:3
```

``` sh
# C.db
YYYY/MM/DD 12:00:01 INFO written number=C:1
YYYY/MM/DD 12:00:02 INFO written number=C:2
YYYY/MM/DD 12:00:03 INFO written number=C:3
```

В итоге обработка 9 сообщений, сгенерированных приложением `gen`, заняла 3 секунды

``` txt
12:00:03 - 12:00:00 = 3 секунды
```
