package helpers

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"time"

	"github.com/dimassfeb-09/pestapasta-be/models"
	"gopkg.in/gomail.v2"
)

type InvoiceData struct {
	Date                 string
	InvoiceNumber        string
	CompanyName          string
	CompanyEmail         string
	CompanyPhone         string
	ClientName           string
	ClientEmail          string
	Items                []InvoiceItem
	Subtotal             float64
	Tax                  float64
	Total                float64
	PaymentAccountName   string
	PaymentAccountNumber string
	PaymentMethod        string
	PaymentQRCodeURL     string
}

type InvoiceItem struct {
	ProductName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

// Define the HTML template as a constant
const emailTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Invoice</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        line-height: 1.6;
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
        background-color: #f4f4f4;
        color: #333;
      }
      .invoice-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        border-bottom: 2px solid #007bff;
        padding-bottom: 15px;
        margin-bottom: 20px;
        color: #007bff;
      }
      .invoice-header h1 {
        margin: 0;
        font-size: 2.5rem;
      }
      .invoice-header p {
        margin: 0;
      }
      .invoice-details {
        background: #ffffff;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
      }
      table {
        width: 100%;
        border-collapse: collapse;
        margin: 20px 0;
      }
      th,
      td {
        border: 1px solid #ddd;
        padding: 8px;
        text-align: left;
      }
      th {
        background-color: #007bff;
        color: #ffffff;
      }
      .total-section {
        text-align: right;
        margin-top: 20px;
        font-weight: bold;
      }
      .total-section p {
        margin: 0;
        font-size: 1rem;
      }
      .total-section h3 {
        color: #007bff;
        font-size: 1.5rem;
        margin-top: 10px;
      }
      .payment-info {
        margin-top: 20px;
        padding-top: 20px;
        background: #ffffff;
        border-top: 3px solid #007bff;
        border-radius: 8px;
        padding: 20px;
      }
      .payment-info h3 {
        margin: 0 0 10px;
        color: #007bff;
      }
      .qris-info,
      .bank-transfer-info {
        padding: 10px;
        background-color: #f9f9f9;
        border: 1px solid #ddd;
        border-radius: 5px;
      }
      .qris-info h4,
      .bank-transfer-info h4 {
        margin-top: 0;
        color: #333;
      }
      .qris-info img {
        display: block;
        margin: 10px auto;
        border: 1px solid #ddd;
        border-radius: 8px;
      }
    </style>
  </head>
  <body>
    <div class="invoice-header">
      <div>
        <h1>INVOICE</h1>
        <p>Invoice Date: {{.Date}}</p>
        <p>Invoice Number: {{.InvoiceNumber}}</p>
      </div>
      <div>
        <h3>FROM:</h3>
        <p>{{.CompanyName}}</p>
        <p>{{.CompanyEmail}}</p>
        <p>{{.CompanyPhone}}</p>
      </div>
    </div>

    <div class="payment-info">
      <h3>Payment Information</h3>
      {{if eq .PaymentMethod "qris"}}
      <div class="qris-info">
        <h4>QRIS Payment</h4>
        <p>Scan the QR code to complete your payment.</p>
        <img
          src="{{.PaymentQRCodeURL}}"
          alt="QRIS QR Code"
          style="height: 120px; width: 120px"
        />
      </div>
      {{else}}
      <div class="bank-transfer-info">
        <h4>Bank Transfer</h4>
        <p>Account Name: {{.PaymentAccountName}}</p>
        <p>Account Number: {{.PaymentAccountNumber}}</p>
      </div>
      {{end}}
    </div>

    <div class="invoice-details">
      <h3>BILL TO:</h3>
      <p>{{.ClientName}}</p>
      <p>{{.ClientEmail}}</p>

      <table>
        <thead>
          <tr>
            <th>Product/Service</th>
            <th>Quantity</th>
            <th>Unit Price</th>
            <th>Total</th>
          </tr>
        </thead>
        <tbody>
          {{range .Items}}
          <tr>
            <td>{{.ProductName}}</td>
            <td>{{.Quantity}}</td>
            <td>Rp{{.UnitPrice}}</td>
            <td>Rp{{.TotalPrice}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>

      <div class="total-section">
        <p>Subtotal: Rp{{.Subtotal}}</p>
        <p>Tax (if applicable): Rp{{.Tax}}</p>
        <h3>Total: Rp{{.Total}}</h3>
      </div>
    </div>
  </body>
</html>
`

// ConvertCheckoutToInvoiceData transforms CheckoutResponse into InvoiceData
func ConvertCheckoutToInvoiceData(data models.CheckoutResponse) InvoiceData {
	var items []InvoiceItem
	subtotal := 0.0

	for _, product := range data.ProductDetails {
		price := float64(product.Price)
		items = append(items, InvoiceItem{
			ProductName: product.Name,
			Quantity:    1, // Assuming quantity is 1; adjust if needed
			UnitPrice:   price,
			TotalPrice:  price,
		})
		subtotal += price
	}

	tax := subtotal * 0.1 // Example: 10% tax
	total := subtotal + tax

	return InvoiceData{
		Date:                 time.Now().Format("2006-01-02"),
		InvoiceNumber:        data.TransactionCode,
		CompanyName:          "Pesta Pasta",
		CompanyEmail:         "support@pestapasta.com",
		CompanyPhone:         "123-456-7890",
		ClientName:           data.Name,
		ClientEmail:          data.Email,
		Items:                items,
		Subtotal:             subtotal,
		Tax:                  tax,
		Total:                total,
		PaymentAccountName:   data.PaymentDetails.PaymentAccountName,
		PaymentAccountNumber: data.PaymentDetails.PaymentAccountNumber,
		PaymentMethod:        data.PaymentDetails.PaymentMethod,
		PaymentQRCodeURL:     data.PaymentDetails.QRImageURL,
	}
}

// RenderTemplate renders the email template with the provided data
func RenderTemplate(data InvoiceData) string {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return rendered.String()
}

func SendMail(checkoutResponse models.CheckoutResponse) {
	d := gomail.NewDialer("smtp.gmail.com", 587, "dimassfeb@gmail.com", "cdjn usfl pdhw jtbe")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	invoiceData := ConvertCheckoutToInvoiceData(checkoutResponse)
	body := RenderTemplate(invoiceData)

	m := gomail.NewMessage()
	m.SetHeader("From", "dimassfeb@gmail.com")
	m.SetAddressHeader("To", checkoutResponse.Email, checkoutResponse.Name)
	m.SetHeader("Subject", "Invoice for Your Purchase")
	m.SetBody("text/html", body)

	if err := gomail.Send(s, m); err != nil {
		log.Printf("Could not send email to %q: %v", checkoutResponse.Email, err)
	}
}
