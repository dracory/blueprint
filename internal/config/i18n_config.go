package config

import "github.com/dracory/env"

// i18nConfig reads i18n configuration from environment variables.
func i18nConfig() i18nSettings {
	// Default Language
	//
	// The default locale used when no language is specified.
	// Must match one of the keys in the language list below.
	// Example: en, fr, de, es
	defaultLanguage := env.GetStringOrDefault(KEY_TRANSLATION_LANGUAGE_DEFAULT, defaultTranslationLanguage)

	// Language List
	//
	// The full list of supported locales mapped to their display names.
	// To add a language, update the languageList map below.
	// Example: {"en": "English", "fr": "French"}
	languageList := map[string]string{
		"en": "English",
		"bg": "Bulgarian",
		"de": "German",
	}

	return i18nSettings{defaultLanguage: defaultLanguage, languageList: languageList}
}

const defaultTranslationLanguage = "en"

type i18nSettings struct {
	defaultLanguage string
	languageList    map[string]string
}
