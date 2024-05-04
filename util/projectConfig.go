package projectUtil

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

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
		UserURN      string
		GPTSecret    string
	}

	GlobalVars struct {
		Ctx context.Context
	}
	LinkedInAuthCfg struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
		Scopes       []string
		Endpoint     oauth2.Endpoint
	}
	ChatGPTQueries struct {
		StyleOf struct {
			Manager         string
			Developer       [10]string
			SoftwareManager string
		}
		PostType struct {
			Guide    string
			Anecdote string
			Comment  string
		}
		Subject struct {
			Technologies [6]string
			GeneralTopic [3]string
		}
	}
}

// var projectConfig TProjectConfig

func InitProjectConfig() TProjectConfig {

	godotenv.Load(".env")
	projectConfig := TProjectConfig{
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
				"https://api.linkedin.com/v2/userinfo",
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
			UserURN      string
			GPTSecret    string
		}{
			os.Getenv("CLIENT_ID"),
			os.Getenv("PRIMARY_SECRET"),
			os.Getenv("ACCESS_TOKEN"),
			os.Getenv("USER_URN"),
			os.Getenv("GPT_SECRET"),
		},
		GlobalVars: struct{ Ctx context.Context }{
			context.Background(),
		},
		LinkedInAuthCfg: struct {
			ClientID     string
			ClientSecret string
			RedirectURL  string
			Scopes       []string
			Endpoint     oauth2.Endpoint
		}{
			os.Getenv("CLIENT_ID"),
			os.Getenv("PRIMARY_SECRET"),
			"http://localhost:8080/redirect",
			[]string{"email", "openid", "profile", "w_member_social"},
			linkedin.Endpoint,
		},
		ChatGPTQueries: struct {
			StyleOf struct {
				Manager         string
				Developer       [10]string
				SoftwareManager string
			}
			PostType struct {
				Guide    string
				Anecdote string
				Comment  string
			}
			Subject struct {
				Technologies [6]string
				GeneralTopic [3]string
			}
		}{
			StyleOf: struct {
				Manager         string
				Developer       [10]string
				SoftwareManager string
			}{"Manager", [10]string{
				"Frontend Developer",
				"Backend Engineer",
				"Full-stack Developer",
				"UI Developer",
				"API Developer",
				"JavaScript Engineer",
				"Database Developer",
				"DevOps Engineer",
				"Web Developer",
				"Software Engineer",
			}, "Software Engineering Manager"},
			PostType: struct {
				Guide    string
				Anecdote string
				Comment  string
			}{"write a guide for: ", "speak in the style of a personal anecdote or experience: ", "in a pensive manner, comment on: "},

			Subject: struct {
				Technologies [6]string
				GeneralTopic [3]string
			}{[6]string{
				"JavaScript",
				"Python",
				"Java",
				"Ruby",
				"Go (Golang)",
				"TypeScript",
			},
				[3]string{
					"Team Management",
					"Learning Programming",
					"Productivity in the office/remote workplace",
				},
			},
		},
	}
	// fmt.Println(projectConfig)
	// fmt.Println(projectConfig.DotEnvVars.GPTSecret)
	return projectConfig
}
