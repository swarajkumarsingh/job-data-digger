package constants

const REDIS_JOBS_LIST_KEY = "myList"

// Server ENV constants
const (
        ENV_PROD  = "prod"
        ENV_UAT   = "uat"
        ENV_DEV   = "dev"
        ENV_LOCAL = "local"
)

var MyMap = map[string]interface{}{
        "name":      "John",
        "age":       30,
        "isStudent": false,
}

var CareerPageLinks = map[string]interface{}{
        "Google": "https://www.google.com/about/careers/applications/jobs/results/?location=India",
        "FinBox": "https://finbox.freshteam.com/jobs",
}
