package db

import (
	"fiber_rest_api/model"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()
	dbhost := os.Getenv("MYSQL_HOST")
	dbuser := os.Getenv("MYSQL_USER")
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DBNAME")
	dbport := os.Getenv("MYSQL_PORT")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpassword, dbhost, dbport, dbname)
	db, err := gorm.Open(mysql.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("db connection failed")
	}
	DB = db
	fmt.Println("db connected successfully")

	/*err = AutoMigrate(db)
	if err != nil {
		fmt.Println(err)
	}*/
}

func AutoMigrate(connection *gorm.DB) error {
	err := connection.Debug().AutoMigrate(
		&model.Cashier{},
		&model.Category{},
		&model.Payment{},
		&model.PaymentType{},
		&model.Product{},
		&model.Discount{},
		&model.Order{},
	)
	if err != nil {
		return err
	}
	return nil
}
