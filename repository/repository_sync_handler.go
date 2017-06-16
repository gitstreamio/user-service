package repository

import (
	"net/http"
	"github.com/gorilla/mux"
	"context"
	"fmt"
	"strings"
)


type RepositorySyncHandler struct {
	repositoryService *RepositoryService
}

func NewRepositorySyncHandler(ctx context.Context) *RepositorySyncHandler{
	repositoryService := NewRepositoryService(ctx)

	return &RepositorySyncHandler{repositoryService}
}

func (srv *RepositorySyncHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]
	gitHubApiToken := r.Header.Get("Authorization")
	err := srv.repositoryService.SyncUserRepositories(login, strings.TrimPrefix(gitHubApiToken, "Bearer "))
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Synchronization of user repositories failed for user: %s", login))
		rw.WriteHeader(500)
	}
}
