package register

import (
	"github.com/dracory/geostore"
	"github.com/dracory/userstore"
)

type registerControllerData struct {
	action             string
	authUser           userstore.UserInterface
	email              string
	firstName          string
	lastName           string
	buinessName        string
	phone              string
	country            string
	timezone           string
	countryList        []geostore.Country
	formErrorMessage   string
	formSuccessMessage string
	formRedirectURL    string
}
