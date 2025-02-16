package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

const baseURL = "http://localhost:8089"

func TestSignUpHandler(t *testing.T) {
	rand.Seed(uint64(time.Now().UnixNano()))
	randomNumber := rand.Intn(9000) + 1000
	payload := map[string]interface{}{
		"name":        "John Doe",
		"email":       fmt.Sprintf("johndoe%d@example.com", randomNumber),
		"password":    "StrongPass123!",
		"cpf":         "12345678909",
		"carPlate":    "ABC1234",
		"isPassenger": true,
		"isDriver":    false,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", baseURL+"/signup", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAccountHandler_NotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", baseURL+"/accounts/invalid-id", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
