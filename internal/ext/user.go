package ext

import (
	"log"
	"strings"

	"github.com/dracory/userstore"
	"github.com/samber/lo"
)

func DisplayNameFull(user userstore.UserInterface) string {
	if user == nil {
		return "n/a"
	}

	displayName := user.GetFirstName() + " " + user.GetLastName()

	if strings.TrimSpace(displayName) == "" {
		return user.GetEmail()
	}

	return displayName
}

func IsClient(user userstore.UserInterface) bool {
	if user == nil {
		return false
	}
	return user.GetMeta("is_client") == "yes"
}

func SetIsClient(user userstore.UserInterface, isClient bool) userstore.UserInterface {
	if user == nil {
		return nil
	}
	value := lo.Ternary(isClient, "yes", "no")
	if err := user.SetMeta("is_client", value); err != nil {
		log.Println("Failed to set is_client meta", err)
	}
	return user
}
