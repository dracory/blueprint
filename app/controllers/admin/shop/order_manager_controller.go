package admin

import (
	"net/http"
)

type orderManagerController struct{}

func NewOrderManagerController() *orderManagerController {
	return &orderManagerController{}
}

func (orderManagerController *orderManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "Order Manager"
}
