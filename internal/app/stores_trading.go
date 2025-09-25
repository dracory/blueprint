package app

import (
	"context"
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/tradingstore"
)

// tradingStoreInitialize initializes the trading store if enabled in the configuration.
func tradingStoreInitialize(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetTradingStoreUsed() {
		return nil
	}

	if store, err := newTradingStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetTradingStore(store)
	}

	return nil
}

func tradingStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetTradingStoreUsed() {
		return nil
	}

	if app.GetTradingStore() == nil {
		return errors.New("trading store is not initialized")
	}

	if err := app.GetTradingStore().AutoMigrateInstruments(context.Background()); err != nil {
		return err
	}
	if err := app.GetTradingStore().AutoMigratePrices(context.Background()); err != nil {
		return err
	}

	return nil
}

// newTradingStore constructs the Trading store without running migrations
func newTradingStore(db *sql.DB) (tradingstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := tradingstore.NewStore(tradingstore.NewStoreOptions{
		DB:                   db,
		InstrumentTableName:  "snv_trading_instruments",
		PriceTableNamePrefix: "snv_trading_prices_",
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("tradingstore.NewStore returned a nil store")
	}
	return st, nil
}
