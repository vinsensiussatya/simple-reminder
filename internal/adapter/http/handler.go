package http

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"simple-reminder/internal/core"
	"simple-reminder/internal/usecase"
	"github.com/google/uuid"
)

type ReminderHandler struct {
	uc *usecase.ReminderUsecase
}

func NewReminderHandler(uc *usecase.ReminderUsecase) *ReminderHandler {
	return &ReminderHandler{uc: uc}
}

func (h *ReminderHandler) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/reminders", h.addReminder).Methods("POST")
	r.HandleFunc("/reminders", h.listReminders).Methods("GET")
	r.HandleFunc("/reminders/{id}", h.deleteReminder).Methods("DELETE")
	return r
}

type addReq struct {
	Message  string    `json:"message"`
	RemindAt time.Time `json:"remind_at"`
}

func (h *ReminderHandler) addReminder(w http.ResponseWriter, r *http.Request) {
	var req addReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rem := &core.Reminder{
		ID: uuid.NewString(),
		Message: req.Message,
		RemindAt: req.RemindAt,
	}
	if err := h.uc.Add(rem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rem)
}

func (h *ReminderHandler) listReminders(w http.ResponseWriter, r *http.Request) {
	rems, err := h.uc.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rems)
}

func (h *ReminderHandler) deleteReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.uc.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
