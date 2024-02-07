package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

const (
	shareEP    = "https://api.linkedin.com/v2/shares"
	meEP       = "https://api.linkedin.com/v2/me"
	userinfoEP = "https://api.linkedin.com/v2/userinfo"
)

var (
	linkedInConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("PRIMARY_SECRET"), //PRIMARY_SECRET
		RedirectURL:  "http://localhost:8080/login",
		Scopes:       []string{"openid", "profile", "w_member_social", "email"},
		Endpoint:     linkedin.Endpoint,
	}
)

func main() {
	if err := envInit(); err != nil {
		log.Fatal("Error loading .env")
		return
	}

	if err := helloHandler(); err != nil {
		log.Fatal("Error loading Hello handler")
		return
	}
	resp, err := http.Get(meEP)
	if err != nil {
		fmt.Errorf("Error with GET request to /me EP")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

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

func helloHandler() error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
		fmt.Println("Server Pinged!")
	}
	http.HandleFunc("/", handler)
	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) error {
	if err := fmt.Fprint(w, "Hello, World!"); err != nil {
		log.Fatal("Error starting server")
	}
	fmt.Println("Server Pinged!")
	return nil
}

func startServer() error {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server dead")
		log.Fatal("Error with ListenAndServe")
	}
	return nil
}
