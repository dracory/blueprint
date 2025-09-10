package testutils

import (
	"log"
	"testing"

	"github.com/dracory/cachestore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDatabase creates an in-memory SQLite database for testing
func SetupTestDatabase(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to database")

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get sql.DB")

	return db, func() {
		sqlDB.Close()
	}
}

// SetupTestAuth creates an auth.Auth instance with in-memory stores for testing
func SetupTestAuth(t *testing.T) (
	userstore.StoreInterface,
	sessionstore.StoreInterface,
	cachestore.StoreInterface,
	func(),
) {
	db, cleanupDB := SetupTestDatabase(t)
	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get SQL DB")

	// Initialize stores with in-memory SQLite
	userStore, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:                 sqlDB,
		UserTableName:      "users",
		AutomigrateEnabled: true,
		DebugEnabled:       false,
	})
	require.NoError(t, err, "Failed to create user store")

	sessionStore, err := sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:                 sqlDB,
		SessionTableName:   "sessions",
		AutomigrateEnabled: true,
		DebugEnabled:       false,
	})
	require.NoError(t, err, "Failed to create session store")

	cacheStore, err := cachestore.NewStore(cachestore.NewStoreOptions{
		DB:                 sqlDB,
		CacheTableName:     "caches",
		AutomigrateEnabled: true,
		DebugEnabled:       false,
	})
	require.NoError(t, err, "Failed to create cache store")

	return userStore, sessionStore, cacheStore, cleanupDB
}

// SetupMailServer creates a mock SMTP server for testing email functionality.
// Returns the server instance and a cleanup function to stop the server.
func SetupMailServer(t *testing.T) (*smtpmock.Server, func()) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false, // Set to true for debugging
		LogServerActivity: true,
		PortNumber:        2525, // Standard test SMTP port
		HostAddress:       "127.0.0.1",
	})

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start mock SMTP server: %v", err)
	}

	return server, func() {
		if err := server.Stop(); err != nil {
			log.Printf("Warning: failed to stop mock SMTP server: %v", err)
		}
	}
}
