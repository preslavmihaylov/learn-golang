package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	uuid "github.com/satori/go.uuid"
)

func main() {
	imageHandler := http.FileServer(http.Dir("./images"))
	http.HandleFunc("/", homeHandler)
	http.Handle("/images/", http.StripPrefix("/images/", imageHandler))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/transform", transformHandler)
	http.HandleFunc("/download/", downloadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("home.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Failed to parse form data")
		return
	}

	img := r.MultipartForm.File["image"][0]
	f, err := img.Open()
	if err != nil {
		panic(err)
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	uuid := uuid.Must(uuid.NewV4())
	err = ioutil.WriteFile(fmt.Sprintf("images/%s.png", uuid.String()), bs, 0666)
	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{Name: "image", Value: uuid.String(), Expires: time.Now().Add(36 * time.Hour)}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/transform", http.StatusMovedPermanently)
}

// TransformedImage ...
type TransformedImage struct {
	Original string
	Option1  string
	Option2  string
	Option3  string
}

func transformHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("image")
	if err != nil {
		panic(err)
	}

	img := c.Value
	orig := fmt.Sprintf("images/%s.png", img)
	opt1 := fmt.Sprintf("images/%s_1.png", img)
	opt2 := fmt.Sprintf("images/%s_2.png", img)
	opt3 := fmt.Sprintf("images/%s_3.png", img)
	timg := TransformedImage{Original: orig, Option1: opt1, Option2: opt2, Option3: opt3}

	transformImg(timg.Original, timg.Option1, 100)
	transformImg(timg.Original, timg.Option2, 300)
	transformImg(timg.Original, timg.Option3, 500)
	t, err := template.ParseFiles("transform.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, timg)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := stripDownloadPrefix(r.URL.Path)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=image.png")
	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, f)
}

func transformImg(input, output string, shapes int) {
	cmd := exec.Command("primitive", "-i", input, "-o", output, "-n", fmt.Sprintf("%d", shapes))
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func stripDownloadPrefix(path string) string {
	return path[10:]
}
