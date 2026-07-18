package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"project/internal/app"
	"project/internal/config"
	"project/internal/helpers"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
	"github.com/dracory/userstore"
)

// CartItem represents an item in the cart
type CartItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       string `json:"price"`
	Quantity    int    `json:"quantity"`
	ImageURL    string `json:"image_url"`
}

// Cart represents the user's shopping cart
type Cart struct {
	Items []CartItem `json:"items"`
}

const (
	maxQuantity = 999     // Maximum quantity per item to prevent abuse
	maxItems    = 100     // Maximum number of items in cart to prevent abuse
	maxBodySize = 1 << 20 // 1MB maximum request body size
)

// validateProductAndGetPrice validates that a product exists and returns its actual price
// This prevents price manipulation by ensuring the price comes from the shop store
func (controller *cartController) validateProductAndGetPrice(ctx context.Context, productID string) (string, string, string, error) {
	if controller.app.GetShopStore() == nil {
		return "", "", "", fmt.Errorf("shop store not available")
	}

	product, err := controller.app.GetShopStore().ProductFindByID(ctx, productID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return "", "", "", fmt.Errorf("product not found")
	}

	// Get product media for image URL
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	mediaQuery.SetLimit(1)
	medias, err := controller.app.GetShopStore().MediaList(ctx, mediaQuery)
	if err != nil {
		controller.app.GetLogger().Warn("Failed to get product media", slog.String("error", err.Error()))
	}

	var imageURL string
	if len(medias) > 0 {
		imageURL = medias[0].GetURL()
	}

	return product.GetTitle(), product.GetPrice(), imageURL, nil
}

// getCartFromCache retrieves cart from cache for guest users
func (controller *cartController) getCartFromCache(ctx context.Context, r *http.Request) Cart {
	if controller.app.GetCacheStore() == nil {
		controller.app.GetLogger().Warn("getCartFromCache: Cache store not available")
		return Cart{Items: []CartItem{}}
	}

	cacheKey := helpers.GenerateCartCacheKey(r)
	controller.app.GetLogger().Info("getCartFromCache: Generated cache key", slog.String("cache_key", cacheKey))

	var cart Cart
	cartData, err := controller.app.GetCacheStore().GetJSON(cacheKey, "")
	if err != nil {
		controller.app.GetLogger().Warn("getCartFromCache: Failed to get cart from cache", slog.String("error", err.Error()))
		return Cart{Items: []CartItem{}}
	}

	if cartData != nil {
		if data, ok := cartData.(map[string]any); ok {
			if items, ok := data["items"].([]any); ok {
				for _, item := range items {
					if itemMap, ok := item.(map[string]any); ok {
						cart.Items = append(cart.Items, CartItem{
							ProductID:   itemMap["product_id"].(string),
							ProductName: itemMap["product_name"].(string),
							Price:       itemMap["price"].(string),
							Quantity:    int(itemMap["quantity"].(float64)),
							ImageURL:    itemMap["image_url"].(string),
						})
					}
				}
			}
		}
	}

	return cart
}

// saveCartToCache saves cart to cache for guest users
func (controller *cartController) saveCartToCache(ctx context.Context, r *http.Request, cart Cart) error {
	if controller.app.GetCacheStore() == nil {
		return fmt.Errorf("cache store not available")
	}

	cacheKey := helpers.GenerateCartCacheKey(r)
	err := controller.app.GetCacheStore().SetJSON(cacheKey, cart, 30*24*60*60) // 30 days in seconds
	if err != nil {
		return err
	}

	return nil
}

// TransferCacheToUser transfers cart from cache to user metadata
func (controller *cartController) TransferCacheToUser(ctx context.Context, r *http.Request, user userstore.UserInterface) error {
	if controller.app.GetCacheStore() == nil {
		return fmt.Errorf("cache store not available")
	}

	// Get cart from cache
	cacheCart := controller.getCartFromCache(ctx, r)
	controller.app.GetLogger().Info("TransferCacheToUser: Retrieved cart from cache", slog.Int("item_count", len(cacheCart.Items)))

	if len(cacheCart.Items) == 0 {
		controller.app.GetLogger().Info("TransferCacheToUser: No items in cache cart, skipping transfer")
		return nil // No items to transfer
	}

	// Get existing user cart
	userCart := controller.getCartFromUser(user)
	controller.app.GetLogger().Info("TransferCacheToUser: Retrieved user cart", slog.Int("item_count", len(userCart.Items)))

	// Merge carts (cache cart takes precedence for quantities)
	// Use slice for deterministic ordering
	var mergedItems []CartItem
	itemMap := make(map[string]int) // Maps productID to index in mergedItems

	// Add user cart items first
	for _, item := range userCart.Items {
		if _, exists := itemMap[item.ProductID]; !exists {
			itemMap[item.ProductID] = len(mergedItems)
			mergedItems = append(mergedItems, item)
		}
	}

	// Add or update with cache cart items
	for _, item := range cacheCart.Items {
		if idx, exists := itemMap[item.ProductID]; exists {
			// Update existing item quantity
			mergedItems[idx].Quantity += item.Quantity
		} else {
			// Add new item
			itemMap[item.ProductID] = len(mergedItems)
			mergedItems = append(mergedItems, item)
		}
	}

	// Create merged cart with deterministic ordering
	mergedCart := Cart{Items: mergedItems}
	controller.app.GetLogger().Info("TransferCacheToUser: Merged cart", slog.Int("total_items", len(mergedCart.Items)))

	// Save to user metadata and sync with cache
	if err := controller.saveCartToUserWithCache(ctx, r, user, mergedCart); err != nil {
		controller.app.GetLogger().Error("TransferCacheToUser: Failed to save merged cart", slog.String("error", err.Error()))
		return err
	}

	controller.app.GetLogger().Info("TransferCacheToUser: Cart saved successfully to user metadata")
	// Keep cache in sync with user meta - both will be used
	return nil
}

// NewCartController creates a new cart controller
func NewCartController(app app.AppInterface) *cartController {
	return &cartController{
		app: app,
	}
}

// ClearCart clears the cart from both user metadata and cache
func (controller *cartController) ClearCart(ctx context.Context, r *http.Request, user userstore.UserInterface) error {
	// Clear from user metadata
	if err := controller.saveCartToUser(user, Cart{Items: []CartItem{}}); err != nil {
		return err
	}

	// Clear from cache
	if controller.app.GetCacheStore() != nil {
		cacheKey := helpers.GenerateCartCacheKey(r)
		controller.app.GetCacheStore().SetJSON(cacheKey, Cart{Items: []CartItem{}}, 30*24*60*60)
	}

	return nil
}

type cartController struct {
	app app.AppInterface
}

// Handler handles cart API requests
func (controller *cartController) Handler(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	// Check if user is authenticated
	authUser := helpers.GetAuthUser(r)

	// Handle different actions based on method
	switch r.Method {
	case http.MethodGet:
		if authUser != nil {
			return controller.handleGetCart(ctx, w, r, authUser)
		}
		return controller.handleGetCartGuest(ctx, w, r)
	case http.MethodPost:
		if authUser != nil {
			return controller.handleAddToCart(ctx, w, r, authUser)
		}
		return controller.handleAddToCartGuest(ctx, w, r)
	case http.MethodDelete:
		if authUser != nil {
			return controller.handleRemoveFromCart(ctx, w, r, authUser)
		}
		return controller.handleRemoveFromCartGuest(ctx, w, r)
	case http.MethodPut:
		if authUser != nil {
			return controller.handleUpdateCart(ctx, w, r, authUser)
		}
		return controller.handleUpdateCartGuest(ctx, w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"status":"error","message":"Method not allowed"}`))
		return ""
	}
}

// handleGetCart retrieves the user's cart from metadata
func (controller *cartController) handleGetCart(ctx context.Context, w http.ResponseWriter, r *http.Request, authUser userstore.UserInterface) string {
	cart := controller.getCartFromUser(authUser)

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Cart retrieved successfully", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleGetCartGuest retrieves the guest cart from cache
func (controller *cartController) handleGetCartGuest(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	cart := controller.getCartFromCache(ctx, r)

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Cart retrieved successfully", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleAddToCart adds an item to the user's cart
func (controller *cartController) handleAddToCart(ctx context.Context, w http.ResponseWriter, r *http.Request, authUser userstore.UserInterface) string {
	var reqBody struct {
		ProductID   string `json:"product_id"`
		ProductName string `json:"product_name"`
		Price       string `json:"price"`
		Quantity    int    `json:"quantity"`
		ImageURL    string `json:"image_url"`
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid request body"}`))
		return ""
	}

	// Validate required fields
	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Product ID is required"}`))
		return ""
	}

	if reqBody.Quantity <= 0 {
		reqBody.Quantity = 1
	}
	if reqBody.Quantity > maxQuantity {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Quantity exceeds maximum limit"}`))
		return ""
	}

	// Validate product exists and get actual price from shop store
	productName, actualPrice, imageURL, err := controller.validateProductAndGetPrice(ctx, reqBody.ProductID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid product"}`))
		return ""
	}

	// Use validated product data instead of client-provided data
	reqBody.ProductName = productName
	reqBody.Price = actualPrice
	if imageURL != "" {
		reqBody.ImageURL = imageURL
	}

	// Get current cart
	cart := controller.getCartFromUser(authUser)

	// Check if cart size limit is reached (for new items)
	itemExists := false
	for i, item := range cart.Items {
		if item.ProductID == reqBody.ProductID {
			// Update quantity
			newQuantity := item.Quantity + reqBody.Quantity
			if newQuantity > maxQuantity {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"Quantity exceeds maximum limit"}`))
				return ""
			}
			cart.Items[i].Quantity = newQuantity
			itemExists = true
			break
		}
	}

	// Add new item if it doesn't exist
	if !itemExists {
		if len(cart.Items) >= maxItems {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status":"error","message":"Cart has reached maximum item limit"}`))
			return ""
		}
		cart.Items = append(cart.Items, CartItem{
			ProductID:   reqBody.ProductID,
			ProductName: reqBody.ProductName,
			Price:       reqBody.Price,
			Quantity:    reqBody.Quantity,
			ImageURL:    reqBody.ImageURL,
		})
	}

	// Save cart to user metadata and sync with cache
	if err := controller.saveCartToUserWithCache(ctx, r, authUser, cart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Failed to save cart"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Item added to cart", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleAddToCartGuest adds an item to the guest cart in cache
func (controller *cartController) handleAddToCartGuest(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	var reqBody struct {
		ProductID   string `json:"product_id"`
		ProductName string `json:"product_name"`
		Price       string `json:"price"`
		Quantity    int    `json:"quantity"`
		ImageURL    string `json:"image_url"`
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid request body"}`))
		return ""
	}

	// Validate required fields
	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Product ID is required"}`))
		return ""
	}

	if reqBody.Quantity <= 0 {
		reqBody.Quantity = 1
	}
	if reqBody.Quantity > maxQuantity {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Quantity exceeds maximum limit"}`))
		return ""
	}

	// Validate product exists and get actual price from shop store
	productName, actualPrice, imageURL, err := controller.validateProductAndGetPrice(ctx, reqBody.ProductID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid product"}`))
		return ""
	}

	// Use validated product data instead of client-provided data
	reqBody.ProductName = productName
	reqBody.Price = actualPrice
	if imageURL != "" {
		reqBody.ImageURL = imageURL
	}

	// Get current cart from cache
	cart := controller.getCartFromCache(ctx, r)

	// Check if cart size limit is reached (for new items)
	itemExists := false
	for i, item := range cart.Items {
		if item.ProductID == reqBody.ProductID {
			// Update quantity
			newQuantity := item.Quantity + reqBody.Quantity
			if newQuantity > maxQuantity {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"Quantity exceeds maximum limit"}`))
				return ""
			}
			cart.Items[i].Quantity = newQuantity
			itemExists = true
			break
		}
	}

	// Add new item if it doesn't exist
	if !itemExists {
		if len(cart.Items) >= maxItems {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status":"error","message":"Cart has reached maximum item limit"}`))
			return ""
		}
		cart.Items = append(cart.Items, CartItem{
			ProductID:   reqBody.ProductID,
			ProductName: reqBody.ProductName,
			Price:       reqBody.Price,
			Quantity:    reqBody.Quantity,
			ImageURL:    reqBody.ImageURL,
		})
	}

	// Save cart to cache
	if err := controller.saveCartToCache(ctx, r, cart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Failed to save cart"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Item added to cart", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleRemoveFromCart removes an item from the user's cart
func (controller *cartController) handleRemoveFromCart(ctx context.Context, w http.ResponseWriter, r *http.Request, authUser userstore.UserInterface) string {
	var reqBody struct {
		ProductID string `json:"product_id"`
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid request body"}`))
		return ""
	}

	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Product ID is required"}`))
		return ""
	}

	// Get current cart
	cart := controller.getCartFromUser(authUser)

	// Remove item from cart
	var updatedItems []CartItem
	for _, item := range cart.Items {
		if item.ProductID != reqBody.ProductID {
			updatedItems = append(updatedItems, item)
		}
	}
	cart.Items = updatedItems

	// Save cart to user metadata and sync with cache
	if err := controller.saveCartToUserWithCache(ctx, r, authUser, cart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Failed to save cart"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Item removed from cart", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleRemoveFromCartGuest removes an item from the guest cart in cache
func (controller *cartController) handleRemoveFromCartGuest(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	var reqBody struct {
		ProductID string `json:"product_id"`
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid request body"}`))
		return ""
	}

	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Product ID is required"}`))
		return ""
	}

	// Get current cart from cache
	cart := controller.getCartFromCache(ctx, r)

	// Remove item from cart
	var updatedItems []CartItem
	for _, item := range cart.Items {
		if item.ProductID != reqBody.ProductID {
			updatedItems = append(updatedItems, item)
		}
	}
	cart.Items = updatedItems

	// Save cart to cache
	if err := controller.saveCartToCache(ctx, r, cart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Failed to save cart"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Item removed from cart", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleUpdateCart updates the quantity of an item in the cart
func (controller *cartController) handleUpdateCart(ctx context.Context, w http.ResponseWriter, r *http.Request, authUser userstore.UserInterface) string {
	var reqBody struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid request body"}`))
		return ""
	}

	if reqBody.ProductID == "" || reqBody.Quantity < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Product ID and quantity are required"}`))
		return ""
	}
	if reqBody.Quantity > maxQuantity {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Quantity exceeds maximum limit"}`))
		return ""
	}

	// Get current cart
	cart := controller.getCartFromUser(authUser)

	// Update item quantity or remove if quantity is 0
	var updatedItems []CartItem
	for _, item := range cart.Items {
		if item.ProductID == reqBody.ProductID {
			if reqBody.Quantity > 0 {
				item.Quantity = reqBody.Quantity
				updatedItems = append(updatedItems, item)
			}
			// If quantity is 0, don't add to updatedItems (removes it)
		} else {
			updatedItems = append(updatedItems, item)
		}
	}
	cart.Items = updatedItems

	// Save cart to user metadata and sync with cache
	if err := controller.saveCartToUserWithCache(ctx, r, authUser, cart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Failed to save cart"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Cart updated", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// handleUpdateCartGuest updates the quantity of an item in the guest cart
func (controller *cartController) handleUpdateCartGuest(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	var reqBody struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Invalid request body"}`))
		return ""
	}

	if reqBody.ProductID == "" || reqBody.Quantity < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Product ID and quantity are required"}`))
		return ""
	}
	if reqBody.Quantity > maxQuantity {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","message":"Quantity exceeds maximum limit"}`))
		return ""
	}

	// Get current cart from cache
	cart := controller.getCartFromCache(ctx, r)

	// Update item quantity or remove if quantity is 0
	var updatedItems []CartItem
	for _, item := range cart.Items {
		if item.ProductID == reqBody.ProductID {
			if reqBody.Quantity > 0 {
				item.Quantity = reqBody.Quantity
				updatedItems = append(updatedItems, item)
			}
			// If quantity is 0, don't add to updatedItems (removes it)
		} else {
			updatedItems = append(updatedItems, item)
		}
	}
	cart.Items = updatedItems

	// Save cart to cache
	if err := controller.saveCartToCache(ctx, r, cart); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Failed to save cart"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	response := api.SuccessWithData("Cart updated", map[string]any{
		"cart": cart,
	})
	w.Write([]byte(response.ToString()))
	return ""
}

// getCartFromUser retrieves the cart from user metadata
func (controller *cartController) getCartFromUser(user userstore.UserInterface) Cart {
	cartJSON := user.GetMeta(config.USER_META_CART)
	if cartJSON == "" {
		return Cart{Items: []CartItem{}}
	}

	var cart Cart
	if err := json.Unmarshal([]byte(cartJSON), &cart); err != nil {
		return Cart{Items: []CartItem{}}
	}

	return cart
}

// saveCartToUser saves the cart to user metadata
func (controller *cartController) saveCartToUser(user userstore.UserInterface, cart Cart) error {
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		controller.app.GetLogger().Error("Failed to marshal cart", slog.String("error", err.Error()))
		return err
	}

	if err := user.SetMeta(config.USER_META_CART, string(cartJSON)); err != nil {
		controller.app.GetLogger().Error("Failed to set cart metadata", slog.String("error", err.Error()))
		return err
	}

	// Persist the user metadata to the database
	if err := controller.app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		controller.app.GetLogger().Error("Failed to update user with cart metadata", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// saveCartToUserWithCache saves the cart to user metadata and syncs with cache
func (controller *cartController) saveCartToUserWithCache(ctx context.Context, r *http.Request, user userstore.UserInterface, cart Cart) error {
	if err := controller.saveCartToUser(user, cart); err != nil {
		return err
	}

	// Sync with cache for guest users or when cache is available
	if controller.app.GetCacheStore() != nil {
		cacheKey := helpers.GenerateCartCacheKey(r)
		if err := controller.app.GetCacheStore().SetJSON(cacheKey, cart, 30*24*60*60); err != nil {
			controller.app.GetLogger().Warn("Failed to sync cart to cache", slog.String("error", err.Error()))
		}
	}

	return nil
}
