package models

import "time"

type UsedFood struct {
	Id         int       `json:"id"`
	MealType   string    `json:"meal_type"`
	FoodId     int       `json:"food_id"`
	UserId     int       `json:"user_id"`
	Quantity   float64   `json:"quantity"`
	Created_at time.Time `json:"created_at"`
}

type UsedFoodWithNutrition struct {
	Id        int       `json:"id"`
	MealType  string    `json:"meal_type"`
	FoodId    int       `json:"food_id"`
	UserId    int       `json:"user_id"`
	Quantity  float64   `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`

	Name     string  `json:"name"`
	Proteins float64 `json:"proteins"`
	Fats     float64 `json:"fats"`
	Carbos   float64 `json:"carbos"`
	Calorie  int     `json:"calorie"`
}

type NutritionTotal struct {
	Proteins float64 `json:"proteins"`
	Fats     float64 `json:"fats"`
	Carbos   float64 `json:"carbos"`
	Calorie  int     `json:"calorie"`
}

type MealData struct {
	Items []*UsedFoodWithNutrition `json:"items"`
	Total NutritionTotal           `json:"total"`
}

type Response struct {
	UsedFood map[string]MealData `json:"usedfood"`
}

type DailyNutrition struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	TargetCal         int     `json:"target_cal"`
	TargetProteins    float64 `json:"target_proteins"`
	TargetFats        float64 `json:"target_fats"`
	TargetCarbos      float64 `json:"target_carbos"`
	TotalProteins     float64 `json:"total_proteins"`
	TotalFats         float64 `json:"total_fats"`
	TotalCarbos       float64 `json:"total_carbos"`
	TotalCalories     int     `json:"total_calories"`
	RemainingCalories int     `json:"remaining_calories"`
}
