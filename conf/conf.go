package conf

import "os"

var ColName string = "links"
var DbName string = "ZipLink"
var ConnectionString string = os.Getenv("MONGO_URL")

var ENV string = getEnv()
var SentryDSN string = os.Getenv("SENTRY_DSN")

// Server ENV constants
const ENV_PROD = "prod"
const ENV_DEV = "dev"

func getEnv() string {
	value := os.Getenv("STAGE")
	if value == "" {
		return ENV_DEV
	}
	return value
}
