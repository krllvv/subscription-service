package main

import (
	"log"
	"net/http"
	"os"
	"subscription-service/config"
	_ "subscription-service/docs"
	"subscription-service/internal/handler"
	"subscription-service/internal/repository/sub/postgres"
	"subscription-service/internal/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title		Subscription Service API
// @version		1.0
// @description	REST service for aggregating data about users' online subscriptions
// @BasePath	/
func main() {
	logger := log.New(os.Stdout, "[SubService] ", log.LstdFlags)
	cfg := config.InitConfig(logger)
	addr := ":" + cfg.ServerPort

	repo, err := postgres.NewSubPostgresRepository(cfg, logger)
	if err != nil {
		logger.Fatalf("Database init error: %v", err)
	}
	srv := service.NewSubService(repo)
	h := handler.NewSubHandler(srv, logger)

	r := mux.NewRouter()
	h.RegisterRoutes(r)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Println("Server starting at " + addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Println("http server error:", err)
	}
}
