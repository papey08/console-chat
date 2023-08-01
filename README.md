# console-chat

## Описание

Данный проект представляет собой сервер, позволяющий вести переписку нескольких 
клиентов прямо в терминале.

## Структура проекта

```text
├── cmd
│   ├── client
│   │   └── main.go // точка входа в приложение
│   └── server
│       └── main.go // консольный клиент
│
├── configs
│   └── config.yml // файл с конфигами
│
├── internal
│   ├── app // слой бизнес-логики (usecase)
│   │   ├── valid // пакет для проверки валидности никнеймов и паролей
│   │   ├── app.go // реализация интерфейса приложения
│   │   └── app_interface.go // интерфейс приложения
│   │
│   ├── model // слой сущностей (entities)
│   │   ├── errs.go
│   │   └── user.go // структура пользователя
│   │
│   ├── ports // сетевой слой (infrastructure)
│   │   ├── ginserver // http-сервер 
│   │   └── wsserver // websocket сервер
│   │
│   └── repo // слой БД
│       └── user_repo // хранилище пользователей
│
├── migrations
│   └── user_repo_init.sql // скрипт для конфигурации user_repo
│
├── Dockerfile
├── README.md
├── docker-compose.yml
├── go.mod
└── go.sum

```

## Бизнес-логика

Новый пользователь должен зарегистрироваться на сервере. Регистрация происходит 
при вводе валидных ника и пароля. Пароли хранятся в базе данных в 
захэшированном виде, то есть админ сервера не будет иметь доступ к аккаунтам 
пользователей. Затем пользователь должен авторизоваться на сервере, введя свои 
ник и пароль, в ответ он получит jwt-токен. JWT-токен нужно будет прислать 
отдельной строкой в чат первым сообщением, из него websocket-сервер расшифрует 
имя пользователя, которым будет подписывать все последующие сообщения от этого 
пользователя в чате.

Также при регистрации данные пользователя попадают во временный кеш, чтобы при 
авторизации этого же пользователя сервер мог быстрее их получить.

## Используемые технологии

* go 1.20
* PostgreSQL — постоянное хранение пользователей
* Redis — временное хранение пользователей
* [Gin Web Framework](https://github.com/gin-gonic/gin)
* Websocket
* Docker

## Запуск сервера

### С помощью Docker

```shell
$ docker-compose up
```

### Локально

Самостоятельно сконфигурировать PostgreSQL (*[скрипт для конфигурации user_repo](https://github.com/papey08/console-chat/blob/master/migrations/user_repo_init.sql)*), 
изменить файл *[config.yml](https://github.com/papey08/console-chat/blob/master/configs/config.yml)*, после чего выполнить команды:

```shell
$ go mod download
$ go run cmd/server/main.go
```

## Запуск клиента

```shell
$ go mod download
$ go run cmd/client/main.go
```

При запуске *[cmd/client/main.go](https://github.com/papey08/console-chat/blob/master/cmd/client/main.go)* 
будет выведена документация. При запуске с флагом `-reg` клиент перейдёт к 
регистрации, при запуске с флагом `-sign` клиент перейдёт к авторизации и 
подключению к чату. Чтобы убедиться в работоспособности, запустите несколько 
клиентов.

## Формат запросов

### Регистрация

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/console-chat/users`
* Формат тела запроса:
```json
{
    "nickname": "papey08",
    "password": "qwerty_123"
}
```
* Формат ответа:
```json
{
    "data": {
        "nickname": "papey08",
        "hashed_password": "21d0c2b75fe758d93ab6dc4911712f5d5667a1d334a9afe92131473fe8c53b40"
    },
    "error": null
}
```

### Авторизация

* Метод: `GET`
* Эндпоинт: `http://localhost:8080/console-chat/users/papey08`
* Формат тела запроса:
```json
{
    "password": "qwerty_123"
}
```
* Формат ответа:
```json
{
    "data": {
        "token_string": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4NDE0NzcsIm5pY2tuYW1lIjoicGFwZXkwOCJ9.h0Gh2bHYuqCkwm26vavTRuT-0qbr6olU6Q-50_yzzLM"
    },
    "error": null
}
```

### Чат

* Адрес: `ws://localhost:8080/console-chat/chat`
* Первое сообщение — полученный токен: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4NDE0NzcsIm5pY2tuYW1lIjoicGFwZXkwOCJ9.h0Gh2bHYuqCkwm26vavTRuT-0qbr6olU6Q-50_yzzLM`
