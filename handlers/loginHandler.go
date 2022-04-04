package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VanjaRo/web-chat/auth"
)

func (api *Api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginUser AuthUser
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		returnErrorResponse(w, "invalid json")
		return
	}

	dbUser, err := api.UserRepository.FindUserByName(loginUser.Username)
	if err != nil {
		returnErrorResponse(w, "user not found")
		return
	}

	if ok, err := auth.ComparePassword(loginUser.Password, dbUser.Password); !ok || err != nil {
		returnErrorResponse(w, "invalid password")
		return
	}

	token, err := auth.GenerateJWTToken(dbUser)
	if err != nil {
		returnErrorResponse(w, "can't generate JWT token")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"token\": \"" + token + "\"}"))

}

func returnErrorResponse(w http.ResponseWriter, errMsg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"error\": \"" + errMsg + "\"}"))
}
