.PHONY: build run logs dockerstop
.SILENT: build run logs dockerstop

dev:
	nodemon --exec go run main.go

install:
	go mod tidy