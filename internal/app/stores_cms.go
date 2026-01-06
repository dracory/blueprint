package app

import (
	"context"
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/cmsstore"
)

func cmsStoreInitialize(app types.RegistryInterface) error {
	if !app.GetConfig().GetCmsStoreUsed() {
		return nil
	}

	if store, err := newCmsStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetCmsStore(store)
	}

	return nil
}

func cmsStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetCmsStoreUsed() {
		return nil
	}

	if app.GetCmsStore() == nil {
		return errors.New("cms store is not initialized")
	}

	if err := app.GetCmsStore().AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

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
