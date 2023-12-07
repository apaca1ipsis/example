# Организация параллельных очередей с помощью партиций kafka и каналов go

[Конфигурация "Одна партиция - один обработчик"](./docs/one-partition-one-handler.md)

[Конфигурация "Две партиции - два обработчика"](./docs/two-partitions-two-handlers.md)

[Конфигурация "Две партиции - четыре обработчика"](./docs/two-partitions-four-handlers.md)

[Конфигурация "Две партиции - четыре обработчика - два consumer'а"](./docs/two-partitions-four-handlers-two-consumers.md)

После каждого эксперимента выполняйте

``` sh
docker-compose down -v
```


run compose
``` sh
docker-compose up -d --build
```

Подождать 5минут+ пока поднимется Кафка....

show select from db
``` sh
docker exec -it kafka-parallel-queues-pg-1  psql -U user -h localhost db -c 'select * from logs'
```

