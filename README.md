## Стек проекта:
- Golang 1.22.4
- Postgres 16.0
- Cleanenv 1.5
- Chi/v5 5.1.0
## Процесс запуска проекта (через docker compose):
### Запуск проекта:
переименовать файл:

.env.example -> .env

КОНСОЛЬ!!! находясь в корневой директории:

docker compose up
  
проект доступен по endpoint:

POST

- http://localhost:8080/signup

- http://localhost:8080/signin

- http://localhost:8080/create (необходим JWT токен)

GET

- http://localhost:8080/notes (необходим JWT токен)


### Примеры запросов можно получить при импорте Postman Collection в корне проекта
