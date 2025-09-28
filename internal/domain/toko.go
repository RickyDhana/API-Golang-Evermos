package domain

import "time"

type Toko struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	NamaToko string `json:"nama_toko"`
	UrlToko  string `gorm:"unique" json:"url_toko"`
	UserID   uint   `json:"user_id"`

	Produk []Produk `gorm:"foreignKey:TokoID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
