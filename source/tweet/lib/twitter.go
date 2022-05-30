package tweet

import (
	"context"
	"encoding/json"
	"fmt"
	twitter "github.com/g8rswimmer/go-twitter/v2"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

type ConfigAndToken struct {
	Config *oauth2.Config
	Token  *oauth2.Token
}

type retweet struct {
	twtclient      *twitter.Client
	configAndToken *ConfigAndToken
	projectID      string
}

// Add implements the go-twitter Authorization interface
func (r retweet) Add(req *http.Request) {
	config := r.configAndToken.Config
	tokenSource := config.TokenSource(req.Context(), r.configAndToken.Token)
	log.Println(tokenSource)
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatal(err)
	}

	if token.AccessToken != r.configAndToken.Token.AccessToken {
		r.configAndToken.Token = token
		twitterOAuthJson, err := json.Marshal(r.configAndToken)
		if err != nil {
			log.Fatalf("Couldn't serialize token: %v", err)
		}

		err = UpdateSecret(r.projectID, "twitter_oauth", twitterOAuthJson)
		if err != nil {
			log.Fatal(err)
		}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
}

func getConfigAndToken(projectID string) (*ConfigAndToken, error) {
	configAndTokenJson, err := AccessSecret(projectID, "twitter_oauth")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(configAndTokenJson))

	var configAndToken *ConfigAndToken
	err = json.Unmarshal(configAndTokenJson, &configAndToken)
	if err != nil {
		return nil, err
	}
	fmt.Println(configAndToken)
	return configAndToken, nil
}

func newClient(projectID string) *retweet {

	configAndToken, err := getConfigAndToken(projectID)
	if err != nil {
		log.Fatal(err)
	}
	retweet := &retweet{
		configAndToken: configAndToken,
		projectID:      projectID,
	}

	twtclient := &twitter.Client{
		Authorizer: retweet,
		Client:     &http.Client{Timeout: 2 * time.Second},
		Host:       "https://api.twitter.com",
	}

	retweet.twtclient = twtclient

	return retweet
}

func Retweet(projectID string) {

	client := newClient(projectID)
	query := "powerpoint"

	opts := twitter.TweetRecentSearchOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	}

	fmt.Println("Callout to tweet recent search callout")

	tweetResponse, err := client.twtclient.TweetRecentSearch(context.Background(), query, opts)
	if err != nil {
		log.Panicf("tweet lookup error: %v", err)
	}

	dictionaries := tweetResponse.Raw.TweetDictionaries()

	fmt.Println("Callout to create tweet callout")

	for k, d := range dictionaries {
		fmt.Println(k + " : " + d.Tweet.Text)

		req := twitter.CreateTweetRequest{
			Text: d.Tweet.Text,
		}

		tweetResponse, err := client.twtclient.CreateTweet(context.Background(), req)
		if err != nil {
			log.Panicf("create tweet error: %v", err)
		}

		fmt.Println(tweetResponse.Tweet.ID)

		break
	}
}
