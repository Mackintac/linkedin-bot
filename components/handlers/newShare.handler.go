package handlers

import (
	"bytes"
	projectUtil "dev/linkedIn/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var projectConfig projectUtil.TProjectConfig = projectUtil.InitProjectConfig()
var linkedInOauthConfig = &oauth2.Config{
	ClientID:     projectConfig.LinkedInAuthCfg.ClientID,
	ClientSecret: projectConfig.LinkedInAuthCfg.ClientSecret,
	RedirectURL:  projectConfig.LinkedInAuthCfg.RedirectURL,
	Scopes:       projectConfig.LinkedInAuthCfg.Scopes,
	Endpoint:     projectConfig.LinkedInAuthCfg.Endpoint,
}

func NewShareHandler() func(w http.ResponseWriter, r *http.Request) {

	newShareHandler := func(w http.ResponseWriter, r *http.Request) {

		token := &oauth2.Token{AccessToken: projectConfig.DotEnvVars.AccessToken}
		httpClient := linkedInOauthConfig.Client(projectConfig.GlobalVars.Ctx, token)
		shareReqBody := map[string]interface{}{
			"author":         "urn:li:person:" + projectConfig.DotEnvVars.UserURN,
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

		req, err := http.NewRequest("POST", projectConfig.Endpoints.LinkedIn.Share, bytes.NewBuffer(jsonShareReqBody))
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
	return newShareHandler
}
