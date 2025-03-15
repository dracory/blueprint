package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/cmsstore"
)

// CmsStoreInitialize initializes the CMS store
func CmsStoreInitialize(db *sql.DB) (cmsstore.StoreInterface, error) {
	cmsStoreInstance, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB: db,

		BlockTableName:    "snv_cms_block",
		PageTableName:     "snv_cms_page",
		TemplateTableName: "snv_cms_template",
		SiteTableName:     "snv_cms_site",

		MenusEnabled:      true,
		MenuItemTableName: "snv_cms_menu_item",
		MenuTableName:     "snv_cms_menu",

		TranslationsEnabled:        true,
		TranslationTableName:       "snv_cms_translation",
		TranslationLanguageDefault: "en",
		TranslationLanguages:       map[string]string{"en": "English", "bg": "Bulgarian", "de": "German"},

		VersioningEnabled:   true,
		VersioningTableName: "snv_cms_version",
	})

	if err != nil {
		return nil, errors.Join(errors.New("cmsstore.NewStore"), err)
	}

	if cmsStoreInstance == nil {
		return nil, errors.New("CmsStore is nil")
	}

	return cmsStoreInstance, nil
}

// CmsStoreAutoMigrate runs migrations for the CMS store
func CmsStoreAutoMigrate(ctx context.Context, store cmsstore.StoreInterface) error {
	if store == nil {
		return errors.New("cmsstore.AutoMigrate: CmsStore is nil")
	}

	err := store.AutoMigrate(ctx)

	if err != nil {
		return errors.Join(errors.New("cmsstore.AutoMigrate"), err)
	}

	return nil
}
