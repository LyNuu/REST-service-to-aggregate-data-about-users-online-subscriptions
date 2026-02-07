package usecase

import (
	"context"
	"errors"
	"time"

	dtomapper "github.com/dranikpg/dto-mapper"

	"online-subscription-service/internal/model"
	"online-subscription-service/internal/model/dao"
	"online-subscription-service/internal/repository"
	l "online-subscription-service/pkg/logger"
)

type SubscriptionService struct {
	repo subscriptionRepository
	l    l.Logger
}

func NewSubscriptionService(repo subscriptionRepository, l l.Logger) *SubscriptionService {
	return &SubscriptionService{repo: repo, l: l}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	var subDao *dao.SubscriptionDB
	if err := dtomapper.Map(&subDao, sub); err != nil {
		s.l.Warn("failed map model -> dao", "error", err.Error())
		return nil, err
	}
	subModel, err := s.repo.Create(ctx, subDao)
	if err != nil {
		s.l.Warn("service err: ", "error", err.Error())
		return nil, err
	}

	var resSub *model.Subscription
	if err := dtomapper.Map(&resSub, subModel); err != nil {
		s.l.Warn("failed map dao -> model", "error", err.Error())
		return nil, err
	}

	return resSub, nil
}

func (s *SubscriptionService) GetAll(ctx context.Context) (*[]model.Subscription, error) {
	var subsModel *[]model.Subscription

	subDao, err := s.repo.GetAll(ctx)
	if err != nil {
		s.l.Warn("service error: ", "error", err.Error())
		return nil, err
	}

	if err := dtomapper.Map(&subsModel, subDao); err != nil {
		s.l.Warn("failed map dao -> model", "error", err.Error())
		return nil, err
	}

	return subsModel, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.l.Info("not found by id", "error", err.Error())
			return nil, ErrNotFound
		}
		s.l.Warn("service error", "error", err.Error())
		return nil, err
	}

	var subModel model.Subscription
	if err := dtomapper.Map(&subModel, sub); err != nil {
		s.l.Warn("failed mapping dao -> model", "error", err.Error())
		return nil, err
	}
	return &subModel, nil
}

func (s *SubscriptionService) Update(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	var subDB dao.SubscriptionDB
	if err := dtomapper.Map(&subDB, sub); err != nil {
		s.l.Warn("failed mapping model -> dao", "error", err.Error())
		return nil, err
	}
	subDao, err := s.repo.Update(ctx, &subDB)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.l.Info("not found by id", "error", err.Error())
			return nil, ErrNotFound
		}
		s.l.Warn("service error", "error", err.Error())
		return nil, err
	}

	var subModel model.Subscription
	if err := dtomapper.Map(&subModel, subDao); err != nil {
		s.l.Warn("failed mapping dao -> model", "error", err.Error())
		return nil, err
	}
	return &subModel, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.l.Info("not found by id", "error", err.Error())
			return ErrNotFound
		}
		s.l.Warn("service error", "error", err.Error())
		return err
	}
	return nil
}

func (s *SubscriptionService) GetTotalCost(ctx context.Context, userID, serviceName string, start, end time.Time) (int, error) {
	if start.After(end) {
		s.l.Info("method GetTotalCost()", "error", errors.New("start date cannot be after end date"))
		return 0, errors.New("start date cannot be after end date")
	}

	return s.repo.GetTotalCost(ctx, userID, serviceName, start, end)
}
