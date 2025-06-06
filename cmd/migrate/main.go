package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"github.com/longgggwwww/hrm-ms-hr/ent"
)

func main() {
	client, err := ent.Open("postgres", "host=192.168.1.24 port=5433 user=root password=123456 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("Migration complete!")
}
