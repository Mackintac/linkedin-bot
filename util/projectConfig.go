package projectUtil

import "os"

type TProjectConfig struct {
	Endpoints struct {
		LinkedIn struct {
			AllShares string
			Share     string
			UserInfo  string
		}
		Server struct {
			NewShare string
			NewQuery string
			UserInfo string
			Redirect string
		}
	}
	DotEnvVars struct {
		ClientId     string
		ClientSecret string
		AccessToken  string
	}
}

var projectConfig TProjectConfig

func InitProjectConfig() TProjectConfig {
	projectConfig = TProjectConfig{
		Endpoints: struct {
			LinkedIn struct {
				AllShares string
				Share     string
				UserInfo  string
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
				UserInfo  string
			}{
				"https://api.linkedin.com/v2/shares",
				"https://api.linkedin.com/v2/ugcPosts",
				"https://api.linkedin.com/v2/me",
			},
			Server: struct {
				NewShare string
				NewQuery string
				UserInfo string
				Redirect string
			}{
				"/newShare",
				"/newQuery",
				"/userInfo",
				"http://localhost:8080/redirect",
			},
		},
		DotEnvVars: struct {
			ClientId     string
			ClientSecret string
			AccessToken  string
		}{
			os.Getenv("CLIENT_ID"),
			os.Getenv("PRIMARY_SECRET"),
			os.Getenv("ACCESS_TOKEN"),
		},
	}

	return projectConfig
}
