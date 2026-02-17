# Lazy Investor

Торговый советник с приоритетом безопасности: детерминированное ядро (данные, фичи, сигналы, риск) + LLM как советник + подтверждение интентов в Telegram. Подробная спецификация — [SPEC.md](SPEC.md).

## Запуск

```bash
cp .env.example .env   # задать APP_NAME, TINKOFF_API_HOST, TINKOFF_API_TOKEN
make run               # или: go run ./cmd/analyst/main.go
```

Сейчас приложение поднимает соединение с Tinkoff API и ждёт graceful shutdown (SIGINT/SIGTERM). Пайплайн (сбор данных → фичи → LLM → Telegram) в разработке.

## Технические детали

- **Go**, структура по [SPEC.md](SPEC.md): `cmd/`, `internal/ports` (интерфейсы), `internal/adapters` (Tinkoff, chart), `internal/application`, `pkg` (config, logging, wire).
- **Рынок**: адаптер Tinkoff Invest API за портом `MarketDataProvider` (свечи, поиск инструмента). Контракт данных — `internal/dto` (Candle, Instrument).
- **DI**: Google Wire в `pkg/wire`; пересборка: `make wire`.
