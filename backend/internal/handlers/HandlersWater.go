package handlers

import (
	"encoding/json"
	"myapp/internal/postgres"
	"net/http"
	"strconv"
)

type WaterHandler struct {
	repo *postgres.WaterRepo
}

func NewWaterHandler(repo *postgres.WaterRepo) *WaterHandler {
	return &WaterHandler{repo: repo}
}

func (wh *WaterHandler) GetWaterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	timestamp := r.URL.Query().Get("time")

	if idStr == "" || timestamp == "" {
		http.Error(w, "id and time are required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	water, err := wh.repo.GetWater(id, timestamp)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(water)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(water)
}

func (wh *WaterHandler) UpsertWaterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserId    int    `json:"user_id"`
		CountMl   int    `json:"count_ml"`
		Timestamp string `json:"timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserId <= 0 || req.CountMl <= 0 || req.Timestamp == "" {
		http.Error(w, "user_id, count_ml, and timestamp are required", http.StatusBadRequest)
		return
	}

	water, err := wh.repo.GetWater(req.UserId, req.Timestamp)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = wh.repo.CreateWater(req.UserId, req.CountMl, req.Timestamp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			water, _ = wh.repo.GetWater(req.UserId, req.Timestamp)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		newCount := water.CountMl + req.CountMl
		err = wh.repo.UpdateWater(water.Id, newCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(water)
}
