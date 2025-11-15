CREATE TABLE "users"(
    id INTEGER PRIMARY KEY,
    name VARCHAR(128) UNIQUE NOT NULL,
    weight DECIMAL(8,2),
    height DECIMAL(10,2),
    age INTEGER,
    sex VARCHAR(24) CHECK (sex in ('мужской', 'женский')),
    target_cal INTEGER,
    target_water INTEGER,
    target_proteins DECIMAL(8,2),
    target_fats DECIMAL(8,2),
    target_carbos DECIMAL(8,2),
    target_breakfast_cal INTEGER,
    target_lunch_cal INTEGER,
    target_dinner_cal INTEGER,
    target_nosh_cal INTEGER,
    target VARCHAR(128) CHECK (target IN ('профицит', 'дефицит', 'поддерживание')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_user_calculations()
RETURNS TRIGGER AS $$
DECLARE
    bmr INTEGER;
    calories INTEGER;
    proteins DECIMAL(8,2);
    fats DECIMAL(8,2);
    carbs DECIMAL(8,2);
BEGIN
    IF NEW.sex = 'мужской' THEN
        bmr := (10 * NEW.weight + 6.25 * NEW.height - 5 * NEW.age + 5)::INTEGER;
    ELSIF NEW.sex = 'женский' THEN
        bmr := (10 * NEW.weight + 6.25 * NEW.height - 5 * NEW.age - 161)::INTEGER;
    ELSE
        bmr := 2000;
    END IF;

    CASE NEW.target
        WHEN 'дефицит' THEN
            calories := GREATEST(bmr - 300, 1200);
        WHEN 'профицит' THEN
            calories := LEAST(bmr + 500, 4000);
        WHEN 'поддерживание' THEN
            calories := bmr;
        ELSE
            calories := bmr;
    END CASE;

    NEW.target_cal := calories;

    NEW.target_breakfast_cal := ROUND(calories * 0.25)::INTEGER;
    NEW.target_lunch_cal := ROUND(calories * 0.35)::INTEGER;
    NEW.target_nosh_cal := ROUND(calories * 0.15)::INTEGER;
    NEW.target_dinner_cal := ROUND(calories * 0.25)::INTEGER;

    NEW.target_water := GREATEST((NEW.weight * 30)::INTEGER, 1000);

    proteins := NEW.weight * 1.2; 
    fats := NEW.weight * 0.8; 
    carbs := GREATEST(
        (calories - (proteins * 4) - (fats * 9)) / 4, 
        0
    )::DECIMAL(8,2);

    NEW.target_proteins := proteins;
    NEW.target_fats := fats;
    NEW.target_carbos := carbs;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS user_calculations_trigger ON users;

CREATE TRIGGER user_calculations_trigger
    BEFORE INSERT OR UPDATE OF weight, height, age, sex, target
    ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_calculations();

