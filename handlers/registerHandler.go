package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VanjaRo/web-chat/auth"
	"github.com/VanjaRo/web-chat/repositories"
)

func (api *Api) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registrationUser AuthUser
	// read the inputs from the frontend
	if err := json.NewDecoder(r.Body).Decode(&registrationUser); err != nil {
		returnErrorResponse(w, "invalid json")
		return
	}
	// check if the user already exists
	_, err := api.UserRepository.FindUserByName(registrationUser.Username)
	if err == nil {
		returnErrorResponse(w, "user already exists")
		return
	}
	// hash the password
	hashedPassword, err := auth.GeneratePassword(registrationUser.Password)
	if err != nil {
		returnErrorResponse(w, "could not hash password")
		return
	}
	// I generate the ID myself to avoid collisions
	// during private rooms creation
	newId := api.GenerateId()
	// create the user
	newUser := &repositories.User{
		ID:       newId,
		Name:     registrationUser.Username,
		Password: hashedPassword,
	}
	if err := api.UserRepository.AddUser(newUser); err != nil {
		returnErrorResponse(w, "could not create user")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
