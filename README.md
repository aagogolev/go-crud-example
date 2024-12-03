# Go CRUD Example

Современный REST API сервис на Go с использованием чистой архитектуры.

## Возможности

- CRUD операции для управления пользователями
- Чистая архитектура (handlers, services, repositories)
- PostgreSQL база данных
- Swagger документация
- Валидация данных
- Middleware для логирования
- Модульные тесты
- Docker поддержка

## Структура проекта
```
go-crud-example/
├── cmd/                      # Точки входа в приложение
│   └── api/                  
│       └── main.go          # Главный файл приложения
│
├── internal/                 # Внутренний код приложения
│   ├── handler/             # HTTP обработчики
│   │   ├── user_handler.go
│   │   └── handler.go
│   │
│   ├── service/             # Бизнес-логика
│   │   ├── user_service.go
│   │   └── service.go
│   │
│   ├── repository/          # Работа с базой данных
│   │   ├── user_repository.go
│   │   └── repository.go
│   │
│   └── model/              # Модели данных
│       └── user.go
│
├── pkg/                    # Общие пакеты
│   ├── config/            # Конфигурация
│   │   └── config.go
│
├── tests/                 # Тесты
│   ├── unit/             # Модульные тесты
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│
├── docs/                 # Документация
│
├── .env.example        # Пример переменных окружения
├── .gitignore         # Игнорируемые Git файлы
├── docker-compose.yml # Docker Compose конфигурация
├── go.mod            # Go модули
├── go.sum           # Зависимости Go модулей
└── README.md        # Документация проекта
```

## Установка

1. Клонировать репозиторий:
```bash
git clone https://github.com/yourusername/go-crud-example.git
```
2. Установить зависимости:
```bash
go mod download
```
3. Создать файл .env и заполнить переменные окружения (см. .env.example)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=users_db
SERVER_PORT=8000
```
4. Запустить PostgreSQL:
```bash
docker-compose up -d postgres
```
5. Запустить приложение:

```bash
go go-crud-example/run cmd/api/main.go
```

