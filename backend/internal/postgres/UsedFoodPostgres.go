package postgres

import (
	"database/sql"
	"fmt"
	"myapp/internal/models"
)

type UsedFoodRepo struct {
	db *sql.DB
}

func NewUsedFoodRepo(db *sql.DB) *UsedFoodRepo {
	return &UsedFoodRepo{db: db}
}

func (ufr *UsedFoodRepo) AddDish(quantity float64, mealtype string, userid, foodid int) (*models.UsedFood, error) {
	if (mealtype != "завтрак" && mealtype != "обед" && mealtype != "ужин" && mealtype != "перекус") ||
		quantity <= 0 || userid <= 0 || foodid <= 0 {
		return nil, fmt.Errorf("invalid data")
	}
	var usedfood models.UsedFood
	err := ufr.db.QueryRow("INSERT INTO usedfood (meal_type, food_id, user_id, quantity) values ($1, $2, $3, $4) RETURNING *",
		mealtype, foodid, userid, quantity).Scan(&usedfood.Id, &usedfood.MealType, &usedfood.FoodId,
		&usedfood.UserId, &usedfood.Quantity, &usedfood.Created_at)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return &usedfood, nil
}

func (ufr *UsedFoodRepo) DeleteDish(id int) (*int, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}
	res, err := ufr.db.Exec("delete from usedfood where id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	affec, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if affec == 0 {
		return nil, fmt.Errorf("dish not found")
	}
	return &id, nil
}

func (ufr *UsedFoodRepo) GetAllUsedFoodByTimestamp(id int, timestamp string, meal_type string) ([]*models.UsedFoodWithNutrition, error) {
	if id <= 0 || timestamp == "" || (meal_type != "завтрак" && meal_type != "ужин" && meal_type != "обед" && meal_type != "перекус") {
		return nil, fmt.Errorf("invalid data")
	}
	query := `SELECT 
    		uf.id,
    		uf.meal_type,
    		uf.food_id,
    		uf.user_id,
			uf.quantity,
			uf.created_at,
			f."name",
			(f.proteins / 100.0) * uf.quantity AS proteins,
			(f.fats / 100.0) * uf.quantity AS fats,
			(f.carbos / 100.0) * uf.quantity AS carbos,
			((f.calorie::DECIMAL / 100.0) * uf.quantity)::INTEGER AS calorie
		FROM usedfood uf
		JOIN food f ON f.id = uf.food_id
		WHERE uf.user_id = $1
		AND uf.meal_type = $2
		AND DATE(uf.created_at) = $3
		ORDER BY uf.created_at DESC;`

	rows, err := ufr.db.Query(query, id, meal_type, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to query used_food: %w", err)
	}
	defer rows.Close()

	var results []*models.UsedFoodWithNutrition
	for rows.Next() {
		var uf models.UsedFoodWithNutrition
		err := rows.Scan(&uf.Id, &uf.MealType, &uf.FoodId, &uf.UserId, &uf.Quantity, &uf.CreatedAt, &uf.Name,
			&uf.Proteins, &uf.Fats, &uf.Carbos, &uf.Calorie)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, &uf)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return results, nil
}

func (ufr *UsedFoodRepo) GetUpperTable(id int, timestamp string) (*models.DailyNutrition, error) {
	if id <= 0 || timestamp == "" {
		return nil, fmt.Errorf("invalid data")
	}
	query := `SELECT
		u.id,
		u.name,
		u.target_cal,
		u.target_proteins,
		u.target_fats,
		u.target_carbos,
		COALESCE(SUM((f.proteins / 100.0) * uf.quantity), 0) AS total_proteins,
		COALESCE(SUM((f.fats / 100.0) * uf.quantity), 0) AS total_fats,
		COALESCE(SUM((f.carbos / 100.0) * uf.quantity), 0) AS total_carbos,
		COALESCE(SUM(((f.calorie::DECIMAL / 100.0) * uf.quantity)::INTEGER), 0) AS total_calories,
		(u.target_cal - COALESCE(SUM(((f.calorie::DECIMAL / 100.0) * uf.quantity)::INTEGER), 0)) AS remaining_calories
	FROM users u
	LEFT JOIN usedfood uf ON u.id = uf.user_id
	LEFT JOIN food f ON uf.food_id = f.id
	WHERE u.id = $1 
	AND DATE(uf.created_at) = $2
	GROUP BY u.id, u.name, u.target_cal, u.target_proteins, u.target_fats, u.target_carbos;`

	var info models.DailyNutrition
	err := ufr.db.QueryRow(query, id, timestamp).Scan(&info.ID, &info.Name, &info.TargetCal, &info.TargetProteins,
		&info.TargetFats, &info.TargetCarbos, &info.TotalProteins, &info.TotalFats, &info.TotalCarbos,
		&info.TotalCalories, &info.RemainingCalories)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return &info, nil
}
