package handlers

import (
	"net/http"

	"github.com/coffeeman1a/saga-citus-go/orders-service/internal/models"
	"github.com/coffeeman1a/saga-citus-go/orders-service/internal/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type OrderHandler struct {
	orderRepo repository.OrderRepository
}

func StartOrderHandler(repo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: repo}
}

func (h *OrderHandler) NewOrder(c *gin.Context) {
	var req *models.NewOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	id, err := h.orderRepo.NewOrder(c.Request.Context(), req.UserID, req.Item, req.Price)
	if err != nil {
		log.WithError(err).Error("failed to create new order")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create a new order"})
		return
	}

	log.WithFields(log.Fields{
		"id":      id,
		"user_id": req.UserID,
		"item":    req.Item,
		"price":   req.Price,
	}).Info("order created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "order created"})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	log.Debug("parsing ID param...")
	idParam := c.Param("id")

	order, err := h.orderRepo.GetOrderByID(c.Request.Context(), idParam)
	if err != nil {
		log.WithError(err).Error("failed to get order from repository")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	if order == nil {
		log.WithField("order_id", idParam).Info("order not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	log.WithField("order_id", idParam).Info("order found and returned")
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req *models.StatusUpdateRquest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	log.Debug("parsing ID param...")
	idParam := c.Param("id")

	updated, err := h.orderRepo.UpdateOrderStatus(c.Request.Context(), idParam, req.Status)
	if err != nil {
		log.WithError(err).Error("failed to get order from repository")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	log.WithFields(log.Fields{
		"id":     idParam,
		"status": req.Status,
	}).Info("order status updated")
	c.JSON(http.StatusOK, updated)
}
