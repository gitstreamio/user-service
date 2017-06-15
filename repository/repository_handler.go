package repository

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"context"
)


type RepositoryHandler struct {
	repositoryService *RepositoryService
}

func NewRepositoryHandler(ctx context.Context) *RepositoryHandler{
	repositoryService := NewRepositoryService(ctx)

	return &RepositoryHandler{repositoryService}
}

func (srv *RepositoryHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]
	repos, err := srv.repositoryService.FindByLogin(login)
	if err != nil {
		log.WithError(err).Error("Could not find repos for login: " + login)
		rw.WriteHeader(404)
		return
	}
	jsonRepos, err := json.Marshal(repos)
	if err != nil {
		log.WithError(err).Error("Could not marshal repos for login: " + login)
		rw.WriteHeader(500)
		return
	}
	rw.Write(jsonRepos)
}
