package handlers

import (
	"encoding/json"
	"myapp/internal/models"
	"myapp/internal/postgres"
	"net/http"
)

type HandlerFood struct {
	foodrepo *postgres.FoodRepo
}

func NewHandlerFood(foodrepo *postgres.FoodRepo) *HandlerFood {
	return &HandlerFood{foodrepo: foodrepo}
}

func (hf *HandlerFood) GetAllFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}
	foods, err := hf.foodrepo.GetAllFoods()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"foods": foods})
}

func (hf *HandlerFood) CreateFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "body required"})
		return
	}
	var food *models.Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	food, err = hf.foodrepo.CreateFood(food.Name, food.Proteins, food.Carbos, food.Fats,
		food.Calorie, food.Types)
	if err != nil {
		if err.Error() == "invalid data for create food" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(food)
}

func (hf *HandlerFood) UpdateFoodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "{error: Method not allowed}", http.StatusMethodNotAllowed)
		return
	}
	food := &models.Food{}
	json.NewDecoder(r.Body).Decode(&food)

	food, err := hf.foodrepo.UpdateFood(food.Id, food.Name, food.Calorie, food.Proteins, food.Carbos, food.Fats, food.Types)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(food)
}

func (hf *HandlerFood) GetByNameFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}
	value := r.URL.Query().Get("name")
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "query required"})
		return
	}
	foods, err := hf.foodrepo.GetAllFoodsByName(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"foods": foods})
}
