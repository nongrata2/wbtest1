package handlers

import (
	"context"
	"firstmod/internal/models"
	"log/slog"
	"net/http"
)

type Storage interface {
	Add(context.Context, models.Order) error
	GetInfo(context.Context, int) (models.Order, error)
	Delete(context.Context, int) error
}

func GetOrderByIDHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func CreateOrderHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteOrderHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
