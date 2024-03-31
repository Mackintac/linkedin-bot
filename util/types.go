	cfg := Config{
		Endpoints: struct {
			LinkedIn struct {
				AllShares string
				Share     string
				Me        string
			}
			Server struct {
				NewShare string
				NewQuery string
				UserInfo string
				Redirect string
			}
		}{
			LinkedIn: struct {
				AllShares string
				Share     string
				Me        string
			}{
				AllShares: "https://api.linkedin.com/v2/shares",
				Share:     "https://api.linkedin.com/v2/ugcPosts",
				Me:        "https://api.linkedin.com/v2/me",
			},
			Server: struct {
				NewShare string
				NewQuery string
				UserInfo string
				Redirect string
			}{
				NewShare: "/newShare",
				NewQuery: "/newQuery",
				UserInfo: "/userInfo",
				Redirect: "http://localhost:8080/redirect",
			},
		},
		OAuth: struct {
			ClientID     string
			ClientSecret string
			RedirectURL  string
		}{
			ClientID:     "YOUR_CLIENT_ID",
			ClientSecret: "YOUR_CLIENT_SECRET",
			RedirectURL:  "YOUR_REDIRECT_URL",
		},
	}

	// Accessing configuration values
	fmt.Println("LinkedIn All Shares:", cfg.Endpoints.LinkedIn.AllShares)
	fmt.Println("Server New Share:", cfg.Endpoints.Server.NewShare)
	fmt.Println("OAuth Client ID:", cfg.OAuth.ClientID)
}