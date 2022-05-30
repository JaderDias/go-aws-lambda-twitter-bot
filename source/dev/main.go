package main

import (
	tweet "example.com/tweet/lib"
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	ctx := context.TODO()
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tweet.Retweet(ctx, cfg)
}
