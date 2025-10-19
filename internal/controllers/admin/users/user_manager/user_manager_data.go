package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/userstore"
	"github.com/spf13/cast"
)

type userManagerControllerData struct {
	request         *http.Request
	action          string
	page            string
	pageInt         int
	perPage         int
	sortOrder       string
	sortBy          string
	formStatus      string
	formEmail       string
	formFirstName   string
	formLastName    string
	formCreatedFrom string
	formCreatedTo   string
	formUserID      string
	userList        []userstore.UserInterface
	userCount       int64
}

func (controller *userManagerController) prepareData(r *http.Request) (data userManagerControllerData, errorMessage string) {
	var err error
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.page = req.GetStringTrimmedOr(r, "page", "0")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(req.GetStringTrimmedOr(r, "per_page", "10"))
	data.sortOrder = req.GetStringTrimmedOr(r, "sort_order", sb.DESC)
	data.sortBy = req.GetStringTrimmedOr(r, "by", userstore.COLUMN_CREATED_AT)
	data.formEmail = req.GetStringTrimmed(r, "email")
	data.formFirstName = req.GetStringTrimmed(r, "first_name")
	data.formLastName = req.GetStringTrimmed(r, "last_name")
	data.formStatus = req.GetStringTrimmed(r, "status")
	data.formCreatedFrom = req.GetStringTrimmed(r, "created_from")
	data.formCreatedTo = req.GetStringTrimmed(r, "created_to")

	userList, userCount, err := controller.fetchUserList(data)

	if err != nil {
		controller.app.GetLogger().Error("Error. At userCreateController > prepareDataAndValidate", slog.String("error", err.Error()))
		return data, "Creating user failed. Please contact an administrator."
	}

	data.userList = userList

	data.userCount = userCount

	return data, ""
}

func (controller *userManagerController) fetchUserList(data userManagerControllerData) ([]userstore.UserInterface, int64, error) {
	if controller.app == nil || controller.app.GetUserStore() == nil {
		return []userstore.UserInterface{}, 0, errors.New("UserStore is not initialized")
	}

	userIDs := []string{}

	if data.formFirstName != "" {
		firstNameUserIDs, err := controller.app.GetBlindIndexStoreFirstName().Search(data.formFirstName, blindindexstore.SEARCH_TYPE_CONTAINS)

		if err != nil {
			controller.app.GetLogger().Error("At userManagerController > prepareData", slog.String("error", err.Error()))
			return []userstore.UserInterface{}, 0, err
		}

		if len(firstNameUserIDs) == 0 {
			return []userstore.UserInterface{}, 0, nil
		}

		userIDs = append(userIDs, firstNameUserIDs...)
	}

	if data.formLastName != "" {
		lastNameUserIDs, err := controller.app.GetBlindIndexStoreLastName().Search(data.formLastName, blindindexstore.SEARCH_TYPE_CONTAINS)

		if err != nil {
			controller.app.GetLogger().Error("At userManagerController > prepareData", slog.String("error", err.Error()))
			return []userstore.UserInterface{}, 0, err
		}

		if len(lastNameUserIDs) == 0 {
			return []userstore.UserInterface{}, 0, nil
		}

		userIDs = append(userIDs, lastNameUserIDs...)
	}

	if data.formEmail != "" {
		emailUserIDs, err := controller.app.GetBlindIndexStoreEmail().Search(data.formEmail, blindindexstore.SEARCH_TYPE_CONTAINS)

		if err != nil {
			controller.app.GetLogger().Error("At userManagerController > prepareData", slog.String("error", err.Error()))
			return []userstore.UserInterface{}, 0, err
		}

		if len(emailUserIDs) == 0 {
			return []userstore.UserInterface{}, 0, nil
		}

		userIDs = append(userIDs, emailUserIDs...)
	}

	query := userstore.NewUserQuery().
		SetSortDirection(data.sortOrder).
		SetOrderBy(data.sortBy).
		SetOffset(data.pageInt * data.perPage).
		SetLimit(data.perPage)

	if len(userIDs) > 0 {
		query.SetIDIn(userIDs)
	}

	if data.formStatus != "" {
		query.SetStatus(data.formStatus)
	}

	if data.formCreatedFrom != "" {
		query.SetCreatedAtGte(data.formCreatedFrom + " 00:00:00")
	}

	if data.formCreatedTo != "" {
		query.SetCreatedAtLte(data.formCreatedTo + " 23:59:59")
	}

	userList, err := controller.app.GetUserStore().UserList(data.request.Context(), query)

	if err != nil {
		controller.app.GetLogger().Error("At userManagerController > prepareData", slog.String("error", err.Error()))
		return []userstore.UserInterface{}, 0, err
	}

	userCount, err := controller.app.GetUserStore().UserCount(data.request.Context(), query)

	if err != nil {
		controller.app.GetLogger().Error("At userManagerController > prepareData", slog.String("error", err.Error()))
		return []userstore.UserInterface{}, 0, err
	}

	return userList, userCount, nil

}
