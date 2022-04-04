package handlers

import (
	"github.com/VanjaRo/web-chat/repositories"
	"github.com/google/uuid"
)

type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Api struct {
	UserRepository *repositories.UserRepository
}

// generate a unique id for the user
func (api *Api) GenerateId() string {
	// check if the id is already in use
	for {
		id := uuid.New().String()
		if _, err := api.UserRepository.FindUserById(id); err != nil {
			return id
		}
	}
}
