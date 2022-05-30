package tweet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	context        context.Context
	s3Client       *s3.Client
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

		_, err = PutObject(r.context, r.s3Client, "my-bucket-included-giraffe", "twitter_secrets", twitterOAuthJson)
		if err != nil {
			log.Fatal(err)
		}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
}

func getConfigAndToken(ctx context.Context, s3Client *s3.Client) (*ConfigAndToken, error) {
	configAndTokenJson, err := GetObject(ctx, s3Client, "my-bucket-included-giraffe", "twitter_secrets")
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

func newClient(ctx context.Context, cfg aws.Config) *retweet {
	s3Client := s3.NewFromConfig(cfg)

	configAndToken, err := getConfigAndToken(ctx, s3Client)
	if err != nil {
		log.Fatal(err)
	}
	retweet := &retweet{
		configAndToken: configAndToken,
		s3Client:       s3Client,
		context:        ctx,
	}

	twtclient := &twitter.Client{
		Authorizer: retweet,
		Client:     &http.Client{Timeout: 2 * time.Second},
		Host:       "https://api.twitter.com",
	}

	retweet.twtclient = twtclient

	return retweet
}

func Retweet(ctx context.Context, cfg aws.Config) {

	client := newClient(ctx, cfg)
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
