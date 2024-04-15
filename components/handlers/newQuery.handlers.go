package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
)

func NewQueryHandler() func(w http.ResponseWriter, r *http.Request) {
	newQueryHandler := func(w http.ResponseWriter, r *http.Request) {
		queryReqBody := map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": map[string]interface{}{
				"role": "user", "content": "Say this is a test!"},
			"temperature": 0.7,
		}

		jsonQueryReqBody, err := json.Marshal(queryReqBody)
		if err != nil {
			log.Fatal("Error Marshalling JSON:", err)
			return
		}
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonQueryReqBody))
		if err != nil {
			log.Fatal("Error creating request with GPT query:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer"+projectConfig.DotEnvVars.GPTSecret)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Error sending GPT query request:", err)
			return
		}

		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			log.Fatal("Error decoding GPT reponse Body:", err)
			return
		}

		if errorMessage, ok := responseBody["errorMessage"].(string); ok {
			fmt.Println("Error Message:", errorMessage)
		}

		defer resp.Body.Close()

	}
	return newQueryHandler
}
