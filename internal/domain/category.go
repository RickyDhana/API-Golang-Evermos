package domain

import "time"

type Category struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Nama      string `gorm:"unique" json:"nama"`
	Slug      string `gorm:"unique" json:"slug"`

	Produk []Produk `gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
