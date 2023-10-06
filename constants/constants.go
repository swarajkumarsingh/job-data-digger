package constants

const GOOGLE_CAREER_PAGE_URL = "https://www.google.com/about/careers/applications/jobs/results/?location=India"

const REDIS_SCRAPE_DATA_KEY = "scrape-data-job-data-digger"

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
