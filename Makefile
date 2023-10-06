install:
	go mod tidy

run:
	docker build -t job-data-scraper . && docker run -p 8080:8080 job-data-scraper

build:
	docker build -t job-data-scraper

start:
	docker run -p 8080:8080 job-data-scraper

dev:
	nodemon --exec go run main.go

deploy: echo "TODO"

test: echo "TODO"

.PHONY: build run logs dockerstop
.SILENT: build run logs dockerstop