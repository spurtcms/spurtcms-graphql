package dbconfig

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgresql db connection
func SetupDB() *gorm.DB {

	er := godotenv.Load()

	if er != nil {

		log.Fatalf("Error loading .env file")
	}

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD")) //Build connection string

	log.Println(dbUri)

	DB, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})

	if err != nil {

		fmt.Println("Status:", err)

		panic(err)
	}
	
	return DB
}


