package contact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"project/internal/links"
	"project/internal/registry"
	"project/internal/tasks/email_admin_new_contact"
	"project/internal/types"

	"github.com/dracory/bs"
	"github.com/dracory/csrf"
	"github.com/dracory/customstore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
)

type formContact struct {
	liveflux.Base
	App             types.RegistryInterface
	UserID          string
	Email           string
	FirstName       string
	LastName        string
	Text            string
	CsrfToken       string
	ErrorMessage    string
	SuccessMessage  string
	RedirectURL     string
	CanUpdateEmail  bool
	CanUpdateFirst  bool
	CanUpdateLast   bool
	CaptchaQuestion string
	CaptchaExpected string
	CaptchaAnswer   string
}

func NewFormContact(registry registry.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&formContact{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*formContact); ok {
		c.App = registry
	}

	return inst
}

func (c *formContact) GetKind() string {
	return "website_contact_form"
}

func (c *formContact) Mount(ctx context.Context, params map[string]string) error {
	c.UserID = strings.TrimSpace(params["user_id"])
	c.CsrfToken = csrf.TokenGenerate("holymoly")

	c.CanUpdateEmail = true
	c.CanUpdateFirst = true
	c.CanUpdateLast = true

	// Initialize simple math captcha (e.g. "6 + 2 =")
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(5) + 1 // 1-5
	b := rand.Intn(5) + 1 // 1-5
	sum := a + b
	sumStr := strconv.Itoa(sum)
	c.CaptchaQuestion = strconv.Itoa(a) + " + " + strconv.Itoa(b) + " ="
	c.CaptchaExpected = hashCaptchaValue(sumStr)
	c.CaptchaAnswer = ""

	if c.UserID == "" || c.App == nil || c.App.GetUserStore() == nil {
		return nil
	}

	user, err := c.App.GetUserStore().UserFindByID(ctx, c.UserID)
	if err != nil || user == nil {
		return nil
	}

	// Mirror previous contact_controller.AnyIndex behavior: use plain user fields
	c.Email = user.Email()
	c.FirstName = user.FirstName()
	c.LastName = user.LastName()

	c.CanUpdateFirst = user.FirstName() == ""
	c.CanUpdateLast = user.LastName() == ""
	c.CanUpdateEmail = user.Email() == ""

	return nil
}

func (c *formContact) Handle(ctx context.Context, action string, data url.Values) error {
	if action != "submit" {
		return nil
	}

	if data == nil {
		data = url.Values{}
	}

	c.Email = strings.TrimSpace(data.Get("email"))
	c.FirstName = strings.TrimSpace(data.Get("first_name"))
	c.LastName = strings.TrimSpace(data.Get("last_name"))
	c.Text = strings.TrimSpace(data.Get("text"))
	c.CsrfToken = strings.TrimSpace(data.Get("csrf_token"))
	c.CaptchaAnswer = strings.TrimSpace(data.Get("captcha_answer"))
	c.CaptchaExpected = strings.TrimSpace(data.Get("captcha_expected"))

	if c.CsrfToken == "" {
		c.ErrorMessage = "CSRF token is required"
		c.SuccessMessage = ""
		c.RedirectURL = links.Website().Contact()
		return nil
	}

	if !csrf.TokenValidate(c.CsrfToken, "holymoly") {
		c.ErrorMessage = "CSRF token is invalid"
		c.SuccessMessage = ""
		c.RedirectURL = links.Website().Contact()
		return nil
	}

	if c.FirstName == "" {
		c.ErrorMessage = "First name is required"
		c.SuccessMessage = ""
		return nil
	}

	if c.LastName == "" {
		c.ErrorMessage = "Last name is required"
		c.SuccessMessage = ""
		return nil
	}

	if c.Email == "" {
		c.ErrorMessage = "Email is required"
		c.SuccessMessage = ""
		return nil
	}

	if c.Text == "" {
		c.ErrorMessage = "Text is required"
		c.SuccessMessage = ""
		return nil
	}

	// Validate captcha
	if c.CaptchaAnswer == "" || c.CaptchaExpected == "" {
		c.ErrorMessage = "Please answer the verification question"
		c.SuccessMessage = ""
		return nil
	}

	if hashCaptchaValue(c.CaptchaAnswer) != c.CaptchaExpected {
		c.ErrorMessage = "Verification answer is incorrect"
		c.SuccessMessage = ""
		return nil
	}

	record := customstore.NewRecord("contact")

	if err := record.SetPayloadMap(map[string]interface{}{
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"email":      c.Email,
		"text":       c.Text,
	}); err != nil {
		if c.App != nil && c.App.GetLogger() != nil {
			c.App.GetLogger().Error("At formContact.Handle", "error", err.Error())
		}
		c.ErrorMessage = "System error occurred. Please try again later."
		c.SuccessMessage = ""
		return nil
	}

	if c.App == nil || c.App.GetCustomStore() == nil {
		c.ErrorMessage = "System error occurred. Please try again later."
		c.SuccessMessage = ""
		return nil
	}

	if err := c.App.GetCustomStore().RecordCreate(record); err != nil {
		c.App.GetLogger().Error("At formContact.Handle", "error", err.Error())
		c.ErrorMessage = "System error occurred. Please try again later."
		c.SuccessMessage = ""
		return nil
	}

	if _, err := email_admin_new_contact.NewEmailToAdminOnNewContactFormSubmittedTaskHandler(c.App).Enqueue(); err != nil {
		c.App.GetLogger().Error("At formContact.Handle. Enqueue EmailToAdminOnNewContactFormSubmittedTask", "error", err.Error())
	}

	if c.UserID != "" && c.App != nil && c.App.GetUserStore() != nil {
		user, err := c.App.GetUserStore().UserFindByID(ctx, c.UserID)
		if err == nil && user != nil {
			if c.CanUpdateFirst {
				user.SetFirstName(c.FirstName)
			}
			if c.CanUpdateLast {
				user.SetLastName(c.LastName)
			}

			if err := c.App.GetUserStore().UserUpdate(context.Background(), user); err != nil {
				c.App.GetLogger().Error("At formContact.Handle", "error", err.Error())
			}
		}
	}

	c.ErrorMessage = ""
	c.SuccessMessage = "Your message has been sent."
	c.RedirectURL = links.Website().Contact()

	return nil
}

func (c *formContact) Render(ctx context.Context) hb.TagInterface {
	required := hb.Sup().
		HTML("required").
		Style("margin-left:5px;color:lightcoral;")

	firstName := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("First name").
			Child(required),
		bs.FormInput().
			Name("first_name").
			Value(c.FirstName).
			AttrIf(!c.CanUpdateFirst, "readonly", "readonly").
			StyleIf(!c.CanUpdateFirst, "background-color:#ccc;"),
	})

	lastName := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Last name").Child(required),
		bs.FormInput().
			Name("last_name").
			Value(c.LastName).
			AttrIf(!c.CanUpdateLast, "readonly", "readonly").
			StyleIf(!c.CanUpdateLast, "background-color:#ccc;"),
	})

	email := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Email").
			Child(required),
		bs.FormInput().
			Name("email").
			Value(c.Email).
			AttrIf(!c.CanUpdateEmail, "readonly", "readonly").
			StyleIf(!c.CanUpdateEmail, "background-color:#ccc;"),
	})

	csrfToken := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormInput().Type(hb.TYPE_HIDDEN).
			Name("csrf_token").
			Value(c.CsrfToken),
	})

	text := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Text").
			Child(required),
		bs.FormTextArea().
			Name("text").
			Value(c.Text).
			HTML(c.Text).
			Style("height:200px;"),
	})

	captchaGroup := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Verify you are human by answering this question").
			Child(required),
		hb.Div().
			Class("input-group").
			Child(hb.Span().
				Class("input-group-text").
				Text(c.CaptchaQuestion)).
			Child(bs.FormInput().
				Name("captcha_answer").
				Value(c.CaptchaAnswer)),
	})

	buttonSubmit := bs.Button().
		Class("btn-primary mb-0").
		Attr("type", "submit").
		Attr(liveflux.DataFluxAction, "submit").
		Attr(liveflux.DataFluxIndicator, ".contact-submit-spinner").
		Child(hb.I().Class("bi bi-rocket me-2")).
		Child(hb.Span().Text("Send")).
		Child(hb.Div().
			Style("display: none").
			Class("contact-submit-spinner ms-2").
			Class("spinner-border spinner-border-sm text-light"))

	formContact := hb.Form().
		ID("FormContact").
		Children([]hb.TagInterface{
			hb.Div().
				Class("row g-4").
				Child(hb.Div().
					Class("col-6").
					Child(firstName)).
				Child(hb.Div().
					Class("col-6").
					Child(lastName)).
				Child(hb.Div().
					Class("col-12").
					Child(email)).
				Child(hb.Div().
					Class("col-12").
					Child(text)).
				Child(hb.Div().
					Class("col-12").
					Child(captchaGroup)),
			hb.Div().
				Class("row mt-3").
				Child(hb.Div().
					Class("col-12").
					Class("d-sm-flex justify-content-end").
					Child(buttonSubmit)),
		}).
		Child(csrfToken).
		Child(bs.FormInput().
			Type(hb.TYPE_HIDDEN).
			Name("captcha_expected").
			Value(c.CaptchaExpected))

	errorMessageJSON, _ := json.Marshal(c.ErrorMessage)
	successMessageJSON, _ := json.Marshal(c.SuccessMessage)

	card := hb.Div().ID("CardContact").
		Class("card bg-transparent border rounded-3").
		Style("text-align:left;").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").Children([]hb.TagInterface{
				formContact,
			}),
		})

	if c.ErrorMessage != "" {
		card = card.Child(hb.Script(`
 			Swal.fire({
 				icon: 'error',
 				title: 'Oops...',
 				text: ` + string(errorMessageJSON) + `,
 			})
 		`))
	}

	if c.SuccessMessage != "" {
		card = card.Child(hb.Script(`
 			Swal.fire({
 				icon: 'success',
 				title: 'Saved',
 				text: ` + string(successMessageJSON) + `,
 			})
 		`))
	}

	if c.RedirectURL != "" {
		card = card.Child(hb.Script(`setTimeout(() => {window.location.href="` + c.RedirectURL + `";}, 5000);`))
	}

	return c.Root(card)
}

func hashCaptchaValue(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:])
}

func init() {
	if err := liveflux.Register(&formContact{}); err != nil {
		log.Printf("Failed to register formContact component: %v", err)
	}
}
