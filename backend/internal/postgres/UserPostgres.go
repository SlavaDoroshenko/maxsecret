package postgres

import (
	"database/sql"
	"fmt"
	"myapp/internal/models"
	"strings"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) GetUser(id int) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(`select * from users
					   where id = $1`, id).Scan(&user.Id, &user.Name, &user.Weight, &user.Height,
		&user.Age, &user.Sex,
		&user.TargetCal, &user.TargetWat, &user.TargetProt, &user.TargetFats, &user.TargetCarbos,
		&user.TargetBreakfast, &user.TargetLunch, &user.TargetDinner, &user.TargetNosh, &user.Target, &user.Created_at)
	if err != nil {
		return nil, fmt.Errorf("error in query: %s", err)
	}
	return &user, nil
}

func (u *UserRepo) UpdateUser(id, age int, weight float64, height float64, target string, sex string) (*models.User, error) {
	var args []any
	count := 1
	var setParts []string

	if weight > 0 {
		setParts = append(setParts, fmt.Sprintf(`weight = $%d`, count))
		args = append(args, weight)
		count++
	}
	if height > 0 {
		setParts = append(setParts, fmt.Sprintf(`height = $%d`, count))
		args = append(args, height)
		count++
	}
	if target == "профицит" || target == "дефицит" || target == "поддерживание" {
		setParts = append(setParts, fmt.Sprintf(`target = $%d`, count))
		args = append(args, target)
		count++
	}
	if sex == "мужской" || sex == "женский" {
		setParts = append(setParts, fmt.Sprintf(`sex = $%d`, count))
		args = append(args, sex)
		count++
	}
	if age > 0 {
		setParts = append(setParts, fmt.Sprintf(`age = $%d`, count))
		args = append(args, age)
		count++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d RETURNING *`,
		strings.Join(setParts, ", "), count)
	args = append(args, id)

	var user models.User
	err := u.db.QueryRow(query, args...).Scan(&user.Id, &user.Name, &user.Weight, &user.Height,
		&user.Age, &user.Sex,
		&user.TargetCal, &user.TargetWat, &user.TargetProt, &user.TargetFats, &user.TargetCarbos,
		&user.TargetBreakfast, &user.TargetLunch, &user.TargetDinner, &user.TargetNosh, &user.Target, &user.Created_at)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	return &user, nil
}

// PASS:Пока не нужная функция
// func (u *UserRepo) DeleteUser(id int) (int, error) {
// 	var deletedID int
// 	err := u.db.QueryRow("DELETE FROM users WHERE id = $1 RETURNING id", id).Scan(&deletedID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return 0, fmt.Errorf("user with id %d not found", id)
// 		}
// 		return 0, fmt.Errorf("error deleting user: %w", err)
// 	}
// 	return deletedID, nil
// }
