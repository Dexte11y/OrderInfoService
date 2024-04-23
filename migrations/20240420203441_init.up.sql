CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),    
    locale VARCHAR(5),
    internal_signature TEXT,
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(255)
);

CREATE TABLE delivery (
    name VARCHAR(255),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(255),
    address TEXT,
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE payment (
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(3),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt TIMESTAMP,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);


CREATE TABLE items (
    chrt_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    track_number VARCHAR(255),
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INTEGER,
    size VARCHAR(50),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER
);

CREATE TABLE testorder (
    id SERIAL PRIMARY KEY,
    item VARCHAR(255),
    amount INTEGER
)