package controllers

import (
	"net/http"
	"strconv"

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

type Galleries struct {
	NewView   *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	service   models.GalleryService
	router    *mux.Router
}

type NewGalleryForm struct {
	Title string `schema:"title"`
}

type EditGalleryForm NewGalleryForm

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		NewView:   views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		service:   gs,
		router:    r,
	}
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {

	loggedInUser := context.User(r.Context())
	galleries, err := g.service.ByUserID(loggedInUser.ID)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var viewData views.Data
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
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, redirectUrl.Path, http.StatusFound)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var viewData views.Data
	viewData.Yield = gallery
	g.ShowView.Render(w, r, viewData)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
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
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, redirectUrl.Path, http.StatusFound)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.service.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong", http.StatusInternalServerError)
		}

		return nil, err
	}

	return gallery, nil
}
