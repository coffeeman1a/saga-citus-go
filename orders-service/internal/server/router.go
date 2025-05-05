package server

import (
	"github.com/coffeeman1a/saga-citus-go/orders-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, orderHandler *handlers.OrderHandler) {
	public := r.Group("/")
	{
		public.POST("/orders/", orderHandler.NewOrder)
		public.GET("/orders/:id", orderHandler.GetOrderByID)
		public.PUT("/status/:id", orderHandler.UpdateOrderStatus)
	}
}
