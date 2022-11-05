package main

import "github.com/dupreehkuda/balance-microservice/internal/api"

func main() {
	srv := api.NewByConfig()
	srv.Run()
}
