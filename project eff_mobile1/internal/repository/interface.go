package repository

import (
    "context"
    "effective_mobile_service/internal/entity"
    "github.com/google/uuid"
)

type SubscriptionRepository interface {
    Create(ctx context.Context, subscription *entity.Subscription) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
    Update(ctx context.Context, subscription *entity.Subscription) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filter entity.SubscriptionFilter) ([]*entity.Subscription, error)
    GetTotalCost(ctx context.Context, filter entity.SubscriptionFilter) (int, error)
}
