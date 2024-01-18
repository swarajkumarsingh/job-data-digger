.PHONY: build run logs dockerstop
.SILENT: build run logs dockerstop

run:
	docker compose build
	docker compose up

build:
	docker build -t job-data-scraper . && docker run -p 8080:8080 job-data-scraper

start:
	docker run -p 8080:8080 job-data-scraper

dev:
	nodemon --exec go run main.go

install:
	go mod tidy

deploy: 
	echo "TODO"

test: 
	echo "TODO"