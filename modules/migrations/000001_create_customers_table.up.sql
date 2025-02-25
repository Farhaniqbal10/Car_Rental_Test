CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(16) UNIQUE NOT NULL,
    phone_number VARCHAR(15) NOT NULL
);
