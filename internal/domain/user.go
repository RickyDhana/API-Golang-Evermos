package domain

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Nama      string    `json:"nama"`
	Email     string    `gorm:"unique_index" json:"email"`
	NoHp      string    `gorm:"unique_index" json:"no_hp"`
	Password  string    `json:"-"` 
	IsAdmin   bool      `json:"is_admin"`
	Pekerjaan string    `json:"pekerjaan"` 
	Tentang   string    `json:"tentang"`   
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
