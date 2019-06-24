package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	err := gg.db.Create(gallery).Error
	if err != nil {
		return fmt.Errorf("failed to create gallery model: %s", err)
	}

	return nil
}
