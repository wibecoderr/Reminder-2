package handler

import (
	"net/http"

	utils "github.com/wibecoderr/Reminder-2.git"
	"github.com/wibecoderr/Reminder-2.git/service"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest

	err := utils.ParseBody(r.Body, &req)
	if err != nil {
		utils.RespondError(w, http.StatusNoContent, err, "Fail to create user")
		return
	}
	userService := service.NewUserService()
	resp, err := userService.CreateUser(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Fail to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, resp)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	err := utils.ParseBody(r.Body, &req)
	if err != nil {
		utils.RespondError(w, http.StatusNoContent, err, "Fail to create user")
		return
	}
	authService := service.NewAuthService()
	resp, err := authService.Login(req)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Login failed")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)

}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := utils.UserContext(r)
	userId := userCtx.UserId
	sessionID := userCtx.SessionId

	authService := service.NewAuthService()
	response, err := authService.Logout(sessionID, userId)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Logout failed")
		return
	}

	utils.RespondJSON(w, http.StatusNoContent, response)

}
