package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex13-quiet-hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	cache := make(map[int]*item)
	lastCacheReset := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Sub(lastCacheReset) > time.Minute*15 {
			lastCacheReset = time.Now()
			cache = make(map[int]*item)
		}

		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		storyCh := make(chan *item, numStories)
		var stories []*item

		currID := 0
		fetchedStories := 0
		for i := 0; i < numStories; i++ {
			getStory(cache, client, storyCh, ids[currID])
			currID++
		}

		for {
			select {
			case st := <-storyCh:
				if st != nil {
					if _, ok := cache[st.ID]; !ok {
						cache[st.ID] = st
					}

					stories = append(stories, st)
					fetchedStories++
				} else {
					getStory(cache, client, storyCh, ids[currID])
					currID++
				}
			}

			if fetchedStories == numStories {
				close(storyCh)
				break
			}
		}

		sort.Sort(byID(stories))
		err = tpl.Execute(w, templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		})
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getStory(cache map[int]*item, client hn.Client, storyCh chan *item, id int) {
	if st, ok := cache[id]; ok {
		storyCh <- st
		return
	}

	go getStoryAsync(client, storyCh, id)
}

func getStoryAsync(client hn.Client, storyCh chan *item, id int) {
	hnItem, err := client.GetItem(id)
	if err != nil {
		storyCh <- nil
		return
	}

	item := parseHNItem(hnItem)
	if isStoryLink(item) {
		storyCh <- &item
		return
	}

	storyCh <- nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []*item
	Time    time.Duration
}

type byID []*item

func (s byID) Len() int {
	return len(s)
}

func (s byID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
