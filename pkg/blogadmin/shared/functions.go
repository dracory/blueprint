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
