package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"project_satu/internal/config"
	"project_satu/internal/domain"
)

type UpdateProfileInput struct {
	Nama      *string `json:"nama"`
	Email     *string `json:"email"`
	NoHp      *string `json:"no_hp"`
	Pekerjaan *string `json:"pekerjaan"`
	Tentang   *string `json:"tentang"`
	Password  *string `json:"password"`
}

func UpdateProfile(userID uint, input UpdateProfileInput) (*domain.User, error) {
	var user domain.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if input.Email != nil && *input.Email != user.Email {
		var tmp domain.User
		if err := config.DB.Where("email = ?", *input.Email).First(&tmp).Error; err == nil {
			return nil, errors.New("email already in use")
		}
		user.Email = *input.Email
	}

	if input.NoHp != nil && *input.NoHp != user.NoHp {
		var tmp domain.User
		if err := config.DB.Where("no_hp = ?", *input.NoHp).First(&tmp).Error; err == nil {
			return nil, errors.New("no_hp already in use")
		}
		user.NoHp = *input.NoHp
	}

	if input.Nama != nil {
		user.Nama = *input.Nama
	}
	if input.Pekerjaan != nil {
		user.Pekerjaan = *input.Pekerjaan
	}
	if input.Tentang != nil {
		user.Tentang = *input.Tentang
	}

	if input.Password != nil && *input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}
