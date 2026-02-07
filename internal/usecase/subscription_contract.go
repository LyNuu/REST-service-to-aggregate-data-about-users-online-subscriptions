package usecase

import (
	"context"
	"time"

	"online-subscription-service/internal/model/dao"
)

type subscriptionRepository interface {
	Create(ctx context.Context, db *dao.SubscriptionDB) (*dao.SubscriptionDB, error)
	GetAll(ctx context.Context) (*[]dao.SubscriptionDB, error)
	GetByID(ctx context.Context, id string) (*dao.SubscriptionDB, error)
	Update(ctx context.Context, sub *dao.SubscriptionDB) (*dao.SubscriptionDB, error)
	Delete(ctx context.Context, id string) error
	GetTotalCost(ctx context.Context, userID, serviceName string, start, end time.Time) (int, error)
}
