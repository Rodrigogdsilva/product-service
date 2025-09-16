package api

import (
	"encoding/json"
	"net/http"
	"product-service/src/config"
	"product-service/src/service"

	"github.com/oklog/ulid/v2"
)

type Handler struct {
	service service.ProductService
	cfg     *config.Config
}

type CreateProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type GetProduct struct {
	ID ulid.ULID `json:"id"`
}

type UpdateProduct struct {
	ID          ulid.ULID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
}

type ReduceStock struct {
	ID       ulid.ULID `json:"id"`
	Quantity int       `json:"quantity"`
}

type DeleteProduct struct {
	ID ulid.ULID `json:"id"`
}

func NewHandler(svc service.ProductService, cfg *config.Config) *Handler {
	return &Handler{
		service: svc,
		cfg:     cfg,
	}
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var createProduct CreateProduct
	if err := json.NewDecoder(r.Body).Decode(&createProduct); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := h.service.Create(r.Context(), createProduct.Name, createProduct.Description, createProduct.Price, createProduct.Stock)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"message": "Product created successfully"})

}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {

	var getProduct GetProduct

	if err := json.NewDecoder(r.Body).Decode(&getProduct); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	product, err := h.service.GetProductByID(r.Context(), getProduct.ID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) HandleReduceStock(w http.ResponseWriter, r *http.Request) {

	var reduceStock ReduceStock

	if err := json.NewDecoder(r.Body).Decode(&reduceStock); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := h.service.ReduceStock(r.Context(), reduceStock.ID, reduceStock.Quantity)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Stock reduced successfully"})
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var updateProduct UpdateProduct
	if err := json.NewDecoder(r.Body).Decode(&updateProduct); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := h.service.Update(r.Context(), updateProduct.ID, updateProduct.Name, updateProduct.Description, updateProduct.Price, updateProduct.Stock)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {

	var deleteProduct DeleteProduct
	if err := json.NewDecoder(r.Body).Decode(&deleteProduct); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := h.service.Delete(r.Context(), deleteProduct.ID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
