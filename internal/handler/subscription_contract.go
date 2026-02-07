package handler

import (
	"context"
	"time"

	"online-subscription-service/internal/model"
)

type subscriptionsService interface {
	Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	GetAll(ctx context.Context) (*[]model.Subscription, error)
	GetByID(ctx context.Context, id string) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	Delete(ctx context.Context, id string) error
	GetTotalCost(ctx context.Context, userID, serviceName string, start, end time.Time) (int, error)
}
