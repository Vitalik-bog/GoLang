package controller

import (
    "context"
    "effective_mobile_service/internal/entity"
    "effective_mobile_service/internal/usecase"
    "github.com/google/uuid"
)

type SubscriptionController struct {
    useCase usecase.SubscriptionUseCase
}

func NewSubscriptionController(useCase usecase.SubscriptionUseCase) *SubscriptionController {
    return &SubscriptionController{useCase: useCase}
}

func (c *SubscriptionController) CreateSubscription(req entity.CreateSubscriptionRequest) (*entity.Subscription, error) {
    return c.useCase.Create(context.Background(), req)
}

func (c *SubscriptionController) GetSubscription(id uuid.UUID) (*entity.Subscription, error) {
    return c.useCase.GetByID(context.Background(), id)
}

func (c *SubscriptionController) UpdateSubscription(id uuid.UUID, req entity.UpdateSubscriptionRequest) (*entity.Subscription, error) {
    return c.useCase.Update(context.Background(), id, req)
}

func (c *SubscriptionController) DeleteSubscription(id uuid.UUID) error {
    return c.useCase.Delete(context.Background(), id)
}

func (c *SubscriptionController) ListSubscriptions(filter entity.SubscriptionFilter) ([]*entity.Subscription, error) {
    return c.useCase.List(context.Background(), filter)
}

func (c *SubscriptionController) GetTotalCost(filter entity.SubscriptionFilter) (*entity.TotalCostResponse, error) {
    return c.useCase.GetTotalCost(context.Background(), filter)
}
