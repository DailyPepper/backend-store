package handlers

import (
	"backend-store/internal/storage"
	"encoding/json"
	"log"
	"net/http"
)

type OrderHander struct {
	storage storage.Storage
}

func NewOrderHandler(storage storage.Storage) *OrderHander {
	return &OrderHander{storage: storage}
}

func (h *OrderHander) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.storage.GetAllOrders()

	if err != nil {
		http.Error(w, "Нет заказов", http.StatusInternalServerError)
		log.Printf("Ошибка получения заказов: %v", err)
		return
	}

	response, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		log.Printf("Ошибка маршалинга: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
