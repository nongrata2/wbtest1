package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"firstmod/internal/models"
	"log/slog"
	"net/http"
)

type Storage interface {
	Add(context.Context, models.Order) error
	GetInfo(context.Context, string) (models.Order, error)
	GetIDs(context.Context) ([]string, error)
	Delete(context.Context, int) error
}

func GetOrderByIDHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orderUID := r.PathValue("orderID")
		if orderUID == "" {
			log.Error("orderUID is missing in URL path")
			http.Error(w, "Order ID is missing", http.StatusBadRequest)
			return
		}
		log.Debug("received request to get order info", "order_uid", orderUID)

		order, err := storage.GetInfo(r.Context(), orderUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("order not found", "order_uid", orderUID)
				http.Error(w, "Order not found", http.StatusNotFound)
				return
			}
			log.Error("failed to get order from storage", "order_uid", orderUID, "error", err)
			http.Error(w, "Failed to retrieve order info", http.StatusInternalServerError)
			return
		}

		responseJSON, err := json.MarshalIndent(order, "", "    ")
		if err != nil {
			log.Error("failed to marshal JSON response", "order_uid", order.OrderUID, "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseJSON)
		if err != nil {
			log.Error("failed to write response", "order_uid", order.OrderUID, "error", err)
		}
	}
}

func CreateOrderHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var order models.Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			log.Error("failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = storage.Add(r.Context(), order)
		if err != nil {
			log.Error("failed to add order to storage", "order_uid", order.OrderUID, "error", err)
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		log.Info("order created successfully", "order_uid", order.OrderUID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Order created successfully", "order_uid": order.OrderUID})
	}
}

func DeleteOrderHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetOrdersIDsHandler(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Warn("received non-GET request for order UIDs list", "method", r.Method)
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		uids, err := storage.GetIDs(r.Context())
		if err != nil {
			log.Error("failed to get order UIDs from storage", "error", err)
			http.Error(w, "Failed to retrieve order IDs", http.StatusInternalServerError)
			return
		}

		response := map[string][]string{"order_uids": uids}
		responseJSON, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			log.Error("failed to marshal UIDs response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseJSON)
		if err != nil {
			log.Error("failed to write UIDs response", "error", err)
		}
		log.Info("successfully retrieved and sent all order UIDs", "count", len(uids))
	}
}
