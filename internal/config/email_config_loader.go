package config

import (
	"github.com/dracory/env"
	"github.com/spf13/cast"
)

// loadMailConfig loads mail configuration directly into the config.
func loadMailConfig(cfg ConfigInterface) {
	cfg.SetMailDriver(env.GetString(KEY_MAIL_DRIVER))
	cfg.SetMailFromAddress(env.GetString(KEY_MAIL_FROM_ADDRESS))
	cfg.SetMailFromName(env.GetString(KEY_MAIL_FROM_NAME))
	cfg.SetMailHost(env.GetString(KEY_MAIL_HOST))
	cfg.SetMailPassword(env.GetString(KEY_MAIL_PASSWORD))
	cfg.SetMailPort(cast.ToInt(env.GetString(KEY_MAIL_PORT)))
	cfg.SetMailUsername(env.GetString(KEY_MAIL_USERNAME))
}
