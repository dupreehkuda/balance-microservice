.PHONY:
run:
	go run cmd/main.go

.PHONY:
compose:
	docker-compose up -d

.PHONY:
compose-down:
	docker-compose down --remove-orphans

.PHONY:
test:
	go test -v -cover -race -count 1 ./...