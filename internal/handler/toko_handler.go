package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"project_satu/internal/config"
	"project_satu/internal/domain"
	"project_satu/internal/middleware"
)

func CreateTokoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		NamaToko string `json:"nama_toko"`
		UrlToko  string `json:"url_toko"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	toko := domain.Toko{
		NamaToko:  req.NamaToko,
		UrlToko:   req.UrlToko,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&toko).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toko)
}

func ListTokoHandler(w http.ResponseWriter, r *http.Request) {
	var tokos []domain.Toko
	if err := config.DB.Find(&tokos).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokos)
}

func UpdateTokoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var toko domain.Toko
	if err := config.DB.First(&toko, id).Error; err != nil {
		http.Error(w, "toko not found", http.StatusNotFound)
		return
	}

	if toko.UserID != user.ID {
		http.Error(w, "forbidden - bukan pemilik toko", http.StatusForbidden)
		return
	}

	var req struct {
		NamaToko string `json:"nama_toko"`
		UrlToko  string `json:"url_toko"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	toko.NamaToko = req.NamaToko
	toko.UrlToko = req.UrlToko
	toko.UpdatedAt = time.Now()

	if err := config.DB.Save(&toko).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toko)
}

func DeleteTokoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var toko domain.Toko
	if err := config.DB.First(&toko, id).Error; err != nil {
		http.Error(w, "toko not found", http.StatusNotFound)
		return
	}

	if toko.UserID != user.ID {
		http.Error(w, "forbidden - bukan pemilik toko", http.StatusForbidden)
		return
	}

	if err := config.DB.Delete(&toko).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
