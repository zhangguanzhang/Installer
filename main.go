package main

import (
	"Installer/router"
	"Installer/service"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// @title Installer API
// @version 1.0
// @description some api of this
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {

	port := flag.Int("port", 8080, "http server port")
	dbInfo := flag.String("db", "127.0.0.1:3306@pxe", "database host with port, like ip:host@database")
	con := flag.String("user", "root:zhangguanzhang", "user:pass connect to db")
	ks := flag.String("ks", "templates/ks.tmpl", "kickstart template file")

	flag.Parse()

	_, err := os.Stat(*ks)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("kickstart template file %s %v", *ks, err)
	}

	if err := service.DBInit(*con, *dbInfo); err != nil {
		log.Fatal(err)
	}

	defer service.DBClose()

	Router := router.InitRouter(*ks)
	if err := Router.Run(":" + strconv.Itoa(*port)); err != nil {
		log.Fatal(err)
	}
}
