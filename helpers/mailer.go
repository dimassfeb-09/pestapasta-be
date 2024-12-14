package helpers

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"time"

	"github.com/dimassfeb-09/pestapasta-be/models"
	"github.com/dimassfeb-09/pestapasta-be/utils"
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
  </head>
  <body style="font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; padding: 20px; background-color: #f4f4f4; color: #333;">
    <div style="font-size: 16px; gap: 20px; align-items: flex-start; border-bottom: 2px solid #007bff; padding-bottom: 15px; margin-bottom: 20px; color: #007bff;">
      <div>
        <h1 style="margin: 0; font-size: 2.5rem;">INVOICE</h1>
        <p style="margin: 0;">Invoice Date: {{.Date}}</p>
        <p style="margin: 0;">Invoice Number: {{.InvoiceNumber}}</p>
      </div>
      <div>
        <h3 style="margin: 0;">FROM:</h3>
        <p style="margin: 0;">{{.CompanyName}}</p>
        <p style="margin: 0;">{{.CompanyEmail}}</p>
        <p style="margin: 0;">{{.CompanyPhone}}</p>
      </div>
    </div>

    <div style="margin-top: 20px; padding-top: 20px; background: #ffffff; border-top: 3px solid #007bff; border-radius: 8px; padding: 20px;">
      <h3 style="margin: 0 0 10px; color: #007bff;">Payment Information</h3>
      {{if eq .PaymentMethod "qris"}}
      <div style="padding: 10px; background-color: #f9f9f9; border: 1px solid #ddd; border-radius: 5px;">
        <h4 style="margin-top: 0; color: #333;">QRIS Payment</h4>
        <p style="margin: 0;">Scan the QR code to complete your payment.</p>
        <img
          src="{{.PaymentQRCodeURL}}"
          alt="QRIS QR Code"
          style="display: block; margin: 10px auto; border: 1px solid #ddd; border-radius: 8px; height: 120px; width: 120px;"
        />
      </div>
      {{else}}
      <div style="padding: 10px; background-color: #f9f9f9; border: 1px solid #ddd; border-radius: 5px;">
        <h4 style="margin-top: 0; color: #333;">Bank Transfer</h4>
        <p style="margin: 0;">Account Name: {{.PaymentAccountName}}</p>
        <p style="margin: 0;">Account Number: {{.PaymentAccountNumber}}</p>
      </div>
      {{end}}
    </div>

    <div style="background: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);">
      <h3 style="margin: 0 0 20px;">BILL TO:</h3>
      <p style="margin: 0;">{{.ClientName}}</p>
      <p style="margin: 0;">{{.ClientEmail}}</p>

      <table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
        <thead>
          <tr>
            <th style="border: 1px solid #ddd; padding: 8px; text-align: left; background-color: #007bff; color: #ffffff;">Product/Service</th>
            <th style="border: 1px solid #ddd; padding: 8px; text-align: left; background-color: #007bff; color: #ffffff;">Quantity</th>
            <th style="border: 1px solid #ddd; padding: 8px; text-align: left; background-color: #007bff; color: #ffffff;">Unit Price</th>
            <th style="border: 1px solid #ddd; padding: 8px; text-align: left; background-color: #007bff; color: #ffffff;">Total</th>
          </tr>
        </thead>
        <tbody>
          {{range .Items}}
          <tr>
            <td style="border: 1px solid #ddd; padding: 8px;">{{.ProductName}}</td>
            <td style="border: 1px solid #ddd; padding: 8px;">{{.Quantity}}</td>
            <td style="border: 1px solid #ddd; padding: 8px;">Rp{{.UnitPrice}}</td>
            <td style="border: 1px solid #ddd; padding: 8px;">Rp{{.TotalPrice}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>

      <div style="text-align: right; margin-top: 20px; font-weight: bold;">
        <p style="margin: 0; font-size: 1rem;">Subtotal: Rp{{.Subtotal}}</p>
        <p style="margin: 0; font-size: 1rem;">Tax (if applicable): Rp{{.Tax}}</p>
        <h3 style="color: #007bff; font-size: 1.5rem; margin-top: 10px;">Total: Rp{{.Total}}</h3>
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
	d := gomail.NewDialer("smtp.gmail.com", 587, utils.GetENV().Email.User, utils.GetENV().Email.Password)
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
