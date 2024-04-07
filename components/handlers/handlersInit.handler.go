package handlers

import "net/http"

func InitHandlers() {

	// LINKEDIN API PROFILE HANDLER FOR THE /ME EP
	userInfoHandler := UserInfoHandler()
	newShareHandler := NewShareHandler()

	http.HandleFunc(projectConfig.Endpoints.Server.NewShare, newShareHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.UserInfo, userInfoHandler)

}
