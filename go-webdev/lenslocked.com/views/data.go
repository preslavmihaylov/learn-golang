package views

import (
	"log"
	"net/http"
	"time"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	AlertMsgGeneric = "Something went wrong. Please try again and contact us if the problem persists."
)

type PublicError interface {
	error
	Public() string
}

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

type Alert struct {
	Level   string
	Message string
}

func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func RedirectWithAlert(w http.ResponseWriter, r *http.Request, url string, code int, a Alert) {
	persistAlert(w, a)
	http.Redirect(w, r, url, code)
}

func persistAlert(w http.ResponseWriter, a Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvlCookie := http.Cookie{
		Name:     "alert_level",
		Path:     "/",
		Value:    a.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	msgCookie := http.Cookie{
		Name:     "alert_msg",
		Path:     "/",
		Value:    a.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, &lvlCookie)
	http.SetCookie(w, &msgCookie)
}

func getPersistedAlert(r *http.Request) *Alert {
	lvlCookie, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}

	msgCookie, err := r.Cookie("alert_msg")
	if err != nil {
		return nil
	}

	return &Alert{
		Level:   lvlCookie.Value,
		Message: msgCookie.Value,
	}
}

func clearPersistedAlert(w http.ResponseWriter) {
	lvlCookie := http.Cookie{
		Name:     "alert_level",
		Path:     "/",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	msgCookie := http.Cookie{
		Name:     "alert_msg",
		Path:     "/",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, &lvlCookie)
	http.SetCookie(w, &msgCookie)
}
