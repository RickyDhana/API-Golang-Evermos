package domain

type OrderItem struct {
    ID        uint `gorm:"primary_key" json:"id"`
    OrderID   uint `json:"order_id"`
    ProdukID  uint `json:"produk_id"`
    Qty       int  `json:"qty"`
    Subtotal  int  `json:"subtotal"`
}
