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
						"text": `Hello LinkedIn fam! 

I wanted to share a little personal anecdote about my journey as a Full-stack Developer working with Swift. When I first started learning Swift, I was blown away by its simplicity and power. The way it seamlessly integrates with Apple's ecosystem and allows for rapid development of iOS and macOS applications is truly impressive.

As I delved deeper into Swift, I realized just how versatile it is. From building sleek user interfaces with SwiftUI to optimizing performance with Combine, Swift has allowed me to create robust and user-friendly applications that exceed expectations.

One of the things I love most about Swift is the strong community surrounding it. Whether I'm attending meetups, participating in online forums, or collaborating with fellow developers on open-source projects, I am constantly inspired by the passion and innovation that this community embodies.

As a Full-stack Developer, mastering Swift has opened up a world of possibilities for me. It has not only enhanced my technical skills but also broadened my career opportunities. I am grateful for the journey so far and excited to see where it will take me next.

I would love to hear about your experiences with Swift and how it has impacted your career as a developer. Feel free to share your thoughts in the comments below!

#FullStackDeveloper #Swift #iOSDevelopment #TechCommunity #CareerGrowth #LinkedInNetworking 
`,
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
