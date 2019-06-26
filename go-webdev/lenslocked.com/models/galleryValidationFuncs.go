package models

type galleryValidationFunc func(g *Gallery) error

func runGalleryValidationFuncs(g *Gallery, fns ...galleryValidationFunc) error {
	for _, fn := range fns {
		err := fn(g)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gv *galleryValidator) requireUserID(g *Gallery) error {
	if g.UserID == 0 {
		return ErrUserIDRequired
	}

	return nil
}

func (gv *galleryValidator) requireTitle(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}

	return nil
}

func (gv *galleryValidator) idGreaterThan(n uint) galleryValidationFunc {
	return galleryValidationFunc(func(g *Gallery) error {
		if g.ID <= n {
			return ErrIDInvalid
		}

		return nil
	})
}
