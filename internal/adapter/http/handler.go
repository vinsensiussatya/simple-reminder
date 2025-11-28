package http

import (
	"encoding/json"
	"net/http"
	"os"
	"simple-reminder/internal/core"
	"simple-reminder/internal/usecase"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReminderHandler struct {
	uc *usecase.ReminderUsecase
}

func NewReminderHandler(uc *usecase.ReminderUsecase) *ReminderHandler {
	return &ReminderHandler{uc: uc}
}

func basicAuthMiddleware(next http.Handler) http.Handler {
	user := os.Getenv("BASIC_AUTH_USER")
	pass := os.Getenv("BASIC_AUTH_PASS")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != user || p != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *ReminderHandler) Router() http.Handler {
	r := mux.NewRouter()
	// API endpoints first!
	r.HandleFunc("/reminders", h.addReminder).Methods("POST")
	r.HandleFunc("/reminders", h.listReminders).Methods("GET")
	r.HandleFunc("/reminders/{id}", h.updateReminder).Methods("PUT")
	r.HandleFunc("/reminders/{id}", h.deleteReminder).Methods("DELETE")
	// Serve frontend static files for everything else
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("internal/adapter/http/static")))

	// Apply middleware by wrapping the router
	return basicAuthMiddleware(r)
}

type addReq struct {
	Message  string    `json:"message"`
	RemindAt time.Time `json:"remind_at"`
	Email    string    `json:"email"`
}

func (h *ReminderHandler) addReminder(w http.ResponseWriter, r *http.Request) {
	var req addReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rem := &core.Reminder{
		ID:       uuid.NewString(),
		Message:  req.Message,
		RemindAt: req.RemindAt,
		Email:    req.Email,
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

func (h *ReminderHandler) updateReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Message  string    `json:"message"`
		RemindAt time.Time `json:"remind_at"`
		Email    string    `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.uc.Update(id, req.Message, req.RemindAt, req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
