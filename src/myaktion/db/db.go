package db

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	err := Connect(os.Getenv("DB_CONNECT"))
	if err != nil {
		panic(err)
	}
}

func Connect(connect string) error {
	dsn := fmt.Sprintf("root:root@tcp(%s)/myaktion?charset=utf8&parseTime=True&loc=Local", connect)
	log.Info("Using database connection string: ", dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")
	}
	log.Info("Starting automatic migration")
	if err := DB.Debug().AutoMigrate(&model.Campaign{}, &model.Donation{}); err != nil {
		return err
	}
	log.Info("Finished automatic migration")
	return nil
}
