package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"os/exec"

	tweet "example.com/tweet/lib"
	pkcecv "github.com/nirasan/go-oauth-pkce-code-verifier"
)

func prepareCommandLineArgument(name string, argument []byte) string {
	trimmed := strings.Trim(string(argument), " \n")
	return fmt.Sprintf("%s=%s", name, trimmed)
}

func getOAuthConfig(clientId string, clientSecret string, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"users.read", "tweet.read", "tweet.write", "follows.read", "follows.write", "offline.access"},
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
			TokenURL: "https://api.twitter.com/2/oauth2/token",
		},
	}
}

func main() {

	err := os.Chdir(filepath.Join("..", "..", "terraform"))
	if err != nil {
		log.Fatalf("Couldn't change directory: %v", err)
	}
	fmt.Println("Checking if there's already a Twitter OAuth Token saved in the secret store")
	twitterOAuthJson, err := exec.Command(
		"aws",
		"s3api",
		"get-object",
		"--bucket",
		"my-bucket-included-giraffe",
		"--key",
		"twitter_secrets",
		"outfile",
	).Output()
	if err == nil {
		twitterOAuthJson, err = exec.Command(
			"cat",
			"outfile",
		).Output()
	} else {
		var clientId string
		var clientSecret string
		var redirectURL string
		fmt.Printf("Enter your Twitter Client Id: ")
		if _, err = fmt.Scan(&clientId); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Enter your Twitter Client Secret: ")
		if _, err = fmt.Scan(&clientSecret); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Enter your Twitter Redirect URL: ")
		if _, err = fmt.Scan(&redirectURL); err != nil {
			log.Fatal(err)
		}

		authConf := getOAuthConfig(clientId, clientSecret, redirectURL)
		var codeVerifier, err = pkcecv.CreateCodeVerifier()
		if err != nil {
			log.Fatalf("Couldn't create code verifier: %v", err)
		}

		authCodeURL := authConf.AuthCodeURL(
			"state",
			oauth2.AccessTypeOffline,
			oauth2.SetAuthURLParam("client_id", authConf.ClientID),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("code_challenge", codeVerifier.CodeChallengeS256()),
		)

		fmt.Println("Visit the URL for the auth dialog: ", authCodeURL)

		// 2. redirect URL contains code
		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		var code string
		fmt.Printf("Enter the code query parameter of the redirected URL: ")
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatalf("Couldn't scan code: %v", err)
		}

		// Use the custom HTTP client when requesting a token.
		ctx := context.Background()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, http.Client{Timeout: 2 * time.Second})

		// 3. Exchange code for token
		twitterToken, err := authConf.Exchange(
			ctx,
			code,
			oauth2.SetAuthURLParam("code_verifier", codeVerifier.Value),
		)
		if err != nil {
			log.Fatalf("Couldn't exchange code for token: %v", err)
		}

		twitterOAuth := &tweet.ConfigAndToken{
			Config: authConf,
			Token:  twitterToken,
		}
		twitterOAuthJson, err = json.Marshal(twitterOAuth)
		if err != nil {
			log.Fatalf("Couldn't serialize token: %v", err)
		}
	}

	err = exec.Command("terraform", "init").Run()
	if err != nil {
		log.Fatalf("Couldn't initiate terraform: %v", err)
	}

	twitterSecrets := prepareCommandLineArgument("twitter_secrets", twitterOAuthJson)
	fmt.Println("terraform", "apply",
		"--var", twitterSecrets,
	)
	cmd := exec.Command("terraform", "apply",
		"--var", twitterSecrets,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		log.Fatalf("Couldn't apply terraform: %v", err)
	}
}
