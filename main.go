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
	userInfoEP  = "https://api.linkedin.com/v2/userinfo"
	redirectEP  = "http://localhost:8080/redirect"
	shareEP     = "https://api.linkedin.com/v2/ugcPosts"
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
	fmt.Println("0")
	url := linkedInOauthConfig.AuthCodeURL("ducksss", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	fmt.Println("1")
	fmt.Println("2")

	fmt.Println("3")

	fmt.Println("4")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	fmt.Println("4")
	fmt.Printf("Code: %s, Verifier: %s\n", code)
	tok, err := linkedInOauthConfig.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("Token: %s: ", tok)
		log.Fatalf("Error exchanging authorization code for access token: %v", err)
	}
	fmt.Printf("\nToken: %s: ", tok)
	fmt.Println(json.MarshalIndent(tok, "", "    "))

	fmt.Println("6")
	client := linkedInOauthConfig.Client(ctx, tok)
	client.Get("...")

	fmt.Println("Server going live")

	if err := startServer(); err != nil {
		log.Fatal("Error starting server")
		return
	}

	go getUserHandler()
}

func envInit() error {
	// .env initialization
	godotenv.Load(".env")
	return nil
}

// func handleInitAuth(w http.ResponseWriter, r *http.Request) {
// 	url := linkedInOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func getAuthCode() {
// 	var code string
// 	if _, err := fmt.Scan(&code); err != nil {
// 		log.Fatal(err)
// 	}
// }

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
