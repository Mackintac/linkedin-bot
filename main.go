package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"

	_ "github.com/joho/godotenv/autoload"
)

type linkedInEndpoints struct {
	AllShares string
	Share     string
	Me        string
}

type serverEndpoints struct {
	NewShare string
	NewQuery string
	UserInfo string
	Redirect string
}

type dotEnvVars struct {
	ClientId     string
	ClientSecret string
	AccessToken  string
}

var APIEndpoints = linkedInEndpoints{
	AllShares: "https://api.linkedin.com/v2/shares",
	Share:     "https://api.linkedin.com/v2/ugcPosts",
	Me:        "https://api.linkedin.com/v2/me",
}

var ServerEndpoints = serverEndpoints{
	NewShare: "/newShare",
	NewQuery: "/newQuery",
	UserInfo: "/userInfo",
	Redirect: "http://localhost:8080/redirect",
}

var DotEnvVars = dotEnvVars{
	ClientId:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("PRIMARY_SECRET"),
	AccessToken:  os.Getenv("ACCESS_TOKEN"),
}

func main() {

	ctx := context.Background()
	if err := envInit(); err != nil {
		log.Fatal("Error loading .env", err)
		return
	}

	if err := handlersInit(); err != nil {
		log.Fatal("Error initializating Server Handlers", err)
		return
	}

	// accessToken := os.Getenv("ACCESS_TOKEN")

	linkedInOauthConfig := &oauth2.Config{
		ClientID:     DotEnvVars.ClientId,
		ClientSecret: DotEnvVars.ClientSecret, //PRIMARY_SECRET
		RedirectURL:  ServerEndpoints.Redirect,
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

	newShareHandler := func(w http.ResponseWriter, r *http.Request) {

		shareReqBody := map[string]interface{}{
			"author":         "urn:li:person:4924372b1",
			"lifecycleState": "PUBLISHED",
			"specificContent": map[string]interface{}{
				"com.linkedin.ugc.ShareContent": map[string]interface{}{
					"shareCommentary": map[string]interface{}{
						"text": "Spending this Thursday learning more about backend systems using GoLang!",
					},
					"shareMediaCategory": "NONE",
				},
			},
			"visibility": map[string]interface{}{
				"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
			},
		}
		jsonShareReqBody, err := json.Marshal(shareReqBody)
		if err != nil {
			fmt.Println("Error Marshalling JSON:", err)
			return
		}

		r.Header.Set("Content-Type", "application/json")

	}
	http.HandleFunc(ServerEndpoints.NewShare, newShareHandler)
	return nil
}

// func getUserHandler() {
// 	resp, err := http.Get(APIEndpoints.Me)
// 	if err != nil {
// 		userErr := fmt.Errorf("error with response from getUser EP: %w", err)
// 		fmt.Println(userErr)
// 	}

// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)

// 	fmt.Println(string(body))
// 	if err != nil {
// 		rBodyErr := fmt.Errorf("error with response body of GET /me EP: %w", err)
// 		fmt.Println(rBodyErr)
// 	}

// }

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
