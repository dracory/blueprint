package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/tradingstore"
)

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
