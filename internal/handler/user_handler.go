package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"project_satu/internal/middleware"
	"project_satu/internal/service"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type responseUser struct {
		ID       uint   `json:"id"`
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		NoHp     string `json:"no_hp"`
		IsAdmin  bool   `json:"is_admin"`
		Pekerjaan string `json:"pekerjaan"`
		Tentang   string `json:"tentang"`
		Created   string `json:"created_at"`
		Modified  string `json:"updated_at"`
	}

	resp := responseUser{
		ID:        user.ID,
		Nama:      user.Nama,
		Email:     user.Email,
		NoHp:      user.NoHp,
		IsAdmin:   user.IsAdmin,
		Pekerjaan: user.Pekerjaan,
		Tentang:   user.Tentang,
		Created:   user.CreatedAt.Format(time.RFC3339),
		Modified:  user.UpdatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func UpdateMeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Nama      *string `json:"nama"`
		Email     *string `json:"email"`
		NoHp      *string `json:"no_hp"`
		Pekerjaan *string `json:"pekerjaan"`
		Tentang   *string `json:"tentang"`
		Password  *string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedUser, err := service.UpdateProfile(user.ID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type responseUser struct {
		ID       uint   `json:"id"`
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		NoHp     string `json:"no_hp"`
		IsAdmin  bool   `json:"is_admin"`
		Pekerjaan string `json:"pekerjaan"`
		Tentang   string `json:"tentang"`
		Created   string `json:"created_at"`
		Modified  string `json:"updated_at"`
	}

	resp := responseUser{
		ID:        updatedUser.ID,
		Nama:      updatedUser.Nama,
		Email:     updatedUser.Email,
		NoHp:      updatedUser.NoHp,
		IsAdmin:   updatedUser.IsAdmin,
		Pekerjaan: updatedUser.Pekerjaan,
		Tentang:   updatedUser.Tentang,
		Created:   updatedUser.CreatedAt.Format(time.RFC3339),
		Modified:  updatedUser.UpdatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
