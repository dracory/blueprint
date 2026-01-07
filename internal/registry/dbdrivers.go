package registry

// Register SQL drivers via blank imports so database/sql recognizes them.
// Keep this file in the app package to ensure it’s linked into all binaries.

import (
	_ "github.com/go-sql-driver/mysql"
)
