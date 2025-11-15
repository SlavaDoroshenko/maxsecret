CREATE TABLE food(
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    proteins DECIMAL(8, 2) NOT NULL,
    fats DECIMAL(8, 2) NOT NULL,
    carbos DECIMAL(8, 2) NOT NULL,
    calorie INTEGER NOT NULL,
    types VARCHAR(128) NOT NULL check (types in ('мл', 'г', 'шт')) DEFAULT 'г'
);

