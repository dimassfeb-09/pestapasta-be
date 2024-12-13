package services

import (
	"log"
	"testing"

	"github.com/dimassfeb-09/pestapasta-be/models"
)

func TestCreateTransactionTest(t *testing.T) {
	transactionBody := models.CreateTransactionMidtransPayload{
		PaymentType: "qris",
		TransactionDetails: struct {
			OrderID     string  "json:\"order_id\""
			GrossAmount float64 "json:\"gross_amount\""
		}{
			OrderID:     "kjasgdausasdasddbu67567yqgwtdasd",
			GrossAmount: 20000,
		},

		ItemDetails: []models.ItemDetails{
			{
				ID:       "1",
				Price:    10000,
				Quantity: 1,
				Name:     "KECOAK",
			},
			{
				ID:       "2",
				Price:    10000,
				Quantity: 1,
				Name:     "KECOAK 2",
			},
		},
		CustomerDetails: struct {
			FirstName string "json:\"first_name\""
			LastName  string "json:\"last_name\""
			Email     string "json:\"email\""
		}{
			FirstName: "Dimas",
			LastName:  "Febriyanto",
			Email:     "dimassfeb@gmail.com",
		},
		QRIS: struct {
			Acquirer string "json:\"acquirer\""
		}{
			Acquirer: "gopay",
		},
	}
	body, errCustom, err := CreateTransaction(transactionBody)
	if err != nil || errCustom != nil {
		log.Println(errCustom.StatusCode)
		log.Fatal(err)
	}

	log.Println(body)
}
