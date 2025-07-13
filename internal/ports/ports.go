package ports

import (
	"context"
	"firstmod/internal/models"
)

type Repository interface {
	Add(ctx context.Context, order models.Order) error
	GetInfo(ctx context.Context, orderUID string) (models.Order, error)
	Delete(ctx context.Context, orderUID string) error
	GetIDs(ctx context.Context) ([]string, error)
}

type CacheRepository interface {
	Get(orderUID string) (models.Order, bool)
	Set(order models.Order)
	Delete(orderUID string)
	GetAllUIDs() []string
	LoadToCacheFromDB(ctx context.Context, db Repository) error
}

type OrderService interface {
	Add(context.Context, models.Order) error
	GetOrder(context.Context, string) (models.Order, error)
	Delete(context.Context, string) error
	GetOrderIDs(context.Context) ([]string, error)
	LoadCacheFromDB(ctx context.Context) error
}
