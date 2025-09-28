package domain

import "time"

type ProductLog struct {
    ID        uint   `gorm:"primary_key" json:"id"`
    ProdukID  uint   `json:"produk_id"`
    Nama      string `json:"nama"`
    Harga     int    `json:"harga"`
    Qty       int    `json:"qty"`
    OrderID   uint   `json:"order_id"`

    CreatedAt time.Time `json:"created_at"`
}
