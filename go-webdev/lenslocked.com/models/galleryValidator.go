package models

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) Create(g *Gallery) error {
	err := runGalleryValidationFuncs(g, gv.requireTitle, gv.requireUserID)
	if err != nil {
		return err
	}

	return gv.GalleryDB.Create(g)
}

func (gv *galleryValidator) Update(g *Gallery) error {
	err := runGalleryValidationFuncs(g, gv.requireTitle, gv.requireUserID)
	if err != nil {
		return err
	}

	return gv.GalleryDB.Update(g)
}
