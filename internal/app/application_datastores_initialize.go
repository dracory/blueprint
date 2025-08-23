package app

// dataStoresInitialize performs phase 1 of store setup. For now it delegates
// to initializeStores to preserve behavior; it will be refactored to create-only.
func (a *Application) dataStoresInitialize() error {

	// BlogStore: create and migrate, then set
	bs, err := newBlogStore(a.db)

	if err != nil {
		return err
	}

	a.SetBlogStore(bs)

	cs, err := newCacheStore(a.db)

	if err != nil {
		return err
	}

	a.SetCacheStore(cs)

	cms, err := newCmsStore(a.db)
	if err != nil {
		return err
	}
	a.SetCmsStore(cms)

	customStore, err := newCustomStore(a.db)
	if err != nil {
		return err
	}
	a.SetCustomStore(customStore)

	es, err := newEntityStore(a.db)
	if err != nil {
		return err
	}

	a.SetEntityStore(es)

	fs, err := newFeedStore(a.db)
	if err != nil {
		return err
	}
	a.SetFeedStore(fs)

	gs, err := newGeoStore(a.db)
	if err != nil {
		return err
	}
	a.SetGeoStore(gs)

	if ls, err := newLogStore(a.db); err != nil {
		return err
	} else {
		a.SetLogStore(ls)
	}

	ms, err := newMetaStore(a.db)
	if err != nil {
		return err
	}
	a.SetMetaStore(ms)

	ss, err := newSessionStore(a.db)
	if err != nil {
		return err
	}
	a.SetSessionStore(ss)

	shopStore, err := newShopStore(a.db)
	if err != nil {
		return err
	}
	a.SetShopStore(shopStore)

	sqlFileStore, err := newSqlFileStorage(a.db)
	if err != nil {
		return err
	}
	a.SetSqlFileStorage(sqlFileStore)

	// Setting store (constructor only)
	st, err := newSettingStore(a.db)
	if err != nil {
		return err
	}
	a.SetSettingStore(st)

	// User store (constructor only, migrate in phase 2)
	us, err := newUserStore(a.db)
	if err != nil {
		return err
	}
	a.SetUserStore(us)

	// Vault store (constructor only)
	vs, err := newVaultStore(a.db)
	if err != nil {
		return err
	}
	a.SetVaultStore(vs)

	// Blind index stores (constructors only)
	be, err := newBlindIndexEmailStore(a.db)
	if err != nil {
		return err
	}
	a.SetBlindIndexStoreEmail(be)

	bfn, err := newBlindIndexFirstNameStore(a.db)
	if err != nil {
		return err
	}
	a.SetBlindIndexStoreFirstName(bfn)

	bln, err := newBlindIndexLastNameStore(a.db)
	if err != nil {
		return err
	}
	a.SetBlindIndexStoreLastName(bln)

	// Task store (constructor only)
	ts, err := newTaskStore(a.db)
	if err != nil {
		return err
	}
	a.SetTaskStore(ts)

	// Stats store (constructor only)
	ss2, err := newStatsStore(a.db)
	if err != nil {
		return err
	}
	a.SetStatsStore(ss2)

	// Trading store (constructor only)
	tr, err := newTradingStore(a.db)
	if err != nil {
		return err
	}
	a.SetTradingStore(tr)

	return nil
}
