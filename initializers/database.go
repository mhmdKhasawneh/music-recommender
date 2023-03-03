package initializers

import (
	"log"
	"os"

	"github.com/mhmdKhasawneh/musicrecommendationapp/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDb() {
	var err error

	dsn := os.Getenv("DB_CONNECTION")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to database.")
	}
}

func MigrateDatabase(){
	Db.AutoMigrate(&models.Recommendation{})
	Db.AutoMigrate(&models.User{})
}
