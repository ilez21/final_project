# TODO Планировщик задач

## Описание

Веб-сервер для управления задачами с поддержкой повторяющихся задач.

Функции:
- Добавление, редактирование и удаление задач
- Отметка задач как выполненных
- Правила повторения (ежедневно, еженедельно, ежемесячно, ежегодно)
- Поиск задач
- Аутентификация по паролю

## Выполненные задания со звёздочкой

1. Переменная окружения TODO_PORT
2. Переменная окружения TODO_DBFILE
3. Все правила повторения (w, m)
4. Поиск задач
5. Аутентификация
6. Docker образ

## Запуск

### Локально

```bash
go run main.go
```

Откройте http://localhost:7540/

### С параметрами

```bash
set TODO_PORT=8080
set TODO_DBFILE=./data/scheduler.db
set TODO_PASSWORD=mypassword
go run main.go
```

## Тесты

Настройте `tests/settings.go`:

```go
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ``
```

Запуск:

```bash
go test ./tests -v
```

Или используйте `run_test.bat` (Windows)

### Тесты с аутентификацией

1. Запустите сервер с паролем:
```bash
set TODO_PASSWORD=test123
go run main.go
```

2. Получите токен:
```bash
curl -X POST http://localhost:7540/api/signin -H "Content-Type: application/json" -d "{\"password\":\"test123\"}"
```

3. Вставьте токен в `tests/settings.go`

## Docker

### Сборка

```bash
docker build -t todo-scheduler .
```

### Запуск

```bash
docker run -d -p 7540:7540 --name todo todo-scheduler
```

### С БД на хосте (Windows)

```bash
docker run -d -p 7540:7540 -v %cd%:/data --name todo todo-scheduler
```

### С БД на хосте (Linux/Mac)

```bash
docker run -d -p 7540:7540 -v $(pwd):/data --name todo todo-scheduler
```

### С аутентификацией

```bash
docker run -d -p 7540:7540 -e TODO_PASSWORD=mypassword --name todo todo-scheduler
```

### Управление

```bash
docker logs todo
docker stop todo
docker start todo
docker rm -f todo
```

## API

### Публичные
- GET / - главная страница
- GET /login.html - вход
- POST /api/signin - аутентификация
- GET /api/nextdate - вычисление даты

### Защищенные
- POST /api/task - добавить задачу
- GET /api/task?id=<id> - получить задачу
- PUT /api/task - обновить задачу
- DELETE /api/task?id=<id> - удалить задачу
- GET /api/tasks?search=<query> - список задач
- POST /api/task/done?id=<id> - выполнить задачу

## Структура

```
final_sprint/
├── pkg/
│   ├── api/      # обработчики API
│   ├── db/       # работа с БД
│   └── server/   # веб-сервер
├── tests/        # тесты
├── web/          # фронтенд
├── main.go
├── Dockerfile
└── scheduler.db
```

## Технологии

- Go 1.21
- SQLite
- JWT
- Docker

## Правила повторения

- `d <число>` - через N дней (1-400)
- `y` - ежегодно
- `w <дни>` - дни недели (1=пн, 7=вс)
- `m <дни> [месяцы]` - дни месяца
