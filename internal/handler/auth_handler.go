package handler

import (
	"encoding/json"
	"net/http"

	"project_satu/internal/service"
)

// ===== Register Handler =====
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		Password string `json:"password"`
		NoHp     string `json:"no_hp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := service.RegisterUser(req.Nama, req.Email, req.Password, req.NoHp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ===== Login Handler =====
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := service.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
