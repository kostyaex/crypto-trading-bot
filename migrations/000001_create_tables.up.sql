-- 000001_create_tables.up.sql

-- Таблица для хранения рыночных данных
CREATE TABLE IF NOT EXISTS market_data (
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    exchange TEXT NOT NULL,
    symbol TEXT NOT NULL,
    time_frame TEXT NOT NULL,
    open_price NUMERIC(20, 8) NOT NULL,
    hight_price NUMERIC(20, 8) NOT NULL,
    low_price NUMERIC(20, 8) NOT NULL,
    close_price NUMERIC(20, 8) NOT NULL,
    volume NUMERIC(20,3) NOT NULL,
    buy_volume NUMERIC(20,3) NOT NULL,
    sell_volume NUMERIC(20,3) NOT NULL
);

CREATE TABLE IF NOT EXISTS market_data_statuss (
    id SERIAL PRIMARY KEY,
    exchange TEXT NOT NULL,
    symbol TEXT NOT NULL,
    time_frame TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    actual_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL,
    UNIQUE(exchange,symbol, time_frame)
);

INSERT INTO market_data_statuss (exchange, symbol, time_frame, active, status)
    VALUES ('binance', 'BTCUSDT', '5m', false, '') ON CONFLICT DO NOTHING;
    

CREATE TABLE IF NOT EXISTS cluster_data (
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    symbol TEXT NOT NULL,
    time_frame TEXT NOT NULL,
    is_buysell BOOLEAN NOT NULL DEFAULT FALSE,
    cluster_price NUMERIC(20, 8) NOT NULL,
    volume NUMERIC(20,3) NOT NULL
);

-- Таблица для хранения информации о биржах
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    exchange TEXT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    api_key TEXT NOT NULL,
    api_secret TEXT NOT NULL
);

-- Таблица для хранения технических индикаторов
CREATE TABLE IF NOT EXISTS indicators (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    indicator_name TEXT NOT NULL,
    value NUMERIC(20, 8) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, indicator_name, timestamp)
);

-- Таблица для хранения торговых стратегий
CREATE TABLE IF NOT EXISTS strategies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    config JSONB NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Таблица для хранения состояний поведенческих деревьев
CREATE TABLE IF NOT EXISTS behavior_trees (
    id SERIAL PRIMARY KEY,
    strategy_id INT NOT NULL,
    state JSONB NOT NULL,
    last_executed TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (strategy_id) REFERENCES strategies(id) ON DELETE CASCADE
);

-- Таблица для хранения торговых операций
CREATE TABLE IF NOT EXISTS trade_operations (
    id SERIAL PRIMARY KEY,
    strategy_id INT NOT NULL,
    symbol TEXT NOT NULL,
    order_type TEXT NOT NULL,
    side TEXT NOT NULL,
    quantity NUMERIC(20, 8) NOT NULL,
    price NUMERIC(20, 8),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL,
    FOREIGN KEY (strategy_id) REFERENCES strategies(id) ON DELETE CASCADE
);

-- Таблица для хранения результатов бэктестинга
CREATE TABLE IF NOT EXISTS backtest_results (
    id SERIAL PRIMARY KEY,
    strategy_id INT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    initial_capital NUMERIC(20, 8) NOT NULL,
    final_capital NUMERIC(20, 8) NOT NULL,
    total_profit NUMERIC(20, 8) NOT NULL,
    drawdown NUMERIC(20, 8) NOT NULL,
    successful_trades INT NOT NULL,
    failed_trades INT NOT NULL,
    config JSONB NOT NULL,
    FOREIGN KEY (strategy_id) REFERENCES strategies(id) ON DELETE CASCADE
);