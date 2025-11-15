CREATE TABLE usedfood(
    id SERIAL PRIMARY KEY,
    meal_type VARCHAR(128) check (meal_type in ('завтрак', 'обед', 'ужин', 'перекус')) NOT NULL,
    food_id INTEGER REFERENCES food(id) NOT NULL,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    quantity DECIMAL(12, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);