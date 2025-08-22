### USER
торговый бот на go может работать в реальном режиме собирая данные с биржи или в режиме бектестинга, выбирая данные из базы данных. Данные свечей проходят ряд этапов. На первом этапе данные кластеризуются в более крупный интервал, на следующем этапе крастеризованные данные группируются в серии, дальше серии анализируются по стратегии, в результате анализа формируются сигналы на покупку/продажи. Сигналы проверяются по набору заданных правил, включающих проверку технических показателей, PNL текущих позиций, текущего баланса и т.п.
Какие подходы программирования лучше всего использовать для этого, как организовать структуру обработки этапов и проверки правил.

### ASSISTANT
Для реализации торгового бота на Go, который может работать как в **реальном времени**, так и в **режиме бектестинга**, с последовательной обработкой данных через несколько этапов (кластеризация, группировка в серии, анализ по стратегии, формирование сигналов, проверка правил), важно использовать **гибкую, модульную и тестируемую архитектуру**. Ниже описаны подходы и структура, которые помогут реализовать такую систему эффективно.

---

### 🔧 **1. Подходы программирования**

#### ✅ **1.1. Реактивное программирование (Reactive Programming)**
- Используйте **каналы (channels)** и **горутины (goroutines)** для асинхронной обработки данных.
- Данные (свечи) "текут" через конвейер этапов.
- Подходит для реального времени: новые свечи поступают → обрабатываются → проходят этапы → формируются сигналы.

#### ✅ **1.2. Паттерн "Конвейер" (Pipeline Pattern)**
- Каждый этап — это **стадия конвейера**, принимающая данные, обрабатывающая их и передающая дальше.
- Пример:
  ```go
  candles -> Clusterizer -> Series -> Strategy -> Signal -> RuleValidator -> Execution
  ```
- Каждый этап — отдельная функция или структура с интерфейсом.

#### ✅ **1.3. Паттерн "Стратегия" (Strategy Pattern)**
- Разные стратегии анализа (например, RSI, MACD, трендовые модели) реализуются как **отдельные модули**, удовлетворяющие общему интерфейсу.
- Можно легко переключаться между стратегиями.

#### ✅ **1.4. Паттерн "Цепочка ответственности" (Chain of Responsibility)**
- Проверка правил (на баланс, PNL, технические индикаторы и т.п.) может быть реализована как цепочка валидаторов.
- Каждый валидатор решает: пропустить сигнал дальше или заблокировать.

#### ✅ **1.5. Инверсия управления (IoC) / Dependency Injection**
- Конфигурация и зависимости (например, источник данных, стратегия, валидаторы) внедряются извне.
- Упрощает тестирование и переключение между режимами (live / backtest).

---

### 🏗️ **2. Архитектура и структура проекта**

```
/pkg/
  /data/
    source.go          // Интерфейс: Source (реалтайм / бектест)
    candle.go          // Структура свечи
  /processing/
    clusterizer.go     // Кластеризация свечей
    series.go          // Группировка в серии
    strategy.go        // Интерфейс стратегии
    signal.go          // Структура сигнала
    validator.go       // Интерфейс валидатора правил
  /indicators/
    rsi.go, macd.go    // Технические индикаторы
  /execution/
    simulator.go       // Для бектеста
    broker.go          // Для реальной торговли
  /storage/
    db.go              // Работа с БД (для бектеста)
/main.go
/cmd/
  /live/
  /backtest/
```

---

### 🔄 **3. Организация обработки этапов (конвейер)**

#### Пример конвейера:
```go
type Pipeline struct {
    Source      DataProvider
    Clusterizer *Clusterizer
    SeriesMaker *SeriesMaker
    Strategy    Strategy
    Validators  []Validator
    Executor    Executor
}

func (p *Pipeline) Run() {
    for candle := range p.Source.Stream() {
        // Этап 1: Кластеризация
        p.Clusterizer.AddCandle(candle)
        if !p.Clusterizer.IsReady() {
            continue
        }
        aggregated := p.Clusterizer.GetAggregated()

        // Этап 2: Группировка в серии
        series := p.SeriesMaker.Add(aggregated)
        if !series.IsComplete() {
            continue
        }

        // Этап 3: Анализ по стратегии
        signal := p.Strategy.Analyze(series)
        if signal == nil {
            continue
        }

        // Этап 4: Проверка правил
        if p.validateSignal(signal) {
            p.Executor.Execute(signal)
        }
    }
}

func (p *Pipeline) validateSignal(signal *Signal) bool {
    for _, v := range p.Validators {
        if !v.Validate(signal) {
            return false
        }
    }
    return true
}
```

---

### 📊 **4. Режимы: Реальный и Бектест**

#### ✅ Общий интерфейс источника данных:
```go
type DataProvider interface {
    Stream() <-chan *Candle
}
```

- **RealTimeProvider**: подключается к API биржи (Binance, Bybit и т.п.)
- **BacktestProvider**: читает свечи из БД или файла (CSV, Parquet), эмулирует временные метки

#### 🔄 Переключение режима:
```go
var source DataProvider
if mode == "live" {
    source = NewRealTimeSource(apiKey, secret)
} else {
    source = NewBacktestSource(db, symbol, start, end)
}
```

---

### 📌 **5. Этапы обработки подробно**

| Этап | Описание | Реализация |
|------|--------|-----------|
| **1. Кластеризация** | Агрегация свечей (например, 1m → 5m) | `Clusterizer` накапливает N свечей, выдает агрегированную |
| **2. Группировка в серии** | Логическая группировка (например, тренд, флэт, коррекция) | `SeriesMaker` использует правила или ML-модель |
| **3. Стратегия** | Анализ серии → сигнал | Интерфейс `Strategy.Analyze(series) *Signal` |
| **4. Валидация сигнала** | Проверка по правилам | Цепочка `Validator`-ов: `BalanceValidator`, `PNLValidator`, `IndicatorValidator` |
| **5. Исполнение** | Отправка ордера или запись в отчет | `Executor` — симулятор или реальный брокер |

---

### ✅ **6. Проверка правил (валидаторы)**

```go
type Validator interface {
    Validate(*Signal) bool
}

type BalanceValidator struct {
    minBalance float64
}

func (v *BalanceValidator) Validate(s *Signal) bool {
    return getCurrentBalance() >= v.minBalance
}

type PNLValidator struct {
    maxDrawdown float64
}

func (v *PNLValidator) Validate(s *Signal) bool {
    return getPNL() > -v.maxDrawdown
}
```

> Можно собирать валидаторы в список, легко добавлять/удалять.

---

### 🧪 **7. Тестирование и гибкость**

- Все этапы тестируются отдельно (unit-тесты).
- Используйте **интерфейсы** для замены компонентов (например, `Strategy`, `DataProvider`).
- Для бектеста можно эмулировать задержки, комиссии, проскальзывание.

---

### 🚀 **8. Производительность и масштабируемость**

- Используйте **буферизованные каналы** для снижения блокировок.
- Храните временные данные в памяти (не в БД на каждом шаге).
- Для бектеста — **параллельный прогон** по нескольким активам/стратегиям.
- Используйте `sync.Pool` для переиспользования объектов (свечи, сигналы).

---

### 📦 **9. Пример интерфейсов**

```go
type Strategy interface {
    Analyze(series *Series) *Signal
}

type DataProvider interface {
    Stream() <-chan *Candle
}

type Executor interface {
    Execute(*Signal)
}
```

---

### 📝 **Вывод: Рекомендации**

| Аспект | Рекомендация |
|-------|-------------|
| **Архитектура** | Конвейер + интерфейсы + DI |
| **Паттерны** | Pipeline, Strategy, Chain of Responsibility |
| **Режимы** | Единый интерфейс `DataProvider` |
| **Тестирование** | Моки для валидаторов, стратегий, источников |
| **Производительность** | Горутины + каналы + кэширование |
| **Гибкость** | Конфигурация через JSON/YAML (стратегии, валидаторы, параметры) |

---

### 🧩 Дополнительно

- Рассмотрите использование **TDD**: сначала напишите тесты для каждого этапа.
- Для хранения исторических данных — **TimescaleDB** (PostgreSQL для временных рядов) или **InfluxDB**.
- Логирование: используйте `zap` или `logrus` для отладки этапов.

---

Если хочешь — могу помочь с **примером реализации конвейера** или **шаблоном проекта на Go**.

### USER
К примеру добавить этап получения на каждом шаге данных о PNL текущих позиций и проверку превышения допустимых значений. 

### ASSISTANT
Отличный вопрос! Добавление **мониторинга PNL (прибыли/убытка) текущих позиций** на каждом шаге — критически важный элемент управления рисками. Ниже я покажу, **как интегрировать этот этап в существующую архитектуру**, где PNL проверяется **на каждом шаге конвейера**, а не только при формировании сигнала.

---

## ✅ Цель
На **каждом шаге обработки** (после получения свечи, после кластеризации, перед формированием сигнала) — получать актуальные данные о:
- Открытых позициях
- Текущем PNL
- Проверять: не превышает ли PNL/Drawdown допустимые пороги

Если PNL ниже порога — **приостанавливать торговлю**, **блокировать новые сигналы**, **возможно, закрывать позиции**.

---

## 🧩 Где встраивать PNL-мониторинг?

Можно реализовать двумя способами:

### Вариант 1: **PNL как часть валидации сигнала** (простой)
→ Проверяется только при формировании сигнала.

### Вариант 2: **PNL-мониторинг как фоновый процесс + состояние в контексте** (рекомендуется)
→ Актуальный PNL доступен на **любом этапе** конвейера.
→ Можно реагировать **до** формирования сигнала (например, при сильном убытке — вообще остановить обработку).

---

## 🛠 Реализация: PNL-мониторинг в Go

### 1. Структуры

```go
// pkg/position/position.go
type Position struct {
    Symbol    string
    Side      string  // "long", "short"
    EntryPrice float64
    Size      float64
    Leverage  float64
}

type PNL struct {
    Unrealized float64 // текущий unrealized PNL
    Realized   float64
    Drawdown   float64 // максимальная просадка
    Timestamp  time.Time
}
```

---

### 2. Интерфейс PNLProvider

```go
// pkg/risk/pnl.go
type PNLProvider interface {
    GetCurrentPNL() (*PNL, error)
    GetOpenPositions() ([]*Position, error)
}
```

---

### 3. Реализация (например, для Binance)

```go
type BinancePNLProvider struct {
    client *binance.Client
}

func (p *BinancePNLProvider) GetCurrentPNL() (*PNL, error) {
    positions, err := p.GetOpenPositions()
    if err != nil {
        return nil, err
    }

    var totalPNL float64
    for _, pos := range positions {
        markPrice := getMarkPrice(pos.Symbol) // можно через API
        var pnl float64
        if pos.Side == "long" {
            pnl = (markPrice - pos.EntryPrice) * pos.Size
        } else {
            pnl = (pos.EntryPrice - markPrice) * pos.Size
        }
        totalPNL += pnl
    }

    return &PNL{
        Unrealized: totalPNL,
        Timestamp: time.Now(),
    }, nil
}
```

---

### 4. Интеграция в Pipeline

Добавим `PNLProvider` в конвейер и будем **проверять PNL на каждом шаге**.

```go
// pkg/pipeline/pipeline.go
type Pipeline struct {
    Source       DataProvider
    PNLProvider  PNLProvider
    Clusterizer  *Clusterizer
    SeriesMaker  *SeriesMaker
    Strategy     Strategy
    Validators   []Validator
    Executor     Executor

    // Риск-контроль
    MaxDrawdown float64
    TradingEnabled bool
}
```

---

### 5. Периодическое обновление PNL (фоном)

```go
func (p *Pipeline) startPNLMonitor() {
    ticker := time.NewTicker(10 * time.Second) // каждые 10 сек
    go func() {
        for range ticker.C {
            pnl, err := p.PNLProvider.GetCurrentPNL()
            if err != nil {
                log.Printf("PNL error: %v", err)
                continue
            }

            // Проверка просадки
            if pnl.Unrealized < -p.MaxDrawdown {
                p.TradingEnabled = false
                log.Printf("🛑 STOP TRADING: Drawdown exceeded! PNL=%.2f, MaxDrawdown=%.2f", pnl.Unrealized, p.MaxDrawdown)
            } else {
                p.TradingEnabled = true // можно снова торговать
            }
        }
    }()
}
```

Вызови это в `Run()`:

```go
func (p *Pipeline) Run() {
    p.startPNLMonitor() // запускаем мониторинг

    for candle := range p.Source.Stream() {
        // 🔒 Проверка разрешения на торговлю
        if !p.TradingEnabled {
            log.Println("Trading paused due to drawdown")
            continue
        }

        // ... остальная обработка: кластеризация, стратегия и т.д.
    }
}
```

---

### 6. Дополнительно: Валидатор PNL

Можно добавить **валидатор**, который проверяет PNL **перед исполнением сигнала**:

```go
type PNLValidator struct {
    pnlProvider PNLProvider
    maxDrawdown float64
}

func (v *PNLValidator) Validate(signal *Signal) bool {
    pnl, err := v.pnlProvider.GetCurrentPNL()
    if err != nil {
        return false
    }
    return pnl.Unrealized >= -v.maxDrawdown
}
```

Добавить в цепочку:
```go
pipeline.Validators = []Validator{
    &BalanceValidator{min: 100},
    &PNLValidator{pnlProvider: provider, maxDrawdown: 500},
    &IndicatorValidator{...},
}
```

---

## 🔄 Как это работает в потоке

```
[Свеча] → Проверка: TradingEnabled? → Кластеризация → Серия → Стратегия → Сигнал → Валидаторы (включая PNL) → Исполнение
                     ↑
             (фоновый PNL-монитор)
```

- Если PNL упал ниже порога → `TradingEnabled = false`
- Все новые свечи **игнорируются** до восстановления
- Можно добавить **автоматическое закрытие позиций** при критической просадке

---

## 💡 Дополнительные улучшения

| Функция | Описание |
|-------|--------|
| **PNL по символу** | Если торгуем несколько пар — отслеживать PNL по каждой |
| **Теплый рестарт** | Сохранять состояние PNL/позиций при перезапуске |
| **Уведомления** | Отправлять алерт в Telegram при большом убытке |
| **Auto-liquidation** | При критическом PNL — отправлять market-ордер на закрытие |

---

## ✅ Пример: Автоматическое закрытие при критической просадке

```go
if pnl.Unrealized < -p.MaxDrawdown*0.8 {
    log.Println("⚠️  High drawdown — closing positions...")
    p.Executor.ForceCloseAll()
}
```

---

## 📦 Вывод

Чтобы **на каждом шаге** получать PNL и проверять риски:

1. ✅ Реализуй `PNLProvider` с интерфейсом
2. ✅ Запусти **фоновую горутину**, обновляющую PNL
3. ✅ Храни флаг `TradingEnabled` в `Pipeline`
4. ✅ Проверяй этот флаг **на входе в конвейер**
5. ✅ Добавь `PNLValidator` в цепочку валидации сигналов
6. ✅ (Опционально) Реализуй автоматическое закрытие позиций

---

Если хочешь — могу скинуть **полный пример на GitHub-подобной структуре** или реализовать **мок для бектеста с симуляцией PNL**.

### USER
Как в режиме бектестинга синхронизировать функцию получения pnl с торговыми данными

### ASSISTANT
Отличный и очень важный вопрос! В **режиме бектестинга** синхронизация **PNL с торговыми данными** критична: PNL должен обновляться **точно в моменты сделок** и **соответствовать временным меткам свечей**, чтобы избежать **look-ahead bias** (просмотра в будущее) и обеспечить реалистичный расчёт.

---

## 🎯 Цель
Синхронизировать обновление PNL с:
- Временем поступления свечей
- Моментами входа/выхода из позиций
- Точным расчётом unrealized PNL на каждом шаге

---

## 🔧 Проблемы при бектесте

1. **Реальный PNL зависит от позиций**, а позиции — от сигналов.
2. Если PNL обновляется асинхронно (как в live), будет **рассинхрон**.
3. Нельзя использовать "текущее время" — нужно **двигаться по историческим данным шаг за шагом**.

---

## ✅ Решение: PNL-менеджер, привязанный к таймстепу

### Основная идея:
> PNL пересчитывается **на каждом шаге конвейера**, на основе **исторических цен** и **истории ордеров**, синхронно с поступлением свечей.

---

## 🏗️ Архитектура для бектеста

```
[Свеча t] → Обновить PNL (по цене t) → Проверить риски → Сформировать сигнал → Исполнить → Обновить позиции
```

---

## 1. Структура BacktestPNLProvider

```go
// pkg/risk/backtest_pnl.go

type Trade struct {
    Time      time.Time
    Symbol    string
    Side      string  // "buy", "sell"
    Price     float64
    Size      float64
    IsClose   bool
}

type BacktestPNLProvider struct {
    trades       []Trade
    positions    map[string]*Position  // открытые позиции
    history      []PNLSnapshot         // история PNL по времени
    currentTime  time.Time
    currentPrice float64               // текущая цена (из свечи)
}
```

---

## 2. Метод обновления PNL на шаге

```go
func (p *BacktestPNLProvider) UpdateAt(candle *Candle) {
    p.currentTime = candle.Time
    p.currentPrice = candle.Close

    // Пересчитываем все открытые позиции
    var totalUnrealized float64
    for _, pos := range p.positions {
        var pnl float64
        if pos.Side == "long" {
            pnl = (p.currentPrice - pos.EntryPrice) * pos.Size
        } else {
            pnl = (pos.EntryPrice - p.currentPrice) * pos.Size
        }
        totalUnrealized += pnl
    }

    // Сохраняем снимок
    snapshot := PNLSnapshot{
        Time:       p.currentTime,
        Unrealized: totalUnrealized,
        Positions:  copyPositions(p.positions),
    }
    p.history = append(p.history, snapshot)
}
```

---

## 3. Интеграция в конвейер бектеста

```go
func (p *Pipeline) RunBacktest() {
    for candle := range p.Source.Stream() {
        // 1. Обновляем PNL по текущей свече
        p.PNLProvider.UpdateAt(candle)

        // 2. Проверяем, разрешена ли торговля
        pnl, _ := p.PNLProvider.GetCurrentPNL()
        if pnl.Unrealized < -p.MaxDrawdown {
            p.TradingEnabled = false
        } else {
            p.TradingEnabled = true
        }

        if !p.TradingEnabled {
            continue
        }

        // 3. Кластеризация и дальше по цепочке
        p.Clusterizer.AddCandle(candle)
        if p.Clusterizer.IsReady() {
            aggregated := p.Clusterizer.GetAggregated()

            series := p.SeriesMaker.Add(aggregated)
            if series.IsComplete() {
                signal := p.Strategy.Analyze(series)
                if signal != nil && p.validateSignal(signal) {
                    // 4. Исполнение изменит позиции → повлияет на будущий PNL
                    p.Executor.Execute(signal)
                    p.PNLProvider.OnTrade(signal.ToTrade()) // обновить позиции
                }
            }
        }
    }
}
```

> 🔁 Важно: `UpdateAt(candle)` вызывается **до** формирования сигнала — чтобы проверить риск **на основе данных "до" решения**.

---

## 4. Executor обновляет позиции

```go
// pkg/execution/simulator.go
type BacktestExecutor struct {
    pnlProvider PNLUpdater
    balance     float64
}

func (e *BacktestExecutor) Execute(signal *Signal) {
    trade := e.executeTrade(signal)
    e.pnlProvider.OnTrade(trade) // передаём изменение позиции
}
```

```go
func (p *BacktestPNLProvider) OnTrade(trade Trade) {
    // Обновляем позиции: открываем, закрываем, увеличиваем
    p.updatePosition(trade)
}
```

---

## 5. Ключевые принципы синхронизации

| Принцип | Как реализован |
|-------|----------------|
| **Нет look-ahead** | PNL считается только по ценам и позициям **до текущего времени** |
| **Шаг за шагом** | Каждая свеча → обновление PNL → принятие решения |
| **История позиций** | Все сделки хранятся, PNL пересчитывается корректно |
| **Цена для PNL** | Берётся из **текущей свечи** (например, Close или Mark) |

---

## 6. Дополнительно: режимы расчёта PNL

| Режим | Как считать |
|------|-----------|
| **Unrealized PNL** | (текущая цена - цена входа) × объём |
| **Realized PNL** | При закрытии: (цена закрытия - цена входа) × объём закрытой части |
| **Floating PNL** | То же, что unrealized — меняется с каждой свечой |

---

## 7. Пример временной шкалы

```
t=10:00 | Свеча [10:00] → PNL = 0 → Сигнал → Куплено 1 BTC по 50k
t=10:01 | Свеча [10:01] → PNL = (50.2k - 50k) × 1 = +200 → Проверка риска → OK
t=10:02 | Свеча [10:02] → PNL = (49.5k - 50k) × 1 = -500 → Превышен drawdown → Остановка
```

---

## ✅ Рекомендации

1. **Не использовать горутины для PNL в бектесте** — всё должно быть синхронным.
2. **Сохраняй историю PNL** — потом можно построить график.
3. **Используй одну временную шкалу** — все компоненты синхронизируются по `candle.Time`.
4. **Мокай PNLProvider** в unit-тестах, чтобы проверять поведение при убытках.

---

## 📦 Вывод

Чтобы синхронизировать PNL с торговыми данными в бектесте:

- ✅ PNL пересчитывается **на каждом шаге** при получении свечи
- ✅ Используется **текущая цена из свечи**
- ✅ PNLProvider **обновляется синхронно**, без задержек
- ✅ Все сделки фиксируются, и позиции обновляются сразу
- ✅ Проверка рисков происходит **до** формирования сигнала

---

Если хочешь, могу скинуть **рабочий пример на Go** с:
- симуляцией 100 свечей
- стратегией
- PNL-менеджером
- проверкой drawdown

Просто скажи!

### USER
как инициализировать структуру Pipeline

### ASSISTANT
Отлично! Давайте посмотрим, **как правильно инициализировать структуру `Pipeline`** в Go, особенно в контексте торгового бота, который поддерживает **режимы live и бектест**, с внедрением зависимостей (DI), конфигурацией и проверкой валидаторов.

---

## 🧱 Структура `Pipeline`

Предположим, у нас есть такая структура (на основе предыдущих обсуждений):

```go
// pkg/pipeline/pipeline.go
type Pipeline struct {
    Source        DataProvider
    PNLProvider   PNLProvider
    Clusterizer   *Clusterizer
    SeriesMaker   *SeriesMaker
    Strategy      Strategy
    Validators    []Validator
    Executor      Executor

    MaxDrawdown     float64
    TradingEnabled  bool
}
```

---

## ✅ Подход: Конструктор с опциями (Functional Options Pattern)

Это **наиболее гибкий и идиоматичный способ** в Go. Он позволяет:
- Гибко настраивать pipeline
- Избежать длинных списков параметров
- Легко добавлять новые опции

---

### 1. Определение функций-опций

```go
// pipeline_options.go

type Option func(*Pipeline) error

func WithSource(source DataProvider) Option {
    return func(p *Pipeline) error {
        if source == nil {
            return fmt.Errorf("source cannot be nil")
        }
        p.Source = source
        return nil
    }
}

func WithPNLProvider(provider PNLProvider) Option {
    return func(p *Pipeline) error {
        if provider == nil {
            return fmt.Errorf("pnl provider cannot be nil")
        }
        p.PNLProvider = provider
        return nil
    }
}

func WithClusterizer(c *Clusterizer) Option {
    return func(p *Pipeline) error {
        p.Clusterizer = c
        return nil
    }
}

func WithStrategy(s Strategy) Option {
    return func(p *Pipeline) error {
        if s == nil {
            return fmt.Errorf("strategy cannot be nil")
        }
        p.Strategy = s
        return nil
    }
}

func WithValidators(validators ...Validator) Option {
    return func(p *Pipeline) error {
        p.Validators = append(p.Validators, validators...)
        return nil
    }
}

func WithExecutor(executor Executor) Option {
    return func(p *Pipeline) error {
        if executor == nil {
            return fmt.Errorf("executor cannot be nil")
        }
        p.Executor = executor
        return nil
    }
}

func WithMaxDrawdown(drawdown float64) Option {
    return func(p *Pipeline) error {
        if drawdown < 0 {
            return fmt.Errorf("drawdown must be non-negative")
        }
        p.MaxDrawdown = drawdown
        return nil
    }
}
```

---

### 2. Конструктор `NewPipeline`

```go
// pipeline.go
func NewPipeline(opts ...Option) (*Pipeline, error) {
    p := &Pipeline{
        Validators:     make([]Validator, 0),
        TradingEnabled: true,
        MaxDrawdown:    1000.0, // значение по умолчанию
    }

    // Применяем опции
    for _, opt := range opts {
        if err := opt(p); err != nil {
            return nil, err
        }
    }

    // Проверяем обязательные поля
    if p.Source == nil {
        return nil, fmt.Errorf("pipeline must have a data source")
    }
    if p.Strategy == nil {
        return nil, fmt.Errorf("pipeline must have a strategy")
    }
    if p.Executor == nil {
        return nil, fmt.Errorf("pipeline must have an executor")
    }

    return p, nil
}
```

---

## 🛠 Пример инициализации: **Режим бектеста**

```go
// cmd/backtest/main.go
source := NewBacktestSource(db, "BTCUSDT", time.Now().Add(-7*24*time.Hour), time.Now())
pnlProvider := NewBacktestPNLProvider()
clusterizer := NewClusterizer("5m")
seriesMaker := NewSeriesMaker()
strategy := NewRSIStrategy(14, 30, 70)
validator := NewBalanceValidator(100)
executor := NewBacktestExecutor()

pipeline, err := pipeline.NewPipeline(
    pipeline.WithSource(source),
    pipeline.WithPNLProvider(pnlProvider),
    pipeline.WithClusterizer(clusterizer),
    pipeline.WithStrategy(strategy),
    pipeline.WithValidators(validator),
    pipeline.WithExecutor(executor),
    pipeline.WithMaxDrawdown(500.0),
)

if err != nil {
    log.Fatalf("Failed to create pipeline: %v", err)
}

pipeline.RunBacktest()
```

---

## 🛠 Пример инициализации: **Live-режим**

```go
// cmd/live/main.go
apiClient := binance.NewClient(apiKey, secretKey)
source := NewRealTimeSource(apiClient, "BTCUSDT", "1m")
pnlProvider := NewBinancePNLProvider(apiClient)
clusterizer := NewClusterizer("5m")
strategy := NewMACDStrategy()
executor := NewLiveExecutor(apiClient)

pipeline, err := pipeline.NewPipeline(
    pipeline.WithSource(source),
    pipeline.WithPNLProvider(pnlProvider),
    pipeline.WithClusterizer(clusterizer),
    pipeline.WithStrategy(strategy),
    pipeline.WithExecutor(executor),
    pipeline.WithMaxDrawdown(1000.0),
    pipeline.WithValidators(
        NewBalanceValidator(10),
        NewLeverageValidator(10),
    ),
)

if err != nil {
    log.Fatalf("Failed to create live pipeline: %v", err)
}

pipeline.Run() // запуск в реальном времени
```

---

## 💡 Альтернатива: Конфигурация через struct

Если хочется использовать **конфиг-структуру** (например, из YAML):

```go
type PipelineConfig struct {
    MaxDrawdown float64 `yaml:"max_drawdown"`
    SourceType  string  `yaml:"source_type"`
    Symbol      string  `yaml:"symbol"`
    Interval    string  `yaml:"interval"`
    Strategy    string  `yaml:"strategy"`
    // ...
}

func NewPipelineFromConfig(cfg *PipelineConfig) (*Pipeline, error) {
    // Здесь создаются зависимости на основе cfg
    // Например: if cfg.SourceType == "backtest" { ... }
    // Затем вызывается NewPipeline с нужными опциями
}
```

---

## ✅ Преимущества подхода с опциями

| Плюс | Описание |
|------|--------|
| **Гибкость** | Можно комбинировать компоненты без изменения сигнатуры |
| **Безопасность** | Проверка `nil` и валидация в опциях |
| **Читаемость** | Код понятен: `WithSource(...)`, `WithStrategy(...)` |
| **Расширяемость** | Легко добавить новую опцию без ломающих изменений |

---

## 📦 Вывод

Чтобы инициализировать `Pipeline`:

1. ✅ Используй **функциональные опции** (`Option func(*Pipeline) error`)
2. ✅ Внедряй зависимости: источник, стратегию, валидаторы, исполнитель
3. ✅ Проверяй обязательные поля в конструкторе
4. ✅ Поддерживай оба режима: live и backtest — через разные реализации интерфейсов
5. ✅ Делай конструктор `NewPipeline(...Option) (*Pipeline, error)`

---

Если хочешь — могу скинуть **полный рабочий пример проекта** с:
- `go.mod`
- структурой пакетов
- примерами `main.go` для live и backtest
- YAML-конфигом

Просто скажи: "Да, скинь шаблон".