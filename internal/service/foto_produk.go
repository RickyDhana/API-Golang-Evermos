package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"project_satu/internal/config"
	"project_satu/internal/domain"
)

type CreateProdukRequest struct {
	NamaProduk string  `json:"nama_produk"`
	Slug       string  `json:"slug"`
	Deskripsi  string  `json:"deskripsi"`
	Stok       int     `json:"stok"`
	Harga      float64 `json:"harga"`
	CategoryID uint    `json:"category_id"`
}

type UpdateProdukRequest struct {
	NamaProduk *string  `json:"nama_produk"`
	Slug       *string  `json:"slug"`
	Deskripsi  *string  `json:"deskripsi"`
	Stok       *int     `json:"stok"`
	Harga      *float64 `json:"harga"`
	CategoryID *uint    `json:"category_id"`
}

func CreateProduk(userID uint, req *CreateProdukRequest) (*domain.Produk, error) {
	var toko domain.Toko
	if err := config.DB.Where("user_id = ?", userID).First(&toko).Error; err != nil {
		return nil, fmt.Errorf("user belum punya toko")
	}

	prod := &domain.Produk{
		NamaProduk: req.NamaProduk,
		Slug:       req.Slug,
		Deskripsi:  req.Deskripsi,
		Stok:       req.Stok,
		Harga:      req.Harga,
		CategoryID: req.CategoryID,
		TokoID:     toko.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := config.DB.Create(prod).Error; err != nil {
		return nil, err
	}
	return prod, nil
}

func GetProdukList(page, perPage int, categoryID uint, minPrice, maxPrice float64, q string) ([]domain.Produk, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	tx := config.DB.Model(&domain.Produk{}).Preload("Foto")

	if categoryID > 0 {
		tx = tx.Where("category_id = ?", categoryID)
	}
	if minPrice > 0 {
		tx = tx.Where("harga >= ?", minPrice)
	}
	if maxPrice > 0 {
		tx = tx.Where("harga <= ?", maxPrice)
	}
	if q != "" {
		qLike := "%" + strings.ToLower(q) + "%"
		tx = tx.Where("LOWER(nama_produk) LIKE ?", qLike)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []domain.Produk
	if err := tx.Offset(offset).Limit(perPage).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

func GetProdukByID(id uint) (*domain.Produk, error) {
	var prod domain.Produk
	if err := config.DB.Preload("Foto").First(&prod, id).Error; err != nil {
		return nil, err
	}
	return &prod, nil
}

func UpdateProduk(prod *domain.Produk, req *UpdateProdukRequest) error {
	if req.NamaProduk != nil {
		prod.NamaProduk = *req.NamaProduk
	}
	if req.Slug != nil {
		prod.Slug = *req.Slug
	}
	if req.Deskripsi != nil {
		prod.Deskripsi = *req.Deskripsi
	}
	if req.Stok != nil {
		prod.Stok = *req.Stok
	}
	if req.Harga != nil {
		prod.Harga = *req.Harga
	}
	if req.CategoryID != nil {
		prod.CategoryID = *req.CategoryID
	}
	prod.UpdatedAt = time.Now()

	return config.DB.Save(prod).Error
}

func DeleteProduk(prod *domain.Produk) error {
	return config.DB.Delete(prod).Error
}

func UploadPhoto(prodID uint, header *multipart.FileHeader) (*domain.FotoProduk, error) {
	filename := fmt.Sprintf("uploads/%d_%s", prodID, filepath.Base(header.Filename))

	foto := &domain.FotoProduk{
		ProdukID:  prodID,
		URL:       filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := config.DB.Create(foto).Error; err != nil {
		return nil, err
	}
	return foto, nil
}
