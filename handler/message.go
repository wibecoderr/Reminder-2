package handler

import (
	"net/http"

	utils "github.com/wibecoderr/Reminder-2.git"
	"github.com/wibecoderr/Reminder-2.git/service"
)

// create , delete , list , show
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	userCTX := utils.UserContext(r)
	userID := userCTX.UserId

	// parse
	var req service.CreateReminderRequest
	err := utils.ParseBody(r.Body, &req)
	if err != nil {
		utils.RespondError(w, http.StatusConflict, nil, "Fail to parse body")
		return
	}
	resp, err :=
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Fail to create reminder")
		return
	}

}
