package services

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dimassfeb-09/pestapasta-be/models"
	"github.com/dimassfeb-09/pestapasta-be/utils"
)

var (
	env = utils.GetENV()
)

func CreateTransaction(trx models.CreateTransactionMidtransPayload) (*models.CreateTransactionMidtransResponse, *models.CreateTransactionMidtransResponseWithError, error) {

	// Create Additional Tax 10%
	taxCount := trx.TransactionDetails.GrossAmount * 0.1
	trx.TransactionDetails.GrossAmount += taxCount
	trx.ItemDetails = append(trx.ItemDetails, models.ItemDetails{
		ID:       "Tax-10%",
		Price:    taxCount,
		Quantity: 1,
		Name:     "Pajak 10%",
	})

	// Marshal the transaction payload
	data, err := json.Marshal(&trx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal transaction payload: %w", err)
	}

	fmt.Println(string(data))

	// Create request
	bytesBuffer := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", utils.EndpointMidtransCharge, bytesBuffer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(utils.GetENV().MidtransKey, "")) // Replace with actual server key

	// Configure HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResponse models.CreateTransactionMidtransResponseWithError
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		log.Printf("Error response: %v", errorResponse)
		return nil, &errorResponse, fmt.Errorf("failed to create transaction, status: %d, response: %v", resp.StatusCode, errorResponse)
	}

	// Unmarshal response body into the success struct
	var transactionResponse models.CreateTransactionMidtransResponse
	if err := json.Unmarshal(body, &transactionResponse); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Return success response
	return &transactionResponse, nil, nil
}

func CheckTransaction(transactionId string) (*models.StatusTransactionMidtransResponse, *models.CreateTransactionMidtransResponseWithError, error) {
	fmt.Println(transactionId)
	url := fmt.Sprintf("https://api.midtrans.com/v2/%s/status", transactionId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(env.MidtransKey, "")) // Replace with actual server key

	// Configure HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResponse models.CreateTransactionMidtransResponseWithError
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		log.Printf("Error response: %v", errorResponse)
		return nil, &errorResponse, fmt.Errorf("failed to create transaction, status: %d, response: %v", resp.StatusCode, errorResponse)
	}

	// Unmarshal response body into the success struct
	var statusTransactionResponse models.StatusTransactionMidtransResponse
	if err := json.Unmarshal(body, &statusTransactionResponse); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Return success response
	return &statusTransactionResponse, nil, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
