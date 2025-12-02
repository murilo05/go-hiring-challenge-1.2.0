package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOKResponse_WithData(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"name": "Test Product",
		"code": "PROD001",
	}

	OKResponse(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Message != "Success" {
		t.Errorf("expected message 'Success', got %s", response.Message)
	}

	if response.Data == nil {
		t.Error("expected data to be present")
	}
}

func TestOKResponse_WithNilData(t *testing.T) {
	w := httptest.NewRecorder()

	OKResponse(w, nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response response
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Error("expected success to be true")
	}
}

func TestOKResponse_WithInvalidData(t *testing.T) {
	w := httptest.NewRecorder()

	invalidData := func() {}

	OKResponse(w, invalidData)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected content in response body")
	}
}

func TestOKResponse_WithChannel(t *testing.T) {
	w := httptest.NewRecorder()

	invalidData := make(chan int)

	OKResponse(w, invalidData)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected content in response body")
	}
}

func TestOKResponse_WithEmptySlice(t *testing.T) {
	w := httptest.NewRecorder()

	data := []string{}

	OKResponse(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response response
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Error("expected data to be present (empty slice)")
	}
}

func TestOKResponse_WithStruct(t *testing.T) {
	w := httptest.NewRecorder()

	type Product struct {
		Code  string  `json:"code"`
		Price float64 `json:"price"`
	}

	data := Product{
		Code:  "PROD001",
		Price: 99.99,
	}

	OKResponse(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Data    Product `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	if response.Data.Code != "PROD001" {
		t.Errorf("expected code PROD001, got %s", response.Data.Code)
	}

	if response.Data.Price != 99.99 {
		t.Errorf("expected price 99.99, got %.2f", response.Data.Price)
	}
}

func TestOKResponse_WithComplexData(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"products": []map[string]interface{}{
			{"code": "PROD001", "price": 99.99},
			{"code": "PROD002", "price": 149.99},
		},
		"total": 2,
		"metadata": map[string]string{
			"page": "1",
		},
	}

	OKResponse(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}
}

func TestErrorResponse_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()

	ErrorResponse(w, http.StatusBadRequest, "Invalid request body")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response errorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("expected success to be false")
	}

	if len(response.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(response.Messages))
	}

	if response.Messages[0] != "Invalid request body" {
		t.Errorf("expected message 'Invalid request body', got %s", response.Messages[0])
	}
}

func TestErrorResponse_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()

	ErrorResponse(w, http.StatusInternalServerError, "Database connection failed")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var response errorResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Messages[0] != "Database connection failed" {
		t.Errorf("expected message 'Database connection failed', got %s", response.Messages[0])
	}
}

func TestErrorResponse_NotFound(t *testing.T) {
	w := httptest.NewRecorder()

	ErrorResponse(w, http.StatusNotFound, "Product not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var response errorResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Messages[0] != "Product not found" {
		t.Errorf("expected message 'Product not found', got %s", response.Messages[0])
	}
}

func TestErrorResponse_MultipleStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"BadRequest", http.StatusBadRequest, "Bad request"},
		{"Unauthorized", http.StatusUnauthorized, "Unauthorized"},
		{"Forbidden", http.StatusForbidden, "Forbidden"},
		{"NotFound", http.StatusNotFound, "Not found"},
		{"InternalServerError", http.StatusInternalServerError, "Internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ErrorResponse(w, tt.statusCode, tt.message)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}

			var response errorResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success {
				t.Error("expected success to be false")
			}

			if response.Messages[0] != tt.message {
				t.Errorf("expected message '%s', got '%s'", tt.message, response.Messages[0])
			}
		})
	}
}

func TestErrorResponse_EmptyMessage(t *testing.T) {
	w := httptest.NewRecorder()

	ErrorResponse(w, http.StatusBadRequest, "")

	var response errorResponse
	json.NewDecoder(w.Body).Decode(&response)

	if len(response.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(response.Messages))
	}

	if response.Messages[0] != "" {
		t.Errorf("expected empty message, got '%s'", response.Messages[0])
	}
}

func TestErrorResponse_LongMessage(t *testing.T) {
	w := httptest.NewRecorder()

	longMessage := "This is a very long error message that contains detailed information about what went wrong in the system"

	ErrorResponse(w, http.StatusInternalServerError, longMessage)

	var response errorResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Messages[0] != longMessage {
		t.Errorf("expected long message, got '%s'", response.Messages[0])
	}
}

func TestNewResponse(t *testing.T) {
	data := map[string]string{"key": "value"}

	response := newResponse(true, "Test message", data)

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Message != "Test message" {
		t.Errorf("expected message 'Test message', got %s", response.Message)
	}

	if response.Data == nil {
		t.Error("expected data to be present")
	}
}

func TestNewResponse_WithNilData(t *testing.T) {
	response := newResponse(false, "Error message", nil)

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Message != "Error message" {
		t.Errorf("expected message 'Error message', got %s", response.Message)
	}
}

func TestNewResponse_WithDifferentTypes(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{"string", "test"},
		{"int", 123},
		{"float", 99.99},
		{"bool", true},
		{"slice", []string{"a", "b"}},
		{"map", map[string]int{"count": 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := newResponse(true, "Success", tt.data)

			if !response.Success {
				t.Error("expected success to be true")
			}

			if response.Data == nil {
				t.Error("expected data to be present")
			}
		})
	}
}

func TestOKResponse_VerifyJSONStructure(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"id":   1,
		"name": "Test",
	}

	OKResponse(w, data)

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if _, ok := result["success"]; !ok {
		t.Error("expected 'success' field in response")
	}

	if _, ok := result["message"]; !ok {
		t.Error("expected 'message' field in response")
	}

	if _, ok := result["data"]; !ok {
		t.Error("expected 'data' field in response")
	}
}

func TestErrorResponse_VerifyJSONStructure(t *testing.T) {
	w := httptest.NewRecorder()

	ErrorResponse(w, http.StatusBadRequest, "Test error")

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if _, ok := result["success"]; !ok {
		t.Error("expected 'success' field in response")
	}

	if _, ok := result["messages"]; !ok {
		t.Error("expected 'messages' field in response")
	}

	if _, ok := result["data"]; ok {
		t.Error("unexpected 'data' field in error response")
	}
}

func TestOKResponse_ContentTypeSetBeforeWrite(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{"test": "value"}

	OKResponse(w, data)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
}

func TestErrorResponse_ContentTypeSetBeforeWrite(t *testing.T) {
	w := httptest.NewRecorder()

	ErrorResponse(w, http.StatusBadRequest, "Error")

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
}
