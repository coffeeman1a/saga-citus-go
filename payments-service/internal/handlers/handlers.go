package handlers

import (
	"net/http"

	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/config"
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/models"
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PaymentHandler struct {
	paymentRepo repository.PaymentRepository
}

func StartPaymentHandler(repo repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{paymentRepo: repo}
}

func (h *PaymentHandler) NewPayment(c *gin.Context) {
	var req *models.NewPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	id, err := h.paymentRepo.NewPayment(c.Request.Context(), req.OrderID, req.Amount)
	if err != nil {
		log.WithError(err).Error("failed to create new payment")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create a new payment"})
		return
	}

	log.WithFields(log.Fields{
		"id":       id,
		"order_id": req.OrderID,
		"status":   config.StatusReserved,
	}).Info("payment created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "payment created"})
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	log.Debug("parsing ID param...")
	idParam := c.Param("id")

	order, err := h.paymentRepo.GetPaymentByID(c.Request.Context(), idParam)
	if err != nil {
		log.WithError(err).Error("failed to get payment from repository")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	if order == nil {
		log.WithField("payment_id", idParam).Info("payment not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "payment not found"})
		return
	}

	log.WithField("payment", idParam).Info("payment found and returned")
	c.JSON(http.StatusOK, order)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	var req *models.StatusUpdateRquest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	log.Debug("parsing ID param...")
	idParam := c.Param("id")

	updated, err := h.paymentRepo.UpdatePaymentStatus(c.Request.Context(), idParam, req.Status)
	if err != nil {
		log.WithError(err).Error("failed to get payment from repository")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	log.WithFields(log.Fields{
		"id":     idParam,
		"status": req.Status,
	}).Info("payment status updated")
	c.JSON(http.StatusOK, updated)
}
