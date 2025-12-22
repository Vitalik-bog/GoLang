package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"
    "effective_mobile_service/internal/entity"
    "effective_mobile_service/internal/repository"
    "github.com/google/uuid"
)

type subscriptionRepository struct {
    db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) repository.SubscriptionRepository {
    return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *entity.Subscription) error {
    query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
    
    _, err := r.db.ExecContext(ctx, query,
        subscription.ID,
        subscription.ServiceName,
        subscription.Price,
        subscription.UserID,
        subscription.StartDate,
        subscription.EndDate,
        subscription.CreatedAt,
        subscription.UpdatedAt,
    )
    
    return err
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE id = $1
    `
    
    var subscription entity.Subscription
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &subscription.ID,
        &subscription.ServiceName,
        &subscription.Price,
        &subscription.UserID,
        &subscription.StartDate,
        &subscription.EndDate,
        &subscription.CreatedAt,
        &subscription.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return &subscription, err
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *entity.Subscription) error {
    query := `
        UPDATE subscriptions
        SET service_name = $1, price = $2, end_date = $3, updated_at = $4
        WHERE id = $5
    `
    
    subscription.UpdatedAt = time.Now()
    
    result, err := r.db.ExecContext(ctx, query,
        subscription.ServiceName,
        subscription.Price,
        subscription.EndDate,
        subscription.UpdatedAt,
        subscription.ID,
    )
    
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return fmt.Errorf("subscription not found")
    }
    
    return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM subscriptions WHERE id = $1`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return fmt.Errorf("subscription not found")
    }
    
    return nil
}

func (r *subscriptionRepository) List(ctx context.Context, filter entity.SubscriptionFilter) ([]*entity.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE 1=1
    `
    
    var conditions []string
    var args []interface{}
    argCounter := 1
    
    if filter.UserID != nil {
        userUUID, err := uuid.Parse(*filter.UserID)
        if err == nil {
            conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCounter))
            args = append(args, userUUID)
            argCounter++
        }
    }
    
    if filter.ServiceName != nil {
        conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argCounter))
        args = append(args, "%"+*filter.ServiceName+"%")
        argCounter++
    }
    
    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }
    
    query += " ORDER BY created_at DESC"
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var subscriptions []*entity.Subscription
    for rows.Next() {
        var subscription entity.Subscription
        err := rows.Scan(
            &subscription.ID,
            &subscription.ServiceName,
            &subscription.Price,
            &subscription.UserID,
            &subscription.StartDate,
            &subscription.EndDate,
            &subscription.CreatedAt,
            &subscription.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        subscriptions = append(subscriptions, &subscription)
    }
    
    return subscriptions, nil
}

func (r *subscriptionRepository) GetTotalCost(ctx context.Context, filter entity.SubscriptionFilter) (int, error) {
    query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
    
    var conditions []string
    var args []interface{}
    argCounter := 1
    
    if filter.UserID != nil {
        userUUID, err := uuid.Parse(*filter.UserID)
        if err == nil {
            conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCounter))
            args = append(args, userUUID)
            argCounter++
        }
    }
    
    if filter.ServiceName != nil {
        conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argCounter))
        args = append(args, "%"+*filter.ServiceName+"%")
        argCounter++
    }
    
    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }
    
    var totalCost int
    err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalCost)
    return totalCost, err
}
