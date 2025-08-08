# Dater
Это аналог Airtable на `Go`, `TypeScript` и `SvelteKit`.

## Функционал
*Пока пусто :)*

## Технологии проекта
Это список технологий и библиотек/фреймворков, используемых в проекте:

### Технологии
| Технология  | Описание                                  |
| ----------- | ----------------------------------------- |
| Go          | Backend-язык программирования             |
| TypeScript  | Язык для типизированного JS               |
| PostgreSQL  | Реляционная база данных                   |
| RESTful API | Основной способ общения клиента и сервера |
| WebSocket   | Двусторонняя связь для real-time          |
| Vite        | Быстрый сборщик для фронтенда             |


## Фреймворки
| Фреймворк    | Описание                       |
| ------------ | ------------------------------ |
| Gin          | HTTP-фреймворк для Go          |
| SvelteKit    | Полноценный фронтенд-фреймворк |
| Tailwind CSS | Утилитарный CSS-фреймворк      |


## Библиотеки
| Библиотека      | Назначение                                      |
| --------------- | ----------------------------------------------- |
| Zod             | Схемы и валидация на фронтенде                  |
| axios           | Библиотека для HTTP-клиентов                    |
| GORM            | ORM для работы с PostgreSQL в Go                |
| zerolog         | Высокопроизводительный логгер для Go            |
| BurntSushi/toml | Парсер и генератор TOML-файлов                  |
| joho/godotenv   | Парсер dotenv-файлов                            |
| pgx             | PostgreSQL-драйвер с расширенными возможностями |
| cors            | Middleware для управления CORS в Gin            |
| bcrypt          | Хэширование паролей                             |
| jwt-go          | Работа с JWT (JSON Web Token)                   |

## TODO для MLP
- Функционал пользователей
- Функционал рабочих пространств (Workspace)
- Функционал баз данных (Base)
- Функционал таблиц (Table)
- Функционал строк (Row)
- Функционал полей (Field)
- WebSocket, RESTful API
- UX и UI дизайн

## Настройка сервера проекта
1. Создать файл `.env` в корне проекта и добавить в него переменные окружения:
   ```env
   DB_PASSWORD=<пароль_пользователя_бд>
   JWT_SECRET=<секретный_ключ_для_jwt>
   ```
2. Добавить конфигурационный файл `config.toml` в директорию `backend/config/`. Об конфигурационном файле можно прочитать [здесь](https://github.com/go-ozzo/ozzo-config).
2. Установить зависимости:
   ```bash
   go mod tidy
   ```
3. Запустить проект:
   ```bash
   go run main.go
   ```

field.type = 'single select', то cell.value = { "option_id": 17 }
field.type = 'multiselect', то cell.value = { "option_ids": [17, 22, 34] }
field.type = 'Many2One', то cell.value = { "relation_id": 1 }
field.type = 'One2Many', то cell.value = { "relation_ids": [1, 2, 3] }
field.type = 'Many2Many', то cell.value = { "relation_ids": [1, 2, 3] }

В domain хранятся модели для бд.
В repo находятся интерфейсы методов для работы с БД у моделей.
В services находятся реализации интерфейсов из repo. Будет содержать бизнес-логику.
В handlers находятся обработчики запросов и будут использовать сервисы.

```
core/
├── config
│   └── config.go
├── domain
│   ├── base.go
│   ├── cell.go
│   ├── field.go
│   ├── record.go
│   ├── record_relation.go
│   ├── select_option.go
│   └── table.go
├── repo
│   ├── interface.go
│   └── postgres
│       ├── base_repo.go
│       ├── cell_repo.go
│       ├── field_repo.go
│       ├── record_relation_repo.go
│       ├── record_repo.go
│       ├── select_option_repo.go
│       └── table_repo.go
├── server
│   ├── db.go
│   ├── logger.go
│   └── server.go
├── service
│   ├── base_service.go
│   ├── cell_service.go
│   ├── field_service.go
│   ├── record_relation_service.go
│   ├── record_service.go
│   ├── select_option_service.go
│   └── table_service.go
└── transport
    ├── api
    │   ├── handlers
    │   │   ├── base_handler.go
    │   │   ├── cell_handler.go
    │   │   ├── field_handler.go
    │   │   ├── record_handler.go
    │   │   ├── record_relation_handler.go
    │   │   ├── select_option_handler.go
    │   │   └── table_handler.go
    │   ├── middleware.go
    │   └── router.go
    └── ws
```