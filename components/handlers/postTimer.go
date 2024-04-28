package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func PostTimer(ep string, epName string) {
	ticker := time.NewTicker(17 * time.Hour)

	for range ticker.C {
		fmt.Println("sending an auto-ping to endpoint: " + epName)
		resp, err := http.Get("http://localhost:8080" + ep)
		if err != nil {
			fmt.Println("Error with auto-ping to endpoint: "+epName, err)
			continue
		}
		fmt.Println("Response status:", resp.Status)
	}
}
