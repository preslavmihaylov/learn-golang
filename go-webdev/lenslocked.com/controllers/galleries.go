package controllers

import (
	"fmt"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

type Galleries struct {
	NewView *views.View
	service models.GalleryService
}

type NewGalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		NewView: views.NewView("bootstrap", "galleries/new"),
		service: gs,
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

	fmt.Fprintln(w, gallery)
}
