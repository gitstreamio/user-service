package user

import (
	"net/http"
	"github.com/gorilla/mux"
	"context"
	"encoding/json"
)


type UserHandler struct {
	userService *UserService
}

func NewUserHandler(ctx context.Context) *UserHandler{
	userService := NewUserService(ctx)
	return &UserHandler{userService}
}

func (srv *UserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]
	user, err := srv.userService.FindByLogin(login)
	if err != nil {
		log.WithError(err).Error("User search failed for user: " + login)
		rw.WriteHeader(500)
		return
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.WithError(err).Error("Could not marshal user data for user: " + login)
		rw.WriteHeader(500)
		return
	}
	rw.Write(jsonUser)
}
