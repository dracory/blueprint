package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/cmsstore"
)

func newCmsStore(db *sql.DB) (cmsstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                   db,
		BlockTableName:       "snv_cms_block",
		PageTableName:        "snv_cms_page",
		TemplateTableName:    "snv_cms_template",
		SiteTableName:        "snv_cms_site",
		MenusEnabled:         true,
		MenuItemTableName:    "snv_cms_menu_item",
		MenuTableName:        "snv_cms_menu",
		TranslationsEnabled:  true,
		TranslationTableName: "snv_cms_translation",
		VersioningEnabled:    true,
		VersioningTableName:  "snv_cms_version",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("cmsstore.NewStore returned a nil store")
	}

	return st, nil
}
