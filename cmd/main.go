package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	reminderhttp "simple-reminder/internal/adapter/http"
	"simple-reminder/internal/adapter/repo/mem"
	"simple-reminder/internal/usecase"
)

func main() {
	repo := mem.NewReminderRepo()
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
