DROP TRIGGER IF EXISTS user_calculations_trigger ON users;
DROP FUNCTION IF EXISTS update_user_calculations();
DROP FUNCTION IF EXISTS calculate_water(DECIMAL);
DROP FUNCTION IF EXISTS calculate_calories(DECIMAL, DECIMAL, INTEGER, VARCHAR, VARCHAR);
drop table "users";