package blogai

// import (
// 	"github.com/dracory/database"
// 	"github.com/dracory/customstore"
// 	_ "modernc.org/sqlite"
// )

// var customStore customstore.StoreInterface

// func Initialize() {
// 	db, err := database.Open(database.Options().
// 		SetDatabaseType("sqlite").
// 		SetDatabaseHost("").
// 		SetDatabasePort("").
// 		SetDatabaseName("blogai.db").
// 		SetCharset(`utf8mb4`).
// 		SetUserName("").
// 		SetPassword(""))

// 	if err != nil {
// 		panic(err)
// 	}

// 	customStore, err := customstore.NewStore(customstore.NewStoreOptions{
// 		DB:        db,
// 		TableName: "blogai_custom_record",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	if customStore == nil {
// 		panic("unexpected nil database")
// 	}

// 	customStore = customStore
// }
