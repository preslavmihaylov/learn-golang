package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var g Gallery
	err := first(gg.db.Where("id = ?", id), &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	err := gg.db.Create(gallery).Error
	if err != nil {
		return fmt.Errorf("failed to create gallery model: %s", err)
	}

	return nil
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	err := gg.db.Save(gallery).Error
	if err != nil {
		return fmt.Errorf("failed to update gallery model: %s", err)
	}

	return nil
}
