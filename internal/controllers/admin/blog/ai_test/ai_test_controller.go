package aitest

import (
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/base/req"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/llm"
)

type AiTestController struct {
	app types.AppInterface
}

func NewAiTestController(app types.AppInterface) *AiTestController {
	return &AiTestController{app: app}
}

func (c *AiTestController) Handler(w http.ResponseWriter, r *http.Request) string {
	if r.Method == http.MethodPost && r.FormValue("action") == "testai" {
		userMsg := req.Value(r, "user_message")
		if userMsg == "" {
			userMsg = "Tell me shortly about blogs."
		}

		model, err := shared.LlmEngine(c.app)
		if err != nil {
			if _, writeErr := w.Write([]byte(hb.Swal(hb.SwalOptions{
				Title:            "Error",
				Text:             err.Error(),
				Icon:             "error",
				Timer:            15000,
				TimerProgressBar: true,
			}).ToHTML())); writeErr != nil {
				return ""
			}
			return ""
		}

		if model == nil {
			if _, writeErr := w.Write([]byte(hb.Swal(hb.SwalOptions{
				Title:            "Error",
				Text:             "model is nil",
				Icon:             "error",
				Timer:            15000,
				TimerProgressBar: true,
			}).ToHTML())); writeErr != nil {
				return ""
			}
			return ""
		}

		response, err := model.GenerateText("You are a helpful assistant.", userMsg, llm.LlmOptions{
			MaxTokens:   128,
			Temperature: 0.7,
		})
		if err != nil {
			if _, writeErr := w.Write([]byte(hb.Swal(hb.SwalOptions{
				Title:            "Error",
				Text:             err.Error(),
				Icon:             "error",
				Timer:            15000,
				TimerProgressBar: true,
			}).ToHTML())); writeErr != nil {
				return ""
			}
			return ""
		}

		if _, writeErr := w.Write([]byte(hb.Swal(hb.SwalOptions{
			Title:            "Success",
			Text:             response,
			Icon:             "success",
			Timer:            15000,
			TimerProgressBar: true,
		}).ToHTML())); writeErr != nil {
			return ""
		}
		return ""
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:      "Test AI Connectivity",
		Content:    c.view(),
		ScriptURLs: []string{cdn.Sweetalert2_11()},
		Styles:     []string{},
	}).ToHTML()
}

func (c *AiTestController) view() *hb.Tag {
	header := hb.Div().
		Class("text-center mb-4").
		Child(
			hb.Heading1().HTML("🤖 Test AI Connectivity").Class("display-4 fw-bold mb-2 text-primary"),
		).
		Child(
			hb.Paragraph().HTML(`<span class='badge bg-info fs-6'>AI Diagnostics</span>`).Class("mb-2"),
		)

	desc := hb.Paragraph().
		HTML(`Welcome to the AI Diagnostics page!<br><br>
		Here you can verify that your AI model is correctly configured and reachable from the application.<br>
		Enter a custom message below and click the button to send a test request to your configured AI provider.<br>
		If everything is set up correctly, you will see a response from the AI.<br>
		If not, you will see an error message with details.<br><br>
		<em>Tip: Use this page after any configuration changes or when troubleshooting AI integration.</em>`).
		Class("mb-4 text-muted fs-5")

	aiIllustration := hb.Div().
		Class("mb-4 text-center").
		Child(
			hb.Img("https://cdn.jsdelivr.net/gh/twitter/twemoji@14.0.2/assets/72x72/1f916.png").
				Alt("AI Bot").
				Style("width:80px;height:80px;filter:drop-shadow(0 4px 16px #2d9cdb44);"),
		)

	inputRow := hb.Form().
		ID("FormTestAI").
		Class("mb-3").
		Attr("hx-post", shared.NewLinks().AiTest(map[string]string{"action": "testai"})).
		Attr("hx-target", "body").
		Attr("hx-swap", "beforeend").
		Attr("hx-indicator", "#ButtonTestAI").
		Child(
			hb.Label().
				For("InputUserMessage").
				Class("form-label fw-semibold mb-2").
				HTML("Message for AI:"),
		).
		Child(
			hb.Div().Class("input-group").
				Child(
					hb.Input().
						Type("text").
						Name("user_message").
						ID("InputUserMessage").
						Placeholder("Enter your message for the AI...").
						Value("what is the meaning of life").
						Class("form-control form-control-lg shadow-sm").
						Attr("required", "required"),
				).
				Child(
					hb.Button().
						Type("submit").
						ID("ButtonTestAI").
						Class("btn btn-primary btn-lg fw-semibold px-5 shadow-sm").
						Child(hb.Span().Class("me-2").HTML("🚦")).
						Text(`Test AI is working`).
						Child(hb.I().Class("spinner-border spinner-border-sm htmx-indicator ms-2")),
				),
		)

	footer := hb.Div().
		Class("text-center mt-4 text-secondary small").
		HTML("<hr><div>Need help? Visit the <a href='https://github.com/yourproject/wiki/ai-troubleshooting' target='_blank'>AI Troubleshooting Guide</a>.</div>")

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Blog",
			URL:  links.Admin().Blog(),
		},
		{
			Name: "AI Tools",
			URL:  shared.NewLinks().AiTools(),
		},
		{
			Name: "Test AI",
			URL:  shared.NewLinks().AiTest(),
		},
	})

	card := hb.Div().
		Class("card shadow border-0 w-100 mb-5 animate__animated animate__fadeInUp").
		Child(
			hb.Div().Class("card-body p-5").
				Child(header).
				Child(aiIllustration).
				Child(desc).
				Child(inputRow).
				Child(footer),
		)

	return hb.Div().
		Class("container min-vh-100 py-4 bg-light animate__animated animate__fadeIn").
		Child(breadcrumbs).
		Child(card)
}
