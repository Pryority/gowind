// db.go
//
// PURPOSE:
// Use ConnectDB to build the connection to the database with configuration parameters

package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/wpcodevo/golang-fiber/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	// Construct the connection string using the provided configuration parameters.
	// The connection string includes the host, username, password, database name, port, and other necessary options.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort)

	// Open a connection to the PostgreSQL database using gorm.Open and the PostgreSQL driver.
	// The connection is assigned to the DB variable.
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Check if the connection was successful. If an error occurs,
		// the program will log an error message and exit.
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	// Execute a query to create the uuid-ossp extension if it doesn't exist in the database.
	// This extension is commonly used for generating UUIDs.
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Set the log mode of the database connection to logger.Info to enable logging of SQL statements.
	/// NOTE: Comment this line to hide/show logs
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")

	// Run database migrations using DB.AutoMigrate to automatically create or update the required
	// database tables based on the provided models.
	DB.AutoMigrate(&models.Note{})

	// Log a success message indicating that the connection to the database was established.
	log.Println("🚀\tConnected Successfully to the Database")
}
