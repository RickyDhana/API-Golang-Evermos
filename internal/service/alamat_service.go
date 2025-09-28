package service

import (
	"errors"

	"project_satu/internal/config"
	"project_satu/internal/domain"
)

type CreateAlamatRequest struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	Detail       string `json:"detail_alamat"`
	KotaSandi    string `json:"kota_sandi"`
}

type UpdateAlamatRequest struct {
	JudulAlamat  *string `json:"judul_alamat"`
	NamaPenerima *string `json:"nama_penerima"`
	NoTelp       *string `json:"no_telp"`
	Detail       *string `json:"detail_alamat"`
	KotaSandi    *string `json:"kota_sandi"`
}

func CreateAlamat(userID uint, req *CreateAlamatRequest) (*domain.Alamat, error) {
	al := &domain.Alamat{
		UserID:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		Detail:       req.Detail,
		KotaSandi:    req.KotaSandi,
	}
	if err := config.DB.Create(al).Error; err != nil {
		return nil, err
	}
	return al, nil
}

func GetAlamatByUser(userID uint) ([]domain.Alamat, error) {
	var list []domain.Alamat
	if err := config.DB.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetAlamatByID(userID, alamatID uint) (*domain.Alamat, error) {
	var al domain.Alamat
	if err := config.DB.Where("id = ? AND user_id = ?", alamatID, userID).First(&al).Error; err != nil {
		return nil, err
	}
	return &al, nil
}

func UpdateAlamat(userID, alamatID uint, req *UpdateAlamatRequest) (*domain.Alamat, error) {
	al, err := GetAlamatByID(userID, alamatID)
	if err != nil {
		return nil, errors.New("alamat not found or not owned by user")
	}

	if req.JudulAlamat != nil {
		al.JudulAlamat = *req.JudulAlamat
	}
	if req.NamaPenerima != nil {
		al.NamaPenerima = *req.NamaPenerima
	}
	if req.NoTelp != nil {
		al.NoTelp = *req.NoTelp
	}
	if req.Detail != nil {
		al.Detail = *req.Detail
	}
	if req.KotaSandi != nil {
		al.KotaSandi = *req.KotaSandi
	}

	if err := config.DB.Save(al).Error; err != nil {
		return nil, err
	}
	return al, nil
}

func DeleteAlamat(userID, alamatID uint) error {
	// ensure owned
	var al domain.Alamat
	if err := config.DB.Where("id = ? AND user_id = ?", alamatID, userID).First(&al).Error; err != nil {
		return errors.New("alamat not found or not owned by user")
	}
	if err := config.DB.Delete(&al).Error; err != nil {
		return err
	}
	return nil
}
