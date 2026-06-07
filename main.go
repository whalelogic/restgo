package main

import (
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	h := &EmployeeHandler{
		rdb: redis.NewClient(&redis.Options{Addr: addr}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /employees", h.List)
	mux.HandleFunc("GET /employees/{id}", h.Get)
	mux.HandleFunc("POST /employees", h.Create)
	mux.HandleFunc("PUT /employees/{id}", h.Update)
	mux.HandleFunc("PATCH /employees/{id}", h.Patch)
	mux.HandleFunc("DELETE /employees/{id}", h.Delete)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
