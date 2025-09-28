package config

const (
	defaultTranslationLanguage = "en"
)

func translationLanguageDefault() string {
	return defaultTranslationLanguage
}

func translationLanguageListDefault() map[string]string {
	return map[string]string{
		"en": "English",
		"bg": "Bulgarian",
		"de": "German",
	}
}
