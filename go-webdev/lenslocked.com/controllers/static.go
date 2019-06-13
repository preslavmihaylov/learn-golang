package controllers

import "github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"

type Static struct {
	HomeView     *views.View
	ContactsView *views.View
	FAQView      *views.View
}

func NewStatic() *Static {
	return &Static{
		HomeView:     views.NewView("bootstrap", "static/home"),
		ContactsView: views.NewView("bootstrap", "static/contacts"),
		FAQView:      views.NewView("bootstrap", "static/faq"),
	}
}
