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

### USER
Как контролировать статус выполнения по разным стратегиям приспользовании веб-интерфейса

### ASSISTANT
Отличный вопрос! Добавление **веб-интерфейса** для мониторинга и управления **несколькими стратегиями** в торговом боте на Go — мощный шаг к удобству, прозрачности и контролю. Ниже я покажу, **как организовать контроль статуса выполнения по разным стратегиям**, используя веб-интерфейс (например, на базе **HTTP + HTML/JS или WebSocket + frontend**).

---

## 🎯 Цель
Пользователь должен видеть в реальном времени:
- Какие стратегии запущены / остановлены
- Их текущий статус (работает, приостановлена, ошибка)
- Сигналы (покупка/продажа)
- PNL, баланс, открытые позиции
- Возможность вручную запустить/остановить стратегию

---

## 🏗️ Архитектура с веб-интерфейсом

```
[Торговые стратегии] → [Менеджер стратегий] → [API (HTTP/WS)] → [Веб-интерфейс]
       ↑                                           ↑
   Статус, PNL, сигналы                     HTML/JS (React, Vue, или простой шаблон)
```

---

## 1. Модель: `StrategyInstance` — экземпляр стратегии

```go
// pkg/strategy/strategy.go
type Status string

const (
    StatusRunning  Status = "running"
    StatusPaused   Status = "paused"
    StatusStopped  Status = "stopped"
    StatusError    Status = "error"
)

type StrategyInstance struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Symbol      string    `json:"symbol"`
    Interval    string    `json:"interval"`
    Status      Status    `json:"status"`
    LastSignal  *Signal   `json:"last_signal,omitempty"`
    PNL         float64   `json:"pnl"`
    Balance     float64   `json:"balance"`
    Position    *Position `json:"position,omitempty"`
    Error       string    `json:"error,omitempty"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## 2. Менеджер стратегий `StrategyManager`

```go
// pkg/strategy/manager.go
type StrategyManager struct {
    strategies map[string]*StrategyInstance
    mutex      sync.RWMutex
    pipelines  map[string]*Pipeline // ID → Pipeline
}

func NewStrategyManager() *StrategyManager {
    return &StrategyManager{
        strategies: make(map[string]*StrategyInstance),
        pipelines:  make(map[string]*Pipeline),
    }
}

func (sm *StrategyManager) AddStrategy(id, name, symbol, interval string, pipeline *Pipeline) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    sm.strategies[id] = &StrategyInstance{
        ID:        id,
        Name:      name,
        Symbol:    symbol,
        Interval:  interval,
        Status:    StatusStopped,
        UpdatedAt: time.Now(),
    }
    sm.pipelines[id] = pipeline
}
```

---

## 3. Обновление статуса в реальном времени

Каждая стратегия (pipeline) должна **обновлять свой статус** через callback:

```go
// Внутри Pipeline
func (p *Pipeline) updateStrategyStatus(manager *StrategyManager, id string, updater func(*StrategyInstance)) {
    manager.mutex.Lock()
    defer manager.mutex.Unlock()

    if s, exists := manager.strategies[id]; exists {
        updater(s)
        s.UpdatedAt = time.Now()
    }
}

// Пример использования при формировании сигнала
p.updateStrategyStatus(manager, "strat-1", func(s *StrategyInstance) {
    s.LastSignal = signal
    s.PNL = getCurrentPNL()
})
```

---

## 4. HTTP API для веб-интерфейса

```go
// cmd/web/server.go
func StartWebServer(manager *StrategyManager, port string) {
    http.HandleFunc("/api/strategies", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            manager.mutex.RLock()
            list := make([]*StrategyInstance, 0, len(manager.strategies))
            for _, s := range manager.strategies {
                list = append(list, s)
            }
            manager.mutex.RUnlock()

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(list)
        }
    })

    http.HandleFunc("/api/strategy/start", func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        // Запуск pipeline...
        manager.updateStatus(id, StatusRunning)
        w.WriteHeader(200)
    })

    http.HandleFunc("/api/strategy/stop", ...)
    
    // Статические файлы (веб-интерфейс)
    http.Handle("/", http.FileServer(http.Dir("./web/")))

    log.Printf("Web UI started at http://localhost:%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

---

## 5. Веб-интерфейс (HTML + JS)

### Пример: `web/index.html`

```html
<!DOCTYPE html>
<html>
<head>
  <title>Trade Bot Monitor</title>
  <meta http-equiv="refresh" content="5"> <!-- автообновление -->
</head>
<body>
  <h1>Стратегии</h1>
  <table border="1">
    <tr><th>ID</th><th>Имя</th><th>Символ</th><th>Статус</th><th>PNL</th><th>Сигнал</th><th>Управление</th></tr>
    <tbody id="strategy-list"></tbody>
  </table>

  <script>
    async function loadStrategies() {
      const res = await fetch('/api/strategies');
      const strategies = await res.json();
      const tbody = document.getElementById('strategy-list');
      tbody.innerHTML = '';

      strategies.forEach(s => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${s.id}</td>
          <td>${s.name}</td>
          <td>${s.symbol}</td>
          <td style="color:${s.status === 'running' ? 'green' : s.status === 'error' ? 'red' : 'orange'}">
            ${s.status}
          </td>
          <td>${s.pnl.toFixed(2)}</td>
          <td>${s.last_signal ? s.last_signal.Side : '-'}</td>
          <td>
            <button onclick="start('${s.id}')">Старт</button>
            <button onclick="stop('${s.id}')">Стоп</button>
          </td>
        `;
        tbody.appendChild(row);
      });
    }

    function start(id) { fetch(`/api/strategy/start?id=${id}`); }
    function stop(id) { fetch(`/api/strategy/stop?id=${id}`); }

    // Автообновление
    setInterval(loadStrategies, 3000);
    loadStrategies();
  </script>
</body>
</html>
```

---

## 6. Улучшения: WebSocket для реального времени

Вместо автообновления через `setInterval` — используй **WebSocket**, чтобы обновлять статус мгновенно.

```go
// WebSocket сервер
var clients = make(map[chan []byte]bool)

func StartWebSocket(manager *StrategyManager) {
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn := make(chan []byte, 10)
        clients[conn] = true

        go func() {
            for msg := range conn {
                w.Write(msg)
            }
        }()

        // При обновлении статуса — отправлять всем
        go func() {
            ticker := time.NewTicker(1 * time.Second)
            for range ticker.C {
                data, _ := json.Marshal(manager.GetAll())
                for client := range clients {
                    select {
                    case client <- 
                    default:
                        close(client)
                        delete(clients, client)
                    }
                }
            }
        }()
    })
}
```

---

## 7. Расширенные возможности

| Функция | Как реализовать |
|--------|----------------|
| 🔍 Детали стратегии | `/api/strategy/{id}` — детали, логи, графики |
| 📈 График PNL | Отправить историю PNL в JSON, построить Chart.js |
| 🔔 Уведомления | WebSocket + toast-уведомления при новых сигналах |
| 📝 Логи | Вести лог действий: "Сигнал BUY", "PNL ниже порога" |
| 🔐 Авторизация | Добавить JWT или базовую аутентификацию |

---

## ✅ Пример: Запуск всего приложения

```go
// main.go
func main() {
    manager := NewStrategyManager()

    // Создаём стратегии
    pipeline1 := NewPipeline(WithStrategy(...), ...)
    manager.AddStrategy("btc-macd", "MACD", "BTCUSDT", "5m", pipeline1)

    // Запускаем веб-интерфейс
    go StartWebServer(manager, "8080")

    // Запускаем стратегии
    go pipeline1.Run()

    // Ждём...
    select {}
}
```

---

## 📦 Вывод

Чтобы контролировать статус стратегий через веб-интерфейс:

1. ✅ Создай `StrategyInstance` с полями: статус, PNL, сигнал, ошибка
2. ✅ Используй `StrategyManager` для хранения и обновления состояния
3. ✅ Экспонируй данные через **HTTP API** (или WebSocket)
4. ✅ Сделай простой веб-интерфейс на HTML/JS (или React/Vue)
5. ✅ Позволь пользователю запускать/останавливать стратегии
6. ✅ Обновляй статус **синхронно** с работой pipeline

---

## 🚀 Дополнительно

Если хочешь, могу:
- Скинуть **рабочий пример проекта** с веб-интерфейсом
- Добавить **авторизацию**
- Подключить **Prometheus + Grafana** для мониторинга
- Сделать **адаптивный UI** под мобильные устройства

Просто скажи: **"Скинь пример с веб-интерфейсом"** — и я подготовлю архив или репозиторий.

### USER
При использовании web сокета как будет выглядеть frontend реализация на js

### ASSISTANT
Отлично! Ниже — **полная и рабочая реализация frontend на чистом JavaScript** (без фреймворков), которая подключается к Go-серверу через **WebSocket**, получает обновления статуса стратегий в реальном времени и отображает их в таблице.

---

## 🧩 Задача
- Подключиться к WebSocket-серверу
- Получать обновления списка стратегий
- Обновлять таблицу без перезагрузки страницы
- Показывать новые сигналы и изменения PNL мгновенно

---

## ✅ Backend (Go): Отправка данных через WebSocket

(Кратко, для контекста)

```go
// Go: WebSocket сервер
var clients = make(map[chan []byte]bool)
var clientsMu sync.Mutex

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn := make(chan []byte, 10)
    clientsMu.Lock()
    clients[conn] = true
    clientsMu.Unlock()

    defer func() {
        clientsMu.Lock()
        delete(clients, conn)
        clientsMu.Unlock()
        close(conn)
    }()

    // Отправляем заголовок
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`
<!DOCTYPE html>
<html>
<head><title>Trade Bot Monitor</title></head>
<body>
  <h1>Стратегии (WebSocket)</h1>
  <table border="1" id="strategy-table">
    <tr><th>ID</th><th>Имя</th><th>Символ</th><th>Статус</th><th>PNL</th><th>Сигнал</th></tr>
    <tbody id="strategy-body"></tbody>
  </table>
  <script src="/static/client.js"></script>
</body>
</html>`))

    // Ждём данные от backend и пересылаем клиенту
    for {
        select {
        case data := <-conn:
            w.(http.Hijacker).Hijack()
            // (реализация отправки фрейма — см. ниже, или используем библиотеку)
        }
    }
}
```

Но проще — **отдельный WebSocket endpoint** и **отдельный HTML**.

---

## ✅ Frontend: `client.js` — WebSocket клиент

```js
// static/client.js

class StrategyMonitor {
    constructor() {
        this.socket = null;
        this.tableBody = document.getElementById('strategy-body');
        this.connect();
    }

    connect() {
        // Попробуем подключиться к ws://localhost:8080/ws
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const url = `${protocol}//${window.location.host}/ws`;
        
        this.socket = new WebSocket(url);

        this.socket.onopen = () => {
            console.log('✅ Подключено к WebSocket');
            this.setStatus('Подключено', 'green');
        };

        this.socket.onmessage = (event) => {
            try {
                const strategies = JSON.parse(event.data);
                this.updateTable(strategies);
            } catch (e) {
                console.error('Ошибка парсинга JSON:', e);
            }
        };

        this.socket.onclose = (event) => {
            console.log('❌ Соединение закрыто', event);
            this.setStatus('Переподключение...', 'orange');
            setTimeout(() => this.connect(), 3000); // reconnect
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket ошибка:', error);
            this.setStatus('Ошибка подключения', 'red');
        };
    }

    setStatus(text, color) {
        const status = document.getElementById('status');
        if (status) status.innerText = text;
        document.body.style.color = color;
    }

    updateTable(strategies) {
        this.tableBody.innerHTML = ''; // очищаем

        strategies.forEach(s => {
            const row = document.createElement('tr');

            // Цвет статуса
            const statusColor = {
                running: 'green',
                paused: 'orange',
                stopped: 'gray',
                error: 'red'
            }[s.status] || 'black';

            const lastSignal = s.last_signal ? `${s.last_signal.side.toUpperCase()} @ ${s.last_signal.price}` : '–';

            row.innerHTML = `
                <td>${this.escape(s.id)}</td>
                <td>${this.escape(s.name)}</td>
                <td>${this.escape(s.symbol)}</td>
                <td style="color:${statusColor}">${this.escape(s.status)}</td>
                <td style="font-weight:bold;">${s.pnl?.toFixed(2) ?? '–'}</td>
                <td>${lastSignal}</td>
            `;

            // Подсветка при новом сигнале
            if (s.last_signal && s.updated_at && this.isRecentlyUpdated(s.updated_at)) {
                row.style.backgroundColor = s.last_signal.side === 'buy' ? '#aaffaa' : '#ffaaaa';
                setTimeout(() => {
                    if (row.parentElement) row.style.backgroundColor = '';
                }, 2000);
            }

            this.tableBody.appendChild(row);
        });
    }

    isRecentlyUpdated(timestamp) {
        const now = new Date();
        const updated = new Date(timestamp);
        return (now - updated) < 3000; // менее 3 секунд назад
    }

    escape(str) {
        if (typeof str !== 'string') return str;
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
}
```

---

## ✅ HTML: `index.html`

```html
<!-- static/index.html -->
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Trade Bot Monitor</title>
  <style>
    body { font-family: Arial, sans-serif; padding: 20px; }
    table { width: 100%; border-collapse: collapse; margin-top: 10px; }
    th, td { padding: 8px; text-align: left; border: 1px solid #ccc; }
    th { background: #f0f0f0; }
    #status { font-weight: bold; margin-bottom: 10px; }
  </style>
</head>
<body>
  <h1>📊 Мониторинг торговых стратегий</h1>
  <p>Статус подключения: <span id="status">Подключение...</span></p>

  <table id="strategy-table">
    <thead>
      <tr>
        <th>ID</th>
        <th>Имя</th>
        <th>Символ</th>
        <th>Статус</th>
        <th>PNL</th>
        <th>Последний сигнал</th>
      </tr>
    </thead>
    <tbody id="strategy-body">
      <!-- Данные будут добавлены через JS -->
    </tbody>
  </table>

  <script src="/static/client.js"></script>
  <script>
    // Запускаем монитор
    new StrategyMonitor();
  </script>
</body>
</html>
```

---

## ✅ Go: Отправка данных клиентам (серверная часть)

```go
// В Go: рассылка всем клиентам
func BroadcastStrategies(manager *StrategyManager) {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        manager.mutex.RLock()
        data, _ := json.Marshal(manager.GetAll()) // []StrategyInstance
        manager.mutex.RUnlock()

        clientsMu.Lock()
        for client := range clients {
            select {
            case client <- data:
            default:
                // Клиент не читает — удаляем
                close(client)
                delete(clients, client)
            }
        }
        clientsMu.Unlock()
    }
}
```

> `clients` — это `map[chan []byte]bool`, где каждый канал получает JSON.

---

## ✅ Как сервер отправляет данные по WebSocket

Чтобы не писать вручную протокол WebSocket, лучше использовать пакет, например:

```bash
go get github.com/gorilla/websocket
```

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // для разработки
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print(err)
        return
    }
    defer conn.Close()

    // Отправляем данные каждую секунду
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        data := getCurrentStrategiesJSON() // []byte
        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            break
        }
    }
}
```

---

## ✅ Результат

Когда стратегия:
- Отправляет **новый сигнал** → frontend мгновенно показывает его
- PNL **падает ниже порога** → статус становится красным
- Стратегия **остановлена** → цвет статуса меняется
- Ошибка в pipeline → отображается в интерфейсе

---

## 🚀 Дополнительные улучшения

| Функция | Реализация |
|--------|-----------|
| 🔔 Уведомления | `new Notification("Новый сигнал: BUY")` |
| 📈 График PNL | Используй Chart.js + WebSocket |
| 📱 Адаптив | Сделай таблицу responsive с `overflow-x: auto` |
| 🔄 Ручное управление | Кнопки "Старт/Стоп" через `fetch('/api/start?id=...')` |

---

## ✅ Вывод

**Frontend с WebSocket на чистом JS** — простой, эффективный и не требует сборки.

🔹 Подключается к Go-серверу  
🔹 Получает обновления в реальном времени  
🔹 Обновляет таблицу без перезагрузки  
🔹 Автоматически переподключается при обрыве  

---

## 📦 Хочешь пример целиком?

Напиши:  
👉 **"Скинь zip с примером"** — и я подготовлю:
- Go-сервер с WebSocket
- HTML + JS
- Пример стратегии
- `go.mod`, структуру папок
- Инструкцию по запуску

Готов к отправке!

### USER
В приложении для построения серий используется SeriesBuilder. У него есть метод AddPoints, который добавляет точки в серии по определённому алгоритму. Алгоритмы реализуются отдельными структурами. Есть также общие методы AddClusteredData и GetActiveSeries, которые не зависят от реализации конкретного алгоритма. Как организовать код?

### ASSISTANT
Отличный вопрос! У вас есть **обобщённый компонент `SeriesBuilder`**, который:
- Использует **разные алгоритмы построения серий** (тренд, флэт, волны и т.п.)
- Имеет **общие методы** (`AddClusteredData`, `GetActiveSeries`)
- Должен быть **гибким, расширяемым и тестируемым**

Ниже — **рекомендуемая организация кода на Go**, с использованием **интерфейсов, композиции и паттерна "Стратегия"**.

---

## 🎯 Цель
Разделить:
- **Общую логику** `SeriesBuilder` (управление сериями, буферизация, доступ к данным)
- **Специфическую логику** построения серий (алгоритмы)
- Обеспечить **лёгкое добавление новых алгоритмов** без изменения основного кода

---

## 🏗️ Архитектура

```
SeriesBuilder
  ├── Алгоритм (Strategy Algorithm) — интерфейс
  ├── Общие методы: AddClusteredData(), GetActiveSeries()
  └── Внутреннее состояние: активные/завершённые серии
        ↓
AlgorithmTrend, AlgorithmFlat, AlgorithmVolatility и т.д.
```

---

## 1. Интерфейс алгоритма

```go
// pkg/series/algorithm.go
type Algorithm interface {
    // Обрабатывает новую точку (например, свечу) и возвращает:
    // - новую серию (если началась)
    // - признак завершения текущей серии
    // - ошибку
    ProcessPoint(point *ClusteredCandle) (*Series, bool, error)
    
    // Возвращает имя алгоритма (для логирования)
    Name() string
}
```

---

## 2. Структура `Series` (серия)

```go
// pkg/series/series.go
type Series struct {
    ID        string
    Type      string // "trend", "flat", etc.
    StartTime time.Time
    EndTime   time.Time
    Points    []*ClusteredCandle
    IsClosed  bool
}
```

---

## 3. Реализация алгоритма (пример: тренд)

```go
// pkg/series/algo_trend.go
type TrendAlgorithm struct {
    threshold float64 // минимальное изменение цены для начала тренда
}

func NewTrendAlgorithm(threshold float64) *TrendAlgorithm {
    return &TrendAlgorithm{threshold: threshold}
}

func (a *TrendAlgorithm) Name() string {
    return "trend"
}

func (a *TrendAlgorithm) ProcessPoint(point *ClusteredCandle) (*Series, bool, error) {
    // Логика: если цена изменилась больше чем на threshold — начинаем новую серию
    // Здесь упрощённо
    if len(a.lastPoints) < 2 {
        a.lastPoints = append(a.lastPoints, point)
        return nil, false, nil
    }

    prev := a.lastPoints[len(a.lastPoints)-1]
    change := (point.Close - prev.Close) / prev.Close

    if math.Abs(change) > a.threshold && change > 0 {
        return &Series{
            ID:        generateID(),
            Type:      "trend_up",
            StartTime: point.Time,
            Points:    []*ClusteredCandle{point},
        }, false, nil
    } else if math.Abs(change) > a.threshold && change < 0 {
        return &Series{
            ID:        generateID(),
            Type:      "trend_down",
            StartTime: point.Time,
            Points:    []*ClusteredCandle{point},
        }, false, nil
    }

    return nil, false, nil
}
```

> Аналогично реализуются `FlatAlgorithm`, `VolatilityAlgorithm` и т.д.

---

## 4. `SeriesBuilder` — обобщённый билдер

```go
// pkg/series/builder.go
type SeriesBuilder struct {
    algorithm     Algorithm
    activeSeries  *Series
    completed     []*Series
    mutex         sync.RWMutex
}

func NewSeriesBuilder(algorithm Algorithm) *SeriesBuilder {
    return &SeriesBuilder{
        algorithm: algorithm,
        completed: make([]*Series, 0),
    }
}
```

---

## 5. Общие методы (не зависят от алгоритма)

```go
// Добавляет агрегированные данные (например, 5-минутные свечи)
func (b *SeriesBuilder) AddClusteredData(candle *ClusteredCandle) {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    // Передаём точку алгоритму
    newSeries, shouldClose, err := b.algorithm.ProcessPoint(candle)
    if err != nil {
        log.Printf("Algorithm error: %v", err)
        return
    }

    // Закрываем текущую серию, если нужно
    if shouldClose && b.activeSeries != nil && !b.activeSeries.IsClosed {
        b.activeSeries.IsClosed = true
        b.activeSeries.EndTime = candle.Time
        b.completed = append(b.completed, b.activeSeries)
        b.activeSeries = nil
    }

    // Начинаем новую серию
    if newSeries != nil {
        b.activeSeries = newSeries
    }

    // Добавляем точку в активную серию
    if b.activeSeries != nil {
        b.activeSeries.Points = append(b.activeSeries.Points, candle)
    }
}

// Возвращает активную серию
func (b *SeriesBuilder) GetActiveSeries() *Series {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    return b.activeSeries
}

// Возвращает завершённые серии
func (b *SeriesBuilder) GetCompletedSeries() []*Series {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    return b.completed
}

// Получить алгоритм (для логирования/мониторинга)
func (b *SeriesBuilder) AlgorithmName() string {
    return b.algorithm.Name()
}
```

---

## 6. Пример использования

```go
// main.go
trendAlg := NewTrendAlgorithm(0.01) // 1%
builder := NewSeriesBuilder(trendAlg)

for _, candle := range clusteredCandles {
    builder.AddClusteredData(candle)

    if series := builder.GetActiveSeries(); series != nil {
        fmt.Printf("Active %s series: %d points\n", series.Type, len(series.Points))
    }
}
```

---

## ✅ Преимущества такой архитектуры

| Плюс | Как реализовано |
|------|-----------------|
| **Расширяемость** | Добавляй новые алгоритмы, реализуя `Algorithm` |
| **Инкапсуляция** | `SeriesBuilder` не знает деталей алгоритма |
| **Тестируемость** | Можно мокать `Algorithm` в unit-тестах |
| **Гибкость** | Можно менять алгоритм "на лету" (если нужно) |
| **Чистота кода** | Общие методы в `SeriesBuilder`, специфика — в алгоритмах |

---

## 📁 Структура проекта

```
/pkg/series/
  ├── builder.go               // SeriesBuilder
  ├── series.go                // структура Series
  ├── algorithm.go             // интерфейс Algorithm
  ├── algo_trend.go            // реализация
  ├── algo_flat.go
  ├── algo_volatility.go
  └── algo_custom_xxx.go
```

---

## 🧪 Тестирование

```go
// series/builder_test.go
func TestSeriesBuilder_TrendStartsOnThreshold(t *testing.T) {
    mockAlg := &MockAlgorithm{
        nextSeries: &Series{Type: "trend_up"},
    }
    builder := NewSeriesBuilder(mockAlg)

    candle := &ClusteredCandle{Time: now, Close: 101.0}
    builder.AddClusteredData(candle)

    require.NotNil(t, builder.GetActiveSeries())
    require.Equal(t, "trend_up", builder.GetActiveSeries().Type)
}
```

---

## 🔄 Возможные улучшения

| Функция | Как добавить |
|--------|-------------|
| Поддержка нескольких алгоритмов одновременно | `SeriesBuilder` принимает `[]Algorithm` |
| Динамическая смена алгоритма | `builder.SetAlgorithm(newAlg)` |
| Сохранение контекста в алгоритме | Добавить `Context()` или `Reset()` в интерфейс |
| Метрики | Считать, сколько серий создано, средняя длина и т.п. |

---

## ✅ Вывод

Чтобы правильно организовать `SeriesBuilder`:

1. ✅ Создай интерфейс `Algorithm` с методом `ProcessPoint`
2. ✅ Реализуй каждый алгоритм как отдельную структуру
3. ✅ `SeriesBuilder` использует алгоритм через интерфейс
4. ✅ Общие методы (`AddClusteredData`, `GetActiveSeries`) — в `SeriesBuilder`
5. ✅ Используй мьютексы для потокобезопасности
6. ✅ Легко тестируй и расширяй

---

Если хочешь — могу скинуть **полный рабочий пример** с:
- 2 алгоритмами
- юнит-тестами
- примером интеграции в pipeline

Просто скажи: **"Скинь пример SeriesBuilder"** — и я подготовлю.

### USER
Реализовать торговый менеджер который по выбранным из базы данных стратегиям запускает выполнение бектеста или реального режима. Запуск будет осуществляться через веб-интерфейс. Также нужна возможность отслеживать статус выполнения,остановить, возможность обновить настройки стратегии без остановки. 

### ASSISTANT
Отлично! Ниже — **полная архитектура и реализация торгового менеджера на Go**, который:

- Загружает стратегии из базы данных
- Запускает **бектест** или **реальный режим**
- Управляется через **веб-интерфейс** (HTTP API)
- Позволяет:
  - Запускать/останавливать стратегии
  - Отслеживать статус
  - Обновлять настройки **без остановки**

---

## 🧱 Архитектура системы

```
[Веб-интерфейс] → HTTP API (Gin/Echo) → TradeManager → [StrategyRunner] → Pipeline
                                         ↑
                                 [База данных: стратегии, настройки]
```

---

## 📁 Структура проекта

```
/cmd/
  /web/
    main.go
/pkg/
  /strategy/
    strategy.go
    repository.go          // Работа с БД
  /pipeline/
    pipeline.go
    pipeline_manager.go    // Управление запуском
  /trading/
    trade_manager.go       // Основной менеджер
  /web/
    handlers.go
    routes.go
  /storage/
    db.go                  // PostgreSQL/SQLite
/web/
  /static/
    index.html
    app.js
```

---

## 1. Модель стратегии в БД

```go
// pkg/strategy/strategy.go
type Strategy struct {
    ID          int
    Name        string
    Type        string  // "rsi", "macd", "ma_cross"
    Mode        string  // "live", "backtest"
    Symbol      string
    Interval    string
    Params      json.RawMessage  // JSON: {"rsi_period": 14, "overbought": 70}
    Enabled     bool
    UpdatedAt   time.Time
}
```

Пример в БД:
```json
{
  "id": 1,
  "name": "RSI Scalper",
  "type": "rsi",
  "mode": "live",
  "symbol": "BTCUSDT",
  "interval": "5m",
  "params": {"period": 14, "low": 30, "high": 70},
  "enabled": true
}
```

---

## 2. Repository: Работа с БД

```go
// pkg/strategy/repository.go
type Repository struct {
    db *sql.DB
}

func (r *Repository) GetAll() ([]*Strategy, error) {
    rows, err := r.db.Query("SELECT id, name, type, mode, symbol, interval, params, enabled, updated_at FROM strategies")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var strategies []*Strategy
    for rows.Next() {
        var s Strategy
        var params []byte
        err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Mode, &s.Symbol, &s.Interval, &params, &s.Enabled, &s.UpdatedAt)
        if err != nil {
            return nil, err
        }
        s.Params = params
        strategies = append(strategies, &s)
    }
    return strategies, nil
}

func (r *Repository) GetByID(id int) (*Strategy, error) { /* ... */ }

func (r *Repository) UpdateParams(id int, params json.RawMessage) error { /* ... */ }
```

---

## 3. Strategy Runner (один экземпляр на стратегию)

```go
// pkg/trading/runner.go
type Runner struct {
    Strategy *Strategy
    Pipeline *pipeline.Pipeline
    Cancel   context.CancelFunc  // Для остановки
    Running  bool
    Err      error
    Mutex    sync.RWMutex
}

func (r *Runner) Start(ctx context.Context) {
    r.Mutex.Lock()
    if r.Running {
        r.Mutex.Unlock()
        return
    }
    r.Running = true
    r.Mutex.Unlock()

    runCtx, cancel := context.WithCancel(ctx)
    r.Cancel = cancel

    // Создаём пайплайн с текущими параметрами
    p, err := pipeline.NewPipelineFromStrategy(r.Strategy)
    if err != nil {
        r.Err = err
        r.Running = false
        return
    }
    r.Pipeline = p

    if r.Strategy.Mode == "backtest" {
        p.RunBacktest()
    } else {
        p.Run(runCtx) // Передаём контекст для остановки
    }

    r.Mutex.Lock()
    r.Running = false
    r.Mutex.Unlock()
}

func (r *Runner) Stop() {
    if r.Cancel != nil {
        r.Cancel()
    }
}

func (r *Runner) IsRunning() bool {
    r.Mutex.RLock()
    defer r.Mutex.RUnlock()
    return r.Running
}

func (r *Runner) UpdateStrategy(s *Strategy) {
    r.Mutex.Lock()
    old := r.Strategy
    r.Strategy = s
    r.Mutex.Unlock()

    // Обновляем пайплайн (без остановки!)
    if r.IsRunning() && r.Pipeline != nil {
        r.Pipeline.UpdateStrategyParams(s.Params) // см. ниже
    }
}
```

> ⚠️ `UpdateStrategyParams` — метод, который **обновляет параметры стратегии "на лету"**, не останавливая пайплайн.

---

## 4. TradeManager — центральный менеджер

```go
// pkg/trading/trade_manager.go
type TradeManager struct {
    runners    map[int]*Runner
    repo       *strategy.Repository
    mu         sync.RWMutex
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewTradeManager(repo *strategy.Repository) *TradeManager {
    ctx, cancel := context.WithCancel(context.Background())
    return &TradeManager{
        runners: make(map[int]*Runner),
        repo:    repo,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (tm *TradeManager) LoadAndStartAll() error {
    strategies, err := tm.repo.GetAll()
    if err != nil {
        return err
    }

    for _, s := range strategies {
        if s.Enabled {
            tm.StartStrategy(s)
        }
    }
    return nil
}

func (tm *TradeManager) StartStrategy(s *Strategy) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    if _, exists := tm.runners[s.ID]; exists {
        return
    }

    runner := &Runner{Strategy: s}
    tm.runners[s.ID] = runner

    go runner.Start(tm.ctx)
}

func (tm *TradeManager) StopStrategy(id int) {
    tm.mu.RLock()
    runner, exists := tm.runners[id]
    tm.mu.RUnlock()

    if exists {
        runner.Stop()
    }
}

func (tm *TradeManager) GetStatus() map[int]Status {
    tm.mu.RLock()
    defer tm.mu.RUnlock()

    status := make(map[int]Status)
    for id, r := range tm.runners {
        status[id] = Status{
            ID:       id,
            Name:     r.Strategy.Name,
            Running:  r.IsRunning(),
            Mode:     r.Strategy.Mode,
            Symbol:   r.Strategy.Symbol,
            Error:    r.Err,
        }
    }
    return status
}

func (tm *TradeManager) UpdateStrategyParams(id int, params json.RawMessage) error {
    s, err := tm.repo.GetByID(id)
    if err != nil {
        return err
    }
    s.Params = params

    tm.repo.UpdateParams(id, params)

    tm.mu.RLock()
    runner, exists := tm.runners[id]
    tm.mu.RUnlock()

    if exists {
        runner.UpdateStrategy(s)
    }
    return nil
}
```

---

## 5. HTTP API (веб-интерфейс)

```go
// pkg/web/handlers.go
func (h *Handler) StartAll(c *gin.Context) {
    err := h.Manager.LoadAndStartAll()
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"status": "all strategies started"})
}

func (h *Handler) StopStrategy(c *gin.Context) {
    id := c.Param("id")
    iid, _ := strconv.Atoi(id)
    h.Manager.StopStrategy(iid)
    c.JSON(200, gin.H{"status": "stopped"})
}

func (h *Handler) GetStatus(c *gin.Context) {
    status := h.Manager.GetStatus()
    c.JSON(200, status)
}

func (h *Handler) UpdateParams(c *gin.Context) {
    id := c.Param("id")
    iid, _ := strconv.Atoi(id)

    var req struct {
        Params json.RawMessage `json:"params"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid JSON"})
        return
    }

    err := h.Manager.UpdateStrategyParams(iid, req.Params)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"status": "updated"})
}
```

---

## 6. Веб-интерфейс (простой HTML + JS)

```html
<!-- /web/static/index.html -->
<h1>Торговый менеджер</h1>
<button onclick="startAll()">Запустить все</button>
<button onclick="getStatus()">Обновить статус</button>

<pre id="status"></pre>

<script>
function getStatus() {
  fetch("/api/status")
    .then(r => r.json())
    .then(data => {
      document.getElementById("status").textContent = JSON.stringify(data, null, 2);
    });
}

function startAll() {
  fetch("/api/start", { method: "POST" })
    .then(() => getStatus());
}

// Автообновление
setInterval(getStatus, 3000);
</script>
```

---

## 7. Настройка маршрутов

```go
// pkg/web/routes.go
func SetupRoutes(r *gin.Engine, manager *trading.TradeManager) {
    handler := &Handler{Manager: manager}

    api := r.Group("/api")
    {
        api.POST("/start", handler.StartAll)
        api.POST("/stop/:id", handler.StopStrategy)
        api.GET("/status", handler.GetStatus)
        api.POST("/update/:id", handler.UpdateParams)
    }

    r.Static("/static", "./web/static")
    r.GET("/", func(c *gin.Context) {
        c.File("./web/static/index.html")
    })
}
```

---

## 8. Запуск (main.go)

```go
// cmd/web/main.go
func main() {
    db := storage.InitDB() // PostgreSQL или SQLite
    repo := strategy.NewRepository(db)
    manager := trading.NewTradeManager(repo)

    // Загружаем и запускаем стратегии
    manager.LoadAndStartAll()

    r := gin.Default()
    web.SetupRoutes(r, manager)

    r.Run(":8080")
}
```

---

## ✅ Функциональность

| Функция | Реализовано |
|-------|-----------|
| Загрузка стратегий из БД | ✅ |
| Запуск live / бектеста | ✅ |
| Веб-интерфейс | ✅ |
| Просмотр статуса | ✅ |
| Остановка стратегии | ✅ |
| Обновление параметров без остановки | ✅ (через `UpdateStrategyParams`) |
| Безопасность (горутины, мьютексы) | ✅ |

---

## 🚀 Дополнительно (по желанию)

- **WebSocket** — для реального обновления статуса
- **Аутентификация** — если веб-интерфейс публичный
- **Логирование в БД** — история PNL, сделок
- **Графики** — Chart.js + WebSocket
- **Конфигурация через UI** — редактирование параметров в браузере

---

## 💡 Подсказка: Как обновлять стратегию "на лету"

Внутри стратегии используй **мьютекс**:

```go
type RSI struct {
    mu     sync.RWMutex
    Period int
    Low    int
    High   int
}

func (s *RSI) Update(params json.RawMessage) {
    var cfg struct {
        Period int `json:"period"`
        Low    int `json:"low"`
        High   int `json:"high"`
    }
    json.Unmarshal(params, &cfg)

    s.mu.Lock()
    s.Period = cfg.Period
    s.Low = cfg.Low
    s.High = cfg.High
    s.mu.Unlock()
}
```

И в `Analyze()` — читай с `s.mu.RLock()`.

---

## 📦 Вывод

Вы построили **масштабируемый торговый менеджер** с:

- Централизованным управлением
- Поддержкой live и backtest
- Веб-интерфейсом
- Возможностью динамического обновления

---

✅ **Следующие шаги:**
1. Добавить WebSocket для живых обновлений
2. Реализовать историю сделок
3. Сделать UI красивее (React/Vue)
4. Добавить уведомления (Telegram)

---

### USER
Почему реквизиты Runner не разместить в Pipeline?

### ASSISTANT
Отличный вопрос! На первый взгляд, действительно, **логично перенести поля `Runner` (например, `Running`, `Cancel`, `Err`, `Mutex`) внутрь `Pipeline`**, ведь `Runner` — это обёртка, управляющая запуском пайплайна. Однако есть **важные архитектурные и практические причины**, почему лучше **оставить `Runner` отдельно**.

Давайте разберёмся.

---

## 🔍 Сравнение: `Pipeline` vs `Runner`

| Компонент | Назначение |
|---------|-----------|
| `Pipeline` | **Обработка данных**: кластеризация, стратегия, валидация, исполнение |
| `Runner`   | **Управление жизненным циклом**: запуск, остановка, статус, обновление |

---

## ✅ Почему `Runner` не должен быть частью `Pipeline`

### 1. **Разделение ответственностей (Single Responsibility Principle)**

- `Pipeline` должен **только обрабатывать данные**, а не управлять своим состоянием.
- Если добавить `Cancel`, `Running`, `Mutex` в `Pipeline`, он станет **одновременно процессором и контроллером** — это нарушает SRP.

> 💡 `Pipeline` — это **машина для обработки**, а `Runner` — **оператор этой машины**.

---

### 2. **`Pipeline` может быть переиспользован в разных контекстах**

- В **бектесте** `Pipeline` запускается **один раз**, синхронно, без остановки.
- В **live-режиме** он работает **вечно**, может быть остановлен, перезапущен.
- Если `Pipeline` сам управляет `context.Cancel`, он **привязывается к конкретному режиму**.

Пример:
```go
// В бектесте нет необходимости в Cancel
pipeline.RunBacktest() // просто цикл по данным

// В live — нужен контекст
ctx, cancel := context.WithCancel(context.Background())
pipeline.Run(ctx) // может быть остановлен
```

Если `Cancel` будет внутри `Pipeline`, то:
- Бектест будет содержать "мёртвый" код
- Нарушится чистота абстракции

---

### 3. **Поддержка нескольких экземпляров одной стратегии**

Допустим, вы хотите запустить **одну и ту же стратегию** на разных парах:
- BTCUSDT
- ETHUSDT

Вы можете создать **один `Pipeline` на каждую пару**, но управлять ими через **единый `Runner` или `Manager`**.

Если `Pipeline` сам управляет `Cancel`, то:
- Каждый пайплайн будет "думать" о себе
- Управление извне станет сложнее
- Невозможно централизованно остановить все

---

### 4. **`Runner` может управлять не только `Pipeline`, а и другими компонентами**

В будущем `Runner` может:
- Логировать статистику
- Отправлять уведомления
- Перезапускать `Pipeline` при панике
- Собирать метрики

Если всё это встроить в `Pipeline`, он станет **"божественным объектом" (God Object)**.

---

### 5. **Гибкость при тестировании**

Когда `Pipeline` — "чистый" процессор:
- Его легко **тестировать отдельно** (подать свечу → проверить сигнал)
- Можно **мокать** `Executor`, `Strategy` и т.д.
- Не нужно запускать контекст, `Cancel`, мьютексы

А `Runner` можно тестировать отдельно:
- Запустил → проверил, что `Running = true`
- Остановил → проверил, что `Cancel` вызван

---

## 🔄 Альтернатива: встроить управление в `Pipeline`?

Можно сделать так:

```go
type Pipeline struct {
    ctx     context.Context
    cancel  context.CancelFunc
    running bool
    mu      sync.RWMutex
    // ...
}

func (p *Pipeline) Start() {
    p.mu.Lock()
    if p.running {
        p.mu.Unlock()
        return
    }
    p.ctx, p.cancel = context.WithCancel(context.Background())
    p.running = true
    p.mu.Unlock()

    go p.run(p.ctx)
}
```

Но тогда:
- `Pipeline` становится **непригодным для бектеста**
- Нельзя запустить один `Pipeline` в двух режимах
- Сложнее управлять извне (например, через веб-API)

---

## ✅ Правильная архитектура: **`Runner` управляет `Pipeline`**

```
[TradeManager]
     ↓
[Runner] → управляет жизненным циклом
     ↓
[Pipeline] → только обработка
     ↓
[Strategy, Executor, ...]
```

- `Runner` создаёт `Pipeline` при старте
- Хранит `Cancel`, статус, ошибки
- Может пересоздать `Pipeline` при обновлении стратегии
- `Pipeline` остаётся **чистым, тестируемым, многоразовым**

---

## 💡 Когда можно встроить управление в `Pipeline`?

Только если:
- У вас **один режим** (только live или только backtest)
- Нет веб-интерфейса
- Нет динамического обновления
- Простой скрипт, а не система

Но в вашем случае — **это полноценная торговая система**, и лучше придерживаться **чистой архитектуры**.

---

## 📦 Вывод

> ❌ **Не стоит размещать реквизиты `Runner` в `Pipeline`**, потому что:
>
> 1. Нарушается **принцип единственной ответственности**
> 2. `Pipeline` теряет **гибкость** между режимами
> 3. Сложнее **тестировать** и **управлять**
> 4. Появляется **избыточность** в бектесте
> 5. Усложняется **централизованное управление**

✅ **Оставьте `Runner` отдельно** — он является **обёрткой управления**, а `Pipeline` — **движком обработки**.

---

Если хочешь — могу показать, как сделать **"перезагрузку" Pipeline без остановки Runner** при обновлении стратегии.

### USER
Как тогда остановить цикл выполнения в pipeline из runner

### ASSISTANT
Отличный вопрос! Если `Runner` управляет `Pipeline`, но **не является его частью**, то как **остановить бесконечный цикл** в `Pipeline` (например, приёма свечей в live-режиме) из `Runner`?

---

## ✅ Ответ: через **`context.Context`**

Это **канонический способ** остановки горутин в Go.  
`Runner` создаёт `context.WithCancel()`, передаёт его в `Pipeline.Run(ctx)`, и при остановке вызывает `cancel()`.

---

## 🔧 Как это работает

```
[Runner] → создаёт context + cancel → передаёт в Pipeline.Run(ctx)
                                      ↓
[Pipeline] → читает из канала ДО ТЕХ ПОР, пока ctx.Done() не сработал
                                      ↓
[Runner.Stop()] → вызывает cancel() → цикл в Pipeline завершается
```

---

## 🧩 Реализация

### 1. Метод `Run` в `Pipeline` принимает `context.Context`

```go
// pkg/pipeline/pipeline.go
func (p *Pipeline) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            log.Println("Pipeline stopped: context cancelled")
            return ctx.Err()

        case candle, ok := <-p.Source.Stream():
            if !ok {
                return errors.New("source closed")
            }

            // Обработка свечи: кластеризация, стратегия, валидация и т.д.
            p.processCandle(candle)
        }
    }
}
```

> ⚠️ Важно: **все блокирующие операции** (каналы, таймеры, API-вызовы) должны учитывать `ctx`.

---

### 2. `Runner` создаёт контекст и может его отменить

```go
// pkg/trading/runner.go
type Runner struct {
    Strategy *Strategy
    Pipeline *pipeline.Pipeline
    cancel   context.CancelFunc
    running  bool
    mu       sync.RWMutex
}

func (r *Runner) Start(parentCtx context.Context) {
    r.mu.Lock()
    if r.running {
        r.mu.Unlock()
        return
    }
    r.running = true
    r.mu.Unlock()

    // Создаём дочерний контекст
    ctx, cancel := context.WithCancel(parentCtx)
    r.cancel = cancel

    // Создаём пайплайн
    p, err := pipeline.NewPipelineFromStrategy(r.Strategy)
    if err != nil {
        log.Printf("Pipeline create error: %v", err)
        r.setRunning(false)
        return
    }
    r.Pipeline = p

    // Запускаем пайплайн в горутине
    go func() {
        defer r.onFinish()
        err := p.Run(ctx) // ← передаём контекст
        if err != nil && err != context.Canceled {
            r.mu.Lock()
            r.Err = err
            r.mu.Unlock()
        }
    }()
}

func (r *Runner) Stop() {
    r.mu.RLock()
    cancel := r.cancel
    r.mu.RUnlock()

    if cancel != nil {
        cancel() // ← это остановит Pipeline
    }
}

func (r *Runner) onFinish() {
    r.mu.Lock()
    r.running = false
    r.cancel = nil
    r.mu.Unlock()
}

func (r *Runner) IsRunning() bool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.running
}
```

---

## 🔄 Что происходит при остановке?

1. Вызов `runner.Stop()`
2. Вызов `r.cancel()` → сигнал в `ctx.Done()`
3. В `Pipeline.Run`:
   ```go
   select {
   case <-ctx.Done():
       return ctx.Err() // выход из цикла
   ```
4. Пайплайн завершается **корректно**, без `panic` или `os.Exit`

---

## 💡 Важные моменты

### 1. **Все каналы и операции — с `ctx`**
Если `Pipeline` делает HTTP-запросы или ждёт из другого канала — тоже используй `ctx`:

```go
resp, err := http.Get("https://api.example.com")
// ❌ плохо — нет таймаута

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req)
// ✅ правильно — прервётся при ctx.Done()
```

---

### 2. **Не используй `for { ... }` без `select`**
```go
for {
    candle := <-ch // ❌ блокирует, не реагирует на ctx
}
```

✅ Правильно:
```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case candle := <-ch:
        // обработка
    }
}
```

---

### 3. **Гарантия завершения**
Даже если `Source.Stream()` — это канал, который **никогда не закрывается**, `ctx.Done()` всё равно прервёт ожидание.

---

## 📦 Пример: Реальный источник с контекстом

```go
// pkg/data/realtime.go
func (s *RealTimeSource) Stream() <-chan *Candle {
    out := make(chan *Candle)
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer close(out)
        for {
            select {
            case <-s.ctx.Done(): // контекст от Runner
                return
            case <-ticker.C:
                candle := s.fetchCandle()
                select {
                case out <- candle:
                case <-s.ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}
```

---

## ✅ Вывод

Чтобы остановить `Pipeline` из `Runner`:

1. ✅ `Runner` создаёт `ctx, cancel := context.WithCancel(...)`
2. ✅ Передаёт `ctx` в `pipeline.Run(ctx)`
3. ✅ При `runner.Stop()` вызывает `cancel()`
4. ✅ `Pipeline` проверяет `ctx.Done()` в `select`
5. ✅ Цикл завершается **мягко и безопасно**

> 🔥 Это **идиоматичный способ** управления жизненным циклом в Go.

---

Если хочешь — могу показать, как **протестировать остановку** с помощью `context.WithTimeout` или моков.

### USER
Как при перехвате сигнала прерывания основного приложения выполнить остановку runner

### ASSISTANT
Отличный вопрос! Когда вы запускаете торговое приложение, важно корректно обработать **сигналы ОС** (например, `SIGINT` — `Ctrl+C`, `SIGTERM` — graceful shutdown в Docker), чтобы:

- Остановить все запущенные `Runner`'ы
- Завершить `Pipeline`'ы
- Закрыть соединения с БД, биржей
- Не потерять данные

---

## ✅ Цель

При получении сигнала (например, `Ctrl+C`) — **мягко остановить все стратегии**, которыми управляет `TradeManager`.

---

## 🧩 Решение: Использовать `signal.Notify` + общий `context`

### Шаги:
1. Создать **главный контекст**, который отменяется при сигнале
2. Передать его в `TradeManager` и `Runner`'ы
3. При сигнале — вызвать `cancel()`, все `Runner`'ы остановят свои `Pipeline`

---

## ✅ Реализация

### 1. Перехват сигнала и отмена контекста

```go
// cmd/web/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "your-project/pkg/storage"
    "your-project/pkg/strategy"
    "your-project/pkg/trading"
    "your-project/pkg/web"
    "github.com/gin-gonic/gin"
)

func main() {
    db := storage.InitDB()
    repo := strategy.NewRepository(db)
    manager := trading.NewTradeManager(repo)

    // === Создаём контекст с возможностью отмены ===
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // гарантированно освобождаем ресурсы

    // === Запускаем менеджер ===
    if err := manager.LoadAndStartAll(); err != nil {
        log.Printf("Warning: failed to start some strategies: %v", err)
    }

    // === Перехватываем сигналы ===
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        sig := <-c
        log.Printf("Received signal: %s. Shutting down gracefully...", sig)
        cancel() // ← отменяем контекст → остановка всех Runner'ов

        // Дополнительная задержка на завершение (опционально)
        time.Sleep(1 * time.Second)
        os.Exit(0)
    }()

    // === Запускаем веб-сервер ===
    r := gin.Default()
    web.SetupRoutes(r, manager)

    go func() {
        if err := r.Run(":8080"); err != nil {
            log.Printf("Server error: %v", err)
        }
    }()

    // === Ждём отмены контекста ===
    <-ctx.Done()

    log.Println("Graceful shutdown...")

    // Можно добавить дополнительную очистку
    manager.StopAll() // явная остановка всех стратегий (если нужно)

    // Закрываем БД и т.п.
    db.Close()

    log.Println("Bye!")
}
```

---

### 2. `TradeManager` реагирует на отмену контекста

```go
// pkg/trading/trade_manager.go
func (tm *TradeManager) StopAll() {
    tm.mu.RLock()
    runners := make([]*Runner, 0, len(tm.runners))
    for _, r := range tm.runners {
        runners = append(runners, r)
    }
    tm.mu.RUnlock()

    for _, runner := range runners {
        if runner.IsRunning() {
            log.Printf("Stopping strategy: %s", runner.Strategy.Name)
            runner.Stop() // вызывает cancel() внутри
        }
    }
}
```

> ⚠️ `runner.Stop()` уже использует `context.Cancel`, так что остановка будет **каскадной**.

---

### 3. `Runner` реагирует на `cancel`

Как мы делали ранее — `Runner.Start(ctx)` принимает контекст, и при `cancel()` извне (или из `runner.Stop()`) пайплайн завершается.

```go
func (r *Runner) Start(parentCtx context.Context) {
    ctx, cancel := context.WithCancel(parentCtx)
    r.cancel = cancel

    go func() {
        defer r.onFinish()
        r.Pipeline.Run(ctx) // ← ctx.Done() остановит цикл
    }()
}
```

---

## 🔄 Что происходит при `Ctrl+C`?

```
Ctrl+C
   ↓
signal.Notify получает os.Interrupt
   ↓
Вызывается cancel() у главного контекста
   ↓
Все Runner'ы получают сигнал через ctx.Done()
   ↓
Pipeline прерывает цикл обработки
   ↓
Runner завершает работу
   ↓
main завершает веб-сервер, закрывает БД, выходит
```

---

## ✅ Дополнительно: Таймаут graceful shutdown

Можно дать системе **максимум N секунд** на завершение:

```go
sig := <-c
log.Printf("Signal received: %s", sig)

ctx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
defer shutdownCancel()

go func() {
    <-ctx.Done()
    if ctx.Err() == context.DeadlineExceeded {
        log.Fatal("Graceful shutdown timeout, forcing exit")
    }
}()

manager.StopAll()

// Дожидаемся завершения
<-ctx.Done()
```

---

## ✅ Проверка: Будет ли остановлен бектест?

Да, если `RunBacktest` тоже учитывает `ctx`:

```go
func (p *Pipeline) RunBacktest(ctx context.Context) {
    for candle := range p.Source.Stream() {
        select {
        case <-ctx.Done():
            return
        default:
            p.processCandle(candle)
        }
    }
}
```

---

## ✅ Вывод

Чтобы остановить `Runner` при сигнале прерывания:

1. ✅ Создай **главный `context`** в `main`
2. ✅ Передай его в `Runner` через `Start(ctx)`
3. ✅ Используй `signal.Notify` для перехвата `SIGINT` / `SIGTERM`
4. ✅ При сигнале вызови `cancel()` → все `Pipeline` остановятся
5. ✅ Дополнительно: вызови `manager.StopAll()` для явной очистки

---

✅ Это **стандартный, надёжный и идиоматичный способ** graceful shutdown в Go.

---

Если хочешь — могу показать, как **протестировать обработку сигнала** с помощью `testify/suite` или `os.Pipe`.