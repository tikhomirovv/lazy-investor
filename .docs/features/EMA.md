# EMA

- [ ] EMA features: добавить checkpoints (gap/spread) + meta для LLM, улучшить события тренда
  - [ ] Checkpoints (фиксированный шаг по барам, последние N точек):
    - [ ] For each checkpoint store: `barsAgo`, `gap`, `spread`
    - [ ] Definitions (via meta.rules, not schema/defs):
      - [ ] `gap = close - EMA(fastPeriod)`
      - [ ] `spread = EMA(fastPeriod) - EMA(slowPeriod)`
  - [ ] Meta (внутри EMA feature, едет вместе с данными в snapshot/LLM):
    - [ ] `interval` (e.g. "1d")
    - [ ] `asOfCandleTime` (time of last closed candle)
    - [ ] `fastPeriod`, `slowPeriod`
    - [ ] `description` (super short plain text)
    - [ ] `rules[]` (short interpretation rules for gap/spread sign and |spread| dynamics)
    - [ ] Example output format (JSON):

          {
            "meta": {
              "interval": "1d",
              "asOfCandleTime": "2026-02-18T00:00:00Z",
              "fastPeriod": 20,
              "slowPeriod": 100,
              "description": "EMA checkpoints for trend regime and strengthening/weakening.",
              "rules": [
                "gap = close - EMA(fastPeriod). gap>0 => price above EMA(fast); gap<0 => below.",
                "spread = EMA(fastPeriod) - EMA(slowPeriod). spread>0 => up-regime; spread<0 => down-regime.",
                "|spread| rising across checkpoints => trend strengthening; falling => weakening."
              ]
            },
            "checkpoints": [
              { "barsAgo": 0, "gap": 1.04, "spread": 3.12 },
              { "barsAgo": 5, "gap": 0.62, "spread": 2.71 }
            ]
          }
  - [ ] Detect crossings/events (leave deterministic):
    - [ ] price vs EMA(fast): `price_crossed_up_*` / `price_crossed_down_*` (verify semantics)
    - [ ] EMA(fast) vs EMA(slow): `ema_fast_crossed_above_slow` / `ema_fast_crossed_below_slow` (or keep current 20/100 names)
  - [ ] Update `.docs/FEATURES.md` spec and add tests (synthetic series for crossings + checkpoints dynamics)
