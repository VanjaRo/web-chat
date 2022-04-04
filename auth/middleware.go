package auth

import (
	"context"
	"net/http"
)

const UserContextKey = "user"

type AnonUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (u *AnonUser) GetId() string {
	return u.Id
}

func (u *AnonUser) GetName() string {
	return u.Name
}

func AuthMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, tokenOk := r.URL.Query()["bearer"]
		if tokenOk && len(token) == 1 {
			user, err := ValidateJWTToken(token[0])
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				handlerFunc(w, r.WithContext(ctx))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("plese login to access the chat"))
		}

	})
}
