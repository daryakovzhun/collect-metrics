CREATE TYPE metric_type AS ENUM ('gauge', 'counter');
COMMENT ON TYPE metric_type IS 'Типы метрик.';

CREATE TABLE metrics (
     id BIGSERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     type metric_type NOT NULL,
     delta BIGINT,
     value DOUBLE PRECISION,
     hash VARCHAR(64),
     created_at timestamp with time zone DEFAULT NOW(),
     updated_at timestamp with time zone DEFAULT NOW()
);

COMMENT ON TABLE metrics IS 'Таблица с метриками.';
COMMENT ON COLUMN metrics.id IS 'Уникальный идентификатор метрики.';
COMMENT ON COLUMN metrics.name IS 'Уникальное название метрики.';
COMMENT ON COLUMN metrics.type IS 'Тип метрики (gauge, counter).';
COMMENT ON COLUMN metrics.delta IS 'Значение метрики в случае передачи counter.';
COMMENT ON COLUMN metrics.value IS 'Значение метрики в случае передачи gauge.';
COMMENT ON COLUMN metrics.hash IS 'Хэш метрики.';
COMMENT ON COLUMN metrics.created_at IS 'Время создания метрики.';
COMMENT ON COLUMN metrics.updated_at IS 'Время обновления метрики.';

CREATE UNIQUE INDEX unique_metric_name ON metrics(name);