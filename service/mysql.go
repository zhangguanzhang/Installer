package service

import (
	"Installer/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

var db *gorm.DB

func DBInit(userPass, hostAndPortAndDBName string) error {
	var err error
	conInfo := strings.Split(hostAndPortAndDBName, "@")
	db, err = gorm.Open("mysql", userPass+"@tcp("+conInfo[0]+")/"+conInfo[1]+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}

	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(0)

	if !db.HasTable(&models.Machine{}) {
		if err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&models.Machine{}).Error; err != nil {
			return err
		}
	}

	db.LogMode(false)

	return db.AutoMigrate(&models.Machine{}).Error
}

func DBClose() error {
	return db.Close()
}

func ReturnColumnNames() []string {

	var (
		data []struct {
			ColumnName string `json:"column_name"`
		}
		result = make([]string, 0)
	)

	db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = ?", "machines").Scan(&data)

	for _, v := range data {
		result = append(result, v.ColumnName)
	}
	return result

}
