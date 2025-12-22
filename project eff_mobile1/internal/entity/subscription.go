package entity

import (
    "time"
    "github.com/google/uuid"
)

type Subscription struct {
    ID         uuid.UUID `json:"id" db:"id"`
    ServiceName string    `json:"service_name" db:"service_name"`
    Price      int       `json:"price" db:"price"`
    UserID     uuid.UUID `json:"user_id" db:"user_id"`
    StartDate  time.Time `json:"start_date" db:"start_date"`
    EndDate    *time.Time `json:"end_date,omitempty" db:"end_date"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionRequest struct {
    ServiceName string    `json:"service_name" binding:"required,min=1,max=100"`
    Price      int       `json:"price" binding:"required,min=0"`
    UserID     string    `json:"user_id" binding:"required"`
    StartDate  string    `json:"start_date" binding:"required"`
    EndDate    *string   `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
    ServiceName *string    `json:"service_name,omitempty" binding:"omitempty,min=1,max=100"`
    Price      *int       `json:"price,omitempty" binding:"omitempty,min=0"`
    EndDate    *string    `json:"end_date,omitempty"`
}

type SubscriptionFilter struct {
    UserID     *string    `form:"user_id"`      // Изменили на *string
    ServiceName *string    `form:"service_name"`
    StartDate  *string    `form:"start_date"`
    EndDate    *string    `form:"end_date"`
}

type TotalCostResponse struct {
    TotalCost int `json:"total_cost"`
}
