package utils

import (
	"fmt"
)

// ConnectionURLBuilder func for building URL connection.
func ConnectionURLBuilder(n string) (string, error) {
	// Define URL to connection.
	var url string

	// Switch given names.
	switch n {
	case "postgres":
		// URL for PostgreSQL connection.
		url = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conf.PostgresHost,
			conf.PostgresPort,
			conf.PostgresUser,
			conf.PostgresPassword,
			conf.PostgresDatabase,
		)
	case "fiber":
		// URL for Fiber connection.
		url = fmt.Sprintf(
			":%d",
			conf.ServerPort,
		)
	default:
		// Return error message.
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	// Return connection URL.
	return url, nil
}
