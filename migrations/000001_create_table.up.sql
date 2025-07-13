-- Создаем таблицу 'orders'
CREATE TABLE IF NOT EXISTS orders (
    order_uid           VARCHAR(255) PRIMARY KEY,
    track_number        VARCHAR(255) NOT NULL,
    entry               VARCHAR(255) NOT NULL,
    locale              VARCHAR(10) NOT NULL,
    internal_signature  VARCHAR(255),
    customer_id         VARCHAR(255) NOT NULL,
    delivery_service    VARCHAR(255) NOT NULL,
    shardkey            VARCHAR(255) NOT NULL,
    sm_id               BIGINT NOT NULL,
    date_created        TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard           VARCHAR(255) NOT NULL
);

-- Создаем таблицу 'delivery_info'
-- Связываем ее с orders по order_uid
CREATE TABLE IF NOT EXISTS delivery_info (
    order_uid   VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(255) NOT NULL,
    zip         VARCHAR(20) NOT NULL,
    city        VARCHAR(255) NOT NULL,
    address     VARCHAR(255) NOT NULL,
    region      VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    CONSTRAINT fk_delivery_order
        FOREIGN KEY (order_uid)
        REFERENCES orders (order_uid)
        ON DELETE CASCADE
);

-- Создаем таблицу 'payments'
-- Связываем ее с orders по transaction (который будет ссылаться на order_uid)
CREATE TABLE IF NOT EXISTS payments (
    transaction_uid VARCHAR(255) PRIMARY KEY, -- Переименовал transaction в transaction_uid для ясности
    request_id      BIGINT,
    currency        VARCHAR(10) NOT NULL,
    provider        VARCHAR(255) NOT NULL,
    amount          INT NOT NULL,
    payment_dt      BIGINT NOT NULL,
    bank            VARCHAR(255) NOT NULL,
    delivery_cost   INT NOT NULL,
    goods_total     INT NOT NULL,
    custom_fee      INT NOT NULL,
    CONSTRAINT fk_payment_order
        FOREIGN KEY (transaction_uid) -- Здесь transaction_uid из payments будет FK к order_uid в orders
        REFERENCES orders (order_uid)
        ON DELETE CASCADE
);

-- Создаем таблицу 'items'
-- Связываем ее с orders по track_number (в Order это TrackNumber)
CREATE TABLE IF NOT EXISTS items (
    id              SERIAL PRIMARY KEY, -- Добавляем ID для уникальности каждой позиции товара
    order_uid       VARCHAR(255) NOT NULL, -- Добавляем order_uid для связи с заказом
    chrt_id         BIGINT NOT NULL,
    track_number    VARCHAR(255) NOT NULL,
    price           INT NOT NULL,
    rid             VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    sale            INT NOT NULL,
    size            VARCHAR(50) NOT NULL,
    total_price     INT NOT NULL,
    nm_id           BIGINT NOT NULL,
    brand           VARCHAR(255) NOT NULL,
    status          INT NOT NULL,
    CONSTRAINT fk_items_order
        FOREIGN KEY (order_uid)
        REFERENCES orders (order_uid)
        ON DELETE CASCADE
);