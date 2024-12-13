package models

type User struct {
	ID       int     `json:"id" gorm:"primary_key"` // Primary key di tabel users
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Orders   []Order `gorm:"foreignKey:UserID"`
}

// Category represents a food category.
type Category struct {
	ID           int    `json:"id" gorm:"primary_key"`
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}

// PastaMenu represents a pasta menu item.
type Menu struct {
	ID          int     `json:"id" gorm:"primary_key"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Rating      int     `json:"rating"`
}

// Order represents an order placed by a customer.
type Order struct {
	ID           int           `json:"id" gorm:"primary_key"` // Primary key di tabel orders
	OrderDate    string        `json:"order_date"`
	Email        string        `json:"email"`
	Name         string        `json:"name"`
	TotalPrice   float64       `json:"total_price"`
	OrderStatus  string        `json:"order_status"`
	Payment      Payment       `json:"payments" gorm:"foreignKey:OrderID"`
	OrderDetails []OrderDetail `json:"order_details" gorm:"foreignKey:OrderID"`
}

// OrderDetail represents details of a single pasta item in an order.
type OrderDetail struct {
	ID            int     `json:"id" gorm:"primary_key"` // Primary key di tabel order_details
	OrderID       int     `json:"order_id"`              // Foreign key ke tabel orders
	MenuID        int     `json:"menu_id"`
	Quantity      int     `json:"quantity"`
	SubtotalPrice float64 `json:"subtotal_price"`
	Notes         string  `json:"notes"`
	Menu          Menu    `json:"menu_detail" gorm:"foreignKey:MenuID"`
}

// Payment represents payment information for an order.
type Payment struct {
	ID                   int    `gorm:"primary_key"`
	OrderID              int    `json:"order_id"`
	PaymentMethod        string `json:"payment_method"`
	PaymentStatus        string `json:"payment_status"`
	PaymentAccountNumber string `json:"payment_account_number,omitempty"`
	PaymentAccountName   string `json:"payment_account_name,omitempty"`
	PaymentQRCodeURL     string `json:"payment_qrcode_url,omitempty"`
	PaymentCreateDate    string `json:"payment_create_date"`
	PaymentExpiredDate   string `json:"payment_expired_date"`
	PaymentTransactionID string `json:"payment_transaction_id,omitempty"`
	TransactionCode      string `json:"transaction_code"`
}

// CheckoutRequest represents the incoming request for checkout.
type CheckoutRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	PaymentMethodID int    `json:"payment_method_id"`
	Products        []struct {
		ID       int    `json:"id"`
		Quantity int    `json:"quantity"`
		Note     string `json:"note"`
	} `json:"products"`
}

type PaymentDetails struct {
	PaymentAccountNumber string `json:"payment_account_number,omitempty"`
	PaymentAccountName   string `json:"payment_account_name,omitempty"`
	PaymentMethod        string `json:"payment_method"`
	PaymentStatus        string `json:"payment_status"`
	PaymentExpiredTime   int64  `json:"payment_expired_time"`
	QRImageURL           string `json:"qr_image_url,omitempty"`
}

// CheckoutResponse represents the response after checkout.
type CheckoutResponse struct {
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	Total           float64        `json:"total"`
	TransactionCode string         `json:"transaction_code"`
	PaymentDetails  PaymentDetails `json:"payment_details"`
	ProductDetails  []Menu         `json:"product_details"`
}

type PaymentMethod struct {
	ID            int    `json:"id" gorm:"primary_key"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Code          string `json:"code"`
}

type ResponseSuccess struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ResponseSuccessWithData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
