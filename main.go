package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"

	_ "github.com/joho/godotenv/autoload"
)

const (
	allSharesEP = "https://api.linkedin.com/v2/shares"
	meEP        = "https://api.linkedin.com/v2/me"
	redirectEP  = "http://localhost:8080/redirect"
	shareEP     = "https://api.linkedin.com/v2/ugcPosts"
)

const (
	newShareEP = "/newShare"
	newQueryEP = "/newQuery"
	userInfoEP = "/userInfo"
)

func main() {

	ctx := context.Background()
	if err := envInit(); err != nil {
		log.Fatal("Error loading .env")
		return
	}

	// accessToken := os.Getenv("ACCESS_TOKEN")
	cID := os.Getenv("CLIENT_ID")
	cSecret := os.Getenv("PRIMARY_SECRET")

	linkedInOauthConfig := &oauth2.Config{
		ClientID:     cID,
		ClientSecret: cSecret, //PRIMARY_SECRET
		RedirectURL:  redirectEP,
		Scopes:       []string{"openid", "profile", "w_member_social", "email"},
		Endpoint:     linkedin.Endpoint,
	}

	// verifier := oauth2.GenerateVerifier()
	url := linkedInOauthConfig.AuthCodeURL("ducksss", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := linkedInOauthConfig.Exchange(ctx, code)

	if err != nil {
		fmt.Printf("Token: %s: ", tok)
		log.Fatalf("Error exchanging authorization code for access token: %v", err)
	}
	fmt.Printf("\nToken: %s: ", tok)
	fmt.Println(json.MarshalIndent(tok, "", "    "))

	client := linkedInOauthConfig.Client(ctx, tok)
	client.Get("...")

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

func handlersInit() error {

	return nil
}

func getUserHandler() {
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
