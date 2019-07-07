package models

import "github.com/jinzhu/gorm"

type GalleryService interface {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryValidator{
		GalleryDB: &galleryGorm{db},
	}
}
