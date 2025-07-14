package service

import (
	"context"
	"encoding/json"
	"firstmod/internal/models"
	"firstmod/internal/ports"
	"log/slog"
)

type OrderService struct {
	db       ports.Repository
	cache    ports.CacheRepository
	producer ports.KafkaProducer
	log      *slog.Logger
}

func NewOrderService(db ports.Repository, cache ports.CacheRepository, log *slog.Logger, producer ports.KafkaProducer) *OrderService {
	return &OrderService{
		db:       db,
		cache:    cache,
		log:      log,
		producer: producer,
	}
}

func (s *OrderService) Add(ctx context.Context, order models.Order) error {
	err := s.db.Add(ctx, order)
	if err != nil {
		return err
	}
	s.cache.Set(order)
	s.log.Debug("order added to cache after DB insert", "orderUID", order.OrderUID)

	orderJSON, err := json.Marshal(order)
	if err != nil {
		s.log.Error("failed to marshal order to JSON for Kafka", "orderUID", order.OrderUID, "error", err)
	} else {
		publishErr := s.producer.Publish(ctx, order.OrderUID, orderJSON)
		if publishErr != nil {
			s.log.Error("failed to publish order to Kafka", "orderUID", order.OrderUID, "error", publishErr)
		} else {
			s.log.Info("order published to Kafka", "orderUID", order.OrderUID)
		}
	}

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (models.Order, error) {
	if order, found := s.cache.Get(orderUID); found {
		s.log.Debug("order retrieved from cache", "orderUID", orderUID)
		return order, nil
	}

	s.log.Debug("order not in cache, fetching from DB", "orderUID", orderUID)
	order, err := s.db.GetInfo(ctx, orderUID)
	if err != nil {
		return order, err
	}

	s.cache.Set(order)
	s.log.Debug("order fetched from DB and added to cache", "orderUID", orderUID)
	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, orderUID string) error {
	err := s.db.Delete(ctx, orderUID)
	if err != nil {
		return err
	}
	s.cache.Delete(orderUID)
	s.log.Debug("order successfully deleted from DB and Cache", "orderUID", orderUID)
	return nil
}

func (s *OrderService) GetOrderIDs(ctx context.Context) ([]string, error) {
	uids := s.cache.GetAllUIDs()
	if len(uids) == 0 {
		s.log.Warn("cache is empty for UIDs, fetching all UIDs from DB", "action", "GetOrderIDs")
		dbUids, err := s.db.GetIDs(ctx)
		if err != nil {
			return nil, err
		}
		return dbUids, nil
	}
	s.log.Debug("retrieved order UIDs from cache")
	return uids, nil
}

func (s *OrderService) LoadCacheFromDB(ctx context.Context) error {
	return s.cache.LoadToCacheFromDB(ctx, s.db)
}
