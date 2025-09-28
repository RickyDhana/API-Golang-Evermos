package domain

import "time"

type Trx struct {
    ID            uint      `gorm:"primary_key" json:"id"`
    UserID        uint      `json:"user_id"`
    AlamatPengiriman string `gorm:"type:text" json:"alamat_pengiriman"` 
    HargaTotal    float64   `json:"harga_total"`
    KodeInvoice   string    `gorm:"unique" json:"kode_invoice"`
    Status        string    `json:"status"` 
    MethodBayar   string    `json:"method_bayar,omitempty"`

    DetailTrx []DetailTrx `gorm:"foreignKey:TrxID" json:"detail_trx"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type DetailTrx struct {
    ID        uint    `gorm:"primary_key" json:"id"`
    TrxID     uint    `json:"trx_id"`
    ProdukID  uint    `json:"produk_id"`
    NamaProduk string `json:"nama_produk"`
    Harga     float64 `json:"harga"`
    Qty       int     `json:"qty"`
    HargaTotal float64 `json:"harga_total"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
