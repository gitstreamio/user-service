package server

import (
	"context"
	"github.com/gorilla/mux"
	"user-service/user"
	"net/http"
	"user-service/repository"
)


func Run() {
	ctx := context.Background()

	router := mux.NewRouter()

	userHandler := user.NewUserHandler(ctx)
	router.Handle("/users/{login}", userHandler)

	userSyncHandler := user.NewUserSyncHandler(ctx)
	router.Handle("/users-sync/{login}", userSyncHandler)

	repoHandler := repository.NewRepositoryHandler(ctx)
	router.Handle("/repos/{login}", repoHandler)

	repoSyncHandler := repository.NewRepositorySyncHandler(ctx)
	router.Handle("/repos-sync/{login}", repoSyncHandler)

	http.ListenAndServe(":2022", router)
}

