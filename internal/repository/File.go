package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FileRepo interface {
	Save(file model.UploadFile) (model.UploadFile, error)
}

func NewFile(db *gorm.DB) FileRepo {
	return &fileRepo{db}
}

type fileRepo struct {
	db *gorm.DB
}

func (r *fileRepo) Save(file model.UploadFile) (model.UploadFile, error) {
	err := r.db.Save(&file).Error
	if err != nil {
		logger.Errorf("save upload file occur error")
		return model.UploadFile{}, err
	}
	return file, nil
}
