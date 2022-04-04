package main

import (
	"encoding/json"
	"net/http"

	"github.com/VanjaRo/web-chat/auth"
	"github.com/VanjaRo/web-chat/repositories"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Api struct {
	UserRepository *repositories.UserRepository
}

func (api *Api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginUser LoginUser
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		returnErrorResponse(w)
		return
	}
	dbUser, err := api.UserRepository.FindUserByUsername(loginUser.Username)
	if err != nil {
		returnErrorResponse(w)
		return
	}

	if ok, err := auth.ComparePassword(dbUser.Password, loginUser.Password); !ok || err != nil {
		returnErrorResponse(w)
		return
	}

	token, err := auth.GenerateJWTToken(dbUser)
	if err != nil {
		returnErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func returnErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"status\": \"error\"}"))
}
