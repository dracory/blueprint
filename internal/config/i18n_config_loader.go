package config

import "github.com/dracory/env"

// loadTranslationConfig loads translation configuration directly into the config.
func loadTranslationConfig(cfg ConfigInterface) {
	defaultLang := env.GetString(KEY_TRANSLATION_LANGUAGE_DEFAULT)
	if defaultLang == "" {
		defaultLang = translationLanguageDefault()
	}

	cfg.SetTranslationLanguageDefault(defaultLang)
	cfg.SetTranslationLanguageList(translationLanguageListDefault())
}
