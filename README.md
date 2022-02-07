# Тестовое задание

1. Описать proto файл с сервисом из 3 методов: добавить пользователя, удалить пользователя, список пользователей
2. Реализовать gRPC сервис на основе proto файла на Go
3. Для хранения данных использовать PostgreSQL
4. На запрос получения списка пользователей данные будут кешироваться в redis на минуту и брать из редиса
5. При добавлении пользователя делать лог в clickHouse
6. Добавление логов в clickHouse делать через очередь Kafka

## Запуск
```
docker-compose up
```
Дождаться появления строки
```
app_1         | {"level":"info","ts":1644227746.196521,"caller":"app/main.go:63","msg":"Service started"}
```

## Использование
Сервис принимает только gRPC запросы. 
Протестировать можно с помощью Goland или e.g. [BloomRPC](https://github.com/bloomrpc/bloomrpc)
```
GRPC localhost:8081/users.Users/ListUsers
###
GRPC localhost:8081/users.Users/AddUser
{
  "name": "test user",
  "mail": "test userasdgfddfgdsad"
}
###
GRPC localhost:8081/users.Users/RemoveUser
{
  "id": "7"
}
```

## Что стоило бы сделать, но пока не сделано

- Тесты
- Пагинация или возвращать stream для метода `ListUsers`
- Валидация аргументов в middlewaare
- Утилитка, которая покажет, что все  работает
- Поддержка gRPC-web
