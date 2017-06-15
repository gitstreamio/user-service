package user

import (
	"net/http"
	"github.com/gorilla/mux"
	"context"
	"fmt"
)

// TODO get token from request
const token string = "1aafed65ecb213cddb83f62d9d8c377f92107ef7"

type UserSyncHandler struct {
	userService *UserService
}

func NewUserSyncHandler(ctx context.Context) *UserSyncHandler{
	userService := NewUserService(ctx)
	return &UserSyncHandler{userService}
}

func (srv *UserSyncHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]
	err := srv.userService.SyncUser(login, token)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("User sync failed for user: %s", login))
		rw.WriteHeader(500)
	}
}
