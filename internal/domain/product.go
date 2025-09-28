package domain

import "time"

type Produk struct {
	ID         uint         `gorm:"primary_key" json:"id"`
	NamaProduk string       `json:"nama_produk"`
	Slug       string       `gorm:"unique" json:"slug"`
	Deskripsi  string       `json:"deskripsi"`
	Stok       int          `json:"stok"`
	Harga      float64      `json:"harga"`
	CategoryID uint         `json:"category_id"`
	TokoID     uint         `json:"toko_id"`

	Foto []FotoProduk `gorm:"foreignKey:ProdukID" json:"foto"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
