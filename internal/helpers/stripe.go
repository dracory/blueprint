package helpers

import (
	"fmt"
	"strings"
	"sync"

	stripe "github.com/stripe/stripe-go/v73"
	stripeSession "github.com/stripe/stripe-go/v73/checkout/session"
)

var (
	stripeKeyOnce sync.Once
	stripeKey     string
)

// GenerateStripePaymentsCheckoutURLOptions options for "Stripe Payments"
type GenerateStripePaymentsCheckoutURLOptions struct {
	CustomerEmail     string
	AmountToPay       float64
	UrlPaymentCancel  string
	UrlPaymentSuccess string
	StripeKeyPrivate  string
	LineItems         []LineItem
	OrderID           string
}

type LineItem struct {
	Name     string
	Quantity int
	Price    float64
}

// GenerateStripePaymentsCheckoutURL - generates checkout URL for "Stripe Payments"
func GenerateStripePaymentsCheckoutURL(options GenerateStripePaymentsCheckoutURLOptions) (string, error) {
	// Initialize stripe key once to prevent race condition with global stripe.Key variable
	stripeKeyOnce.Do(func() {
		stripeKey = options.StripeKeyPrivate
		stripe.Key = stripeKey
	})

	priceCurrency := "GBP"

	// Build line items from cart items
	lineItems := []*stripe.CheckoutSessionLineItemParams{}
	if len(options.LineItems) > 0 {
		for _, item := range options.LineItems {
			productName := item.Name
			if options.OrderID != "" {
				productName = fmt.Sprintf("%s (Order: %s)", item.Name, options.OrderID)
			}
			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(strings.ToLower(priceCurrency)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(int64(item.Price * 100)),
				},
				Quantity: stripe.Int64(int64(item.Quantity)),
			})
		}
	} else {
		// Fallback to single line item if no cart items provided
		productTitle := fmt.Sprintf("Order: %s", options.OrderID)
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(strings.ToLower(priceCurrency)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(productTitle),
				},
				UnitAmount: stripe.Int64(int64(options.AmountToPay * 100)),
			},
			Quantity: stripe.Int64(1),
		})
	}

	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(options.CustomerEmail),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  lineItems,
		SuccessURL: stripe.String(options.UrlPaymentSuccess),
		CancelURL:  stripe.String(options.UrlPaymentCancel),
	}

	checkoutSession, err := stripeSession.New(params)

	if err != nil {
		return "", err
	}

	return checkoutSession.URL, nil
}
