package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
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

func (gg *galleryGorm) ByUserID(userID uint) ([]Gallery, error) {
	var gs []Gallery
	err := all(gg.db.Where("user_id = ?", userID), &gs)
	if err != nil {
		return nil, err
	}

	return gs, nil
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

func (gg *galleryGorm) Delete(id uint) error {
	delGallery := Gallery{Model: gorm.Model{ID: id}}
	err := gg.db.Delete(delGallery).Error
	if err != nil {
		return fmt.Errorf("failed to delete gallery model: %s", err)
	}

	return nil
}
