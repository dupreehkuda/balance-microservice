.PHONY: run
run:
	go run cmd/main.go

.PHONY: compose
compose:
	docker-compose up -d

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans