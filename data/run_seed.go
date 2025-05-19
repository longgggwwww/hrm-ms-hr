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

	// Initialize database client
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

// initDBClient initializes and returns an Ent database client.
func initDBClient() (*ent.Client, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, errMissingEnv("DB_URL")
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// runSeeders executes all seed functions.
func runSeeders(ctx context.Context, client *ent.Client) error {
	if err := seeds.SeedOrganizations(ctx, client); err != nil {
		return wrapError("Org", err)
	}

	if err := seeds.SeedDepartments(ctx, client); err != nil {
		return wrapError("Department", err)
	}

	if err := seeds.SeedPositions(ctx, client); err != nil {
		return wrapError("Position", err)
	}

	return nil
}

// errMissingEnv creates an error for missing environment variables.
func errMissingEnv(varName string) error {
	return fmt.Errorf("environment variable %s is not set", varName)
}

// wrapError wraps a seeding error with additional context.
func wrapError(entity string, err error) error {
	return fmt.Errorf("failed to seed %s: %w", entity, err)
}
