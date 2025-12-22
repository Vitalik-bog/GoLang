package usecase

import (
    "context"
    "fmt"
    "time"
    "effective_mobile_service/internal/entity"
    "effective_mobile_service/internal/repository"
    "github.com/google/uuid"
)

type SubscriptionUseCase interface {
    Create(ctx context.Context, req entity.CreateSubscriptionRequest) (*entity.Subscription, error)
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
    Update(ctx context.Context, id uuid.UUID, req entity.UpdateSubscriptionRequest) (*entity.Subscription, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filter entity.SubscriptionFilter) ([]*entity.Subscription, error)
    GetTotalCost(ctx context.Context, filter entity.SubscriptionFilter) (*entity.TotalCostResponse, error)
}

type subscriptionUseCase struct {
    repo repository.SubscriptionRepository
}

func NewSubscriptionUseCase(repo repository.SubscriptionRepository) SubscriptionUseCase {
    return &subscriptionUseCase{repo: repo}
}

func (uc *subscriptionUseCase) Create(ctx context.Context, req entity.CreateSubscriptionRequest) (*entity.Subscription, error) {
    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }
    
    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        return nil, fmt.Errorf("invalid start date format: %w", err)
    }
    
    var endDate *time.Time
    if req.EndDate != nil {
        parsedEndDate, err := time.Parse("01-2006", *req.EndDate)
        if err != nil {
            return nil, fmt.Errorf("invalid end date format: %w", err)
        }
        endDate = &parsedEndDate
    }
    
    subscription := &entity.Subscription{
        ID:         uuid.New(),
        ServiceName: req.ServiceName,
        Price:      req.Price,
        UserID:     userID,
        StartDate:  startDate,
        EndDate:    endDate,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    // Здесь должен быть вызов репозитория
    // err = uc.repo.Create(ctx, subscription)
    // if err != nil { ... }
    
    return subscription, nil
}

func (uc *subscriptionUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
    // Заглушка
    return &entity.Subscription{
        ID:         id,
        ServiceName: "Test Service",
        Price:      100,
        UserID:     uuid.New(),
        StartDate:  time.Now(),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

func (uc *subscriptionUseCase) Update(ctx context.Context, id uuid.UUID, req entity.UpdateSubscriptionRequest) (*entity.Subscription, error) {
    // Заглушка
    return &entity.Subscription{
        ID:         id,
        ServiceName: "Updated Service",
        Price:      200,
        UserID:     uuid.New(),
        StartDate:  time.Now(),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

func (uc *subscriptionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
    return nil
}

func (uc *subscriptionUseCase) List(ctx context.Context, filter entity.SubscriptionFilter) ([]*entity.Subscription, error) {
    // Заглушка
    return []*entity.Subscription{
        {
            ID:         uuid.New(),
            ServiceName: "Service 1",
            Price:      100,
            UserID:     uuid.New(),
            StartDate:  time.Now(),
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        },
    }, nil
}

func (uc *subscriptionUseCase) GetTotalCost(ctx context.Context, filter entity.SubscriptionFilter) (*entity.TotalCostResponse, error) {
    // Заглушка
    return &entity.TotalCostResponse{TotalCost: 1000}, nil
}
