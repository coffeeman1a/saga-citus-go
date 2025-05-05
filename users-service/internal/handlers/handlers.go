package handlers

import (
	"net/http"

	"github.com/coffeeman1a/saga-citus-go/users-service/internal/models"
	"github.com/coffeeman1a/saga-citus-go/users-service/internal/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func StartUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

func (h *UserHandler) NewUser(c *gin.Context) {
	var req *models.NewUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	id, err := h.userRepo.NewUser(c.Request.Context(), req.Email)
	if err != nil {
		log.WithError(err).Error("failde to create new user")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create a new user"})
		return
	}

	log.WithFields(log.Fields{
		"id":    id,
		"email": req.Email,
	}).Info("User created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	log.Debug("Parsing ID param...")
	idParam := c.Param("id")

	user, err := h.userRepo.GetUserByID(c.Request.Context(), idParam)
	if err != nil {
		log.WithError(err).Error("failed to get user from repository")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	if user == nil {
		log.WithField("user_id", idParam).Info("User not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	log.WithField("user_id", idParam).Info("User found and returned")
	c.JSON(http.StatusOK, user)
}
