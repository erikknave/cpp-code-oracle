package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/erikknave/go-code-oracle/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	if DB != nil {
		return
	}
	var err error
	// DB, err = gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	DB, err = initMariaDB()

	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	DB.AutoMigrate(&types.User{})
	DB.AutoMigrate(&types.ChatMessage{})
	DB.AutoMigrate(&types.UserSearchResults{})
	DB.AutoMigrate(&types.UserAgentType{})
}

func initMariaDB() (*gorm.DB, error) {
	dbUser := os.Getenv("MARIA_DB_USER")
	dbPass := os.Getenv("MARIA_DB_PASSWORD")
	dbName := os.Getenv("MARIA_DB_NAME")
	dbHost := os.Getenv("MARIA_DB_HOST")
	dbPort := os.Getenv("MARIA_DB_PORT")
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	var db *gorm.DB
	var err error

	for attempts := 1; attempts <= 20; attempts++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		fmt.Printf("Attempt %d: Unable to connect to database. Retrying in 3 seconds...\n", attempts)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 20 attempts: %v", err)
}
