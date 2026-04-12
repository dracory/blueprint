package shared

import (
	"project/internal/registry"

	"github.com/dracory/hb"
	"github.com/dracory/llm"
)

func LlmEngine(registry registry.RegistryInterface) (llm.LlmInterface, error) {
	provider := llm.ProviderOpenRouter
	if registry.GetConfig().IsEnvTesting() {
		provider = llm.ProviderMock
	}
	return llm.JSONModel(provider, llm.LlmOptions{
		ApiKey: registry.GetConfig().GetOpenRouterApiKey(),
		Model:  llm.OPENROUTER_MODEL_GEMINI_2_5_FLASH_LITE,
	})
}

func ErrorPopup(errorMessage string) hb.TagInterface {
	return hb.Swal(hb.SwalOptions{
		Title:            "Error",
		Text:             errorMessage,
		Icon:             "error",
		Timer:            10000,
		TimerProgressBar: true,
	})
}

func SuccessPopup(successMessage string) hb.TagInterface {
	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             successMessage,
		Icon:             "success",
		Timer:            10000,
		TimerProgressBar: true,
	})
}

func SuccessPopupWithRedirect(successMessage string, redirectUrl string, redirectSeconds int) hb.TagInterface {
	if redirectUrl != "" {
		return hb.Swal(hb.SwalOptions{
			Title:            "Success",
			Text:             successMessage,
			Icon:             "success",
			Timer:            redirectSeconds * 1000,
			TimerProgressBar: true,
			RedirectURL:      redirectUrl,
			RedirectSeconds:  redirectSeconds,
		})
	}

	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             successMessage,
		Icon:             "success",
		Timer:            redirectSeconds * 1000,
		TimerProgressBar: true,
	})
}

// Slugify converts a string to a URL-friendly slug
func Slugify(text string) string {
	result := []rune{}
	for _, r := range text {
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+32)
		} else if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result = append(result, r)
		} else if r == ' ' || r == '-' || r == '_' {
			result = append(result, '-')
		}
	}
	return string(result)
}
