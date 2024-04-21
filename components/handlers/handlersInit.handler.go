package handlers

import "net/http"

func InitHandlers(customQuery string) {

	// LINKEDIN API PROFILE HANDLER FOR THE /ME EP
	userInfoHandler := UserInfoHandler()
	newShareHandler := NewShareHandler()
	newQueryHandler := NewQueryHandler(customQuery)

	http.HandleFunc(projectConfig.Endpoints.Server.NewShare, newShareHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.UserInfo, userInfoHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.NewQuery, newQueryHandler)

}
