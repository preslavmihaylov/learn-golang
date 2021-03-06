package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

const (
	IndexGalleriesRoute = "index_galleries"
	ShowGalleryRoute    = "show_gallery"
	EditGalleryRoute    = "edit_gallery"
)

const (
	maxMultipartMem = 1 << 20 // 1 MB
)

type Galleries struct {
	NewView      *views.View
	ShowView     *views.View
	EditView     *views.View
	IndexView    *views.View
	service      models.GalleryService
	imageService models.ImageService
	router       *mux.Router
}

type NewGalleryForm struct {
	Title string `schema:"title"`
}

type EditGalleryForm NewGalleryForm

func NewGalleries(gs models.GalleryService, is models.ImageService, r *mux.Router) *Galleries {
	return &Galleries{
		NewView:      views.NewView("bootstrap", "galleries/new"),
		ShowView:     views.NewView("bootstrap", "galleries/show"),
		EditView:     views.NewView("bootstrap", "galleries/edit"),
		IndexView:    views.NewView("bootstrap", "galleries/index"),
		service:      gs,
		imageService: is,
		router:       r,
	}
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data

	loggedInUser := context.User(r.Context())
	galleries, err := g.service.ByUserID(loggedInUser.ID)
	if err != nil {
		viewData.SetAlert(err)
		g.IndexView.Render(w, r, viewData)
		return
	}

	viewData.Yield = galleries
	g.IndexView.Render(w, r, viewData)
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data

	form := NewGalleryForm{}
	err := parseForm(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		g.NewView.Render(w, r, viewData)
		return
	}

	loggedInUser := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: loggedInUser.ID,
	}

	err = g.service.Create(&gallery)
	if err != nil {
		viewData.SetAlert(err)
		g.NewView.Render(w, r, viewData)
		return
	}

	redirectUrl, err := g.router.Get(EditGalleryRoute).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Something went wrong...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	a := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "You have successfully created a gallery!",
	}

	views.RedirectWithAlert(w, r, redirectUrl.Path, http.StatusFound, a)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	var viewData views.Data
	viewData.Yield = gallery
	g.ShowView.Render(w, r, viewData)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	usr := context.User(r.Context())
	if gallery.UserID != usr.ID {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var viewData views.Data
	viewData.Yield = gallery
	g.EditView.Render(w, r, viewData)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	usr := context.User(r.Context())
	if gallery.UserID != usr.ID {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var viewData views.Data
	viewData.Yield = gallery

	var form EditGalleryForm
	err = parseForm(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		g.EditView.Render(w, r, viewData)
		return
	}

	gallery.Title = form.Title
	err = g.service.Update(gallery)
	if err != nil {
		viewData.SetAlert(err)
		g.EditView.Render(w, r, viewData)
		return
	}

	viewData.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery updated successfully!",
	}

	g.EditView.Render(w, r, viewData)
}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	usr := context.User(r.Context())
	if gallery.UserID != usr.ID {
		http.Error(w, "You do not have permission to delete this gallery", http.StatusForbidden)
		return
	}

	var viewData views.Data
	err = g.service.Delete(gallery.ID)
	if err != nil {
		viewData.SetAlert(err)
		viewData.Yield = gallery
		g.EditView.Render(w, r, viewData)
		return
	}

	redirectUrl, err := g.router.Get(IndexGalleriesRoute).URL()
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Something went wrong...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	a := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "You have successfully deleted the gallery...",
	}

	views.RedirectWithAlert(w, r, redirectUrl.Path, http.StatusFound, a)
}

func (g *Galleries) ImageViaLink(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	usr := context.User(r.Context())
	if gallery.UserID != usr.ID {
		http.Error(w, "gallery not found", http.StatusForbidden)
		return
	}

	var viewData views.Data
	viewData.Yield = gallery
	err = r.ParseForm()
	if err != nil {
		viewData.SetAlert(err)
		g.EditView.Render(w, r, viewData)
		return
	}

	imgURLs := r.PostForm["files"]
	err = g.downloadImagesViaURLs(gallery.ID, imgURLs)
	if err != nil {
		views.RedirectWithAlert(w, r, "/", http.StatusFound, views.Alert{
			Level:   views.AlertLvlError,
			Message: "There was a problem getting some of the images...",
		})
		return
	}

	url, err := g.router.Get(EditGalleryRoute).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		views.RedirectWithAlert(w, r, "/", http.StatusFound, views.Alert{
			Level:   views.AlertLvlError,
			Message: "Something went wrong...",
		})
		return
	}

	views.RedirectWithAlert(w, r, url.Path, http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "You have successfully added your images",
	})
}

func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	usr := context.User(r.Context())
	if gallery.UserID != usr.ID {
		http.Error(w, "gallery not found", http.StatusForbidden)
		return
	}

	var viewData views.Data
	viewData.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		viewData.SetAlert(err)
		g.EditView.Render(w, r, viewData)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			viewData.SetAlert(err)
			g.EditView.Render(w, r, viewData)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("wasn't able to close image file: %s", err)
			}
		}()

		err = g.imageService.Create(gallery.ID, file, f.Filename)
		if err != nil {
			viewData.SetAlert(err)
			g.EditView.Render(w, r, viewData)
			return
		}
	}

	url, err := g.router.Get(EditGalleryRoute).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Something went wrong...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	a := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "You have successfully uploaded an image...",
	}

	views.RedirectWithAlert(w, r, url.Path, http.StatusFound, a)
}

func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Gallery not found...",
		}

		views.RedirectWithAlert(w, r, "/", http.StatusFound, a)
		return
	}

	usr := context.User(r.Context())
	if gallery.UserID != usr.ID {
		http.Error(w, "You do not have permission to edit this gallery or image", http.StatusForbidden)
		return
	}

	filename := mux.Vars(r)["filename"]
	img := models.Image{
		Filename:  filename,
		GalleryID: gallery.ID,
	}

	err = g.imageService.Delete(&img)
	if err != nil {
		var viewData views.Data
		viewData.Yield = gallery
		viewData.SetAlert(err)
		g.EditView.Render(w, r, viewData)
		return
	}

	url, err := g.router.Get(EditGalleryRoute).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		a := views.Alert{
			Level:   views.AlertLvlError,
			Message: "Something went wrong...",
		}

		views.RedirectWithAlert(w, r, "/galleries", http.StatusFound, a)
		return
	}

	a := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "You have successfully deleted an image...",
	}

	views.RedirectWithAlert(w, r, url.Path, http.StatusFound, a)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}

	gallery, err := g.service.ByID(uint(id))
	if err != nil {
		return nil, err
	}

	images, err := g.imageService.ByGalleryID(gallery.ID)
	if err != nil {
		return nil, err
	}

	gallery.Images = images

	return gallery, nil
}

func (g *Galleries) downloadImagesViaURLs(galleryID uint, imgURLs []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error)
	defer close(errChan)

	for _, imgURL := range imgURLs {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				errChan <- err
			}
			defer resp.Body.Close()

			_, filename := filepath.Split(url)
			err = g.imageService.Create(galleryID, resp.Body, filename)
			if err != nil {
				errChan <- err
			}
		}(imgURL)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	default:
	}

	return nil
}
