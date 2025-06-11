# Vira Logger

Пакет `vira-logger` предоставляет структурированный логгер с поддержкой цветного вывода, JSON-форматирования, ротации логов и контекстного логирования.

## Особенности

- Поддержка нескольких уровней логирования (DEBUG, INFO, WARN, ERROR, FATAL)
- Вывод в консоль с цветовой подсветкой или в файл
- Поддержка JSON-формата для структурированного логирования
- Ротация лог-файлов (по размеру, возрасту, сжатие)
- Контекстное логирование (request_id, user_id и др.)
- Информация о месте вызова (файл:строка)
- Потокобезопасность

## Установка

```bash
go get github.com/skrolikov/vira-logger
```

## Использование

### Инициализация логгера

```go
import "github.com/skrolikov/vira-logger/logger"

func main() {
    // Конфигурация логгера
    cfg := logger.Config{
        Level:      logger.INFO,
        JsonOutput: false,
        ShowCaller: true,
        Color:      true,
        OutputFile: "app.log",
        MaxSizeMB:  100,   // 100MB
        MaxBackups: 3,     // 3 backup файла
        MaxAgeDays: 30,    // хранить 30 дней
        Compress:   true,  // сжимать старые логи
    }

    log := logger.New(cfg)
    
    log.Info("Приложение запущено")
}
```

### Уровни логирования

```go
log.Debug("Отладочная информация")
log.Info("Информационное сообщение")
log.Warn("Предупреждение")
log.Error("Ошибка")
log.Fatal("Критическая ошибка, приложение завершится") // Вызывает os.Exit(1)
```

### Контекстное логирование

```go
// Добавление полей к логгеру
logWithFields := log.WithFields(map[string]any{
    "service": "auth",
    "version": "1.0",
})

logWithFields.Info("Пользователь аутентифицирован")

// Использование контекста
ctx := context.WithValue(context.Background(), "request_id", "abc123")
ctx = context.WithValue(ctx, "user_id", "user123")

logWithCtx := log.WithContext(ctx)
logWithCtx.Info("Запрос обработан")
```

## Конфигурация

Параметры `Config`:

- `Level` - минимальный уровень логирования (DEBUG, INFO, WARN, ERROR, FATAL)
- `JsonOutput` - вывод в JSON-формате (true/false)
- `ShowCaller` - показывать место вызова (файл:строка)
- `Color` - цветной вывод в консоль (только для не-JSON)
- `OutputFile` - путь к файлу для логирования (пустая строка = stdout)
- `MaxSizeMB` - максимальный размер файла перед ротацией (MB)
- `MaxBackups` - количество резервных копий
- `MaxAgeDays` - максимальный возраст файлов (дни)
- `Compress` - сжимать старые файлы (gzip)

## Формат вывода

### Текстовый формат (по умолчанию)

```
[INFO] 2023-10-01T15:04:05Z main.go:42 Приложение запущено | service=auth version=1.0 request_id=abc123
```

### JSON формат

```json
{
  "time": "2023-10-01T15:04:05Z",
  "level": "INFO",
  "message": "Приложение запущено",
  "caller": "main.go:42",
  "service": "auth",
  "version": "1.0",
  "request_id": "abc123"
}
```

## Лучшие практики

1. Для production используйте JSON-формат и файловый вывод
2. DEBUG уровень только для разработки
3. Добавляйте request_id для трассировки запросов
4. Настройте ротацию логов, чтобы избежать переполнения диска
5. Для микросервисов добавляйте идентификатор сервиса через WithFields

## Производительность

- Логгер использует sync.Mutex для потокобезопасности
- Форматирование сообщений происходит только если уровень логирования позволяет
- Для файлового вывода используется lumberjack с эффективной ротацией