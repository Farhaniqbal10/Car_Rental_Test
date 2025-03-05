CREATE TABLE booking (
    booking_id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    cars_id INT NOT NULL,
    start_period TIMESTAMP NOT NULL,
    end_period TIMESTAMP NOT NULL,
    total_cost BIGINT NOT NULL,
    finished BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE CASCADE,
    FOREIGN KEY (cars_id) REFERENCES cars(cars_id) ON DELETE CASCADE
);
