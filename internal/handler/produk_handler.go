package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"project_satu/internal/config"
	"project_satu/internal/domain"
	"project_satu/internal/middleware"
	"project_satu/internal/service"
)

func CreateProdukHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req service.CreateProdukRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	prod, err := service.CreateProduk(user.ID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prod)
}

func ListProdukHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")
	categoryStr := r.URL.Query().Get("category_id")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)
	categoryID := uint(0)
	if categoryStr != "" {
		if v, err := strconv.Atoi(categoryStr); err == nil {
			categoryID = uint(v)
		}
	}
	minPrice, _ := strconv.ParseFloat(minPriceStr, 64)
	maxPrice, _ := strconv.ParseFloat(maxPriceStr, 64)

	products, total, err := service.GetProdukList(page, perPage, categoryID, minPrice, maxPrice, q)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	totalPages := (total + perPage - 1) / perPage

	res := map[string]interface{}{
		"data": products,
		"meta": map[string]interface{}{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func GetProdukHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	id := uint(id64)

	prod, err := service.GetProdukByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
}

// PUT /product/{id}
func UpdateProdukHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	id := uint(id64)

	var prod domain.Produk
	if err := config.DB.First(&prod, id).Error; err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	var toko domain.Toko
	if err := config.DB.First(&toko, prod.TokoID).Error; err != nil {
		http.Error(w, "toko not found", http.StatusBadRequest)
		return
	}
	if toko.UserID != user.ID {
		http.Error(w, "forbidden - not owner", http.StatusForbidden)
		return
	}

	var req service.UpdateProdukRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := service.UpdateProduk(&prod, &req); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&prod)
}

// DELETE /product/{id}
func DeleteProdukHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	id := uint(id64)

	var prod domain.Produk
	if err := config.DB.First(&prod, id).Error; err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	var toko domain.Toko
	if err := config.DB.First(&toko, prod.TokoID).Error; err != nil {
		http.Error(w, "toko not found", http.StatusBadRequest)
		return
	}
	if toko.UserID != user.ID {
		http.Error(w, "forbidden - not owner", http.StatusForbidden)
		return
	}

	if err := service.DeleteProduk(&prod); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /product/{id}/upload
func UploadProdukPhotoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	prodID := uint(id64)

	var prod domain.Produk
	if err := config.DB.First(&prod, prodID).Error; err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	var toko domain.Toko
	if err := config.DB.First(&toko, prod.TokoID).Error; err != nil {
		http.Error(w, "toko not found", http.StatusBadRequest)
		return
	}
	if toko.UserID != user.ID {
		http.Error(w, "forbidden - not owner", http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file missing: "+err.Error(), http.StatusBadRequest)
		return
	}
	file.Close()

	foto, err := service.UploadPhoto(prodID, header)
	if err != nil {
		http.Error(w, "upload error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foto)
}
