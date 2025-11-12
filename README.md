# URL checker

Программа проверяет доступность url ссылок

## Что уже реализовано

- Проверка доступности url
- Формирование pdf отчетов
- Unit тесты

## Использование

```bash
cd /URL-check

#Запуск тестов
go test ./...

#Запуск программы
go run cmd/main.go

```

## Необходимо создать .env файл с переменными
```env
PORT="8080"
GIN_MODE="release" или debug
```

```go
//Структура
cmd/
 └── main.go          — точка входа, запуск сервера
internal/
 ├── api/             — обработчики HTTP-запросов
 ├── config/          — сохранение/загрузка данных и проверки URL
 └── struct.go        — общие структуры и переменные
.env                  — настройки окружения
.gitignore
```

## Описание примянемых решений в коде

Для хранения данных была использована мапа и для постоянного хранения data.json 
При завершении работы программы, данные сохраняются в файл data.json

Graceful Shutdown работает по принципу: закрывает доступ к новым запросам и ожидает завершения текущих в течении 60 секунд

## Эндпоинты

```go
POST /check // ручка для проверки доступности url
GET /report // ручка для формирования pdf отчета
```

```bash
# По конкретным ID
curl "http://localhost:8080/report?id=1&id=2" --output report.pdf

# Все проверки
curl "http://localhost:8080/report" --output all.pdf
```

Пример запросов и ответов к эндпоинту
POST /check

Запрос
```json
{
  "links": ["https://google.com", "https://youtube.com"]
}
```

Ответ
```json
{
  "statuses": [
    {"url": "https://google.com", "status": "Available"},
    {"url": "youtube.com", "status": "Not available"}
  ],
  "id": 1
}
```