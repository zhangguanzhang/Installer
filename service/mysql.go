package service

import (
	"Installer/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:zhangguanzhang@tcp(127.0.0.1:3306)/pxe?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(0)

	if !db.HasTable(&models.Machines{}) {
		if err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&models.Machines{}).Error; err != nil {
			panic(err)
		}
	}
}
