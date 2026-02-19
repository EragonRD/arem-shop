package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

const testShopID = "11111111-1111-1111-1111-111111111111"

type fakeProductHandlerService struct {
	listData interface{}
	listErr  error

	getByIDData interface{}
	getByIDErr  error

	createData interface{}
	createErr  error

	updateData interface{}
	updateErr  error

	deleteErr error
}

func (f *fakeProductHandlerService) List(_ context.Context, _ string, _ models.UserRole) (interface{}, error) {
	return f.listData, f.listErr
}

func (f *fakeProductHandlerService) GetByID(_ context.Context, _, _ string, _ models.UserRole) (interface{}, error) {
	return f.getByIDData, f.getByIDErr
}

func (f *fakeProductHandlerService) Create(_ context.Context, _ string, _ models.UserRole, _ dto.CreateProductRequest) (interface{}, error) {
	return f.createData, f.createErr
}

func (f *fakeProductHandlerService) Update(_ context.Context, _, _ string, _ models.UserRole, _ dto.UpdateProductRequest) (interface{}, error) {
	return f.updateData, f.updateErr
}

func (f *fakeProductHandlerService) Delete(_ context.Context, _, _ string) error {
	return f.deleteErr
}

type testEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

func newProductHandlerTestRouter(handler *ProductHandler, withShop, withClaims bool, role models.UserRole) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if withShop {
			c.Set("shop_id", testShopID)
		}
		if withClaims {
			c.Set("auth_claims", &utils.Claims{
				Role:   string(role),
				ShopID: testShopID,
			})
		}
		c.Next()
	})

	router.GET("/products", handler.List)
	router.GET("/products/:id", handler.GetByID)
	router.POST("/products", handler.Create)
	router.PUT("/products/:id", handler.Update)
	router.DELETE("/products/:id", handler.Delete)
	return router
}

func decodeEnvelope(t *testing.T, body []byte) testEnvelope {
	t.Helper()
	var envelope testEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("failed to decode envelope: %v body=%s", err, string(body))
	}
	return envelope
}

func TestProductHandler_GetByID_SuccessEnvelope(t *testing.T) {
	service := &fakeProductHandlerService{
		getByIDData: map[string]interface{}{
			"id":   "product-1",
			"name": "Laptop",
		},
	}
	handler := NewProductHandler(service)
	router := newProductHandlerTestRouter(handler, true, true, models.RoleSuperAdmin)

	req := httptest.NewRequest(http.MethodGet, "/products/product-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	envelope := decodeEnvelope(t, rec.Body.Bytes())
	if !envelope.Success {
		t.Fatalf("expected success=true, got false with error=%s", envelope.Error)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("failed to decode data payload: %v", err)
	}
	if data["id"] != "product-1" {
		t.Fatalf("expected product id product-1, got %#v", data["id"])
	}
}

func TestProductHandler_List_MissingShopContext_Returns401(t *testing.T) {
	service := &fakeProductHandlerService{
		listData: []map[string]interface{}{},
	}
	handler := NewProductHandler(service)
	router := newProductHandlerTestRouter(handler, false, true, models.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}

	envelope := decodeEnvelope(t, rec.Body.Bytes())
	if envelope.Success {
		t.Fatalf("expected success=false, got true")
	}
	if envelope.Error != "shop context not found" {
		t.Fatalf("unexpected error message: %s", envelope.Error)
	}
}

func TestProductHandler_Create_AdminPurchasePriceForbidden_Returns403(t *testing.T) {
	service := &fakeProductHandlerService{
		createErr: services.ErrPurchasePriceForbidden,
	}
	handler := NewProductHandler(service)
	router := newProductHandlerTestRouter(handler, true, true, models.RoleAdmin)

	body := []byte(`{
		"name":"Laptop",
		"description":"Test product",
		"category":"Laptops",
		"purchasePrice":800,
		"sellingPrice":999.99,
		"stock":10,
		"imageURL":"https://example.com/laptop.jpg"
	}`)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", rec.Code, rec.Body.String())
	}

	envelope := decodeEnvelope(t, rec.Body.Bytes())
	if envelope.Success {
		t.Fatalf("expected success=false, got true")
	}
	if envelope.Error != services.ErrPurchasePriceForbidden.Error() {
		t.Fatalf("unexpected error message: %s", envelope.Error)
	}
}

func TestProductHandler_Update_InvalidProductID_Returns400(t *testing.T) {
	service := &fakeProductHandlerService{
		updateErr: services.ErrInvalidProductID,
	}
	handler := NewProductHandler(service)
	router := newProductHandlerTestRouter(handler, true, true, models.RoleAdmin)

	body := []byte(`{
		"name":"Laptop",
		"description":"Updated",
		"category":"Laptops",
		"sellingPrice":999.99,
		"stock":10,
		"imageURL":"https://example.com/laptop.jpg"
	}`)
	req := httptest.NewRequest(http.MethodPut, "/products/not-a-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}

	envelope := decodeEnvelope(t, rec.Body.Bytes())
	if envelope.Success {
		t.Fatalf("expected success=false, got true")
	}
}

func TestProductHandler_Delete_ProductHasTransactions_Returns409(t *testing.T) {
	service := &fakeProductHandlerService{
		deleteErr: services.ErrProductHasTransactions,
	}
	handler := NewProductHandler(service)
	router := newProductHandlerTestRouter(handler, true, true, models.RoleSuperAdmin)

	req := httptest.NewRequest(http.MethodDelete, "/products/product-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", rec.Code, rec.Body.String())
	}

	envelope := decodeEnvelope(t, rec.Body.Bytes())
	if envelope.Success {
		t.Fatalf("expected success=false, got true")
	}
	if envelope.Error != services.ErrProductHasTransactions.Error() {
		t.Fatalf("unexpected error message: %s", envelope.Error)
	}
}

func TestProductHandler_Delete_Success_Returns200Envelope(t *testing.T) {
	service := &fakeProductHandlerService{}
	handler := NewProductHandler(service)
	router := newProductHandlerTestRouter(handler, true, true, models.RoleSuperAdmin)

	req := httptest.NewRequest(http.MethodDelete, "/products/product-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	envelope := decodeEnvelope(t, rec.Body.Bytes())
	if !envelope.Success {
		t.Fatalf("expected success=true, got false with error=%s", envelope.Error)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("failed to decode data payload: %v", err)
	}
	if data["message"] != "product deleted" {
		t.Fatalf("expected message=product deleted, got %#v", data["message"])
	}
}
