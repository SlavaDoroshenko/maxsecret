package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"myapp/internal/handlers"
	"myapp/internal/postgres"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", `host=db port=5432 user=postgres 
						password=12345 dbname=maxhack sslmode=disable`)
	if err != nil {
		log.Fatalf("DB connect error:%s", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("DB ping error:%s", err)
	}

	//User
	repo := postgres.NewUserRepo(db)
	userhandler := handlers.NewHandlerUser(repo)

	//Food
	repofood := postgres.NewFoodRepo(db)
	foodhandler := handlers.NewHandlerFood(repofood)

	//UsedFood
	repoUsedFood := postgres.NewUsedFoodRepo(db)
	usedFoodHandler := handlers.NewUsedFoodHandler(repoUsedFood)

	//Water
	repoWater := postgres.NewWaterRepo(db)
	waterHandler := handlers.NewWaterHandler(repoWater)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"health": "success"})
	})
	mux.HandleFunc("/v1/user", userhandler.GetUserHandler)
	mux.HandleFunc("/v1/user/update", userhandler.UpdateUserHandler)
	mux.HandleFunc("/v1/food/all", foodhandler.GetAllFoodHandler)
	mux.HandleFunc("/v1/food/create", foodhandler.CreateFoodHandler)
	mux.HandleFunc("/v1/food/update", foodhandler.UpdateFoodHandler)
	mux.HandleFunc("/v1/food/name", foodhandler.GetByNameFoodHandler)
	mux.HandleFunc("/v1/usedfood/add", usedFoodHandler.AddDishHandler)
	mux.HandleFunc("/v1/usedfood/delete", usedFoodHandler.DeleteDishHandler)
	mux.HandleFunc("/v1/usedfood/get", usedFoodHandler.GetUsedFoodHandler)
	mux.HandleFunc("/v1/usedfood/getinfo", usedFoodHandler.GetUpperFoodHandler)
	mux.HandleFunc("/v1/water/get", waterHandler.GetWaterHandler)
	mux.HandleFunc("/v1/water/update", waterHandler.UpsertWaterHandler)
	log.Fatal(http.ListenAndServe(":3000", mux))
}
