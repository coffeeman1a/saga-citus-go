package server

import (
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, paymentHandler *handlers.PaymentHandler) {
	public := r.Group("/")
	{
		public.POST("/payments/", paymentHandler.NewPayment)
		public.GET("/payments/:id", paymentHandler.GetPaymentByID)
		public.PUT("/status/:id", paymentHandler.UpdatePaymentStatus)
	}
}
