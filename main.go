package main

import (
	"bytes"
	"context"
	"dev/linkedIn/components/handlers"
	projectUtil "dev/linkedIn/util"
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
	UserInfo  string
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
	UserInfo:  "https://api.linkedin.com/v2/userinfo",
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
	Scopes:       []string{"email", "openid", "profile", "w_member_social"},
	Endpoint:     linkedin.Endpoint,
}

var ctx = context.Background()

func main() {
	handlers.NewShareHandler()
	projectConfig := projectUtil.InitProjectConfig()
	fmt.Println(projectConfig.Endpoints.LinkedIn.AllShares)
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

	// LINKEDIN API newShare handler
	newShareHandler := func(w http.ResponseWriter, r *http.Request) {

		token := &oauth2.Token{AccessToken: DotEnvVars.AccessToken}
		httpClient := linkedInOauthConfig.Client(ctx, token)
		shareReqBody := map[string]interface{}{
			"author":         "urn:li:person:",
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

		fmt.Printf("req: %+v\n", req)
		req.Header.Set("Authorization", "Bearer"+token.AccessToken)
		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
		req.Header.Set("Content-Type", "application/json")

		// log.Printf("req: %+v\n", req)
		// fmt.Printf("body: %s\n", jsonShareReqBody)
		// fmt.Printf("-------------BODY------------ %s \n", bytes.NewBuffer(jsonShareReqBody))

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatal("Error sending request:", err)
			return
		}

		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			log.Fatal("Error decoding response body:", err)
			return
		}

		if errorMessage, ok := responseBody["errorMessage"].(string); ok {
			fmt.Println("Error Message:", errorMessage)
		}

		defer resp.Body.Close()

	}

	// LINKEDIN API PROFILE HANDLER FOR THE /ME EP
	getInfoHandler := func(w http.ResponseWriter, r *http.Request) {

		token := &oauth2.Token{AccessToken: DotEnvVars.AccessToken}
		httpClient := linkedInOauthConfig.Client(ctx, token)

		resp, err := httpClient.Get(APIEndpoints.UserInfo)
		if err != nil {

			log.Fatal("Error with Get Request:", err)
			return
		}
		var responseBody map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			log.Fatal("Error decoding response body:", err)
			return
		}

		fmt.Println("Response Body:", responseBody)

		if errorMessage, ok := responseBody["errorMessage"].(string); ok {
			fmt.Println("Error Message:", errorMessage)
		}

		fmt.Println("Response Status:", resp.Body)
		defer resp.Body.Close()
		fmt.Println("Response Status:", resp.Status)

	}

	http.HandleFunc(ServerEndpoints.NewShare, newShareHandler)
	http.HandleFunc(ServerEndpoints.UserInfo, getInfoHandler)
	return nil
}

// STARTS THE SERVER
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
