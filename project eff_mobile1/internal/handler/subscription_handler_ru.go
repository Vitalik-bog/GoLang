package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "effective_mobile_service/internal/entity"
    "effective_mobile_service/internal/usecase"
    "effective_mobile_service/internal/locale/ru"
)

// Русская версия обработчика
type SubscriptionHandlerRU struct {
    useCase usecase.SubscriptionUseCase
}

func NewSubscriptionHandlerRU(useCase usecase.SubscriptionUseCase) *SubscriptionHandlerRU {
    return &SubscriptionHandlerRU{useCase: useCase}
}

func (h *SubscriptionHandlerRU) CreateSubscription(c *gin.Context) {
    var req entity.CreateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   ru.InvalidJSON,
            "details": err.Error(),
        })
        return
    }

    // Валидация
    if req.ServiceName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Название сервиса обязательно"})
        return
    }
    if req.Price <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Цена должна быть больше 0"})
        return
    }
    if req.UserID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
        return
    }
    if req.StartDate == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Дата начала обязательна"})
        return
    }

    subscription, err := h.useCase.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   ru.DatabaseError,
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":      ru.SubscriptionCreated,
        "subscription": subscription,
    })
}

func (h *SubscriptionHandlerRU) GetSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": ru.InvalidSubscriptionID})
        return
    }

    subscription, err := h.useCase.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": ru.DatabaseError})
        return
    }

    if subscription == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": ru.SubscriptionNotFound})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":      "Информация о подписке",
        "subscription": subscription,
    })
}

func (h *SubscriptionHandlerRU) ListSubscriptions(c *gin.Context) {
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": ru.DatabaseError})
        return
    }

    if len(subscriptions) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": ru.SubscriptionListEmpty,
            "count":   0,
            "data":    []string{},
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Список подписок",
        "count":   len(subscriptions),
        "data":    subscriptions,
    })
}

func (h *SubscriptionHandlerRU) GetTotalCost(c *gin.Context) {
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": ru.DatabaseError})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":    ru.TotalCostCalculated,
        "total_cost": totalCost.TotalCost,
        "currency":   "RUB",
        "filters": gin.H{
            "user_id":      filter.UserID,
            "service_name": filter.ServiceName,
            "start_date":   filter.StartDate,
            "end_date":     filter.EndDate,
        },
    })
}

func (h *SubscriptionHandlerRU) UpdateSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": ru.InvalidSubscriptionID})
        return
    }

    var req entity.UpdateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": ru.InvalidJSON})
        return
    }

    subscription, err := h.useCase.Update(c.Request.Context(), id, req)
    if err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": ru.SubscriptionNotFound})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": ru.DatabaseError})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":      ru.SubscriptionUpdated,
        "subscription": subscription,
    })
}

func (h *SubscriptionHandlerRU) DeleteSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": ru.InvalidSubscriptionID})
        return
    }

    if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": ru.SubscriptionNotFound})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": ru.DatabaseError})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": ru.SubscriptionDeleted,
        "id":      id.String(),
    })
}

func (h *SubscriptionHandlerRU) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "healthy",
        "message": ru.HealthCheckSuccess,
        "service": "Сервис управления подписками",
        "version": "1.0",
    })
}
