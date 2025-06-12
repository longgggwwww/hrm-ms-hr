package main

import (
	"context"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/longgggwwww/hrm-ms-hr/ent"
)

func main() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	cli, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()
	if err := cli.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("Migration complete!")
}
