package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "effective_mobile_service/internal/entity"
    "effective_mobile_service/internal/usecase"
)

type SubscriptionHandler struct {
    useCase usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(useCase usecase.SubscriptionUseCase) *SubscriptionHandler {
    return &SubscriptionHandler{useCase: useCase}
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body entity.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} entity.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
    var req entity.CreateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    subscription, err := h.useCase.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, subscription)
}

// GetSubscription godoc
// @Summary Get subscription
// @Description Get subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} entity.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
        return
    }
    
    subscription, err := h.useCase.GetByID(c.Request.Context(), id)
    if err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    c.JSON(http.StatusOK, subscription)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription data
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body entity.UpdateSubscriptionRequest true "Update data"
// @Success 200 {object} entity.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
        return
    }
    
    var req entity.UpdateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    subscription, err := h.useCase.Update(c.Request.Context(), id, req)
    if err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
        return
    }
    
    if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    c.Status(http.StatusNoContent)
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description Get list of subscriptions with filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param start_date query string false "Start date (MM-YYYY)"
// @Param end_date query string false "End date (MM-YYYY)"
// @Success 200 {array} entity.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
    var filter entity.SubscriptionFilter
    
    if userID := c.Query("user_id"); userID != "" {
        filter.UserID = &userID
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, subscriptions)
}

// GetTotalCost godoc
// @Summary Get total cost
// @Description Calculate total cost of subscriptions with filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param start_date query string false "Start date (MM-YYYY)"
// @Param end_date query string false "End date (MM-YYYY)"
// @Success 200 {object} entity.TotalCostResponse
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
    var filter entity.SubscriptionFilter
    
    if userID := c.Query("user_id"); userID != "" {
        filter.UserID = &userID
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, totalCost)
}
