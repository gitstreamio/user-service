package user

import (
	"net/http"
	"github.com/gorilla/mux"
	"context"
	"fmt"
	"strings"
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
	gitHubApiToken := r.Header.Get("Authorization")
	err := srv.userService.SyncUser(login, strings.TrimPrefix(gitHubApiToken, "Bearer "))
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("User sync failed for user: %s", login))
		rw.WriteHeader(500)
	}
}
