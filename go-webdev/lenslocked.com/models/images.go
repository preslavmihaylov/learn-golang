package models

import (
	"fmt"
	"net/url"
	"path/filepath"
)

type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}

	return temp.String()
}

func (i *Image) RelativePath() string {
	galleryID := fmt.Sprintf("%v", i.GalleryID)

	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}
