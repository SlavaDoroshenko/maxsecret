package models

type Food struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Proteins float64 `json:"proteins"`
	Fats     float64 `json:"fats"`
	Carbos   float64 `json:"carbos"`
	Calorie  int     `json:"calorie"`
	Types    string  `json:"types"`
}
