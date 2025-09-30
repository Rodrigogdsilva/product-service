package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-service/src/config"
	"product-service/src/domain"
	"product-service/src/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreate_Success(t *testing.T) {
	// Arrange: Cria o mock do serviço e o handler.
	mockService := new(service.ProductServiceMock)
	handler := NewHandler(mockService, &config.Config{})

	requestBody := `{"name": "New Product", "description": "A great product", "price": 99.99, "stock": 10}`
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Mock: Diz ao mock para esperar uma chamada ao método 'Create' com os parâmetros específicos e retornar nil (sem erro).
	mockService.On("Create", mock.Anything, "New Product", "A great product", 99.99, 10).Return(nil)

	// Act: Chama o handler.
	handler.HandleCreate(rr, req)

	// Assert: Verifica se o status code é 201 Created e se o mock foi chamado conforme esperado.
	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestHandleCreate_ServiceError(t *testing.T) {
	// Arrange: Cria o mock do serviço e o handler.
	mockService := new(service.ProductServiceMock)
	handler := NewHandler(mockService, &config.Config{})

	requestBody := `{"name": "Invalid Product", "price": -10}` // Dados que causarão um erro
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Mock: Diz ao mock para esperar uma chamada ao método 'Create' e retornar um erro específico.
	mockService.On("Create", mock.Anything, "Invalid Product", "", -10.0, 0).Return(domain.ErrInvalidPrice)

	// Act: Chama o handler.
	handler.HandleCreate(rr, req)

	// Assert: Verifica se o status code é 400 Bad Request e se o mock foi chamado conforme esperado.
	assert.Equal(t, http.StatusBadRequest, rr.Code) // Esperamos um status 400 Bad Request
	mockService.AssertExpectations(t)

	// Verifica se a resposta JSON do erro está correta
	var errResponse ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)
	assert.Equal(t, "INVALID_INPUT", errResponse.Code)
	assert.Equal(t, domain.ErrInvalidPrice.Error(), errResponse.Message)
}

func TestHandleList_Success(t *testing.T) {
	// Arrange: Cria o mock do serviço e o handler.
	mockService := new(service.ProductServiceMock)
	handler := NewHandler(mockService, &config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	rr := httptest.NewRecorder()

	// Mock: Mock para retornar uma lista de produtos.
	expectedProducts := []*domain.Product{{Name: "Test Product 1"}, {Name: "Test Product 2"}}
	mockService.On("ListProducts", mock.Anything).Return(expectedProducts, nil)

	// Act: Chama o handler.
	handler.HandleList(rr, req)

	// Assert: Verifica se o status code é 200 OK e se a resposta contém os produtos esperados.
	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)

	var products []*domain.Product
	json.Unmarshal(rr.Body.Bytes(), &products)
	assert.Len(t, products, 2)
	assert.Equal(t, "Test Product 1", products[0].Name)
}
