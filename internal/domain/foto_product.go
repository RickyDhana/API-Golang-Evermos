package domain

import "time"

type FotoProduk struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ProdukID  uint      `gorm:"index" json:"produk_id"`
	URL       string    `gorm:"size:500" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
