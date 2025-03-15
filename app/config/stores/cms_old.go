package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/base/database"
	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/cms"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/ui"
)

// CmsInitialize initializes the CMS
func CmsInitialize(db *sql.DB, blockEditorDefinitions []blockeditor.BlockDefinition, blockEditorRenderer func(blocks []ui.BlockInterface) string, translationLanguageDefault string, translationLanguageList map[string]string) (*cms.Cms, error) {
	cmsInstance, err := cms.NewCms(cms.Config{
		Database:                   sb.NewDatabase(db, database.DatabaseType(db)),
		Prefix:                     "cms_",
		TemplatesEnable:            true,
		PagesEnable:                true,
		MenusEnable:                true,
		BlocksEnable:               true,
		BlockEditorDefinitions:     blockEditorDefinitions,
		BlockEditorRenderer:        blockEditorRenderer,
		EntitiesAutomigrate:        true,
		Shortcodes:                 []cms.ShortcodeInterface{},
		TranslationsEnable:         true,
		TranslationLanguageDefault: translationLanguageDefault,
		TranslationLanguages:       translationLanguageList,
	})

	if err != nil {
		return nil, errors.Join(errors.New("cms.NewCms"), err)
	}

	if cmsInstance == nil {
		return nil, errors.New("cmsInstance is nil")
	}

	return cmsInstance, nil
}

// CmsAutoMigrate runs migrations for the CMS
// Note: No need to implement this as it's migrated during initialization
func CmsAutoMigrate(_ context.Context, _ *cms.Cms) error {
	// No need. Migrated during initialize
	return nil
}
