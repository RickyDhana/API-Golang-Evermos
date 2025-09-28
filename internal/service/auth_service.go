package service

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"project_satu/internal/config"
	"project_satu/internal/domain"
)

var jwtKey = []byte("secret_key") 

func RegisterUser(nama, email, password, noHp string) (*domain.User, error) {
	var existing domain.User
	if err := config.DB.Where("email = ? OR no_hp = ?", email, noHp).First(&existing).Error; err == nil {
		return nil, errors.New("email atau no_hp sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Nama:     nama,
		Email:    email,
		Password: string(hashedPassword),
		NoHp:     noHp,
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	toko := &domain.Toko{
		NamaToko: nama + " Store",
		UrlToko:  nama + "-store", 
		UserID:   user.ID,
	}
	if err := config.DB.Create(toko).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func LoginUser(email, password string) (string, error) {
	var user domain.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("password salah")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   user.Email,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
