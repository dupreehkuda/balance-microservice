.PHONY:
run:
	go run cmd/main.go

dc:
	docker-compose up

rebuild:
	docker-compose build --no-cache