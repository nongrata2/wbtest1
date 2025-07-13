package cache

import (
	"context"
	"firstmod/internal/models"
	"log/slog"
	"sync"
)

type Cache struct {
	data map[string]models.Order
	mu   sync.RWMutex
	log  *slog.Logger
}

type Storage interface {
	GetInfo(context.Context, string) (models.Order, error)
	GetIDs(context.Context) ([]string, error)
}

func (c *Cache) NewCache(log *slog.Logger) *Cache {
	return &Cache{
		data: make(map[string]models.Order),
		log:  log,
	}
}

func (c *Cache) Get(orderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, found := c.data[orderUID]
	if found {
		c.log.Debug("order found in cache", "orderUID", orderUID)
	} else {
		c.log.Debug("order is not found in cache", "orderUID", orderUID)
	}
	return order, found
}

func (c *Cache) Set(order models.Order) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.data[order.OrderUID] = order
	c.log.Debug("order added in cache", "orderUID", order.OrderUID)
}

func (c *Cache) Delete(orderUID string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	delete(c.data, orderUID)
	c.log.Debug("order was removed from cache", "orderUID", orderUID)
}

func (c *Cache) LoadFromDB(ctx context.Context, db Storage) error {
	uids, err := db.GetIDs(ctx)
	if err != nil {
		c.log.Error("Error getting IDs from DB", "error", err)
		return err
	}
	for _, uid := range uids {
		order, err := db.GetInfo(ctx, uid)
		if err != nil {
			c.log.Warn("Failed to get info about order", "orderUID", uid, "error", err)
			continue
		}
		c.Set(order)
	}
	c.log.Error("Cache loaded successfully")
	return nil
}
