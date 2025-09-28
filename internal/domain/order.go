package domain

import "time"

type Order struct {
    ID        uint   `gorm:"primary_key" json:"id"`
    UserID    uint   `json:"user_id"`
    TokoID    uint   `json:"toko_id"`
    Status    string `json:"status"` 
    Total     int    `json:"total"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
