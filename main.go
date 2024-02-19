package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

const (
	allSharesEP = "https://api.linkedin.com/v2/shares"
	meEP        = "https://api.linkedin.com/v2/me"
	userInfoEP  = "https://api.linkedin.com/v2/userinfo"
	redirectEP  = "https://www.linkedin.com/developers/tools/oauth/redirect"
	shareEP     = "https://api.linkedin.com/v2/ugcPosts"
)

var linkedInOauthConfig *oauth2.Config

func main() {

	ctx := context.Background()
	if err := envInit(); err != nil {
		log.Fatal("Error loading .env")
		return
	}

	// accessToken := os.Getenv("ACCESS_TOKEN")

	linkedInOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("PRIMARY_SECRET"), //PRIMARY_SECRET

		RedirectURL: redirectEP,
		Scopes:      []string{"openid", "profile", "w_member_social", "email"},
		Endpoint:    linkedin.Endpoint,
	}

	fmt.Println(linkedInOauthConfig.ClientID)
	http.HandleFunc("/", handleInitAuth)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	tok, err := linkedInOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := linkedInOauthConfig.Client(ctx, tok)
	_ = client

	getUser()

	fmt.Println("Server going live")

	if err := startServer(); err != nil {
		log.Fatal("Error starting server")
		return
	}

}

func envInit() error {
	// .env initialization
	godotenv.Load(".env")
	return nil
}

func handleInitAuth(w http.ResponseWriter, r *http.Request) {
	url := linkedInOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// func getAuthCode() {
// 	var code string
// 	if _, err := fmt.Scan(&code); err != nil {
// 		log.Fatal(err)
// 	}
// }

func getUser() {
	resp, err := http.Get(meEP)
	if err != nil {
		userErr := fmt.Errorf("error with response from getUser EP: %w", err)
		fmt.Println(userErr)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))
	if err != nil {
		rBodyErr := fmt.Errorf("error with response body of GET /me EP: %w", err)
		fmt.Println(rBodyErr)
	}

}

func startServer() error {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server dead")
		log.Fatal("Error with ListenAndServe", err)
	}
	return nil
}
