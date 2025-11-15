package handlers

import (
	"encoding/json"
	"math"
	"myapp/internal/models"
	"myapp/internal/postgres"
	"net/http"
	"strconv"
	"sync"
)

type UsedFoodHandler struct {
	repo *postgres.UsedFoodRepo
}

func NewUsedFoodHandler(repo *postgres.UsedFoodRepo) *UsedFoodHandler {
	return &UsedFoodHandler{repo: repo}
}

func (ufh *UsedFoodHandler) AddDishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "body required"})
		return
	}
	var usedfood *models.UsedFood
	err := json.NewDecoder(r.Body).Decode(&usedfood)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	usedfood, err = ufh.repo.AddDish(usedfood.Quantity, usedfood.MealType, usedfood.UserId, usedfood.FoodId)
	if err != nil {
		if err.Error() == "invalid data" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usedfood)
}

func (ufh *UsedFoodHandler) DeleteDishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "body required"})
		return
	}
	type request struct {
		Id *int `json:"id"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	req.Id, err = ufh.repo.DeleteDish(*req.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"delete success id": req.Id})
}

func (ufh *UsedFoodHandler) GetUsedFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	timestamp := r.URL.Query().Get("time")
	if id == "" || timestamp == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "invalid query"})
		return
	}
	ids, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "id not integer"})
		return
	}

	var (
		breakfast []*models.UsedFoodWithNutrition
		lunch     []*models.UsedFoodWithNutrition
		dinner    []*models.UsedFoodWithNutrition
		snack     []*models.UsedFoodWithNutrition

		errBreakfast error
		errLunch     error
		errDinner    error
		errSnack     error
	)

	wg := &sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		breakfast, errBreakfast = ufh.repo.GetAllUsedFoodByTimestamp(ids, timestamp, `завтрак`)
	}()

	go func() {
		defer wg.Done()
		lunch, errLunch = ufh.repo.GetAllUsedFoodByTimestamp(ids, timestamp, `обед`)
	}()

	go func() {
		defer wg.Done()
		dinner, errDinner = ufh.repo.GetAllUsedFoodByTimestamp(ids, timestamp, `ужин`)
	}()

	go func() {
		defer wg.Done()
		snack, errSnack = ufh.repo.GetAllUsedFoodByTimestamp(ids, timestamp, `перекус`)
	}()
	wg.Wait()
	if errBreakfast != nil || errDinner != nil || errLunch != nil || errSnack != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "trouble with get data"})
		return
	}
	calculateTotal := func(foods []*models.UsedFoodWithNutrition) models.NutritionTotal {
		var total models.NutritionTotal
		for _, item := range foods {
			total.Proteins += item.Proteins
			total.Fats += item.Fats
			total.Carbos += item.Carbos
			total.Calorie += item.Calorie
		}
		total.Proteins = roundFloat(total.Proteins, 2)
		total.Fats = roundFloat(total.Fats, 2)
		total.Carbos = roundFloat(total.Carbos, 2)
		return total
	}

	usedFoodMap := make(map[string]models.MealData)

	usedFoodMap["breakfast"] = models.MealData{
		Items: breakfast,
		Total: calculateTotal(breakfast),
	}

	usedFoodMap["lunch"] = models.MealData{
		Items: lunch,
		Total: calculateTotal(lunch),
	}

	usedFoodMap["dinner"] = models.MealData{
		Items: dinner,
		Total: calculateTotal(dinner),
	}

	usedFoodMap["snack"] = models.MealData{
		Items: snack,
		Total: calculateTotal(snack),
	}

	response := models.Response{
		UsedFood: usedFoodMap,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (ufh *UsedFoodHandler) GetUpperFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	timestamp := r.URL.Query().Get("time")
	if id == "" || timestamp == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "invalid query"})
		return
	}
	ids, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "id not integer"})
		return
	}
	info, err := ufh.repo.GetUpperTable(ids, timestamp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "id not integer"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}

func roundFloat(val float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(val*multiplier) / multiplier
}
