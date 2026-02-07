package dao

import "time"

type SubscriptionDB struct {
	ID          string     `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"price"`
	UserID      string     `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}
