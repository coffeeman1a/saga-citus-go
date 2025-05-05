package server

import (
	"github.com/coffeeman1a/saga-citus-go/users-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, usersHandler *handlers.UserHandler) {
	public := r.Group("/")
	{
		public.POST("/register", usersHandler.NewUser)
		public.GET("/users/:id", usersHandler.GetUserByID)
	}
}
