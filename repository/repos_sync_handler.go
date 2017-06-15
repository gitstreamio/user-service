package repository

import (
	"net/http"
	"github.com/gorilla/mux"
	"context"
	"fmt"
)

// TODO get token from request
const token string = "1aafed65ecb213cddb83f62d9d8c377f92107ef7"

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
	err := srv.repositoryService.SyncUserRepositories(login, token)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Synchronization of user repositories failed for user: %s", login))
		rw.WriteHeader(500)
	}
}
