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
	ShowGalleryRoute = "show_gallery"
)

type Galleries struct {
	NewView  *views.View
	ShowView *views.View
	EditView *views.View
	service  models.GalleryService
	router   *mux.Router
}

type NewGalleryForm struct {
	Title string `schema:"title"`
}

type EditGalleryForm NewGalleryForm

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		NewView:  views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		EditView: views.NewView("bootstrap", "galleries/edit"),
		service:  gs,
		router:   r,
	}
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data

	form := NewGalleryForm{}
	err := parseForm(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		g.NewView.Render(w, viewData)
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
		g.NewView.Render(w, viewData)
		return
	}

	redirectUrl, err := g.router.Get(ShowGalleryRoute).URL("id", strconv.Itoa(int(gallery.ID)))
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
	g.ShowView.Render(w, viewData)
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
	g.EditView.Render(w, viewData)
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
		g.EditView.Render(w, viewData)
		return
	}

	gallery.Title = form.Title
	err = g.service.Update(gallery)
	if err != nil {
		viewData.SetAlert(err)
		g.EditView.Render(w, viewData)
		return
	}

	viewData.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery updated successfully!",
	}

	g.EditView.Render(w, viewData)
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
