package conf

import "os"

var ENV string = os.Getenv("STAGE")
var SentryDSN string = os.Getenv("SENTRY_DSN")

// Server ENV constants
const ENV_DEV = "dev"
const ENV_PROD = "prod"