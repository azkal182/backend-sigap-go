package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/your-org/go-backend-starter/internal/infrastructure/database" // Import to register migrations
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse command line flags
	command := flag.String("command", "up", "Migration command: up, down, status, or to")
	version := flag.String("version", "", "Target version for 'to' command")
	flag.Parse()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Execute migration command
	switch *command {
	case "up":
		if err := database.MigrateUp(database.DB); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("\n✅ Migrations completed successfully")

	case "down":
		if err := database.MigrateDown(database.DB); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("\n✅ Migration rolled back successfully")

	case "status":
		status, err := database.GetMigrationStatus(database.DB)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

		fmt.Println("\n📊 Migration Status:")
		fmt.Println("===================")
		for _, s := range status {
			applied := "❌ Not Applied"
			if s["applied"].(bool) {
				applied = "✅ Applied"
			}
			fmt.Printf("  %s - %s: %s\n", s["version"], s["name"], applied)
		}
		fmt.Println()

	case "to":
		if *version == "" {
			log.Fatal("Version is required for 'to' command. Use -version flag")
		}
		if err := database.MigrateToVersion(database.DB, *version); err != nil {
			log.Fatalf("Failed to migrate to version %s: %v", *version, err)
		}
		fmt.Printf("\n✅ Migrated to version %s successfully\n", *version)

	default:
		fmt.Fprintf(os.Stderr, "❌ Unknown command: %s\n\n", *command)
		fmt.Fprintf(os.Stderr, "Usage: %s -command [up|down|status|to] [-version VERSION]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  up      - Apply all pending migrations\n")
		fmt.Fprintf(os.Stderr, "  down    - Rollback the last migration\n")
		fmt.Fprintf(os.Stderr, "  status  - Show migration status\n")
		fmt.Fprintf(os.Stderr, "  to      - Migrate to specific version (requires -version flag)\n")
		os.Exit(1)
	}
}
