package initializers

import (
	"fmt"
	"log"

	// "github.com/farhan-nahid/email-service/initializers"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type DatabaseConfiguration struct {
	DBName   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func ConnectToDatabase() {
	err := godotenv.Load()
	
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	// dbUser 	   := os.Getenv("postgres") 
	// dbPassword := os.Getenv("DB_PASS")
	// dbHost 	   := os.Getenv("DB_HOST")
	// DBName     := os.Getenv("DB_NAME")
	// dbPort     := os.Getenv("DB_PORT")

	dbUser := "postgres"
	dbPassword := "postgres"
	dbHost := "localhost"
	DBName := "postgres"
	dbPort := "5432"

	log.Println("DB_USER: ", "postgres", "DB_PASS: ", "postgres", "DB_HOST: ", "localhost", "DB_NAME: ", "postgres", "DB_PORT: ", "5432")

	log.Println("Attempting to connect to db")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, DBName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Println(err.Error())
		panic("Failed to Connect Database !")
	} else {
		log.Println("Database Connected Successfully !")
	}

	DB = db
}
