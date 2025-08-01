package shared

import "project/internal/links"

type Links struct{}

func NewLinks() *Links {
	return &Links{}
}

func (*Links) Home(params map[string]string) string {
	params["controller"] = CONTROLLER_HOME
	return links.Admin().Shop(params)
}

func (*Links) Discounts(params map[string]string) string {
	params["controller"] = CONTROLLER_DISCOUNTS
	return links.Admin().Shop(params)
}

func (*Links) Orders(params map[string]string) string {
	params["controller"] = CONTROLLER_ORDERS
	return links.Admin().Shop(params)
}

func (*Links) ProductCreate(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCT_CREATE
	return links.Admin().Shop(params)
}

func (*Links) ProductDelete(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCT_DELETE
	return links.Admin().Shop(params)
}

func (*Links) Products(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCTS
	return links.Admin().Shop(params)
}

func (*Links) ProductUpdate(params map[string]string) string {
	params["controller"] = CONTROLLER_PRODUCT_UPDATE
	return links.Admin().Shop(params)
}

func (*Links) TestLlm(params map[string]string) string {
	params["controller"] = CONTROLLER_TEST_LLM
	return links.Admin().Shop(params)
}
