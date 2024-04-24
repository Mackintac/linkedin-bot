package handlers

import (
	"net/http"
)

func InitHandlers() {

	// LINKEDIN API PROFILE HANDLER FOR THE /ME EP
	userInfoHandler := UserInfoHandler()
	newQueryHandler, queryHolderChannel := NewQueryHandler()
	newShareHandler := NewShareHandler(queryHolderChannel)

	http.HandleFunc(projectConfig.Endpoints.Server.UserInfo, userInfoHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.NewQuery, newQueryHandler)
	http.HandleFunc(projectConfig.Endpoints.Server.NewShare, newShareHandler)

	go PostTimer(projectConfig.Endpoints.Server.NewQuery, "query")
	go PostTimer(projectConfig.Endpoints.Server.NewShare, "share")
}
