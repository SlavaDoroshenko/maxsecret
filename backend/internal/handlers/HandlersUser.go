package handlers

import (
	"encoding/json"
	"myapp/internal/models"
	"myapp/internal/postgres"
	"net/http"
	"strconv"
)

type HandlerUser struct {
	repo *postgres.UserRepo
}

func NewHandlerUser(repo *postgres.UserRepo) *HandlerUser {
	return &HandlerUser{repo: repo}
}

func (hu HandlerUser) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "{error: Method not allowed}", http.StatusMethodNotAllowed)
		return
	}
	value := r.URL.Query().Get("id")
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}
	user, err := hu.repo.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (hu HandlerUser) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "{error: Method not allowed}", http.StatusMethodNotAllowed)
		return
	}
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	age := 0
	if user.Age != nil {
		age = *user.Age
	}

	weight := 0.0
	if user.Weight != nil {
		weight = *user.Weight
	}

	height := 0.0
	if user.Height != nil {
		height = *user.Height
	}

	target := ""
	if user.Target != nil {
		target = *user.Target
	}

	sex := ""
	if user.Sex != nil {
		sex = *user.Sex
	}

	updatedUser, err := hu.repo.UpdateUser(user.Id, age, weight, height, target, sex)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

// PASS:Пока не нужная функция
// func (hu HandlerUser) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if r.Method != http.MethodDelete {
// 		http.Error(w, "{error: Method not allowed}", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	if r.Body == nil {
// 		http.Error(w, "{error: Method not allowed}", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	var id int
// 	err := json.NewDecoder(r.Body).Decode(&id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{"error":err.Error()})
// 		return
// 	}

// }
