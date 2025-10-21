package utils

// Ensure SQL drivers are registered when utils is imported.
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)
