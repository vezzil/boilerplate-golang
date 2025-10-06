package stripe

import (
	"fmt"
	"sync"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeManager handles Stripe payment operations
type StripeManager struct {
	stripeKey     string
	webhookSecret string
}

var (
	instance *StripeManager
	once     sync.Once
)

// NewStripeManager creates a new instance of StripeManager
func NewStripeManager(apiKey, webhookSecret string) *StripeManager {
	once.Do(func() {
		stripe.Key = apiKey
		instance = &StripeManager{
			stripeKey:     apiKey,
			webhookSecret: webhookSecret,
		}
	})
	return instance
}

// CreateCheckoutSession creates a new Stripe Checkout session
func (sm *StripeManager) CreateCheckoutSession(amount int64, currency, successURL, cancelURL, orderID, customerEmail string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Order #" + orderID),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"order_id": orderID,
		},
	}

	if customerEmail != "" {
		params.CustomerEmail = stripe.String(customerEmail)
	}

	return session.New(params)
}

// CreatePaymentIntent creates a new payment intent
func (sm *StripeManager) CreatePaymentIntent(amount int64, currency, orderID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amount),
		Currency:           stripe.String(currency),
		PaymentMethodTypes: []*string{stripe.String("card")},
		Metadata: map[string]string{
			"order_id": orderID,
		},
	}

	return paymentintent.New(params)
}

// GetPaymentIntent retrieves a payment intent by ID
func (sm *StripeManager) GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	return paymentintent.Get(id, nil)
}

// HandleWebhook handles Stripe webhook events
func (sm *StripeManager) HandleWebhook(payload []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, sm.webhookSecret)
}

// CreateRefund creates a refund for a payment intent
func (sm *StripeManager) CreateRefund(paymentIntentID string, amount int64) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}
	if amount > 0 {
		params.Amount = stripe.Int64(amount)
	}

	refund, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	return refund, nil
}
