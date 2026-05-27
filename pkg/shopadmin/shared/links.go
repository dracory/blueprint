package shared

import "project/internal/links"

type Links struct {
	baseURL string
}

func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/shop"
	}
	return &Links{baseURL: baseURL}
}

func (l *Links) Home(params map[string]string) string {
	params["controller"] = CONTROLLER_HOME
	return links.Admin().Shop(params)
}

func (l *Links) Products(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCTS
	return links.Admin().Shop(params)
}

func (l *Links) ProductCreate(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCT_CREATE
	return links.Admin().Shop(params)
}

func (l *Links) ProductDelete(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCT_DELETE
	return links.Admin().Shop(params)
}

func (l *Links) ProductUpdate(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCT_UPDATE
	return links.Admin().Shop(params)
}

func (l *Links) Categories(params map[string]string) string {
	params["controller"] = CONTROLLER_CATEGORIES
	return links.Admin().Shop(params)
}

func (l *Links) CategoryCreate(params map[string]string) string {
	params["controller"] = CONTROLLER_CATEGORY_CREATE
	return links.Admin().Shop(params)
}

func (l *Links) CategoryUpdate(params map[string]string) string {
	params["controller"] = CONTROLLER_CATEGORY_UPDATE
	return links.Admin().Shop(params)
}

func (l *Links) Discounts(params map[string]string) string {
	params["controller"] = CONTROLLER_DISCOUNTS
	return links.Admin().Shop(params)
}

func (l *Links) Orders(params map[string]string) string {
	params["controller"] = CONTROLLER_ORDERS
	return links.Admin().Shop(params)
}
