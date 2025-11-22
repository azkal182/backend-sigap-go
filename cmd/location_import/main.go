// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/joho/godotenv"
// 	"github.com/your-org/go-backend-starter/internal/domain/entity"
// 	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
// )

// const (
// 	locationsBaseDir  = "data/locations"
// 	provincesFileName = "provinces.json"
// 	regenciesFileName = "regencies.json"
// 	districtsFileName = "districts.json"
// 	villagesFileName  = "villages.json"
// )

// func main() {
// 	// Load environment variables
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using environment variables")
// 	}

// 	// Connect to database
// 	if err := database.Connect(); err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	ctx := context.Background()

// 	// Import in hierarchical order
// 	if err := importProvinces(ctx); err != nil {
// 		log.Fatalf("Failed to import provinces: %v", err)
// 	}
// 	if err := importRegencies(ctx); err != nil {
// 		log.Fatalf("Failed to import regencies: %v", err)
// 	}
// 	if err := importDistricts(ctx); err != nil {
// 		log.Fatalf("Failed to import districts: %v", err)
// 	}
// 	if err := importVillages(ctx); err != nil {
// 		log.Fatalf("Failed to import villages: %v", err)
// 	}

// 	log.Println("Location data imported successfully")
// }

// func openJSON(path string) (*os.File, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return f, nil
// }

// func importProvinces(ctx context.Context) error {
// 	path := filepath.Join(locationsBaseDir, provincesFileName)
// 	log.Printf("Importing provinces from %s", path)

// 	f, err := openJSON(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	var records []entity.Province
// 	if err := json.NewDecoder(f).Decode(&records); err != nil {
// 		return err
// 	}

// 	for _, rec := range records {
// 		var existing entity.Province
// 		// Check by ID; if not found, insert
// 		result := database.DB.WithContext(ctx).Where("id = ?", rec.ID).First(&existing)
// 		if result.Error == nil {
// 			continue
// 		}
// 		if result.Error != nil && result.Error.Error() != "record not found" {
// 			// For gorm v2, ErrRecordNotFound is preferred, but we compare string to avoid extra import
// 			log.Printf("Failed to check province id=%d: %v", rec.ID, result.Error)
// 			continue
// 		}
// 		if err := database.DB.WithContext(ctx).Create(&rec).Error; err != nil {
// 			log.Printf("Failed to insert province id=%d: %v", rec.ID, err)
// 		}
// 	}

// 	log.Printf("Imported %d provinces (idempotent)", len(records))
// 	return nil
// }

// func importRegencies(ctx context.Context) error {
// 	path := filepath.Join(locationsBaseDir, regenciesFileName)
// 	log.Printf("Importing regencies from %s", path)

// 	f, err := openJSON(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	var records []entity.Regency
// 	if err := json.NewDecoder(f).Decode(&records); err != nil {
// 		return err
// 	}

// 	for _, rec := range records {
// 		var existing entity.Regency
// 		result := database.DB.WithContext(ctx).Where("id = ?", rec.ID).First(&existing)
// 		if result.Error == nil {
// 			continue
// 		}
// 		if result.Error != nil && result.Error.Error() != "record not found" {
// 			log.Printf("Failed to check regency id=%d: %v", rec.ID, result.Error)
// 			continue
// 		}
// 		if err := database.DB.WithContext(ctx).Create(&rec).Error; err != nil {
// 			log.Printf("Failed to insert regency id=%d: %v", rec.ID, err)
// 		}
// 	}

// 	log.Printf("Imported %d regencies (idempotent)", len(records))
// 	return nil
// }

// func importDistricts(ctx context.Context) error {
// 	path := filepath.Join(locationsBaseDir, districtsFileName)
// 	log.Printf("Importing districts from %s", path)

// 	f, err := openJSON(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	var records []entity.District
// 	if err := json.NewDecoder(f).Decode(&records); err != nil {
// 		return err
// 	}

// 	for _, rec := range records {
// 		var existing entity.District
// 		result := database.DB.WithContext(ctx).Where("id = ?", rec.ID).First(&existing)
// 		if result.Error == nil {
// 			continue
// 		}
// 		if result.Error != nil && result.Error.Error() != "record not found" {
// 			log.Printf("Failed to check district id=%d: %v", rec.ID, result.Error)
// 			continue
// 		}
// 		if err := database.DB.WithContext(ctx).Create(&rec).Error; err != nil {
// 			log.Printf("Failed to insert district id=%d: %v", rec.ID, err)
// 		}
// 	}

// 	log.Printf("Imported %d districts (idempotent)", len(records))
// 	return nil
// }

// func importVillages(ctx context.Context) error {
// 	path := filepath.Join(locationsBaseDir, villagesFileName)
// 	log.Printf("Importing villages from %s", path)

// 	f, err := openJSON(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	var records []entity.Village
// 	if err := json.NewDecoder(f).Decode(&records); err != nil {
// 		return err
// 	}

// 	for _, rec := range records {
// 		var existing entity.Village
// 		result := database.DB.WithContext(ctx).Where("id = ?", rec.ID).First(&existing)
// 		if result.Error == nil {
// 			continue
// 		}
// 		if result.Error != nil && result.Error.Error() != "record not found" {
// 			log.Printf("Failed to check village id=%d: %v", rec.ID, result.Error)
// 			continue
// 		}
// 		if err := database.DB.WithContext(ctx).Create(&rec).Error; err != nil {
// 			log.Printf("Failed to insert village id=%d: %v", rec.ID, err)
// 		}
// 	}

// 	log.Printf("Imported %d villages (idempotent)", len(records))
// 	return nil
// }

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	locationsBaseDir  = "data/locations"
	provincesFileName = "provinces.json"
	regenciesFileName = "regencies.json"
	districtsFileName = "districts.json"
	villagesFileName  = "villages.json"
	batchSize         = 1000 // Adjust based on your data size
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	// Import in hierarchical order
	if err := importProvinces(ctx); err != nil {
		log.Fatalf("Failed to import provinces: %v", err)
	}
	if err := importRegencies(ctx); err != nil {
		log.Fatalf("Failed to import regencies: %v", err)
	}
	if err := importDistricts(ctx); err != nil {
		log.Fatalf("Failed to import districts: %v", err)
	}
	if err := importVillages(ctx); err != nil {
		log.Fatalf("Failed to import villages: %v", err)
	}

	log.Println("Location data imported successfully")
}

func importProvinces(ctx context.Context) error {
	path := filepath.Join(locationsBaseDir, provincesFileName)
	log.Printf("Importing provinces from %s", path)

	var records []entity.Province
	if err := loadJSON(path, &records); err != nil {
		return err
	}

	// Use UPSERT for idempotent import (ON CONFLICT DO NOTHING)
	if err := bulkUpsert(ctx, records); err != nil {
		return err
	}

	log.Printf("Imported %d provinces", len(records))
	return nil
}

func importRegencies(ctx context.Context) error {
	path := filepath.Join(locationsBaseDir, regenciesFileName)
	log.Printf("Importing regencies from %s", path)

	var records []entity.Regency
	if err := loadJSON(path, &records); err != nil {
		return err
	}

	if err := bulkUpsert(ctx, records); err != nil {
		return err
	}

	log.Printf("Imported %d regencies", len(records))
	return nil
}

func importDistricts(ctx context.Context) error {
	path := filepath.Join(locationsBaseDir, districtsFileName)
	log.Printf("Importing districts from %s", path)

	var records []entity.District
	if err := loadJSON(path, &records); err != nil {
		return err
	}

	if err := bulkUpsert(ctx, records); err != nil {
		return err
	}

	log.Printf("Imported %d districts", len(records))
	return nil
}

func importVillages(ctx context.Context) error {
	path := filepath.Join(locationsBaseDir, villagesFileName)
	log.Printf("Importing villages from %s", path)

	var records []entity.Village
	if err := loadJSON(path, &records); err != nil {
		return err
	}

	if err := bulkUpsert(ctx, records); err != nil {
		return err
	}

	log.Printf("Imported %d villages", len(records))
	return nil
}

// Generic helper to load JSON
func loadJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}

// Generic bulk upsert with batching
func bulkUpsert(ctx context.Context, records interface{}) error {
	// Use transaction for better performance
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clauses: ON CONFLICT (id) DO NOTHING for idempotency
		// This prevents duplicate key errors and skips existing records
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true, // Skip if exists
		}).CreateInBatches(records, batchSize).Error
	})
}

// Alternative: If you want to UPDATE existing records instead
func bulkUpsertWithUpdate(ctx context.Context, records interface{}) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// This will update all fields except ID on conflict
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true, // Update all fields if exists
		}).CreateInBatches(records, batchSize).Error
	})
}

// Alternative approach using raw SQL for maximum performance
func bulkUpsertRaw(ctx context.Context, tableName string, records []map[string]interface{}) error {
	if len(records) == 0 {
		return nil
	}

	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Build batch insert with ON CONFLICT
		// This is database-specific, example for PostgreSQL:
		// INSERT INTO table (...) VALUES (...) ON CONFLICT (id) DO NOTHING

		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}
			batch := records[i:end]

			if err := tx.Table(tableName).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(batch).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
