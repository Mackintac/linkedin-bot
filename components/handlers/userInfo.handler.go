package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func UserInfoHandler() func(w http.ResponseWriter, r *http.Request) {

	userInfoHandler := func(w http.ResponseWriter, r *http.Request) {

		token := &oauth2.Token{AccessToken: projectConfig.DotEnvVars.AccessToken}
		httpClient := linkedInOauthConfig.Client(projectConfig.GlobalVars.Ctx, token)

		resp, err := httpClient.Get(projectConfig.Endpoints.LinkedIn.UserInfo)
		if err != nil {

			log.Fatal("Error with Get Request:", err)
			return
		}
		defer resp.Body.Close()
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
		fmt.Println("Response Status:", resp.Status)

	}
	return userInfoHandler
}
