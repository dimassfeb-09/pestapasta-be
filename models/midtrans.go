package models

type CreateTransactionMidtransResponse struct {
	StatusCode        string   `json:"status_code"`
	StatusMessage     string   `json:"status_message"`
	TransactionID     string   `json:"transaction_id"`
	OrderID           string   `json:"order_id"`
	MerchantID        string   `json:"merchant_id"`
	GrossAmount       string   `json:"gross_amount"`
	Currency          string   `json:"currency"`
	PaymentType       string   `json:"payment_type"`
	TransactionTime   string   `json:"transaction_time"`
	TransactionStatus string   `json:"transaction_status"`
	FraudStatus       string   `json:"fraud_status"`
	Actions           []Action `json:"actions"`
	Acquirer          string   `json:"acquirer"`
	QRString          string   `json:"qr_string"`
	ExpiryTime        string   `json:"expiry_time"`
}

type CreateTransactionMidtransResponseWithError struct {
	ID            string `json:"id"`
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

type ItemDetails struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Name     string  `json:"name"`
}

type CreateTransactionMidtransPayload struct {
	PaymentType        string `json:"payment_type"`
	TransactionDetails struct {
		OrderID     string  `json:"order_id"`
		GrossAmount float64 `json:"gross_amount"`
	} `json:"transaction_details"`
	ItemDetails     []ItemDetails `json:"item_details"`
	CustomerDetails struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"customer_details"`
	QRIS struct {
		Acquirer string `json:"acquirer"`
	} `json:"qris"`
}

type StatusTransactionMidtransResponse struct {
	StatusCode        string `json:"status_code"`
	TransactionID     string `json:"transaction_id"`
	GrossAmount       string `json:"gross_amount"`
	Currency          string `json:"currency"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	StatusMessage     string `json:"status_message"`
	MerchantID        string `json:"merchant_id"`
	TransactionTime   string `json:"transaction_time"`
	ExpiryTime        string `json:"expiry_time"`
}
