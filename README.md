# Планировщик задач (TODO-лист)

Веб-сервер на Go, реализующий функциональность планировщика задач с поддержкой периодических задач и аутентификацией.

## Описание проекта

Приложение предоставляет API для управления задачами:
- Создание, редактирование и удаление задач с заголовком, комментарием и датой
- Настройка правил повторения: ежедневно (`d`), еженедельно (`w`), ежемесячно (`m`), ежегодно (`y`)
- Автоматический перенос периодических задач при отметке о выполнении
- Получение списка ближайших задач с сортировкой по дате
- Поиск задач по подстроке в заголовке/комментарии или по дате
- JWT-аутентификация для защиты доступа


## Выполненные задания со звёздочкой

Расширенные правила повторения** — реализованы `w` (дни недели 1-7) и `m` (дни месяца 1-31, -1, -2 с опциональными месяцами)  
Аутентификация** — JWT-токены, защита API endpoints через middleware  
Поиск — параметр `search` в `/api/tasks` (подстрока в title/comment или дата в формате DD.MM.YYYY)  
Docker — сборка (Dockerfile в репозитории)  
CI/CD — GitHub Actions workflow для автоматического запуска тестов при push

## Запуск локально

### Базовый запуск
```bash
go run .

Сервер запустится на http://localhost:7540

С аутентификацией
http://localhost:7540/login.html

Тесты:
go test ./tests

Отдельные тесты:
go test -run ^TestApp$ ./tests
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests

Docker:

Сборка
docker build -t go_final .

Запуск
docker run -d \
  -p 7540:7540 \
  -v $(pwd)/data:/data \
  -e TODO_PASSWORD=your_password \
  go_final

Структура проекта:

  go_final/
├── main.go                 # Точка входа
├── Dockerfile              # Сборка образа
├── go.mod                  # Зависимости
├── README.md              # Документация
├── web/                   # Фронтенд (HTML/CSS/JS)
├── pkg/
│   ├── api/              # HTTP handlers, JWT auth
│   └── db/               # Работа с SQLite
├── tests/                # Тесты
│   └── settings.go       # Конфигурация тестов
└── .github/
    └── workflows/
        └── go.yml        # CI/CD pipeline
