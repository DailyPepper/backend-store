package handlers

import (
	"backend-store/internal/storage"
	"encoding/json"
	"log"
	"net/http"
)

type ProductHandler struct {
	storage storage.Storage
}

func NewProductHandler(storage storage.Storage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// var product storage.Storage

	// if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
	// 	http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }

}

func (h *ProductHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.storage.GetAllProduct()

	if err != nil {
		http.Error(w, "Нет такого товара", http.StatusInternalServerError)
		log.Printf("Ошибка получения заказов: %v", err)
		return
	}

	response, err := json.Marshal(products)

	if err != nil {
		http.Error(w, "Ошибка форматиорвания ответа", http.StatusInternalServerError)
		log.Printf("Ошибка маршалинга: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
