package user

import (
	"net/http"
	"github.com/gorilla/mux"
	"context"
	"fmt"
)


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
	token := ""   // TODO get token from request
	err := srv.userService.SyncUser(login, token)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("User sync failed for user: %s", login))
		rw.WriteHeader(500)
	}
}
