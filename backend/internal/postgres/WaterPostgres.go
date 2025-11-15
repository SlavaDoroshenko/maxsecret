package postgres

import (
	"database/sql"
	"fmt"
	"myapp/internal/models"
)

type WaterRepo struct {
	db *sql.DB
}

func NewWaterRepo(db *sql.DB) *WaterRepo {
	return &WaterRepo{db: db}
}

func (wr *WaterRepo) GetWater(id int, timestamp string) (*models.Water, error) {
	if id <= 0 || timestamp == "" {
		return nil, fmt.Errorf("invalid data")
	}
	var water models.Water
	err := wr.db.QueryRow(`select w.id, w.count_ml, user_id, w.created_at, u.target_water from water w
							join users u on u.id = w.user_id
							where w.user_id = $1 and w.created_at = $2;`, id, timestamp).Scan(&water.Id,
		&water.CountMl, &water.UserId, &water.Created_at, &water.Target)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return &water, nil
}

func (wr *WaterRepo) CreateWater(userId int, countMl int, timestamp string) error {
	query := `INSERT INTO water (count_ml, user_id, created_at) VALUES ($1, $2, $3)`
	_, err := wr.db.Exec(query, countMl, userId, timestamp)
	if err != nil {
		return fmt.Errorf("failed to create water: %w", err)
	}
	return nil
}

func (wr *WaterRepo) UpdateWater(id int, newCountMl int) error {
	query := `UPDATE water SET count_ml = $1 WHERE id = $2`
	result, err := wr.db.Exec(query, newCountMl, id)
	if err != nil {
		return fmt.Errorf("failed to update water: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no water entry found with id %d", id)
	}
	return nil
}
