package main

import (
	"bytes"
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

var linkedInOauthConfig = &oauth2.Config{
	ClientID:     DotEnvVars.ClientId,
	ClientSecret: DotEnvVars.ClientSecret, //PRIMARY_SECRET
	RedirectURL:  ServerEndpoints.Redirect,
	Scopes:       []string{"openid", "profile", "w_member_social", "email"},
	Endpoint:     linkedin.Endpoint,
}

var ctx = context.Background()

func main() {

	if err := envInit(); err != nil {
		log.Fatal("Error loading .env", err)
		return
	}

	fmt.Println(DotEnvVars.AccessToken)

	if err := handlersInit(); err != nil {
		log.Fatal("Error initializating Server Handlers", err)
		return
	}

	// verifier := oauth2.GenerateVerifier()

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

		token := &oauth2.Token{AccessToken: DotEnvVars.AccessToken}

		httpClient := linkedInOauthConfig.Client(ctx, token)

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
			log.Fatal("Error Marshalling JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", APIEndpoints.Share, bytes.NewBuffer(jsonShareReqBody))
		if err != nil {
			log.Fatal("Error Creating Request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// fmt.Printf("req: %+v\n", req)
		// fmt.Printf("ctx: %+v\n", ctx)
		// fmt.Printf("client: %+v\n", httpClient)
		// fmt.Printf("token: %+v\n", token)
		// fmt.Printf("body: %s\n", jsonShareReqBody)

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatal("Error sending request:", err)
			return
		}
		var responseBody map[string]interface{}

		// Decode JSON response body
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			log.Fatal("Error decoding response body:", err)
			return
		}

		// Print response body for inspection
		fmt.Println("Response Body:", responseBody)

		// Extract error details if available
		if errorMessage, ok := responseBody["errorMessage"].(string); ok {
			fmt.Println("Error Message:", errorMessage)
		}

		fmt.Println("Response Status:", resp.Body)
		defer resp.Body.Close()
		fmt.Println("Response Status:", resp.Status)

	}

	newerShareHandler := func(w http.ResponseWriter, r *http.Request) {
		accessToken := DotEnvVars.AccessToken

		// Use the access token to create an OAuth 2.0 token
		token := &oauth2.Token{AccessToken: accessToken}

		// Create an OAuth 2.0 HTTP client
		httpClient := linkedInOauthConfig.Client(context.Background(), token)

		// Use the HTTP client to make authenticated API requests
		url := "https://api.linkedin.com/v2/ugcPosts"
		postData := []byte(`{
		"author": "urn:li:person:4924372b1",
		"lifecycleState": "PUBLISHED",
		"specificContent": {
			"com.linkedin.ugc.ShareContent": {
				"shareCommentary": {
					"text": "Your share text here"
				}
			}
		},
		"visibility": {
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC"
		}
	}`)

		// Send POST request to LinkedIn API
		resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(postData))
		if err != nil {
			// Handle error by returning HTTP 500 Internal Server Error
			http.Error(w, "Failed to send request to LinkedIn API", http.StatusInternalServerError)
			log.Printf("Error sending request to LinkedIn API: %v", err)
			return
		}
		defer resp.Body.Close()

		// Check if the response status code is not successful (not 2xx)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// Handle error by returning HTTP status code received from LinkedIn API
			http.Error(w, fmt.Sprintf("LinkedIn API returned status code %d", resp.StatusCode), resp.StatusCode)
			log.Printf("LinkedIn API returned non-successful status code: %d", resp.StatusCode)
			return
		}

		// Decode response body
		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			// Handle error by returning HTTP 500 Internal Server Error
			http.Error(w, "Failed to decode response from LinkedIn API", http.StatusInternalServerError)
			log.Printf("Error decoding response body from LinkedIn API: %v", err)
			return
		}

		// Print response body for inspection
		fmt.Println("Response Body:", responseBody)

		// Write response to client
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}
	http.HandleFunc(ServerEndpoints.NewShare, newShareHandler)
	http.HandleFunc(ServerEndpoints.UserInfo, newerShareHandler)
	return nil
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
