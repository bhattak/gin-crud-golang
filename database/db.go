package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

// InitDB initializes a MySQL database connection and creates a users table if it doesn't exist.
func InitDB() (*sql.DB, error) {
	godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PWD")
	dbName := os.Getenv("DB_NAME")
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)
	var err error

	db, err = sql.Open("mysql", dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createUserTable := `
        CREATE TABLE IF NOT EXISTS users (
            id INT NOT NULL AUTO_INCREMENT,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) NOT NULL,
            password VARCHAR(100) NOT NULL,
            PRIMARY KEY (id));`
	createAuthTable := `
		CREATE TABLE IF NOT EXISTS super_user (
			id INT NOT NULL AUTO_INCREMENT,
			email VARCHAR(100) NOT NULL,
			password VARCHAR(100) NOT NULL,
			PRIMARY KEY (id));
	`
	if _, err = db.Exec(createUserTable); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	if _, err = db.Exec(createAuthTable); err != nil {
		log.Fatalf("Failed to create super_user table: %v", err)
	}

	return db, nil
}

// func initDb() (*gorm.DB, error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "dbUser", "dbPwd", "dbHost", "dbPort", "dbName")
// 	db, err := gorm.Open(sql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to database: %v", err)
// 	}

// 	// Migrate the schema
// 	err = db.AutoMigrate(&model.User{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to migrate database schema: %v", err)
// 	}

// 	return db, nil
// }
