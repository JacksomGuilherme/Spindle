package database

import (
	"fmt"

	"github.com/JacksomGuilherme/Spindle/configs"
	"github.com/JacksomGuilherme/Spindle/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnection(config *configs.Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%v:%v@/%v?charset=utf8&parseTime=True&loc=Local", config.DBUser, config.DBPassword, config.DBName)

	dataBase, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dataBase.AutoMigrate(&entity.User{})

	return dataBase, nil
}
