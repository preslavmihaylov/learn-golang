package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

func main() {
	twID := flag.String("id", "", "id of contest tweet")
	flag.Parse()

	if *twID == "" {
		log.Fatalf("twitter ID is required")
	}

	client, err := ConnectToTwitterAPI()
	if err != nil {
		log.Fatalf("failed to connect to twitter api: %s", err)
	}

	tweet, err := GetTweetByID(client, *twID)
	if err != nil {
		log.Fatalf("failed to get tweet by id: %s", err)
	}

	fmt.Println("Text of competition tweet")
	fmt.Println(tweet.Text())

	fmt.Println("\nList of participants:")
	retweets, err := GetRetweetsByID(client, "1176928103526014979")
	if err != nil {
		log.Fatalf("failed to get retweets for tweet: %s", err)
	}

	winner := PickAWinner(retweets)
	for _, retw := range retweets {
		fmt.Printf("\t%s\n", retw.User().Name())
	}

	fmt.Printf("\nThe winner is %s (%d)! Congratulations!\n", winner.User().Name(), winner.User().Id())
}

func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("failed to open CREDENTIALS file (https://github.com/kurrik/twittergo-examples): %s", err)
	}

	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}

	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)

	return client, nil
}

func ConnectToTwitterAPI() (*twittergo.Client, error) {
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
		user   *twittergo.User
	)
	client, err = LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("Could not parse CREDENTIALS file: %v\n", err)
	}

	req, err = http.NewRequest("GET", "/1.1/account/verify_credentials.json", nil)
	if err != nil {
		return nil, fmt.Errorf("Could not parse request: %v\n", err)
	}

	resp, err = client.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("Could not send request: %v\n", err)
	}

	user = &twittergo.User{}
	err = resp.Parse(user)
	if err != nil {
		return nil, fmt.Errorf("Problem parsing response: %v\n", err)
	}

	return client, nil
}

func GetTweetByID(client *twittergo.Client, id string) (*twittergo.Tweet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/1.1/statuses/show.json?id=%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not parse request: %v", err)
	}

	resp, err := client.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("Could not send request: %v", err)
	}

	tweet := twittergo.Tweet{}
	err = resp.Parse(&tweet)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tweet: %s", err)
	}

	return &tweet, nil
}

func GetRetweetsByID(client *twittergo.Client, id string) ([]*twittergo.Tweet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/1.1/statuses/retweets/%s.json", id), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not parse request: %v", err)
	}

	resp, err := client.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("Could not send request: %v", err)
	}

	tweets := []*twittergo.Tweet{}
	err = resp.Parse(&tweets)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tweet: %s", err)
	}

	return tweets, nil
}

func PickAWinner(tweets []*twittergo.Tweet) *twittergo.Tweet {
	if len(tweets) == 0 {
		return nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	winner := rng.Intn(len(tweets))

	return tweets[winner]
}
