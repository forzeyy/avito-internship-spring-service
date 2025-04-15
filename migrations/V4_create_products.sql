CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES receipts(id),
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);