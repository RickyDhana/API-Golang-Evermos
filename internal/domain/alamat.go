package domain

import "time"

type Alamat struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	JudulAlamat  string    `gorm:"size:100" json:"judul_alamat"`   
	NamaPenerima string    `gorm:"size:100" json:"nama_penerima"`
	NoTelp       string    `gorm:"size:30" json:"no_telp"`
	Detail       string    `gorm:"type:text" json:"detail_alamat"`
	KotaSandi    string    `gorm:"size:100" json:"kota_sandi"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
