# Lazy Investor

Торговый советник с приоритетом безопасности: детерминированное ядро (данные, фичи, сигналы, риск) + LLM как советник + подтверждение интентов в Telegram. Подробная спецификация — [SPEC.md](SPEC.md).

## Запуск

```powershell
copy .env.example .env   # задать APP_NAME, TINKOFF_API_HOST, TINKOFF_API_TOKEN; опционально TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID
make run                 # или: go run ./cmd/analyst/main.go
```

**Stage 0**: при старте и по расписанию (интервал в секундах в `config.yml` → `scheduler.intervalSeconds`) приложение загружает дневные свечи по инструментам из конфига, считает метрики (доходность, волатильность, max drawdown), формирует текстовый отчёт и при включённом Telegram отправляет его в чат (и опционально PNG-график). Без LLM и интентов. При отсутствии `TELEGRAM_BOT_TOKEN` или `TELEGRAM_CHAT_ID` отчёт только пишется в лог.

## Выгрузка свечей в CSV

Одна и та же логика доступна из **CLI** и из **Telegram-бота**.

### CLI (cmd/export-candles)

Отдельный бинарник: инструмент (ISIN или тикер), таймфрейм и период → CSV в stdout или файл.

```powershell
go run ./cmd/export-candles -instrument SBER -timeframe 1d -last 30
go run ./cmd/export-candles -instrument RU0009029540 -timeframe 1d -from 2024-01-01 -to 2024-12-31 -output candles.csv
```

Флаги: `-instrument` (обязательный), `-timeframe` (1m, 5m, 15m, 1h, 1d, 1w, 1month), период: `-last N` (последние N дней) или `-from` и `-to` (YYYY-MM-DD), `-output` (путь к файлу; по умолчанию stdout). Требуется `.env` с Tinkoff API (как для основного приложения).

### Telegram: команда /candles

Если в `config.yml` включено `telegram.handleCommands: true`, приложение в фоне слушает входящие сообщения и обрабатывает команду `/candles`. Пользователь пишет, например:

```
/candles SBER 1d 30
```

Бот присылает в ответ файл CSV со свечами (инструмент SBER, таймфрейм 1d, последние 30 дней). Опционально задайте `telegram.allowedChatID` (числовой ID чата), чтобы на команды отвечал только этот чат.

## Технические детали

- **Go**, структура по [SPEC.md](SPEC.md): `cmd/`, `internal/ports` (интерфейсы), `internal/adapters` (Tinkoff, chart, telegram), `internal/application` (Stage 0 pipeline, metrics, report), `pkg` (config, logging, wire).
- **Рынок**: адаптер Tinkoff Invest API за портом `MarketDataProvider` (свечи, поиск инструмента). Контракт данных — `internal/dto` (Candle, Instrument).
- **Telegram**: порт `TelegramNotifier`, адаптер на [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api); при пустых env — no-op. **Chat ID**: для личной переписки с ботом укажите свой числовой user ID (узнать: @userinfobot или @getmyid_bot); для отправки в группу — ID группы (отрицательное число). В обоих случаях бот просто шлёт сообщения в указанный чат.
- **DI**: Google Wire в `pkg/wire`; пересборка: `cd pkg/wire && go run -mod=mod github.com/google/wire/cmd/wire .` или `make wire` (если в Makefile добавлена эта команда).
