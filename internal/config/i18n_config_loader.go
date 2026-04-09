package config

import "github.com/dracory/env"

// loadTranslationConfig loads i18n configuration directly into the config.
func loadTranslationConfig(cfg ConfigInterface) {
	// Default Language
	//
	// The default locale used when no language is specified.
	// Must match one of the keys in the language list below.
	// Example: en, fr, de, es
	defaultLang := env.GetString(KEY_TRANSLATION_LANGUAGE_DEFAULT)
	if defaultLang == "" {
		defaultLang = translationLanguageDefault()
	}

	// Language List
	//
	// The full list of supported locales mapped to their display names.
	// To add a language, update translationLanguageListDefault() in defaults.go.
	// Example: {"en": "English", "fr": "French"}
	cfg.SetTranslationLanguageDefault(defaultLang)
	cfg.SetTranslationLanguageList(translationLanguageListDefault())
}
