package registry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/cmsstore"
)

func cmsStoreInitialize(registry RegistryInterface) error {
	if !registry.GetConfig().GetCmsStoreUsed() {
		return nil
	}

	if store, err := newCmsStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetCmsStore(store)
	}

	return nil
}

func cmsStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetCmsStoreUsed() {
		return nil
	}

	if registry.GetCmsStore() == nil {
		return errors.New("cms store is not initialized")
	}

	if err := registry.GetCmsStore().AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

// newCmsStore constructs the CMS store without running migrations
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
