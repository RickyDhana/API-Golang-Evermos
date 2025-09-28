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

// admin only
func CreateKategoriHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil || !user.IsAdmin {
		http.Error(w, "forbidden - admin only", http.StatusForbidden)
		return
	}

	var req struct {
		Nama string `json:"nama"`
		Slug string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	kategori := domain.Category{
		Nama:      req.Nama,
		Slug:      req.Slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&kategori).Error; err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kategori)
}

func ListKategoriHandler(w http.ResponseWriter, r *http.Request) {
	var kategoris []domain.Category
	if err := config.DB.Find(&kategoris).Error; err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kategoris)
}

func UpdateKategoriHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil || !user.IsAdmin {
		http.Error(w, "forbidden - admin only", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var kategori domain.Category
	if err := config.DB.First(&kategori, id).Error; err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	var req struct {
		Nama string `json:"nama"`
		Slug string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	kategori.Nama = req.Nama
	kategori.Slug = req.Slug
	kategori.UpdatedAt = time.Now()

	if err := config.DB.Save(&kategori).Error; err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kategori)
}

func DeleteKategoriHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil || !user.IsAdmin {
		http.Error(w, "forbidden - admin only", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var kategori domain.Category
	if err := config.DB.First(&kategori, id).Error; err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&kategori).Error; err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
