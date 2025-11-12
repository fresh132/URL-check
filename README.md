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

## Эндроинты

```go
POST /check // ручка для проверки доступности url
POST /report // ручка для формирования pdf отчета
```

```bash
# Проверка ссылок
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{"links": ["https://google.com", "https://httpstat.us/404"]}'

# Получение отчёта
curl -X POST http://localhost:8080/report \
  -H "Content-Type: application/json" \
  -d '{"links_list": ["1"]}' \
  --output report.pdf
```

```json
//Пример запросов к обоим эндпоинтам
//POST /check
{
  "links": ["https://google.com", "https://youtube.com"]
}
// Ответ
{
  "statuses": [
    {"url": "https://google.com", "status": "Available"},
    {"url": "youtube.com", "status": "Not available"}
  ],
  "id": 1
}

//GET /report
{
  "links_list": ["1", "2"] // или [""] тогда будет сформирован pdf отчет по всем ссылкам
}
// Ответ
// PDF file 