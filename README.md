# Lazy Investor

Торговый советник с приоритетом безопасности: детерминированное ядро (данные, фичи, сигналы, риск) + LLM как советник + подтверждение интентов в Telegram. Подробная спецификация — [SPEC.md](SPEC.md).

## Запуск

```powershell
copy .env.example .env   # задать APP_NAME, TINKOFF_API_HOST, TINKOFF_API_TOKEN; опционально TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID
make run                 # или: go run ./cmd/analyst/main.go
```

**Stage 0**: при старте и по расписанию (интервал в секундах в `config.yml` → `scheduler.intervalSeconds`) приложение загружает дневные свечи по инструментам из конфига, считает метрики (доходность, волатильность, max drawdown), формирует текстовый отчёт и при включённом Telegram отправляет его в чат (и опционально PNG-график). Без LLM и интентов. При отсутствии `TELEGRAM_BOT_TOKEN` или `TELEGRAM_CHAT_ID` отчёт только пишется в лог.

## Технические детали

- **Go**, структура по [SPEC.md](SPEC.md): `cmd/`, `internal/ports` (интерфейсы), `internal/adapters` (Tinkoff, chart, telegram), `internal/application` (Stage 0 pipeline, metrics, report), `pkg` (config, logging, wire).
- **Рынок**: адаптер Tinkoff Invest API за портом `MarketDataProvider` (свечи, поиск инструмента). Контракт данных — `internal/dto` (Candle, Instrument).
- **Telegram**: порт `TelegramNotifier`, адаптер на [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api); при пустых env — no-op. **Chat ID**: для личной переписки с ботом укажите свой числовой user ID (узнать: @userinfobot или @getmyid_bot); для отправки в группу — ID группы (отрицательное число). В обоих случаях бот просто шлёт сообщения в указанный чат.
- **DI**: Google Wire в `pkg/wire`; пересборка: `cd pkg/wire && go run -mod=mod github.com/google/wire/cmd/wire .` или `make wire` (если в Makefile добавлена эта команда).
