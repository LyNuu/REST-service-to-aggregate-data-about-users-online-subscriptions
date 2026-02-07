package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"online-subscription-service/internal/model/dao"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *dao.SubscriptionDB) (*dao.SubscriptionDB, error) {
	var resSub dao.SubscriptionDB
	err := r.pool.QueryRow(ctx, `
			INSERT INTO subscriptions(user_id, service_name, price, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, user_id, service_name, price, start_date, end_date`,
		sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate).
		Scan(&resSub.ID,
			&resSub.UserID,
			&resSub.ServiceName,
			&resSub.Price,
			&resSub.StartDate,
			&resSub.EndDate)

	if err != nil {
		return nil, err
	}

	return &resSub, nil
}

func (r *SubscriptionRepository) GetAll(ctx context.Context) (*[]dao.SubscriptionDB, error) {
	tx, err := r.pool.Begin(ctx)
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []dao.SubscriptionDB
	for rows.Next() {
		var sub dao.SubscriptionDB
		err = rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &subs, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*dao.SubscriptionDB, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var sub dao.SubscriptionDB
	err = tx.QueryRow(ctx, `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
		`, id).Scan(&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *dao.SubscriptionDB) (*dao.SubscriptionDB, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	var subDB dao.SubscriptionDB
	err = tx.QueryRow(ctx, `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE id = $5
		RETURNING id, service_name, price, user_id, start_date, end_date
		`, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.ID).
		Scan(&subDB.ID,
			&subDB.ServiceName,
			&subDB.Price,
			&subDB.UserID,
			&subDB.StartDate,
			&subDB.EndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &subDB, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	res, err := tx.Exec(ctx, `
		DELETE FROM subscriptions
		WHERE id = $1
		`, id)

	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, userID, serviceName string, start, end time.Time) (int, error) {
	var total int
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE user_id = $1
          AND start_date >= $2
          AND start_date <= $3
          AND (service_name = $4 OR $4 = '')
    `

	err := r.pool.QueryRow(ctx, query, userID, start, end, serviceName).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
