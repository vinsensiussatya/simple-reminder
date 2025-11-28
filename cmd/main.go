package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	reminderhttp "simple-reminder/internal/adapter/http"
	"simple-reminder/internal/adapter/repo/mem"
	pgrepo "simple-reminder/internal/adapter/repo/postgres"
	"simple-reminder/internal/usecase"
	"simple-reminder/internal/core"
)

func main() {
	var repo core.ReminderRepo

	dbURL := os.Getenv("DB_URL")
	if dbURL != "" && len(dbURL) >= 11 && dbURL[:11] == "postgres://" {
		ctx := context.Background()
		pgRepo, err := pgrepo.NewReminderRepo(ctx)
		if err != nil {
			log.Fatalf("failed to connect to Postgres: %v", err)
		}
		repo = pgRepo
		log.Println("Using Postgres repository")
	} else {
		repo = mem.NewReminderRepo()
		log.Println("Using in-memory repository")
	}

	uc := usecase.NewReminderUsecase(repo)
	handler := reminderhttp.NewReminderHandler(uc)

	r := handler.Router()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
