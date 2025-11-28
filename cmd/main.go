package main

import (
	"log"
	"net/http"

	reminderhttp "simple-reminder/internal/adapter/http"
	"simple-reminder/internal/adapter/repo/mem"
	"simple-reminder/internal/usecase"
)

func main() {
	repo := mem.NewReminderRepo()
	uc := usecase.NewReminderUsecase(repo)
	handler := reminderhttp.NewReminderHandler(uc)

	r := handler.Router()
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
