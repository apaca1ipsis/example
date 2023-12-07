# Конфигурация "Одна партиция - один обработчик"

Запустите контейнер с kafka

``` sh
docker-compose up -d
```

Создайте топик numbers с одной партицией

``` sh
docker exec -it kafka kafka-topics --create --topic numbers --bootstrap-server localhost:9092 --partitions 1
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

Запустите приложение `service` в отдельном терминале

``` sh
go run cmd/service/main.go
```

Запустите приложение `gen` в отдельном терминале

``` sh
go run cmd/gen/main.go
```

Приложение `gen` напишет в терминал о том, что сообщения отправлены в kafka

``` sh
YYYY/MM/DD 12:00:00 INFO messages sent
```

Последим обработается сообщение С:3. В терминале С.db видим последнюю запись журнала

``` sh
YYYY/MM/DD 12:00:09 INFO written number=C:3
```

Все 9 сообщений, сгенерированных приложением `gen`, находились в одной очереди (в единственной партиции kafka) и обрабатывались последовательно в течение 9 секунд

``` txt
12:00:09 - 12:00:00 = 9 секунд
```


show select from db
``` sh
docker exec -it kafka-parallel-queues-pg-1  psql -U user -h localhost db -c 'select * from logs'
```

run compose
``` sh
docker-compose up -d --build
```


