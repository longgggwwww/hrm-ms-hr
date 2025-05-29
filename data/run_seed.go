package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/longgggwwww/hrm-ms-hr/data/seeds"
	"github.com/longgggwwww/hrm-ms-hr/ent"
)

func main() {
	log.Println("Starting seeding process...")

	client, err := initDBClient()
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer client.Close()

	// Run seeders
	if err := runSeeders(context.Background(), client); err != nil {
		log.Fatalf("Seeding process failed: %v", err)
	}

	log.Println("Seeding process completed successfully.")
}

func initDBClient() (*ent.Client, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("environment variable DB_URL is not set")
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func runSeeders(ctx context.Context, client *ent.Client) error {
	seeders := []struct {
		name string
		fn   func(context.Context, *ent.Client) error
	}{
		{"Organization", seeds.SeedOrganizations},
		{"Department", seeds.SeedDepartments},
		{"Position", seeds.SeedPositions},
		{"Label", seeds.SeedLabels},
	}

	for _, seeder := range seeders {
		if err := seeder.fn(ctx, client); err != nil {
			return fmt.Errorf("failed to seed %s: %w", seeder.name, err)
		}
	}
	return nil
}
