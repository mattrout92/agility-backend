package main

import (
	"context"
	"net/http"

	"github.com/mattrout92/agility-backend/cors"
	"github.com/mattrout92/agility-backend/handlers"

	"github.com/mattrout92/agility-backend/store"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/mattrout92/agility-backend/logger"
)

func main() {
	r := mux.NewRouter()

	storer := store.Connect(&store.Config{Ctx: context.Background()})

	svc := handlers.Service{Store: storer}

	r.Path("/login").Methods("POST", "OPTIONS").Handler(cors.Middleware(logger.Handler(http.HandlerFunc(svc.Login))))
	r.Path("/dogs").Methods("POST", "OPTIONS").Handler(cors.Middleware(logger.Handler(http.HandlerFunc(svc.AddDog))))
	r.Path("/dogs").Methods("GET", "OPTIONS").Handler(cors.Middleware(logger.Handler(http.HandlerFunc(svc.GetDogs))))
	r.Path("/dogs").Methods("DELETE", "OPTIONS").Handler(cors.Middleware(logger.Handler(http.HandlerFunc(svc.DeleteDog))))
	r.Path("/dogs").Methods("PUT", "OPTIONS").Handler(cors.Middleware(logger.Handler(http.HandlerFunc(svc.UpdateDog))))

	port := ":8080"

	logger.Trace("starting server", logger.Data{"port": port})

	if err := http.ListenAndServe(port, r); err != nil {
		logger.Error(err)
	}
}
