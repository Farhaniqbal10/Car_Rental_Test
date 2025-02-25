CREATE TABLE cars (
    cars_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rent_price_daily BIGINT NOT NULL,
    stock INT NOT NULL
);
