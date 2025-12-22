package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"effective_mobile_service/internal/entity"
	"effective_mobile_service/internal/usecase"
	"effective_mobile_service/pkg/logger"
)

type SubscriptionHandler struct {
	useCase usecase.SubscriptionUseCase
	log     logger.Logger
}

func NewSubscriptionHandler(useCase usecase.SubscriptionUseCase, log logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{useCase: useCase, log: log}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req entity.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.useCase.Create(c.Request.Context(), req)
	if err != nil {
		h.log.Error("Failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	h.log.Info("Subscription created", "id", subscription.ID)
	c.JSON(http.StatusCreated, subscription)
}
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid subscription ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to get subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
		return
	}

	if subscription == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid subscription ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req entity.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.useCase.Update(c.Request.Context(), id, req)
	if err != nil {
		h.log.Error("Failed to update subscription", "id", id, "error", err)
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		}
		return
	}

	h.log.Info("Subscription updated", "id", id)
	c.JSON(http.StatusOK, subscription)
}
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid subscription ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		h.log.Error("Failed to delete subscription", "id", id, "error", err)
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		}
		return
	}

	h.log.Info("Subscription deleted", "id", id)
	c.Status(http.StatusNoContent)
}
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	var filter entity.SubscriptionFilter

	if userID := c.Query("user_id"); userID != "" {
		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		filter.UserID = &parsedUUID
	}

	if serviceName := c.Query("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	subscriptions, err := h.useCase.List(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("Failed to list subscriptions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	var filter entity.SubscriptionFilter

	if userID := c.Query("user_id"); userID != "" {
		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		filter.UserID = &parsedUUID
	}

	if serviceName := c.Query("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	totalCost, err := h.useCase.GetTotalCost(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("Failed to calculate total cost", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total cost"})
		return
	}

	c.JSON(http.StatusOK, totalCost)
}
