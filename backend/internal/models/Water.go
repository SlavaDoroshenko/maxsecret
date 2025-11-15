package models

import "time"

type Water struct {
	Id         int       `json:"id"`
	CountMl    int       `json:"count_ml"`
	UserId     int       `json:"user_id"`
	Created_at time.Time `json:"created_at"`
	Target     int       `json:"target_water"`
}
