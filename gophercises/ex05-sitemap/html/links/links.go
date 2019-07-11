package links

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html document: %s", err)
	}

	return innerLinks(doc)
}

func innerLinks(parent *html.Node) ([]Link, error) {
	lnks := []Link{}
	for node := parent.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode && node.Data == "a" {
			var lnk Link
			var err error

			lnk.Href, err = findHrefAttr(node.Attr)
			if err != nil {
				// return nil, err
				continue // ignore bad links
			}

			lnk.Text = strings.Trim(innerHTML(node), " \n\t")
			lnks = append(lnks, lnk)
		} else {
			iLinks, err := innerLinks(node)
			if err != nil {
				return nil, err
			}

			lnks = append(lnks, iLinks...)
		}
	}

	return lnks, nil
}

func innerHTML(parent *html.Node) string {
	text := ""
	for node := parent.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.TextNode {
			text += node.Data
		}

		text += innerHTML(node)
	}

	return text
}

func findHrefAttr(attributes []html.Attribute) (string, error) {
	for _, attr := range attributes {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}

	return "", fmt.Errorf("no href attribute found")
}
