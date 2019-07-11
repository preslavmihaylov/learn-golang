package sitemap

import (
	"fmt"
	"io"
	"log"
	"net/http"
	gourl "net/url"
	"strings"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex05-sitemap/html/links"
)

func Build(w io.Writer, domain string) error {
	lnk, err := sanitizeURL(domain, domain)
	if err != nil {
		return fmt.Errorf("invalid domain given")
	}

	fmt.Fprintln(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	fmt.Fprintln(w, "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">")

	traverseURL(w, domain, lnk, map[string]bool{lnk: true})

	fmt.Fprintln(w, "</urlset>")

	return nil
}

func traverseURL(w io.Writer, domain, currentLnk string, linksVisited map[string]bool) {
	fmt.Fprintln(w, "\t<url>\n\t\t<loc>"+currentLnk+"</loc>\n\t</url>")
	fmt.Println("Parsing " + currentLnk + "...")
	linksVisited[currentLnk] = true

	resp, err := http.Get(currentLnk)
	if err != nil {
		log.Fatalf("failed to get page %s: %s", currentLnk, err)
	}

	lnks, err := links.Parse(resp.Body)
	if err != nil {
		resp.Body.Close()
		log.Fatalf("failed to parse links from page %s: %s", currentLnk, err)
	}

	resp.Body.Close()
	for _, lnk := range lnks {
		l, err := sanitizeURL(domain, lnk.Href)
		if err != nil {
			continue
		}

		if !strings.HasPrefix(l, "https://"+domain) || linksVisited[l] {
			continue
		}

		traverseURL(w, domain, l, linksVisited)
	}
}

func sanitizeURL(baseHost, url string) (string, error) {
	if strings.HasPrefix(baseHost, "https://") {
		return "", fmt.Errorf("invalid host provided")
	}

	if strings.HasPrefix(url, "/") {
		url = "https://" + baseHost + url
	} else if strings.HasPrefix(url, baseHost) {
		url = "https://" + url
	}

	u, err := gourl.ParseRequestURI(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %s", url, err)
	}

	return (u.Scheme + "://" + u.Host + u.Path), nil
}
