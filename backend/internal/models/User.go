package models

import "time"

type User struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Weight          *float64  `json:"weight"`
	Height          *float64  `json:"height"`
	Sex             *string   `json:"sex"`
	Age             *int      `json:"age"`
	TargetCal       *int      `json:"target_cal"`
	TargetWat       *int      `json:"target_water"`
	TargetProt      *float64  `json:"target_prot"`
	TargetCarbos    *float64  `json:"target_carb"`
	TargetFats      *float64  `json:"target_fats"`
	TargetBreakfast *int      `json:"target_breakfast"`
	TargetDinner    *int      `json:"target_dinner"`
	TargetLunch     *int      `json:"target_lunch"`
	TargetNosh      *int      `json:"target_nosh"`
	Target          *string   `json:"target"`
	Created_at      time.Time `json:"created_at"`
}
