package handlers

import "net/http"

func InitHandlers() {

	// LINKEDIN API PROFILE HANDLER FOR THE /ME EP
	userInfoHandler := UserInfoHandler()
	newQueryHandler, queryHolderChannel := NewQueryHandler()
	newShareHandler := NewShareHandler(queryHolderChannel)

	http.HandleFunc(projectConfig.Endpoints.Server.NewShare, newShareHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.UserInfo, userInfoHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.NewQuery, newQueryHandler)

}
