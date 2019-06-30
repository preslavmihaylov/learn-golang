package models

import (
	"fmt"
	"path/filepath"
)

type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	galleryID := fmt.Sprintf("%v", i.GalleryID)

	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}
