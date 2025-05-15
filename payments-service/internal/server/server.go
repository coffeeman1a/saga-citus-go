package server

import (
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/config"
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/handlers"
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router *gin.Engine
}

func StartServer(repo *repository.PGXRepository) *Server {
	log.Debug("Creating router...")
	r := gin.Default()
	orderHandler := handlers.StartPaymentHandler(repo)
	log.Debug("Setting up routes...")
	setupRoutes(r, orderHandler)
	return &Server{router: r}
}

func (s *Server) Run() error {
	log.Printf("Starting server on %v", config.Port)
	return s.router.Run(":" + config.Port)
}
