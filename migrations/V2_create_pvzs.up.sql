-- +migrate Up
CREATE TABLE IF NOT EXISTS pvzs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city VARCHAR(50) NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
    reg_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);