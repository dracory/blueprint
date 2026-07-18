package helpers

import (
	"fmt"
	"log/slog"
	"project/internal/app"
)

const PAYMENT_SUCCESS = "success"
const PAYMENT_CANCELED = "canceled"
const IDEMPOTENCY_CACHE_KEY_PREFIX = "payment_begin_"

// PaymentCodeData represents the data required for processing a payment.
type PaymentCodeData struct {
	PaymentSuccessOrCanceled string // "success" or "canceled"
	OrderID                  string
	BuyerID                  string
	Amount                   string
	DiscountID               string
	IdempotencyKey           string // Key to prevent duplicate order creation
}

// generatePaymentCodeKey generates a payment code key by concatenating the input
// string with "_order_payment_code".
//
// code: a string representing the base code.
// Returns a string representing the generated payment code key.
func generatePaymentCodeKey(code string) string {
	return code + "_order_payment_code"
}

// DeletePaymentCodeData removes the PaymentCodeData associated with the given code from the cache store.
//
// Parameters:
// - app: the app interface
// - code: string representing the code of the PaymentCodeData to be removed.
//
// Returns an error in case of any issues while removing the data.
func DeletePaymentCodeData(app app.AppInterface, code string) error {
	key := generatePaymentCodeKey(code)
	return app.GetCacheStore().Remove(key)
}

// GetPaymentCodeData retrieves the payment code data based on the provided code.
//
// Parameters:
// - app: the app interface
// - code: string representing the payment code.
// Returns a PaymentCodeData struct and an error if applicable.
func GetPaymentCodeData(app app.AppInterface, code string) (data PaymentCodeData, err error) {
	key := generatePaymentCodeKey(code)
	app.GetLogger().Info("GetPaymentCodeData: Retrieving payment code data", slog.String("code", code), slog.String("cache_key", key))
	dataAny, err := app.GetCacheStore().GetJSON(key, nil)
	if err != nil {
		app.GetLogger().Error("GetPaymentCodeData: Failed to retrieve from cache", slog.String("error", err.Error()))
		return data, err
	}
	if dataAny == nil {
		app.GetLogger().Error("GetPaymentCodeData: Cache returned nil")
		return data, fmt.Errorf("cache returned nil")
	}

	app.GetLogger().Info("GetPaymentCodeData: Cache returned data", slog.String("data_type", fmt.Sprintf("%T", dataAny)))

	// Convert dataAny to map[string]string manually
	dataMap := make(map[string]string)
	dataMapTyped, ok := dataAny.(map[string]string)
	if ok {
		dataMap = dataMapTyped
	} else {
		// Try map[string]any as fallback
		dataMapAny, okAny := dataAny.(map[string]any)
		if !okAny {
			app.GetLogger().Error("GetPaymentCodeData: Data is not map[string]string or map[string]any", slog.String("data_type", fmt.Sprintf("%T", dataAny)))
			return data, nil
		}
		// Convert map[string]any to map[string]string
		for k, v := range dataMapAny {
			if strVal, ok := v.(string); ok {
				dataMap[k] = strVal
			}
		}
	}

	app.GetLogger().Info("GetPaymentCodeData: Data map keys", slog.Int("key_count", len(dataMap)))

	// Safely access map keys with checks
	if paymentStatus, exists := dataMap["payment_status"]; exists {
		data.PaymentSuccessOrCanceled = paymentStatus
	}
	if orderID, exists := dataMap["order_id"]; exists {
		data.OrderID = orderID
	}
	if buyerID, exists := dataMap["buyer_id"]; exists {
		data.BuyerID = buyerID
	}
	if amount, exists := dataMap["amount"]; exists {
		data.Amount = amount
	}
	if discountID, exists := dataMap["discount_id"]; exists {
		data.DiscountID = discountID
	}
	if idempotencyKey, exists := dataMap["idempotency_key"]; exists {
		data.IdempotencyKey = idempotencyKey
	}

	app.GetLogger().Info("GetPaymentCodeData: Parsed data", slog.String("payment_status", data.PaymentSuccessOrCanceled), slog.String("order_id", data.OrderID), slog.String("idempotency_key", data.IdempotencyKey))

	return data, nil
}

// SetPaymentCodeData generates a unique key from the given code
// and saves the data to the cache store with an expiration time of 24 hours.
//
// Parameters:
// - app: the app interface
// - code: a string representing the payment code.
// - data: a struct containing the payment data.
// Returns an error if the cache store fails to save the data.
func SetPaymentCodeData(app app.AppInterface, code string, data PaymentCodeData) error {
	hoursValid := 24 // 24 hours to account for timezone issues and give users more time
	key := generatePaymentCodeKey(code)
	app.GetLogger().Info("SetPaymentCodeData: Storing payment code data", slog.String("code", code), slog.String("cache_key", key), slog.Int("expiration_seconds", hoursValid*3600))

	values := map[string]string{
		"payment_status":  data.PaymentSuccessOrCanceled,
		"order_id":        data.OrderID,
		"buyer_id":        data.BuyerID,
		"amount":          data.Amount,
		"discount_id":     data.DiscountID,
		"idempotency_key": data.IdempotencyKey,
	}

	// Cache store subtracts expiration time, so use negative value to add time
	err := app.GetCacheStore().SetJSON(key, values, int64(hoursValid*3600))
	if err != nil {
		app.GetLogger().Error("SetPaymentCodeData: Failed to store payment code data", slog.String("error", err.Error()))
	}
	return err
}
