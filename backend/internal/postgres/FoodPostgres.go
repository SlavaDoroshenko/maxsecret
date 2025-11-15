package postgres

import (
	"database/sql"
	"fmt"
	"myapp/internal/models"
	"strings"
)

type FoodRepo struct {
	db *sql.DB
}

func NewFoodRepo(db *sql.DB) *FoodRepo {
	return &FoodRepo{db: db}
}

func (fr *FoodRepo) GetAllFoods() ([]models.Food, error) {
	rows, err := fr.db.Query("SELECT * FROM food")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var food models.Food
		err := rows.Scan(&food.Id, &food.Name, &food.Proteins, &food.Fats, &food.Carbos,
			&food.Calorie, &food.Types)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		foods = append(foods, food)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return foods, nil
}

func (fr *FoodRepo) GetAllFoodsByName(name string) ([]models.Food, error) {
	rows, err := fr.db.Query("SELECT * FROM food where name ILIKE $1 || '%'", name)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var food models.Food
		err := rows.Scan(&food.Id, &food.Name, &food.Proteins, &food.Fats, &food.Carbos,
			&food.Calorie, &food.Types)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		foods = append(foods, food)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return foods, nil
}

func (fr *FoodRepo) CreateFood(name string, proteins, carbos, fats float64, calories int, types string) (*models.Food, error) {
	if name == "" ||
		proteins <= 0 ||
		carbos <= 0 ||
		fats <= 0 ||
		calories <= 0 ||
		types != "шт" && types != "мл" && types != "г" {
		return nil, fmt.Errorf("invalid data for create food")
	}
	var food models.Food
	err := fr.db.QueryRow(`insert into food (name, proteins, fats, carbos, calorie, types) values ($1, $2, $3, $4, $5, $6)
						returning *`, name, proteins, fats, carbos, calories, types).Scan(&food.Id, &food.Name, &food.Proteins, &food.Fats, &food.Carbos, &food.Calorie, &food.Types)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return &food, nil
}

func (u *FoodRepo) UpdateFood(id int, name string, calorie int, proteins, carbos, fats float64, types string) (*models.Food, error) {
	var args []any
	count := 1
	var setParts []string

	if calorie > 0 {
		setParts = append(setParts, fmt.Sprintf(`calorie = $%d`, count))
		args = append(args, calorie)
		count++
	}
	if carbos > 0 {
		setParts = append(setParts, fmt.Sprintf(`carbos = $%d`, count))
		args = append(args, carbos)
		count++
	}
	if fats > 0 {
		setParts = append(setParts, fmt.Sprintf(`fats = $%d`, count))
		args = append(args, fats)
		count++
	}
	if proteins > 0 {
		setParts = append(setParts, fmt.Sprintf(`proteins = $%d`, count))
		args = append(args, proteins)
		count++
	}
	if name != "" {
		setParts = append(setParts, fmt.Sprintf(`name = $%d`, count))
		args = append(args, name)
		count++
	}
	if types != "" && (types == "мл" || types == "г" || types == "шт") {
		setParts = append(setParts, fmt.Sprintf(`types = $%d`, count))
		args = append(args, types)
		count++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`UPDATE food SET %s WHERE id = $%d RETURNING *`,
		strings.Join(setParts, ", "), count)
	args = append(args, id)

	var food models.Food
	err := u.db.QueryRow(query, args...).Scan(&food.Id, &food.Name, &food.Proteins, &food.Fats,
		&food.Carbos, &food.Calorie, &food.Types)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	return &food, nil
}
