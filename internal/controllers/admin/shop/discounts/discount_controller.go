package admin

import (
	"context"
	"errors"
	"net/http"
	"project/internal/controllers/admin/shop/shared"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/cdn"
	crud "github.com/dracory/crud/v2"
	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type discountController struct {
	registry registry.RegistryInterface
}

func NewDiscountController(registry registry.RegistryInterface) *discountController {
	return &discountController{registry: registry}
}

func (discountController *discountController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	discountsCrud, err := crud.New(crud.Config{
		EntityNameSingular: "Discount",
		EntityNamePlural:   "Discounts",
		Endpoint:           shared.NewLinks().Discounts(map[string]string{}),
		ColumnNames: []string{
			"Title",
			"Status",
			"Type",
			"Amount",
			"Period Valid",
			"Discount Code",
			"Created",
		},
		CreateFields: []crud.FieldInterface{
			crud.NewField(crud.FieldOptions{
				Label: "Title",
				Name:  "title",
				Type:  crud.FORM_FIELD_TYPE_STRING,
			}),
		},
		UpdateFields: []crud.FieldInterface{
			crud.NewField(crud.FieldOptions{
				Label: "Status",
				Name:  "status",
				Type:  crud.FORM_FIELD_TYPE_SELECT,
				Options: []crud.FieldOption{
					{
						Key:   "",
						Value: "",
					},
					{
						Key:   shopstore.DISCOUNT_STATUS_DRAFT,
						Value: shopstore.DISCOUNT_STATUS_DRAFT,
					},
					{
						Key:   shopstore.DISCOUNT_STATUS_INACTIVE,
						Value: shopstore.DISCOUNT_STATUS_INACTIVE,
					},
					{
						Key:   shopstore.DISCOUNT_STATUS_ACTIVE,
						Value: shopstore.DISCOUNT_STATUS_ACTIVE,
					},
				},
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Title",
				Name:  "title",
				Type:  crud.FORM_FIELD_TYPE_STRING,
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Type",
				Name:  "type",
				Type:  crud.FORM_FIELD_TYPE_SELECT,
				Options: []crud.FieldOption{
					{
						Key:   "",
						Value: "",
					},
					{
						Key:   shopstore.DISCOUNT_TYPE_AMOUNT,
						Value: shopstore.DISCOUNT_TYPE_AMOUNT,
					},
					{
						Key:   shopstore.DISCOUNT_TYPE_PERCENT,
						Value: shopstore.DISCOUNT_TYPE_PERCENT,
					},
				},
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Amount",
				Name:  "amount",
				Type:  crud.FORM_FIELD_TYPE_NUMBER,
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Discount Code",
				Name:  "code",
				Type:  crud.FORM_FIELD_TYPE_STRING,
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Starts",
				Name:  "starts_at",
				Type:  crud.FORM_FIELD_TYPE_DATETIME,
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Ends",
				Name:  "ends_at",
				Type:  crud.FORM_FIELD_TYPE_DATETIME,
			}),
			crud.NewField(crud.FieldOptions{
				Label: "Description",
				Name:  "description",
				Type:  crud.FORM_FIELD_TYPE_HTMLAREA,
			}),
		},
		FuncRows:            discountController.FuncRows,
		FuncCreate:          discountController.FuncCreate,
		FuncFetchReadData:   discountController.FuncFetchReadData,
		FuncFetchUpdateData: discountController.FuncFetchUpdateData,
		FuncTrash:           discountController.FuncTrash,
		FuncUpdate:          discountController.FuncUpdate,
		FuncLayout:          discountController.FuncLayout,
		HomeURL:             links.Admin().Home(),
	})

	if err != nil {
		return "Error: " + err.Error()
	}

	discountsCrud.Handler(w, r)
	return ""
}

func (discountController *discountController) FuncLayout(w http.ResponseWriter, r *http.Request, title string, content string, styleURLs []string, style string, scriptURLs []string, script string) string {
	scriptURLs = append([]string{
		cdn.Jquery_3_6_4(),
	}, scriptURLs...)

	return layouts.NewAdminLayout(discountController.registry, r, layouts.Options{
		Title:      title + " | Admin",
		Content:    hb.Wrap().HTML(content),
		StyleURLs:  styleURLs,
		ScriptURLs: scriptURLs,
		Scripts:    []string{script},
		Styles: []string{
			`nav#Toolbar {border-bottom: 4px solid red;}`,
			style,
		},
	}).ToHTML()
}

func (discountController *discountController) FuncRows(r *http.Request) ([]crud.Row, error) {
	if discountController.registry.GetShopStore() == nil {
		return nil, errors.New("shop store not configured")
	}

	discounts, err := discountController.registry.GetShopStore().DiscountList(context.Background(), shopstore.NewDiscountQuery())

	if err != nil {
		return nil, err
	}

	rows := lo.Map(discounts, func(discount shopstore.DiscountInterface, _ int) crud.Row {
		return crud.Row{
			ID: discount.GetID(),
			Data: []string{
				discount.GetTitle(),
				discount.GetStatus(),
				discount.GetType(),
				cast.ToString(discount.GetAmount()),
				discount.GetStartsAtCarbon().Format("d M Y") + " - " + discount.GetEndsAtCarbon().Format("d M Y"),
				discount.GetCode(),
				discount.GetCreatedAtCarbon().Format("d M Y"),
			},
		}
	})

	return rows, nil
}

func (discountController *discountController) FuncUpdate(r *http.Request, entityID string, data map[string]string) error {
	if discountController.registry.GetShopStore() == nil {
		return errors.New("shop store not configured")
	}

	discount, err := discountController.registry.GetShopStore().DiscountFindByID(context.Background(), entityID)

	if err != nil {
		return err
	}

	if discount == nil {
		return errors.New("discount not found")
	}

	amountStr := data["amount"]
	startsAt := data["starts_at"]
	endsAt := data["ends_at"]
	title := data["title"]
	code := data["code"]
	status := data["status"]
	discountType := data["type"]

	if title == "" {
		return errors.New("title is required")
	}

	if status == "" {
		return errors.New("status is required")
	}

	if code == "" {
		return errors.New("code is required")
	}

	if discountType == "" {
		return errors.New("discount type is required")
	}

	if startsAt == "" {
		return errors.New("starts_at is required")
	}

	if endsAt == "" {
		return errors.New("ends_at is required")
	}

	if amountStr == "" {
		amountStr = "0"
	}

	amount := cast.ToFloat64(amountStr)
	startsAt = carbon.Parse(startsAt).ToDateTimeString(carbon.UTC)
	endsAt = carbon.Parse(endsAt).ToDateTimeString(carbon.UTC)

	discount.SetTitle(title)
	discount.SetDescription(data["description"])
	discount.SetStatus(status)
	discount.SetAmount(amount)
	discount.SetType(discountType)
	discount.SetCode(code)
	discount.SetStartsAt(startsAt)
	discount.SetEndsAt(endsAt)

	err = discountController.registry.GetShopStore().DiscountUpdate(context.Background(), discount)

	if err != nil {
		return err
	}

	return nil
}

func (discountController *discountController) FuncFetchReadData(r *http.Request, discountID string) ([]crud.KeyValue, error) {
	if discountController.registry.GetShopStore() == nil {
		return nil, errors.New("shop store not configured")
	}

	discount, err := discountController.registry.GetShopStore().DiscountFindByID(context.Background(), discountID)

	if err != nil {
		return nil, err
	}

	if discount == nil {
		return nil, errors.New("discount not found")
	}

	data := []crud.KeyValue{
		{Key: "Title", Value: discount.GetTitle()},
		{Key: "Status", Value: discount.GetStatus()},
		{Key: "Description", Value: discount.GetDescription()},
		{Key: "Type", Value: discount.GetType()},
		{Key: "Amount", Value: cast.ToString(discount.GetAmount())},
		{Key: "Starts At", Value: discount.GetStartsAtCarbon().Format("d M Y")},
		{Key: "Ends At", Value: discount.GetEndsAtCarbon().Format("d M Y")},
		{Key: "Created", Value: discount.GetCreatedAtCarbon().Format("d M Y")},
		{Key: "Updated", Value: discount.GetUpdatedAtCarbon().Format("d M Y")},
	}

	return data, nil
}

func (discountController *discountController) FuncFetchUpdateData(r *http.Request, discountID string) (map[string]string, error) {
	if discountController.registry.GetShopStore() == nil {
		return nil, errors.New("shop store not configured")
	}

	discount, err := discountController.registry.GetShopStore().DiscountFindByID(context.Background(), discountID)

	if err != nil {
		return nil, err
	}

	if discount == nil {
		return nil, errors.New("discount not found")
	}

	return map[string]string{
		"title":       discount.GetTitle(),
		"status":      discount.GetStatus(),
		"amount":      cast.ToString(discount.GetAmount()),
		"description": discount.GetDescription(),
		"type":        discount.GetType(),
		"code":        discount.GetCode(),
		"starts_at":   discount.GetStartsAtCarbon().ToDateTimeString(),
		"ends_at":     discount.GetEndsAtCarbon().ToDateTimeString(),
		"created_at":  discount.GetCreatedAtCarbon().ToDateTimeString(),
		"updated_at":  discount.GetUpdatedAtCarbon().ToDateTimeString(),
	}, nil
}

func (discountController *discountController) FuncCreate(r *http.Request, data map[string]string) (discountID string, err error) {
	if discountController.registry.GetShopStore() == nil {
		return "", errors.New("shop store not configured")
	}

	discount := shopstore.NewDiscount()
	discount.SetTitle(data["title"])
	discount.SetStatus(shopstore.DISCOUNT_STATUS_DRAFT)
	discount.SetAmount(0.00)

	err = discountController.registry.GetShopStore().DiscountCreate(context.Background(), discount)

	if err != nil {
		return "", err
	}

	return discount.GetID(), nil
}

func (discountController *discountController) FuncTrash(r *http.Request, discountID string) error {
	if discountController.registry.GetShopStore() == nil {
		return errors.New("shop store not configured")
	}

	err := discountController.registry.GetShopStore().DiscountSoftDeleteByID(context.Background(), discountID)
	return err
}
