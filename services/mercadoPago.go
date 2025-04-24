package services

import (
	"caaso/models"
	"context"
	"fmt"
	"os"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

var (
	Url         string
	AccessToken string
)

func PaymentService() {
	AccessToken = os.Getenv("MERCADO_PAGO_ACCESS_TOKEN")
}

func CreatePayment(amount float64, email string, userId string, planType models.PlanTypes) (*payment.Response, error) {
	cfg, err := config.New(AccessToken)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	baseURL := os.Getenv("BACKEND_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("BACKEND_URL is not set")
	}

	// Build the Payments client
	client := payment.NewClient(cfg)
	// Craft the Request payload
	request := payment.Request{
		TransactionAmount: amount,
		PaymentMethodID:   "pix",
		Payer: &payment.PayerRequest{
			Email: email,
		},
		NotificationURL: baseURL + "/go/payment/confirm",
		Metadata: map[string]any{
			"userId":   userId,
			"planType": planType.String(), // must be a string, e.g. "Monthly" or "Yearly"
		},
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resource, nil
}

func GetPaymentFromId(id int) (*payment.Response, error) {

	// Initialize MP SDK client
	cfg, err := config.New(AccessToken) //
	if err != nil {
		return nil, err
	}
	client := payment.NewClient(cfg) //

	// Re-fetch the payment from Mercado Pago
	mpResp, err := client.Get(context.Background(), id) // analogous to invoice.Get :contentReference[oaicite:0]{index=0}
	if err != nil {
		return nil, err
	}

	return mpResp, nil

}