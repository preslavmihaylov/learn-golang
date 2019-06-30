package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null;index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

func (g *Gallery) ImagesSplitN(n int) [][]Image {
	ret := make([][]Image, n)

	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}

	for i, img := range g.Images {
		ret[i%n] = append(ret[i%n], img)
	}

	return ret
}
