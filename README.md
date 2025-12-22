# Документация для запуска бэкенда Movies-service-back

## Быстрый старт

### 1. Настройка окружения
Создайте файл `.env` в корне проекта:

```env
POSTGRES_HOST=movies-db
POSTGRES_DB=movies
POSTGRES_USER=postgres
POSTGRES_PASSWORD=123
```

### 2. Запуск проекта
```bash
docker compose up --build -d
```

## Доступ к сервисам

- **API:** http://localhost:8080
- **PostgreSQL:** localhost:5432
- **Ollama:** localhost:11434

## Основные команды

```bash
# Запуск
docker compose up -d

# Запуск с пересборкой
docker compose up --build -d

# Остановка
docker compose down

# Просмотр логов
docker compose logs -f

# Статус контейнеров
docker compose ps
```

## Конфигурация

Измените настройки в `config.yml` при необходимости. Файл автоматически монтируется в контейнер.

**Примечание:** Файл `.env` содержит чувствительные данные, не коммитьте его в Git.