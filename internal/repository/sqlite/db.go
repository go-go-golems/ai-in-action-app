package sqlite

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBManager manages the database connection
type DBManager struct {
	db *gorm.DB
}

// NewDBManager creates a new database manager
func NewDBManager(dbPath string) (*DBManager, error) {
	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create a new DBManager
	manager := &DBManager{
		db: db,
	}

	// Auto migrate the schema
	if err := manager.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return manager, nil
}

// Migrate automatically migrates the database schema
func (m *DBManager) Migrate() error {
	log.Println("Migrating database schema...")

	// Auto migrate all models
	err := m.db.AutoMigrate(
		&EventModel{},
		&TimerModel{},
		&NoteModel{},
		&QuestionModel{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// GetDB returns the GORM database instance
func (m *DBManager) GetDB() *gorm.DB {
	return m.db
}

// Close closes the database connection
func (m *DBManager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
