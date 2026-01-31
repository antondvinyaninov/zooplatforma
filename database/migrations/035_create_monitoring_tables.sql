-- Миграция для системы мониторинга и логирования ошибок

-- Таблица для логирования ошибок
CREATE TABLE IF NOT EXISTS error_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT NOT NULL,           -- Название сервиса (Main Backend, Admin Backend, etc.)
    endpoint TEXT NOT NULL,          -- URL endpoint где произошла ошибка
    method TEXT NOT NULL,            -- HTTP метод (GET, POST, etc.)
    error_message TEXT NOT NULL,     -- Сообщение об ошибке
    stack_trace TEXT,                -- Stack trace (опционально)
    user_id INTEGER,                 -- ID пользователя (если авторизован)
    ip_address TEXT,                 -- IP адрес клиента
    user_agent TEXT,                 -- User Agent браузера
    request_body TEXT,               -- Тело запроса (для POST/PUT)
    response_body TEXT,              -- Тело ответа
    status_code INTEGER,             -- HTTP статус код
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_error_logs_service ON error_logs(service);
CREATE INDEX IF NOT EXISTS idx_error_logs_created_at ON error_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_error_logs_user_id ON error_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_error_logs_status_code ON error_logs(status_code);

-- Таблица для логирования всех запросов (опционально, для статистики)
CREATE TABLE IF NOT EXISTS request_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER,        -- Время ответа в миллисекундах
    user_id INTEGER,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Индексы для request_logs
CREATE INDEX IF NOT EXISTS idx_request_logs_service ON request_logs(service);
CREATE INDEX IF NOT EXISTS idx_request_logs_created_at ON request_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_request_logs_status_code ON request_logs(status_code);

-- Таблица для метрик производительности
CREATE TABLE IF NOT EXISTS performance_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT NOT NULL,
    metric_name TEXT NOT NULL,       -- Название метрики (cpu_usage, memory_usage, etc.)
    metric_value REAL NOT NULL,      -- Значение метрики
    unit TEXT,                       -- Единица измерения (%, MB, ms, etc.)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для performance_metrics
CREATE INDEX IF NOT EXISTS idx_performance_metrics_service ON performance_metrics(service);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_created_at ON performance_metrics(created_at);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_name ON performance_metrics(metric_name);

-- Таблица для алертов
CREATE TABLE IF NOT EXISTS system_alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alert_type TEXT NOT NULL,        -- Тип алерта (error_spike, service_down, slow_response, etc.)
    severity TEXT NOT NULL,          -- Серьезность (low, medium, high, critical)
    service TEXT,                    -- Сервис (если применимо)
    message TEXT NOT NULL,           -- Сообщение алерта
    details TEXT,                    -- Дополнительные детали (JSON)
    is_resolved BOOLEAN DEFAULT 0,   -- Решен ли алерт
    resolved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для system_alerts
CREATE INDEX IF NOT EXISTS idx_system_alerts_type ON system_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_system_alerts_severity ON system_alerts(severity);
CREATE INDEX IF NOT EXISTS idx_system_alerts_resolved ON system_alerts(is_resolved);
CREATE INDEX IF NOT EXISTS idx_system_alerts_created_at ON system_alerts(created_at);

-- Представление для быстрого доступа к статистике ошибок
CREATE VIEW IF NOT EXISTS error_stats_24h AS
SELECT 
    service,
    COUNT(*) as error_count,
    COUNT(DISTINCT user_id) as affected_users,
    MIN(created_at) as first_error,
    MAX(created_at) as last_error
FROM error_logs
WHERE created_at > datetime('now', '-24 hours')
GROUP BY service
ORDER BY error_count DESC;

-- Представление для статистики запросов
CREATE VIEW IF NOT EXISTS request_stats_hourly AS
SELECT 
    service,
    strftime('%Y-%m-%d %H:00:00', created_at) as hour,
    COUNT(*) as total_requests,
    SUM(CASE WHEN status_code >= 400 THEN 1 ELSE 0 END) as error_count,
    AVG(response_time_ms) as avg_response_time,
    MAX(response_time_ms) as max_response_time
FROM request_logs
WHERE created_at > datetime('now', '-24 hours')
GROUP BY service, hour
ORDER BY hour DESC, service;

SELECT 'Monitoring tables created successfully' as status;
