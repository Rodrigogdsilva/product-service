package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"product-service/src/config"
	"product-service/src/domain"
	"product-service/src/service"

	"github.com/google/uuid"
)

type Handler struct {
	service service.ProductService
	cfg     *config.Config
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type UpdateProductRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
}

type GetProductRequest struct {
	ID uuid.UUID `json:"id"`
}

type ReduceStockRequest struct {
	ID       uuid.UUID `json:"id"`
	Quantity int       `json:"quantity"`
}

type DeleteProductResquest struct {
	ID uuid.UUID `json:"id"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHandler(svc service.ProductService, cfg *config.Config) *Handler {
	return &Handler{
		service: svc,
		cfg:     cfg,
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %v", err)

	if errors.Is(err, domain.ErrProductNotFound) {
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Code: "PRODUCT_NOT_FOUND", Message: err.Error()})
		return
	}
	if errors.Is(err, domain.ErrParametersMissing) || errors.Is(err, domain.ErrInvalidPrice) || errors.Is(err, domain.ErrInvalidStock) {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "INVALID_INPUT", Message: err.Error()})
		return
	}

	if errors.Is(err, domain.ErrFailedCreatingProduct) {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Code: "FAILED_CREATING_PRODUCT", Message: err.Error()})
		return
	}

	// Erro genérico
	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_SERVER_ERROR", Message: "An unexpected error occurred"})
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST_BODY", Message: "Invalid request body"})
		return
	}

	err := h.service.Create(r.Context(), req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"message": "Product created successfully"})
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {

	var getProduct GetProductRequest

	if err := json.NewDecoder(r.Body).Decode(&getProduct); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	product, err := h.service.GetProductByID(r.Context(), getProduct.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) HandleReduceStock(w http.ResponseWriter, r *http.Request) {

	var reduceStock ReduceStockRequest

	if err := json.NewDecoder(r.Body).Decode(&reduceStock); err != nil {
		h.handleError(w, err)
		return
	}

	err := h.service.ReduceStock(r.Context(), reduceStock.ID, reduceStock.Quantity)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Stock reduced successfully"})
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, err)
		return
	}

	productToUpdate := &domain.Product{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	err := h.service.Update(r.Context(), productToUpdate)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {

	var deleteProduct DeleteProductResquest
	if err := json.NewDecoder(r.Body).Decode(&deleteProduct); err != nil {
		h.handleError(w, err)
		return
	}

	err := h.service.Delete(r.Context(), deleteProduct.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}
