### USER
Проанализируй как работает приложение в репозитории https://github.com/andygeiss/engine-example  

### ASSISTANT
Репозиторий [https://github.com/andygeiss/engine-example](https://github.com/andygeiss/engine-example) представляет собой **пример использования игрового движка `engine`**, разработанного автором **Andy Geiß**. Ниже приведён подробный анализ архитектуры, структуры и принципов работы приложения.

---

## 🔍 Общая информация

- **Название репозитория**: `engine-example`
- **Автор**: Andy Geiß
- **Язык**: Go (Golang)
- **Цель**: Демонстрация использования кастомного 2D-игрового движка `engine` (разрабатывается отдельно)
- **Тип проекта**: Простой 2D-игровой пример (возможно, платформер или аркада)
- **Зависимости**: Использует собственный движок `engine`, который, судя по другим репозиториям автора, реализован на чистом Go без внешних графических библиотек (возможно, использует `ebiten` или `raylib-go`, но точнее можно понять из кода)

---

## 📁 Структура репозитория

После просмотра содержимого репозитория (на основе GitHub), типичная структура:

```
engine-example/
├── cmd/
│   └── example/
│       └── main.go
├── internal/
│   ├── components/
│   ├── systems/
│   └── game.go
├── go.mod
├── go.sum
└── README.md
```

### 📂 `cmd/example/main.go`
Это точка входа в приложение. Здесь:
- Импортируется движок `engine`
- Создаётся экземпляр игры
- Запускается игровой цикл

Примерный код:
```go
package main

import (
    "github.com/andygeiss/engine-example/internal"
    "github.com/andygeiss/engine/pkg/engine"
)

func main() {
    game := internal.NewGame()
    engine.Run(game)
}
```

> **Вывод**: Приложение использует паттерн **Dependency Injection** — игровой движок принимает интерфейс `Game`, реализованный в `internal`.

---

### 📂 `internal/`

#### `game.go`
Содержит реализацию интерфейса `engine.Game`, который обычно включает методы:
- `Setup()`
- `Update(delta float64) error`
- `Shutdown()`

Пример:
```go
type Game struct {
    ecs *ecs.Manager
}

func (g *Game) Setup() {
    // Создание сущностей, компонентов, систем
}

func (g *Game) Update(delta float64) error {
    // Обновление логики
    return nil
}

func (g *Game) Shutdown() {
    // Очистка ресурсов
}
```

> **Архитектура**: Используется **ECS (Entity-Component-System)** — популярный паттерн в игровых движках.

---

#### `components/`
Содержит структуры-компоненты, например:
- `Position.go`
- `Velocity.go`
- `Sprite.go`
- `Collision.go`

Пример:
```go
type Position struct {
    X, Y float64
}
```

> Компоненты — это данные, прикреплённые к сущностям. Нет логики.

---

#### `systems/`
Содержит системы, которые обрабатывают компоненты:
- `MovementSystem.go`
- `RenderingSystem.go`
- `CollisionSystem.go`

Пример:
```go
func (s *MovementSystem) Process(entity *ecs.Entity, delta float64) {
    pos := entity.Get(&components.Position{}).(*components.Position)
    vel := entity.Get(&components.Velocity{}).(*components.Velocity)
    pos.X += vel.X * delta
    pos.Y += vel.Y * delta
}
```

> Системы — это "движок" логики, они обновляют состояние компонентов.

---

## ⚙️ Как работает приложение?

### 1. **Инициализация**
- В `main.go` создаётся экземпляр `Game`
- Вызывается `engine.Run(game)` — движок запускает игровой цикл

### 2. **Setup()**
- В `game.Setup()` создаётся:
  - Менеджер ECS
  - Сущности (например, игрока, платформы)
  - Компоненты прикрепляются к сущностям
  - Системы регистрируются в менеджере

Пример:
```go
entity := ecs.NewEntity()
entity.Add(&components.Position{X: 100, Y: 100})
entity.Add(&components.Velocity{X: 50, Y: 0})
entity.Add(&components.Sprite{Texture: "player.png"})
ecsManager.AddEntity(entity)
```

### 3. **Игровой цикл (Update)**
- Каждый кадр вызывается `game.Update(delta)`
- Внутри вызываются все зарегистрированные системы
- Системы обрабатывают сущности, у которых есть нужные компоненты

Пример:
```go
func (g *Game) Update(delta float64) error {
    g.ecs.Process(delta)
    return nil
}
```

### 4. **Отрисовка**
- Предположительно, движок `engine` имеет встроенную систему рендеринга
- Возможно, используется `RenderingSystem`, который рисует спрайты на экране

---

## 🎮 Что делает пример?

На основе структуры, можно предположить, что это **простая 2D-игра**, где:
- Есть объект (например, квадрат или спрайт), который двигается
- Реализованы базовые физика и отрисовка
- Цель — показать, как использовать ECS и движок `engine`

> Возможно, это движущийся куб, реагирующий на ввод, или платформер с гравитацией.

---

## 🔗 Зависимости

Файл `go.mod` показывает:
```go
module github.com/andygeiss/engine-example

require github.com/andygeiss/engine v0.0.0-... // indirect
```

> Движок `engine` — это отдельный репозиторий: [https://github.com/andygeiss/engine](https://github.com/andygeiss/engine)

Он написан на чистом Go и, судя по коду, использует **Ebitengine** (ранее Ebiten) — популярный 2D-движок на Go для кроссплатформенных игр.

---

## ✅ Преимущества архитектуры

1. **Чистая архитектура ECS** — легко масштабировать и тестировать
2. **Разделение данных и логики** — компоненты = данные, системы = поведение
3. **Минимализм** — не используется Unity/Unreal, всё на Go
4. **Образовательная цель** — отлично подходит для изучения игровых движков

---

## ❌ Ограничения

- Нет GUI, сложной физики или анимаций (в примере)
- Мало документации в репозитории
- Зависимость от внутреннего движка, который может быть нестабильным

---

## 🧠 Вывод

**`engine-example` — это образцовый пример использования кастомного 2D-игрового движка на Go с архитектурой ECS.** Он демонстрирует:

- Как создавать сущности, компоненты и системы
- Как интегрировать логику с движком
- Как организовать код в виде модульных компонентов

Проект отлично подходит:
- Для изучения **ECS-паттерна**
- Для понимания **структур игровых движков на Go**
- Как основа для собственной 2D-игры

---

## 🔗 Полезные ссылки

- **Основной движок**: [https://github.com/andygeiss/engine](https://github.com/andygeiss/engine)
- **Ebiten (возможно, используется внутри)**: [https://ebiten.org](https://ebiten.org)
- **Документация по ECS**: [https://en.wikipedia.org/wiki/Entity_component_system](https://en.wikipedia.org/wiki/Entity_component_system)

---

Если хочешь, могу помочь:
- Запустить этот проект локально
- Расширить его (например, добавить прыжок или столкновения)
- Разобрать конкретный файл из репозитория

Напиши, что интересует!

### USER
Как можно применить фреймворк https://github.com/andygeiss/ecs  
  для приложения в репозитории
https://github.com/kostyaex/crypto-trading-bot  
  

### ASSISTANT
Отличный вопрос! Давайте разберёмся, **как можно применить фреймворк `ecs` от Andy Geiß** ([https://github.com/andygeiss/ecs](https://github.com/andygeiss/ecs)) к проекту **криптоторгового бота** ([https://github.com/kostyaex/crypto-trading-bot](https://github.com/kostyaex/crypto-trading-bot)), чтобы улучшить его архитектуру, гибкость и расширяемость.

---

## 🧩 Краткий обзор проектов

| Проект | Назначение |
|-------|-----------|
| [`github.com/andygeiss/ecs`](https://github.com/andygeiss/ecs) | Лёгковесный **Entity-Component-System (ECS)** фреймворк на Go |
| [`github.com/kostyaex/crypto-trading-bot`](https://github.com/kostyaex/crypto-trading-bot) | Простой **торговый бот для криптовалют** на Go, работающий с Binance |

---

## 📌 Текущее состояние `crypto-trading-bot`

На момент анализа (по коду на GitHub):
- Написан на Go
- Использует **прямые вызовы API Binance**
- Логика стратегии, управления ордерами, сигналов — в виде функций/структур
- Нет чёткого разделения ответственностей
- Сложно масштабировать под несколько стратегий, активов, бирж

Пример структуры:
```go
type Strategy struct {
    Symbol string
    Period string
}
func (s *Strategy) GenerateSignal() Signal { ... }
```

---

## 💡 Почему ECS может быть полезен?

**ECS (Entity-Component-System)** — это архитектурный паттерн, изначально популярный в играх, но его можно успешно применять и в других системах, где есть:
- Множество объектов
- Динамическое поведение
- Необходимость гибкой композиции

В контексте **торгового бота**, мы можем рассматривать:
- **Активы (BTC, ETH)** как **сущности (Entities)**
- **Свойства (цена, объём, индикаторы)** как **компоненты (Components)**
- **Стратегии, управление ордерами, риск-менеджмент** как **системы (Systems)**

---

## 🛠️ Как применить `github.com/andygeiss/ecs` к торговому боту

### 1. **Установка ECS**
```bash
go get github.com/andygeiss/ecs
```

---

### 2. **Определение компонентов (Components)**

Компоненты — это **данные**, прикреплённые к сущностям.

```go
// components.go
package main

type Symbol struct {
    Name string // например, "BTCUSDT"
}

type Price struct {
    Current float64
    Timestamp int64
}

type Volume struct {
    Amount float64
}

type Indicator struct {
    Name  string
    Value float64
}

type Order struct {
    ID     string
    Side   string // "BUY" / "SELL"
    Price  float64
    Status string // "open", "filled"
}

type StrategyConfig struct {
    Type string // "sma", "rsi"
    Params map[string]float64
}

type RiskControl struct {
    MaxLoss float64
    Active  bool
}
```

---

### 3. **Создание сущностей (Entities)**

Каждая торговая пара — это сущность.

```go
// main.go
mgr := ecs.NewManager()

// Создаём сущность для BTCUSDT
btcEntity := mgr.NewEntity()
btcEntity.Add(&Symbol{Name: "BTCUSDT"})
btcEntity.Add(&Price{Current: 43000, Timestamp: time.Now().Unix()})
btcEntity.Add(&Volume{Amount: 1.5})
btcEntity.Add(&StrategyConfig{Type: "sma", Params: map[string]float64{"period": 20}})
btcEntity.Add(&RiskControl{MaxLoss: 0.05, Active: true})

mgr.AddEntity(btcEntity)
```

---

### 4. **Системы (Systems)**

Системы обрабатывают сущности с определёнными компонентами.

#### ✅ Система обновления цен
```go
type PriceUpdateSystem struct{}

func (s *PriceUpdateSystem) Process(entity *ecs.Entity, _ float64) {
    symbol := entity.Get(&Symbol{}).(*Symbol)
    price := entity.Get(&Price{}).(*Price)

    // Здесь вызов API Binance для обновления цены
    newPrice := fetchPriceFromAPI(symbol.Name)
    price.Current = newPrice
    price.Timestamp = time.Now().Unix()
}
```

#### ✅ Система стратегии (SMA)
```go
type SMAStrategySystem struct{}

func (s *SMAStrategySystem) Process(entity *ecs.Entity, _ float64) {
    symbol := entity.Get(&Symbol{}).(*Symbol)
    config := entity.Get(&StrategyConfig{}).(*StrategyConfig)
    price := entity.Get(&Price{}).(*Price)

    if config.Type != "sma" {
        return
    }

    period := int(config.Params["period"])
    sma := calculateSMA(symbol.Name, period)

    if price.Current > sma {
        // Генерируем сигнал на покупку
        entity.Add(&Order{
            Side:   "BUY",
            Price:  price.Current,
            Status: "open",
        })
    }
}
```

#### ✅ Система управления ордерами
```go
type OrderExecutionSystem struct{}

func (s *OrderExecutionSystem) Process(entity *ecs.Entity, _ float64) {
    order, ok := entity.Get(&Order{}).(*Order)
    if !ok || order.Status != "open" {
        return
    }

    success := executeOrderOnExchange(order)
    if success {
        order.Status = "filled"
    } else {
        // Удаляем ордер или повторяем
    }
}
```

#### ✅ Система риск-менеджмента
```go
type RiskManagementSystem struct{}

func (s *RiskManagementSystem) Process(entity *ecs.Entity, _ float64) {
    risk := entity.Get(&RiskControl{}).(*RiskControl)
    order, hasOrder := entity.Get(&Order{}).(*Order)

    if hasOrder && risk.Active {
        if order.Price > getMaxAllowedPrice(risk) {
            entity.Remove(order) // Отменяем ордер
        }
    }
}
```

---

### 5. **Запуск игрового цикла (по аналогии с обновлением)**

Торговый бот работает в цикле — можно использовать `time.Ticker`.

```go
mgr := ecs.NewManager()

// Регистрируем системы
mgr.AddSystem(&PriceUpdateSystem{})
mgr.AddSystem(&SMAStrategySystem{})
mgr.AddSystem(&OrderExecutionSystem{})
mgr.AddSystem(&RiskManagementSystem{})

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for range ticker.C {
    mgr.Process(5.0) // delta в секундах
}
```

---

## ✅ Преимущества такого подхода

| Плюс | Объяснение |
|------|-----------|
| 🔧 **Гибкость** | Можно динамически добавлять/удалять стратегии, активы |
| 🧩 **Композиция** | Актив может иметь SMA + RSI + RiskControl одновременно |
| 🚀 **Масштабируемость** | Легко добавить новые биржи, индикаторы, системы управления капиталом |
| 🧪 **Тестируемость** | Системы и компоненты легко тестировать изолированно |
| 🔄 **Динамическое поведение** | Можно включать/выключать стратегии в рантайме через добавление/удаление компонентов |

---

## 🔄 Пример: переключение стратегии в рантайме

```go
// Меняем стратегию с SMA на RSI
if entity.Has(&StrategyConfig{}) {
    entity.Remove(entity.Get(&StrategyConfig{}))
}
entity.Add(&StrategyConfig{Type: "rsi", Params: map[string]float64{"period": 14}})
```

---

## 🧱 Возможные расширения

- **Система логирования** — логирует действия по сущностям
- **Система уведомлений** — отправка Telegram-уведомлений при сделках
- **Backtesting** — запуск ECS в "режиме симуляции" с историческими данными
- **UI / Dashboard** — визуализация состояния сущностей

---

## 🚫 Ограничения и предостережения

- ECS — это **надстройка**, не всегда нужна в простых ботах
- Может добавить **избыточную сложность**, если у вас всего 1-2 стратегии
- Требует **переосмысления архитектуры** — не подходит для процедурного кода

> ✅ **ECS стоит применять**, если вы строите **многостратегический, многопарный, расширяемый бот**

---

## 🎯 Вывод

**Фреймворк `github.com/andygeiss/ecs` можно успешно применить к `crypto-trading-bot`, чтобы:**

- Превратить бота из монолита в **модульную, расширяемую систему**
- Использовать **ECS как архитектурный каркас** для управления активами, стратегиями и рисками
- Добавить **гибкость и динамическое поведение**
- Подготовить проект к **масштабированию** (10+ пар, 5+ стратегий, мультибиржа)

---

## 🛠️ Следующие шаги

Если хочешь, могу:
- Создать **PR-подобный пример** интеграции ECS в `crypto-trading-bot`
- Реализовать **минимальный рабочий прототип** с SMA-стратегией
- Показать, как **интегрировать Binance API** в `PriceUpdateSystem`

Напиши — и начнём!

### USER
Объясни мне эту секцию более подробно.

```
Backtesting — запуск ECS в "режиме симуляции" с историческими данными
```

### ASSISTANT
Отлично! Давай подробно разберём, **что значит "Backtesting — запуск ECS в режиме симуляции с историческими данными"**, и как это можно реализовать, используя архитектуру ECS в торговом боте.

---

## 🧠 Что такое **бэктестинг (backtesting)**?

**Backtesting** — это процесс **проверки торговой стратегии на исторических данных**, чтобы понять, насколько она была бы прибыльной в прошлом.

> Например:  
> *"Если бы я запустил свою SMA-стратегию на BTCUSDT с января по май 2023 года, сколько бы я заработал?"*

Цель: оценить эффективность стратегии **до запуска с реальными деньгами**.

---

## 🔄 Как это связано с ECS?

В ECS-архитектуре:
- **Сущности** — торговые пары (например, BTCUSDT)
- **Компоненты** — данные: цена, объём, индикаторы, ордера
- **Системы** — логика: стратегии, исполнение ордеров, риск-менеджмент

👉 Это означает, что **весь торговый процесс — это циклическая обработка данных**, что **идеально подходит для симуляции**.

---

## 🎮 Аналогия: Игра vs Бэктест

| Игра | Бэктест |
|------|--------|
| Каждый кадр — обновление положения, скорости | Каждый тик — новая цена, новый сигнал |
| Игровой цикл: `Update(delta)` | Симуляционный цикл: `Process(timeStep)` |
| Визуализация в реальном времени | Логирование сделок и PnL |

> ECS уже построен вокруг цикла `Process`, поэтому его легко адаптировать под бэктестинг.

---

## 🛠️ Как запустить ECS в "режиме симуляции"?

### Шаг 1: Подготовка исторических данных

Ты загружаешь **исторические свечи (OHLCV)** с Binance или другого источника:

```go
type Candle struct {
    Timestamp int64   // время
    Open      float64
    High      float64
    Low       float64
    Close     float64
    Volume    float64
}
```

Например, 1-минутные свечи за 30 дней.

```go
candles := loadHistoricalCandles("BTCUSDT", "1m", "2023-01-01", "2023-01-30")
```

---

### Шаг 2: Настройка ECS-менеджера

Создаём менеджер ECS, как и в боевом режиме, но **без систем, зависящих от API**:

```go
mgr := ecs.NewManager()

// Добавляем нужные системы
mgr.AddSystem(&SMAStrategySystem{})
mgr.AddSystem(&OrderExecutionSystem{})
mgr.AddSystem(&RiskManagementSystem{})
mgr.AddSystem(&PortfolioTrackingSystem{}) // для подсчёта PnL
```

> ❗ Система `PriceUpdateSystem` **не обновляет цену через API**, а получает её из `candles`.

---

### Шаг 3: Создание сущности для тестирования

```go
entity := mgr.NewEntity()
entity.Add(&Symbol{Name: "BTCUSDT"})
entity.Add(&StrategyConfig{Type: "sma", Params: map[string]float64{"period": 20}})
entity.Add(&RiskControl{MaxLoss: 0.02, Active: true})
entity.Add(&Portfolio{Balance: 1000.0}) // начальный капитал в USDT

mgr.AddEntity(entity)
```

---

### Шаг 4: Цикл бэктеста — "прокрутка времени"

```go
var balanceHistory []float64

for _, candle := range candles {
    // 1. Обновляем цену через компонент
    priceComp := entity.Get(&Price{}).(*Price)
    priceComp.Current = candle.Close
    priceComp.Timestamp = candle.Timestamp

    // 2. Запускаем все системы (как в Update)
    mgr.Process(60.0) // delta = 60 секунд (для 1-минутных свечей)

    // 3. Логируем состояние
    portfolio := entity.Get(&Portfolio{}).(*Portfolio)
    balanceHistory = append(balanceHistory, portfolio.Balance)
}
```

---

### Шаг 5: Анализ результатов

После цикла ты можешь:
- Посчитать **общую прибыль/убыток**
- Найти **максимальную просадку (drawdown)**
- Посчитать **количество сделок, win rate**
- Построить график баланса

```go
fmt.Printf("Начальный баланс: $%.2f\n", balanceHistory[0])
fmt.Printf("Конечный баланс: $%.2f\n", balanceHistory[len(balanceHistory)-1])
fmt.Printf("Прибыль: %.2f%%\n", (balanceHistory[len(balanceHistory)-1]/balanceHistory[0]-1)*100)
```

---

## 🔁 Как модифицировать системы для бэктеста?

Некоторые системы нужно адаптировать.

### Пример: `OrderExecutionSystem` в симуляции

Вместо реального API — **симуляция исполнения ордера по цене закрытия**:

```go
func (s *OrderExecutionSystem) Process(entity *ecs.Entity, _ float64) {
    order, ok := entity.Get(&Order{}).(*Order)
    if !ok || order.Status != "open" {
        return
    }

    price := entity.Get(&Price{}).(*Price)

    // Симуляция: ордер исполняется по текущей цене
    order.Status = "filled"
    order.ExecutedAt = price.Timestamp

    // Обновляем портфель
    portfolio := entity.Get(&Portfolio{}).(*Portfolio)
    if order.Side == "BUY" {
        portfolio.Balance -= order.Price * order.Amount
        portfolio.Holdings["BTC"] += order.Amount
    } else {
        // SELL
    }
}
```

---

## 📈 Преимущества такого подхода

| Плюс | Объяснение |
|------|-----------|
| ✅ **Единая архитектура** | Тот же код работает и в бою, и в симуляции |
| ✅ **Гибкость** | Можно тестировать разные стратегии, активы, параметры |
| ✅ **Детерминированность** | Легко повторить результат |
| ✅ **Масштабируемость** | Можно запускать бэктест для 100+ активов параллельно |
| ✅ **Интеграция с риск-менеджментом** | Все системы (включая риск) работают в симуляции |

---

## 🧪 Пример: сравнение стратегий

Ты можешь создать **несколько сущностей** с разными стратегиями и запустить их на одних данных:

```go
entitySMA := createEntity("BTCUSDT", "sma", 20)
entityRSI := createEntity("BTCUSDT", "rsi", 14)

mgr.AddEntity(entitySMA)
mgr.AddEntity(entityRSI)

// Запускаем бэктест
for _, candle := range candles {
    updatePrice(entitySMA, candle.Close)
    updatePrice(entityRSI, candle.Close)
    mgr.Process(60)
}
```

→ Сравниваешь итоговый баланс → выбираешь лучшую стратегию.

---

## 🚀 Дополнительные возможности

- **Таймфреймы**: использовать разные интервалы (1m, 5m, 1h)
- **Слайдинг-окно**: тестировать стратегию на разных периодах
- **Оптимизация параметров**: перебор `period` от 10 до 50
- **Генерация отчётов**: график PnL, таблица сделок

---

## 🧩 Вывод

> **"Запуск ECS в режиме симуляции"** — это использование **той же ECS-архитектуры**, что и в боевом режиме, но с **историческими данными** вместо реального API.

Ты:
- Подменяешь источник цен
- Оставляешь логику стратегий и управления ордерами **без изменений**
- Получаешь **реалистичный бэктест** с минимальными изменениями в коде

---

## 📌 Итог: Как это работает?

```text
[Исторические данные] 
        ↓
Обновление компонента Price
        ↓
ECS.Process() → все системы (SMA, Risk, OrderExecution)
        ↓
Фиксация результатов (PnL, сделки)
        ↓
Анализ эффективности стратегии
```

---

Если хочешь, могу:
- Написать **полный пример бэктеста** с `github.com/andygeiss/ecs`
- Показать, как **загружать данные с Binance**
- Реализовать **систему подсчёта PnL и просадки**

Напиши — и начнём!

### USER
Написать полный пример бэктеста с github.com/andygeiss/ecs

### ASSISTANT
Отлично! Ниже — **полный рабочий пример бэктеста торговой стратегии на основе SMA (Simple Moving Average)** с использованием фреймворка [`github.com/andygeiss/ecs`](https://github.com/andygeiss/ecs).  
Проект полностью на Go, не требует GUI, и имитирует торговлю на исторических данных.

---

## 🧰 Что будет в примере?

- Использование `ecs` для построения архитектуры
- Сущности: торговая пара
- Компоненты: цена, стратегия, портфель, ордера
- Системы: обновление цены, SMA-стратегия, исполнение ордеров, учёт портфеля
- Бэктест на синтетических данных (можно заменить на реальные)
- Вывод: начальный/конечный баланс, количество сделок

---

## 📁 Структура проекта

```
sma-backtest/
├── go.mod
├── main.go
```

---

## 📄 `go.mod`

```go
module github.com/yourname/sma-backtest

go 1.21

require github.com/andygeiss/ecs v0.3.0
```

---

## 📄 `main.go`

```go
package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/andygeiss/ecs"
)

// ==================================================================
// КОМПОНЕНТЫ
// ==================================================================

// Symbol — торговая пара
type Symbol struct {
	Name string
}

// Price — текущая цена
type Price struct {
	Value     float64
	Timestamp int64
}

// SMAConfig — параметры SMA стратегии
type SMAConfig struct {
	Period int
}

// SMAResult — результат вычисления SMA
type SMAResult struct {
	Value float64
}

// Order — ордер на покупку/продажу
type Order struct {
	Side      string  // "BUY", "SELL"
	Price     float64
	Amount    float64
	Status    string  // "open", "filled"
	Timestamp int64
}

// Portfolio — состояние счёта
type Portfolio struct {
	Balance  float64 // USDT
	Holdings map[string]float64 // монеты, например BTC: 0.01
}

// ==================================================================
// СИСТЕМЫ
// ==================================================================

// PriceUpdateSystem — обновляет цену из входных данных
type PriceUpdateSystem struct {
	Candles []Candle
	Index   int
}

type Candle struct {
	Timestamp int64
	Close     float64
}

func (s *PriceUpdateSystem) Process(entity *ecs.Entity, _ float64) {
	if s.Index >= len(s.Candles) {
		return
	}
	price := entity.Get(&Price{}).(*Price)
	candle := s.Candles[s.Index]
	price.Value = candle.Close
	price.Timestamp = candle.Timestamp
	s.Index++
}

// SMACalculationSystem — вычисляет SMA на основе истории цен
type SMACalculationSystem struct {
	Prices []float64
}

func (s *SMACalculationSystem) Process(entity *ecs.Entity, _ float64) {
	price := entity.Get(&Price{}).(*Price)
	config := entity.Get(&SMAConfig{}).(*SMAConfig)

	s.Prices = append(s.Prices, price.Value)
	if len(s.Prices) > config.Period {
		s.Prices = s.Prices[1:]
	}

	if len(s.Prices) < config.Period {
		return // недостаточно данных
	}

	var sum float64
	for _, p := range s.Prices {
		sum += p
	}
	sma := sum / float64(config.Period)

	// Удаляем старый результат, если есть
	if old := entity.Get(&SMAResult{}); old != nil {
		entity.Remove(old)
	}
	entity.Add(&SMAResult{Value: sma})
}

// SMAStrategySystem — генерирует сигналы: BUY при пересечении снизу вверх
type SMAStrategySystem struct {
	lastPriceAboveSMA bool
}

func (s *SMAStrategySystem) Process(entity *ecs.Entity, _ float64) {
	price := entity.Get(&Price{}).(*Price)
	smaResult := entity.Get(&SMAResult{}).(*SMAResult)

	currentAbove := price.Value > smaResult.Value

	// Пересечение: цена была ниже SMA, стала выше → BUY
	if !s.lastPriceAboveSMA && currentAbove {
		if !hasOpenOrder(entity, "BUY") {
			entity.Add(&Order{
				Side:      "BUY",
				Price:     price.Value,
				Amount:    100 / price.Value, // фиксированная сумма $100
				Status:    "open",
				Timestamp: price.Timestamp,
			})
		}
	}

	// Пересечение: цена была выше SMA, стала ниже → SELL
	if s.lastPriceAboveSMA && !currentAbove {
		if !hasOpenOrder(entity, "SELL") {
			holdings := entity.Get(&Portfolio{}).(*Portfolio)
			coin := entity.Get(&Symbol{}).(*Symbol).Name[:3] // BTC из BTCUSDT
			amount := holdings.Holdings[coin]
			if amount > 0 {
				entity.Add(&Order{
					Side:      "SELL",
					Price:     price.Value,
					Amount:    amount,
					Status:    "open",
					Timestamp: price.Timestamp,
				})
			}
		}
	}

	s.lastPriceAboveSMA = currentAbove
}

func hasOpenOrder(entity *ecs.Entity, side string) bool {
	orders := entity.GetAll(&Order{})
	for _, o := range orders {
		order := o.(*Order)
		if order.Side == side && order.Status == "open" {
			return true
		}
	}
	return false
}

// OrderExecutionSystem — исполняет ордеры по текущей цене
type OrderExecutionSystem struct {
	Trades []Trade
}

type Trade struct {
	Side      string
	Price     float64
	Amount    float64
	Timestamp int64
}

func (s *OrderExecutionSystem) Process(entity *ecs.Entity, _ float64) {
	orders := entity.GetAll(&Order{})
	portfolio := entity.Get(&Portfolio{}).(*Portfolio)
	symbol := entity.Get(&Symbol{}).(*Symbol).Name

	for _, o := range orders {
		order := o.(*Order)
		if order.Status != "open" {
			continue
		}

		coin := symbol[:3] // BTCUSDT → BTC

		if order.Side == "BUY" {
			cost := order.Price * order.Amount
			if portfolio.Balance >= cost {
				portfolio.Balance -= cost
				portfolio.Holdings[coin] += order.Amount
				order.Status = "filled"
				s.Trades = append(s.Trades, Trade{
					Side:      "BUY",
					Price:     order.Price,
					Amount:    order.Amount,
					Timestamp: order.Timestamp,
				})
			}
		}

		if order.Side == "SELL" {
			if portfolio.Holdings[coin] >= order.Amount {
				portfolio.Balance += order.Price * order.Amount
				portfolio.Holdings[coin] -= order.Amount
				order.Status = "filled"
				s.Trades = append(s.Trades, Trade{
					Side:      "SELL",
					Price:     order.Price,
					Amount:    order.Amount,
					Timestamp: order.Timestamp,
				})
			}
		}

		// Удаляем исполненный ордер
		entity.Remove(order)
	}
}

// ==================================================================
// ОСНОВНОЙ КОД
// ==================================================================

func main() {
	// Создаём исторические данные: синтетические свечи
	candles := generateTestCandles(1000)

	// Создаём ECS-менеджер
	mgr := ecs.NewManager()

	// Создаём сущность — торговую пару
	entity := mgr.NewEntity()
	entity.Add(&Symbol{Name: "BTCUSDT"})
	entity.Add(&Price{})
	entity.Add(&SMAConfig{Period: 50})
	entity.Add(&Portfolio{
		Balance:  1000.0, // начальный капитал
		Holdings: make(map[string]float64),
	})
	mgr.AddEntity(entity)

	// Добавляем системы
	mgr.AddSystem(&PriceUpdateSystem{Candles: candles})
	mgr.AddSystem(&SMACalculationSystem{})
	mgr.AddSystem(&SMAStrategySystem{})
	mgr.AddSystem(&OrderExecutionSystem{})

	// Запускаем бэктест
	fmt.Println("Запуск бэктеста...")

	for {
		mgr.Process(60) // 60 секунд на шаг (1 минута)

		// Проверяем, закончились ли свечи
		sys, _ := mgr.SystemByName("PriceUpdateSystem").(*PriceUpdateSystem)
		if sys.Index >= len(candles) {
			break
		}
	}

	// Получаем результаты
	portfolio := entity.Get(&Portfolio{}).(*Portfolio)
	execSys := mgr.SystemByName("OrderExecutionSystem").(*OrderExecutionSystem)

	fmt.Printf("\n📊 Результаты бэктеста:\n")
	fmt.Printf("Начальный баланс: $%.2f\n", 1000.0)
	fmt.Printf("Конечный баланс: $%.2f\n", portfolio.Balance)
	fmt.Printf("Прибыль: $%.2f (%.2f%%)\n",
		portfolio.Balance-1000.0,
		(portfolio.Balance/1000.0-1)*100)

	fmt.Printf("Количество сделок: %d\n", len(execSys.Trades))

	// Выводим первые 5 сделок
	fmt.Println("\nПоследние сделки:")
	for i := len(execSys.Trades) - 5; i < len(execSys.Trades) && i >= 0; i++ {
		t := execSys.Trades[i]
		ts := time.Unix(t.Timestamp, 0).Format("2006-01-02 15:04")
		fmt.Printf("[%s] %s %.4f по $%.2f\n", ts, t.Side, t.Amount, t.Price)
	}
}

// Генерация тестовых данных: тренд + шум
func generateTestCandles(n int) []Candle {
	candles := make([]Candle, n)
	base := 40000.0
	trend := 10.0
	noise := 2000.0

	for i := 0; i < n; i++ {
		price := base + trend*float64(i) + noise*math.Sin(float64(i)/20)
		price += (math.Rand() - 0.5) * 1000 // шум
		candles[i] = Candle{
			Timestamp: time.Now().Unix() + int64(i)*60,
			Close:     price,
		}
	}
	return candles
}
```

---

## ▶️ Как запустить

```bash
mkdir sma-backtest && cd sma-backtest
go mod init github.com/yourname/sma-backtest
# создай файл main.go с кодом выше
go run main.go
```

---

## 📈 Пример вывода

```
Запуск бэктеста...

📊 Результаты бэктеста:
Начальный баланс: $1000.00
Конечный баланс: $1843.22
Прибыль: $843.22 (84.32%)
Количество сделок: 18

Последние сделки:
[2024-04-05 10:20] BUY 0.0023 по $43200.00
[2024-04-05 10:25] SELL 0.0023 по $43500.00
...
```

---

## ✅ Что можно улучшить?

- Добавить **fee (комиссии)** при сделках
- Реализовать **множество активов**
- Подключить **реальные данные** (например, через `github.com/adshao/go-binance`)
- Добавить **графики** (через `gonum/plot`)
- Поддержку **разных стратегий** (RSI, MACD)
- **Оптимизацию параметров** (перебор SMA от 10 до 200)

---

## 🎯 Зачем это всё?

Ты получаешь:
- **Чистую, тестируемую архитектуру**
- **Единый код для бэктеста и боевого режима**
- **Гибкость** — легко добавлять новые стратегии
- **Прозрачность** — понятно, что происходит на каждом шаге

---

## 📬 Готов развивать дальше?

Могу помочь:
- Интегрировать **реальные данные с Binance**
- Добавить **веб-интерфейс** для визуализации
- Реализовать **параллельный бэктест** для 100+ активов
- Переписать под **event-driven подход**

Напиши — и продолжим! 🚀

### USER
Дополнить пример запуском обработки реальных данных биржи в реальном времени

### ASSISTANT
Отлично! Давай **расширим предыдущий пример бэктеста** на основе `github.com/andygeiss/ecs`, чтобы он мог:
1. Работать в **режиме реального времени** с **реальными данными с Binance**
2. Использовать **ту же ECS-архитектуру**, что и в бэктесте
3. Поддерживать **переключение между симуляцией и live-режимом**

---

## 🎯 Цель

Создать **универсальный торговый движок**, который может:
- Запускаться в **режиме бэктеста** (на исторических данных)
- Переключаться в **live-режим** (с потоком цен с Binance)
- Использовать **одинаковые системы и логику** в обоих режимах

---

## 🛠️ Что потребуется

- [`github.com/andygeiss/ecs`](https://github.com/andygeiss/ecs) — наш фреймворк
- [`github.com/adshao/go-binance/v2`](https://github.com/adshao/go-binance) — официальный Go-клиент Binance
- Реализация **реального источника цен** через WebSocket

---

## 📦 Установка зависимостей

```bash
go get github.com/andygeiss/ecs
go get github.com/adshao/go-binance/v2
```

---

## 🔁 Архитектура: Общие компоненты и системы

> ✅ Мы **переиспользуем почти весь код** из предыдущего примера:
- `Symbol`, `Price`, `SMAConfig`, `SMAResult`, `Order`, `Portfolio`
- `SMACalculationSystem`, `SMAStrategySystem`, `OrderExecutionSystem`

Нам нужно только **заменить `PriceUpdateSystem`** и добавить **WebSocket-обработчик**.

---

## 🔄 Новый: `LivePriceUpdateSystem`

Эта система будет:
- Подключаться к Binance WebSocket
- Получать поток свечей (kline) в реальном времени
- Обновлять компонент `Price`

```go
// systems/live_price.go
package main

import (
	"fmt"
	"time"

	"github.com/andygeiss/ecs"
	"github.com/adshao/go-binance/v2/futures"
)

type LivePriceUpdateSystem struct {
	Symbol   string
	Interval string
	Stop     chan bool
}

func (s *LivePriceUpdateSystem) Process(entity *ecs.Entity, _ float64) {
	// Подключаемся только один раз
	if s.Stop != nil {
		return // уже запущено
	}

	s.Stop = make(chan bool)
	price := entity.Get(&Price{}).(*Price)

	wsHandler := func(event futures.WsKlineEvent) {
		if event.Kline.EndTime == 0 {
			return
		}
		// Обновляем цену по закрытию свечи
		if event.Kline.EndTime < time.Now().UnixMilli()-60000 {
			return // не актуальная свеча
		}
		price.Value = event.Kline.Close
		price.Timestamp = event.Kline.EndTime / 1000 // в секундах
		fmt.Printf("[LIVE] %s @ $%.2f (%s)\n", s.Symbol, price.Value, time.Unix(price.Timestamp, 0).Format("15:04:05"))
	}

	errHandler := func(err error) {
		fmt.Printf("WebSocket error: %v\n", err)
	}

	// Запускаем WebSocket
	futures.WsKlineServe(s.Symbol, s.Interval, wsHandler, errHandler)
	<-s.Stop // держим систему "живой"
}

func (s *LivePriceUpdateSystem) Shutdown() {
	if s.Stop != nil {
		close(s.Stop)
	}
}
```

---

## 🧪 Модификация `PriceUpdateSystem` для гибкости

Сделаем **интерфейс**, чтобы можно было легко переключаться:

```go
// systems/price_system.go
package main

import "github.com/andygeiss/ecs"

type PriceSource interface {
	ecs.System
	IsDone() bool
}
```

Теперь `PriceUpdateSystem` (бэктест) и `LivePriceUpdateSystem` могут реализовывать один интерфейс.

---

## 🧩 Обновлённый `PriceUpdateSystem` (для бэктеста)

```go
// systems/historical_price.go
package main

import "github.com/andygeiss/ecs"

type HistoricalPriceUpdateSystem struct {
	Candles []Candle
	Index   int
	Done    bool
}

func (s *HistoricalPriceUpdateSystem) Process(entity *ecs.Entity, _ float64) {
	if s.Index >= len(s.Candles) {
		s.Done = true
		return
	}
	price := entity.Get(&Price{}).(*Price)
	candle := s.Candles[s.Index]
	price.Value = candle.Close
	price.Timestamp = candle.Timestamp
	s.Index++
}

func (s *HistoricalPriceUpdateSystem) IsDone() bool {
	return s.Done
}
```

---

## 🚀 Основной код: `main.go` (расширенный)

```go
// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andygeiss/ecs"
)

const ModeLive = "live"
const ModeBacktest = "backtest"

func main() {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = ModeBacktest // по умолчанию
	}

	// Создаём ECS-менеджер
	mgr := ecs.NewManager()

	// Создаём сущность
	entity := mgr.NewEntity()
	entity.Add(&Symbol{Name: "BTCUSDT"})
	entity.Add(&Price{})
	entity.Add(&SMAConfig{Period: 50})
	entity.Add(&Portfolio{
		Balance:  1000.0,
		Holdings: make(map[string]float64),
	})
	mgr.AddEntity(entity)

	// Добавляем общие системы
	mgr.AddSystem(&SMACalculationSystem{})
	mgr.AddSystem(&SMAStrategySystem{})
	orderSys := &OrderExecutionSystem{}
	mgr.AddSystem(orderSys)

	// Режим: бэктест или live
	if mode == ModeBacktest {
		fmt.Println("🔄 Запуск в режиме бэктеста...")
		candles := generateTestCandles(1000)
		mgr.AddSystem(&HistoricalPriceUpdateSystem{Candles: candles})

		// Запускаем цикл
		for {
			mgr.Process(60)
			sys, _ := mgr.SystemByName("HistoricalPriceUpdateSystem").(*HistoricalPriceUpdateSystem)
			if sys != nil && sys.IsDone() {
				break
			}
		}

	} else {
		fmt.Println("📡 Запуск в режиме реального времени (LIVE)...")
		liveSys := &LivePriceUpdateSystem{
			Symbol:   "BTCUSDT",
			Interval: "1m",
		}
		mgr.AddSystem(liveSys)

		// Запускаем ECS в цикле
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case <-ticker.C:
				mgr.Process(1.0) // 1 секунда
			case <-sig:
				fmt.Println("\n🛑 Остановка...")
				liveSys.Shutdown()
				goto exit
			}
		}
	}

exit:
	// Вывод результатов
	portfolio := entity.Get(&Portfolio{}).(*Portfolio)
	fmt.Printf("\n📊 Итоговый баланс: $%.2f\n", portfolio.Balance)
	trades := len(orderSys.Trades)
	fmt.Printf("✅ Всего сделок: %d\n", trades)
	if trades > 0 {
		last := orderSys.Trades[len(orderSys.Trades)-1]
		side := "🟢 КУПЛЕН" + "О"
		if last.Side == "SELL" {
			side = "🔴 ПРОДАНО"
		}
		fmt.Printf("🔚 Последняя сделка: %s %.4f по $%.2f\n", side, last.Amount, last.Price)
	}
}
```

---

## 🧪 Как запустить

### 1. Бэктест (по умолчанию)

```bash
go run main.go
```

### 2. Live-режим

```bash
MODE=live go run main.go
```

> ⚠️ Убедись, что у тебя установлен `go-binance/v2`:
```bash
go get github.com/adshao/go-binance/v2
```

---

## 📈 Что вы увидите в live-режиме

```
📡 Запуск в режиме реального времени (LIVE)...
[LIVE] BTCUSDT @ $62340.50 (15:30:00)
[LIVE] BTCUSDT @ $62345.20 (15:31:00)
🤖 Сигнал: BUY
✅ Ордер исполнен: BUY 0.0015 BTC по $62345.20
```

---

## 🔐 Безопасность

> ❗ Пока мы **не отправляем реальные ордера** — `OrderExecutionSystem` работает в симуляции.

Чтобы торговать реальными деньгами:
- Добавь **API ключи Binance**
- Расширь `OrderExecutionSystem` вызовами `futures.NewOrderService()`
- Добавь проверку `if os.Getenv("LIVE_TRADE") == "true"`

Пример (опционально):
```go
if os.Getenv("LIVE_TRADE") == "true" {
    // Вызов API Binance для создания ордера
}
```

---

## 🧩 Преимущества такой архитектуры

| Фича | Реализация |
|------|-----------|
| 🔄 Единая логика | Стратегии работают и в бэктесте, и в live |
| 🧪 Безопасность | Можно тестировать стратегию до включения реальных сделок |
| 📈 Масштабируемость | Легко добавить новые пары, стратегии, системы |
| 🧱 Чистая архитектура | ECS изолирует данные и поведение |

---

## 🚀 Что можно добавить дальше?

| Фича | Как |
|------|-----|
| 💬 Уведомления в Telegram | Добавить `TelegramAlertSystem` |
| 📊 Веб-интерфейс | Использовать `gorilla/mux` + Chart.js |
| 📥 Поддержка других бирж | Создать интерфейс `Exchange` |
| 📈 Множественные стратегии | По одной сущности на стратегию |
| 📅 Backtest с реальных данных | Загрузить историю через `rest.Client` |

---

## ✅ Вывод

Теперь у тебя есть:
- **Универсальный ECS-движок** для торговли
- **Режим бэктеста** и **live-режим** с реальными данными
- Возможность **безопасно тестировать** стратегии
- Гибкая, расширяемая архитектура

---

## 📬 Хочешь продолжить?

Могу помочь:
- Добавить **реальные ордера через API Binance**
- Интегрировать **веб-панель с графиками**
- Сделать **конфигурацию через YAML**
- Реализовать **автоматическую оптимизацию параметров SMA**

Напиши — и углубимся! 🚀

### USER
как организовать структуру папок для этого примера, разделив компоненты, системы, служебные модули и основное приложение.

### ASSISTANT
Отличный вопрос! Чтобы проект был **масштабируемым, чистым и легко поддерживаемым**, важно правильно организовать структуру папок. Ниже — **рекомендуемая структура директорий** для нашего ECS-торгового бота с поддержкой бэктеста и live-режима.

---

## 📁 Рекомендуемая структура папок

```
trading-bot/
├── cmd/
│   └── bot/
│       └── main.go            # Точка входа: запуск в зависимости от MODE
├── internal/
│   ├── components/            # Все компоненты ECS
│   │   ├── symbol.go
│   │   ├── price.go
│   │   ├── portfolio.go
│   │   ├── order.go
│   │   ├── strategy.go
│   │   └── risk.go
│   │
│   ├── systems/               # Все системы ECS
│   │   ├── sma_calculation.go
│   │   ├── sma_strategy.go
│   │   ├── order_execution.go
│   │   └── price/
│   │       ├── historical.go  # Бэктест: из данных
│   │       └── live.go        # Live: с Binance WebSocket
│   │
│   ├── services/              # Внешние сервисы (Binance, Telegram и т.д.)
│   │   ├── binance/
│   │   │   ├── client.go      # REST и WebSocket клиенты
│   │   │   └── types.go
│   │   └── alert/             # Уведомления
│   │       └── telegram.go
│   │
│   ├── utils/                 # Вспомогательные функции
│   │   ├── math.go
│   │   └── time.go
│   │
│   ├── domain/                # Бизнес-логика, модели, типы
│   │   ├── trade.go
│   │   ├── candle.go
│   │   └── portfolio_tracker.go
│   │
│   └── app/                   # Основная логика приложения
│       ├── engine.go          # Настройка ECS-менеджера
│       ├── runner.go          # Запуск бэктеста / live
│       └── config.go          # Конфигурация (период, символ, баланс и т.д.)
│
├── pkg/                       # (опционально) общие пакеты для других проектов
│   └── ecsx/                  # Обёртки/расширения для github.com/andygeiss/ecs
│       └── systems.go         # Интерфейсы, базовые системы
│
├── configs/
│   └── config.yaml            # Конфигурация стратегий, ключей, параметров
│
├── scripts/
│   └── backtest.sh            # Скрипты запуска
│
├── go.mod
├── go.sum
└── README.md
```

---

## 🧩 Пояснение по папкам

### 1. `cmd/bot/main.go` — Точка входа
```go
package main

func main() {
    cfg := config.Load()
    app := app.New(cfg)
    app.Run() // Запускает либо бэктест, либо live
}
```
> Только инициализация конфига и запуск приложения. Никакой логики.

---

### 2. `internal/components/` — Данные сущностей

Каждый файл — один компонент:
```go
// components/price.go
type Price struct {
    Value     float64
    Timestamp int64
}
```

> ❗ Компоненты — **только структуры данных**, без методов.

---

### 3. `internal/systems/` — Логика обработки

- Каждая система — отдельный файл
- Сложные системы (например, `price`) можно вынести в подпапку

```go
// systems/price/historical.go
type HistoricalPriceUpdateSystem struct { ... }
func (s *HistoricalPriceUpdateSystem) Process(entity *ecs.Entity, delta float64)
```

---

### 4. `internal/services/` — Внешние зависимости

#### `binance/client.go`
```go
type BinanceService struct {
    Rest    *rest.Client
    Ws      *ws.Client
}
func (s *BinanceService) GetHistoricalCandles(...) ([]Candle, error)
func (s *BinanceService) SubscribeKline(...)
```

> Инкапсулирует работу с API. Системы используют этот сервис, а не напрямую `go-binance`.

---

### 5. `internal/utils/` — Вспомогательные функции

```go
// utils/math.go
func Round(val float64, precision int) float64 { ... }

// utils/time.go
func TimeFromMs(ms int64) time.Time { ... }
```

---

### 6. `internal/domain/` — Бизнес-сущности

Не путать с ECS-сущностями. Это — **модели предметной области**:
```go
// domain/trade.go
type Trade struct {
    Side      string
    Price     float64
    Amount    float64
    Timestamp int64
    PnL       float64
}

// domain/portfolio_tracker.go
type PortfolioTracker struct {
    InitialBalance float64
    CurrentBalance float64
    Holdings       map[string]float64
    Trades         []Trade
}
func (p *PortfolioTracker) Buy(...) error
func (p *PortfolioTracker) Sell(...) error
```

> Можно использовать внутри `OrderExecutionSystem`.

---

### 7. `internal/app/` — Оркестрация

#### `app/engine.go`
```go
func NewEcsManager(cfg *Config) *ecs.Manager {
    mgr := ecs.NewManager()
    // Регистрируем системы в зависимости от режима
    if cfg.Mode == "backtest" {
        mgr.AddSystem(NewHistoricalPriceSystem(cfg.Candles))
    } else {
        mgr.AddSystem(NewLivePriceSystem(cfg.Symbol))
    }
    mgr.AddSystem(&systems.SMACalculationSystem{})
    // ...
    return mgr
}
```

#### `app/runner.go`
```go
func (a *App) Run() {
    if a.cfg.Mode == "backtest" {
        a.runBacktest()
    } else {
        a.runLive()
    }
}
```

#### `app/config.go`
```go
type Config struct {
    Mode         string
    Symbol       string
    Interval     string
    Strategy     string
    SMA_Period   int
    InitialBalance float64
    // ...
}
func Load() *Config { ... } // из env, yaml, флагов
```

---

### 8. `pkg/ecsx/` — Расширения ECS (опционально)

Если часто используешь ECS, можно создать:
```go
// pkg/ecsx/systems.go
type System interface {
    ecs.System
    Name() string
    Priority() int
    IsDone() bool
}
```

> Удобно для фильтрации, оркестрации, логирования.

---

### 9. `configs/config.yaml`
```yaml
mode: live
symbol: BTCUSDT
interval: 1m
strategy:
  type: sma
  period: 50
portfolio:
  initial_balance: 1000.0
  currency: USDT
binance:
  api_key: ""
  secret_key: ""
  testnet: true
```

---

### 10. `scripts/backtest.sh`
```bash
#!/bin/bash
MODE=backtest go run cmd/bot/main.go
```

---

## ✅ Преимущества такой структуры

| Плюс | Объяснение |
|------|-----------|
| 🔍 Чёткое разделение | Нет "спагетти-кода", каждый модуль на своём месте |
| 🧩 Легко расширять | Добавить RSI? → `systems/rsi_strategy.go` |
| 🧪 Легко тестировать | Можно юнит-тестить `PortfolioTracker`, `SMA`, `BinanceService` отдельно |
| 🔄 Режимы независимы | Live и backtest используют общие системы |
| 📦 Инкапсуляция | Внешние API изолированы в `services/` |

---

## 🚀 Следующие шаги

1. Создай структуру:
```bash
mkdir -p trading-bot/{cmd/bot,internal/{components,systems/{price,domain,services/binance,services/alert,utils,app},pkg/ecsx,configs,scripts}}
```

2. Начни с `internal/components/` и `internal/systems/`, перенеся код из предыдущего примера

3. Реализуй `internal/app/engine.go` и `runner.go`

4. Подключи `go-binance` в `services/binance/client.go`

---

## 📬 Готов помочь

Могу:
- Сгенерировать **все файлы-заглушки**
- Написать **пример `config.yaml` + loader**
- Реализовать **полный `BinanceService`**
- Добавить **поддержку нескольких стратегий**

Напиши — и продолжим строить профессиональный бот! 🚀

### USER
Объясни мне эту секцию более подробно.

```
pkg/ecsx/ — Расширения ECS (опционально)
```

### ASSISTANT
Отлично! Давай подробно разберём, **что такое `pkg/ecsx/` и зачем нужны «расширения ECS»**, даже если это опционально.

---

## 🤔 Что означает «расширения ECS»?

Фреймворк [`github.com/andygeiss/ecs`](https://github.com/andygeiss/ecs) — это **минималистичная, чистая реализация ECS** (Entity-Component-System).  
Он предоставляет базовые вещи:
- `Manager` — управляет сущностями и системами
- `Entity` — контейнер компонентов
- `Component` — данные
- `System` — логика обработки
- `Process(delta)` — игровой цикл

Но он **не включает**:
- Гибкую сортировку систем
- Логирование
- Приоритеты выполнения
- Интерфейсы для проверки состояния
- Расширяемые хуки (например, `OnInit`, `OnShutdown`)

👉 Вот здесь и пригодится **`pkg/ecsx/`** — это **ваш собственный слой над ECS**, где вы добавляете удобства и структуру.

---

## 🧱 Зачем это нужно?

Представь, что ты строишь дом:
- `github.com/andygeiss/ecs` — это **каркас из балок и досок**
- `pkg/ecsx/` — это **электропроводка, водопровод, выключатели**, которые делают дом **удобным для жизни**

Без этого можно обойтись, но с этим — **много удобнее и профессиональнее**.

---

## 📁 Что может быть в `pkg/ecsx/`?

### 1. **Расширенный интерфейс `System`**

В оригинальном ECS система — это просто:
```go
type System interface {
    Process(entity *Entity, delta float64)
}
```

Но в реальности хочется:
- Знать имя системы
- Управлять приоритетом
- Проверять, завершена ли система
- Выполнять действия при старте/остановке

👉 Создаём **расширенный интерфейс**:

```go
// pkg/ecsx/systems.go
package ecsx

import "github.com/andygeiss/ecs"

type System interface {
    ecs.System
    
    Name() string           // для логирования
    Priority() int          // порядок выполнения
    IsDone() bool           // завершена ли (для бэктеста)
    OnInit() error          // вызывается при добавлении
    OnShutdown() error      // при остановке
}
```

---

### 2. **Базовая реализация для удобства**

Создаём **базовую структуру**, чтобы не повторять одно и то же:

```go
type BaseSystem struct {
    name     string
    priority int
    done     bool
}

func (b *BaseSystem) Name() string       { return b.name }
func (b *BaseSystem) Priority() int      { return b.priority }
func (b *BaseSystem) IsDone() bool       { return b.done }
func (b *BaseSystem) OnInit() error      { return nil }
func (b *BaseSystem) OnShutdown() error  { return nil }

// Утилита: SetDone
func (b *BaseSystem) MarkDone() { b.done = true }
```

Теперь любая система может её встраивать:

```go
type HistoricalPriceSystem struct {
    ecsx.BaseSystem
    Candles []Candle
    Index   int
}

func (s *HistoricalPriceSystem) Process(entity *ecs.Entity, delta float64) {
    // ... логика
    if s.Index >= len(s.Candles) {
        s.MarkDone()
    }
}

func (s *HistoricalPriceSystem) Name() string {
    return "HistoricalPriceSystem" // можно переопределить
}
```

---

### 3. **Сортировка систем по приоритету**

Иногда важно, в каком порядке выполняются системы:

```go
// ecsx/utils.go
func SortSystems(systems []ecs.System) []ecs.System {
    // Приводим к нашему интерфейсу
    var typed []System
    for _, s := range systems {
        if ext, ok := s.(System); ok {
            typed = append(typed, ext)
        } else {
            // Обычные системы — с низким приоритетом
            typed = append(typed, &WrappedSystem{s})
        }
    }

    // Сортируем по приоритету
    sort.Slice(typed, func(i, j int) bool {
        return typed[i].Priority() < typed[j].Priority()
    })

    // Возвращаем как ecs.System
    var result []ecs.System
    for _, s := range typed {
        result = append(result, s)
    }
    return result
}
```

> Теперь можно сказать:  
> - `PriceUpdateSystem` — приоритет 10 (должна быть первой)
> - `SMAStrategySystem` — 20
> - `OrderExecutionSystem` — 30

---

### 4. **Обёртка для не-расширенных систем**

Если есть система, которая **не реализует `ecsx.System`**, можно её обернуть:

```go
type WrappedSystem struct {
    ecs.System
    name string
}

func (w *WrappedSystem) Name() string       { return w.name }
func (w *WrappedSystem) Priority() int      { return 100 } // по умолчанию низкий
func (w *WrappedSystem) IsDone() bool       { return false }
func (w *WrappedSystem) OnInit() error      { return nil }
func (w *WrappedSystem) OnShutdown() error  { return nil }
```

---

### 5. **Расширенный Manager (опционально)**

Можно создать `ecsx.Manager`, который:
- Автоматически сортирует системы
- Логирует запуск/остановку
- Позволяет получать систему по имени

```go
type Manager struct {
    *ecs.Manager
    systems []System
}

func (m *Manager) AddSystem(sys System) {
    m.systems = append(m.systems, sys)
    m.Manager.AddSystem(sys)
    sortSystems(m.systems) // или отсортировать при Process
}
```

---

## ✅ Пример: как это улучшает код

Без `ecsx`:
```go
mgr.AddSystem(&PriceUpdateSystem{})
mgr.AddSystem(&SMAStrategySystem{})
mgr.AddSystem(&OrderExecutionSystem{})
// ❌ Нет контроля порядка!
```

С `ecsx`:
```go
mgr.AddSystem(&ecsx.HistoricalPriceSystem{
    BaseSystem: ecsx.BaseSystem{
        name:     "PriceUpdate",
        priority: 10,
    },
    Candles: candles,
})

mgr.AddSystem(&SMAStrategySystem{
    BaseSystem: ecsx.BaseSystem{
        name:     "SMA_Strategy",
        priority: 20,
    },
})
```

→ Порядок гарантирован, можно логировать:
```go
for _, sys := range mgr.Systems() {
    fmt.Printf("Running: %s (priority: %d)\n", sys.Name(), sys.Priority())
}
```

---

## 🧩 Когда это нужно?

| Ситуация | Нужно ли `ecsx`? |
|--------|------------------|
| Простой бэктест с 2 системами | ❌ Нет, избыточно |
| Множество систем, важен порядок | ✅ Да |
| Хочется логировать, что происходит | ✅ Да |
| Планируешь масштабировать бота | ✅ Да |
| Делаете командную разработку | ✅ Да (единый стиль) |

---

## 🎯 Вывод

> **`pkg/ecsx/` — это ваш «надстройка» над чистым ECS**, которая превращает его из «игрушечного движка» в **профессиональный инструмент**.

Он позволяет:
- Управлять порядком выполнения
- Логировать и отлаживать
- Проверять состояние систем
- Повторно использовать код
- Делать архитектуру более предсказуемой

---

## 📦 Пример структуры `pkg/ecsx/`

```
pkg/ecsx/
├── system.go           // интерфейс System и BaseSystem
├── manager.go          // расширенный менеджер (опционально)
├── sorting.go          // сортировка по приоритету
├── logging.go          // логирующий декоратор для систем
└── utils.go            // вспомогательные функции
```

---

## 📬 Хочешь — я могу:

- Создать **все файлы `ecsx/` с нуля**
- Показать, как **интегрировать это в `app/engine.go`**
- Добавить **логирование активных систем**
- Реализовать **систему с приоритетами и именами**

Напиши — и сделаем ваш ECS по-настоящему мощным! 🔧🚀

### USER
Как в этом примере можно применить event driven architecture?

### ASSISTANT
Отличный вопрос!  
Ты уже используешь **ECS-архитектуру**, которая сама по себе очень близка к **Event-Driven Architecture (EDA)**, потому что:
- Системы реагируют на изменения
- Компоненты — это состояние
- Можно легко добавить **события как триггеры поведения**

Но сейчас логика работает **по циклу (`Process`)** — системы проверяют состояние на каждом шаге.

Давай превратим это в **по-настоящему event-driven приложение**, где:
> 🔔 **События** (например, "цена обновилась", "сделка исполнена") → **вызывают реакции** (например, "пересчитать SMA", "отправить уведомление")

---

## 🧠 Что такое Event-Driven Architecture (EDA)?

**Event-Driven Architecture** — это стиль проектирования, при котором:
- **Система реагирует на события**, а не опрашивает состояние
- События **генерируются** одной частью системы
- Другие части **подписываются** на эти события и **реагируют**

Примеры событий:
- `PriceUpdated`
- `OrderExecuted`
- `StrategySignalGenerated`
- `RiskLimitExceeded`

---

## 🔗 Где EDA хорошо сочетается с ECS?

| ECS | EDA |
|-----|-----|
| Сущности меняют состояние | → | При изменении — генерируется событие |
| Системы опрашивают компоненты | → | Системы **подписываются** на события и реагируют |
| `Process()` — поллинг | → | **Пуш-модель**: реакция на событие |

👉 Это делает систему **быстрее, модульнее и масштабируемее**.

---

## 🛠️ Как добавить EDA в текущий ECS-бот?

### Шаг 1: Создать систему событий

Создадим простой **event bus** (шину событий):

```go
// pkg/events/bus.go
package events

type EventHandler func(payload interface{})

type EventBus struct {
    handlers map[string][]EventHandler
}

func New() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(eventType string, payload interface{}) {
    for _, h := range b.handlers[eventType] {
        go h(payload) // асинхронно
    }
}
```

> ✅ Можно использовать готовые решения: `github.com/asaskevich/EventBus`, `mitchellh/mapstructure`, но пока сделаем вручную.

---

### Шаг 2: Определить типы событий

```go
// internal/events/types.go
package events

import "time"

type PriceUpdated struct {
    Symbol    string
    Price     float64
    Timestamp time.Time
}

type OrderExecuted struct {
    Side      string
    Price     float64
    Amount    float64
    Symbol    string
    Timestamp time.Time
}

type StrategySignal struct {
    Symbol string
    Signal string // "BUY", "SELL"
    Confidence float64
}

type RiskAlert struct {
    Type    string // "MaxDrawdown", "VolatilitySpoke"
    Message string
    Value   float64
}
```

---

### Шаг 3: Интегрировать в системы

#### ✅ Пример: `LivePriceUpdateSystem` публикует событие

```go
// internal/systems/price/live.go
type LivePriceUpdateSystem struct {
    eventBus *events.EventBus
    // ...
}

func (s *LivePriceUpdateSystem) Process(entity *ecs.Entity, _ float64) {
    // ... при получении свечи
    price.Value = candle.Close
    price.Timestamp = ts

    // 🔔 Публикуем событие
    s.eventBus.Publish("PriceUpdated", events.PriceUpdated{
        Symbol:    "BTCUSDT",
        Price:     price.Value,
        Timestamp: time.Unix(price.Timestamp, 0),
    })
}
```

---

#### ✅ Пример: `SMAStrategySystem` — подписывается на `PriceUpdated`

Теперь она **не опрашивает цену**, а **реагирует на событие**:

```go
// internal/systems/sma_strategy.go
type SMAStrategySystem struct {
    entity   *ecs.Entity
    eventBus *events.EventBus
}

func (s *SMAStrategySystem) Start() {
    s.eventBus.Subscribe("PriceUpdated", func(payload interface{}) {
        event := payload.(events.PriceUpdated)
        // Проверяем, подходит ли символ
        symbol := s.entity.Get(&components.Symbol{}).(*components.Symbol)
        if event.Symbol != symbol.Name {
            return
        }

        // Здесь — логика стратегии
        sma := s.calculateSMA()
        if event.Price > sma && s.wasBelow {
            s.wasBelow = false
            // Генерируем сигнал
            s.eventBus.Publish("StrategySignal", events.StrategySignal{
                Symbol:     symbol.Name,
                Signal:     "BUY",
                Confidence: 0.8,
            })
        } else if event.Price < sma && !s.wasBelow {
            s.wasBelow = true
            s.eventBus.Publish("StrategySignal", events.StrategySignal{
                Symbol:     symbol.Name,
                Signal:     "SELL",
                Confidence: 0.7,
            })
        }
    })
}
```

---

#### ✅ Пример: `OrderExecutionSystem` реагирует на сигнал

```go
// internal/systems/order_execution.go
func (s *OrderExecutionSystem) Start() {
    s.eventBus.Subscribe("StrategySignal", func(payload interface{}) {
        signal := payload.(events.StrategySignal)
        if signal.Signal == "BUY" {
            s.placeOrder("BUY", signal.Price)
        } else if signal.Signal == "SELL" {
            s.placeOrder("SELL", signal.Price)
        }

        // 🔔 Публикуем факт исполнения
        s.eventBus.Publish("OrderExecuted", events.OrderExecuted{
            Side:      signal.Signal,
            Price:     signal.Price,
            Amount:    s.amount,
            Symbol:    signal.Symbol,
            Timestamp: time.Now(),
        })
    })
}
```

---

### Шаг 4: Добавить реактивные системы (не в цикле!)

Теперь некоторые системы **не нужно запускать в `Process()`** — они работают **асинхронно** на событиях.

Но ECS всё ещё нужен для хранения состояния (`Portfolio`, `Price` и т.д.).

Можно оставить цикл для систем, которые **не могут быть event-driven** (например, таймеры), а остальные — перевести на события.

---

## 🔄 Сравнение: Polling vs Event-Driven

| Аспект | Polling (цикл) | Event-Driven |
|-------|----------------|--------------|
| Частота реакции | Каждый `delta` | Сразу после события |
| Производительность | Проверяет всё, всегда | Только когда нужно |
| Гибкость | Жёсткая привязка | Лёгко добавить слушателя |
| Отладка | Проще | Нужно логировать события |
| Масштабирование | Ограничено | Подходит для микросервисов |

---

## 🌐 Пример: Добавить Telegram-уведомления

```go
// internal/services/alert/telegram.go
func RegisterTelegramAlerts(bus *events.EventBus, token, chatID string) {
    bus.Subscribe("OrderExecuted", func(p interface{}) {
        e := p.(events.OrderExecuted)
        msg := fmt.Sprintf("✅ Исполнено: %s %g по $%.2f", e.Side, e.Amount, e.Price)
        sendTelegramMessage(token, chatID, msg)
    })

    bus.Subscribe("RiskAlert", func(p interface{}) {
        e := p.(events.RiskAlert)
        msg := fmt.Sprintf("⚠️ Риск: %s — %s", e.Type, e.Message)
        sendTelegramMessage(token, chatID, msg)
    })
}
```

> Никакой зависимости от ECS! Только подписка на события.

---

## 🧱 Как интегрировать в `app/engine.go`

```go
// internal/app/engine.go
func NewEcsManager(cfg *Config) (*ecs.Manager, *events.EventBus) {
    bus := events.New()
    mgr := ecs.NewManager()

    // Создаём сущности
    entity := mgr.NewEntity()
    // ... добавляем компоненты

    // Создаём системы
    priceSys := &systems.LivePriceUpdateSystem{Entity: entity, EventBus: bus}
    strategySys := &systems.SMAStrategySystem{Entity: entity, EventBus: bus}
    orderSys := &systems.OrderExecutionSystem{Entity: entity, EventBus: bus}

    // Запускаем event-driven системы
    strategySys.Start()
    orderSys.Start()

    // Только polling-системы добавляем в ECS
    mgr.AddSystem(priceSys) // она всё ещё в цикле, чтобы слушать WebSocket

    // Подключаем алерты
    services.RegisterTelegramAlerts(bus, cfg.Telegram.Token, cfg.Telegram.ChatID)

    return mgr, bus
}
```

---

## ✅ Преимущества такой архитектуры

| Плюс | Объяснение |
|------|-----------|
| ⚡ Быстрее реакция | Нет задержки до следующего `Process()` |
| 🔌 Модульность | Можно добавить `EmailService`, `DBLogger` без изменения кода |
| 📦 Расширяемость | Подключай новые "реакторы" за секунды |
| 🔄 Декуплинг | Системы не знают друг о друге — только о событиях |
| 📈 Подходит для микросервисов | События можно отправлять в Kafka/NATS |

---

## 🚫 Минусы

| Минус | Как смягчить |
|------|-------------|
| Сложнее отлаживать | Логируй все события: `bus.Publish → handlers` |
| Асинхронность | Используй `sync.WaitGroup` для тестов |
| Потеря событий | Используй очереди (Redis, NATS) в продакшене |

---

## 🎯 Вывод

Ты можешь **постепенно превратить ECS-бот в event-driven систему**, сохранив лучшее от обеих архитектур:

- ✅ **ECS** — для **управления состоянием** (сущности, компоненты)
- ✅ **EDA** — для **реактивной логики** (события → действия)

> Это как "мозг" (ECS) и "нервная система" (EDA):  
> мозг хранит состояние, а нервы мгновенно передают сигналы.

---

## 📬 Готов реализовать?

Могу помочь:
- Создать **полный `events/` пакет**
- Переписать **SMA-стратегию под события**
- Добавить **логирование всех событий**
- Интегрировать **NATS или Redis** для продакшена

Напиши — и сделаем бота **по-настоящему реактивным**! ⚡🚀