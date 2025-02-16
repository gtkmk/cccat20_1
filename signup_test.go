package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type SignUpRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CPF         string `json:"cpf"`
	CarPlate    string `json:"carPlate"`
	IsPassenger bool   `json:"isPassenger"`
	IsDriver    bool   `json:"isDriver"`
}

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/signup", SignUpHandler)

	tests := []struct {
		name     string
		input    SignUpRequest
		expected int
	}{
		{
			"Valid Signup",
			SignUpRequest{"John Doe", "johndoe@example.com", "Password1", "12345678901", "ABC1234", true, false},
			http.StatusOK,
		},
		{
			"Invalid Email",
			SignUpRequest{"John Doe", "invalid-email", "Password1", "12345678901", "ABC1234", true, false},
			http.StatusUnprocessableEntity,
		},
		{
			"Weak Password",
			SignUpRequest{"John Doe", "johndoe@example.com", "123456", "12345678901", "ABC1234", true, false},
			http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code)
		})
	}
}

func TestGetAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/accounts/:accountId", GetAccountHandler)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/non-existing-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
