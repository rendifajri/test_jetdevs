package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// var DatabaseKey map[string]interface{}
var sqlDb *sql.DB

var ShareConnection *gorm.DB

func Connection() *gorm.DB {
	return poolDbConnection()
	//return singletonDbConnection()
}

func poolDbConnection() *gorm.DB {
	if sqlDb == nil {
		sqlDb = connectToDb()
		if sqlDb == nil {
			return nil
		}
	}
	log.Println("create connection by gorm framework")
	gormDb, err := gorm.Open("mysql", sqlDb)
	if err != nil {
		log.Fatalln("error on creating gorm connection ", err)
		return nil
	}
	log.Println("gorm connection is created")
	return gormDb
}

func connectToDb() *sql.DB {
	if sqlDb != nil {
		return sqlDb
	}
	log.Println("create pool database connection")
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	log.Println(dbURL)
	sqlDb, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Println("failed to connect database", err)
		return nil
	}
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(1 * time.Hour)
	log.Println("pool database connection is created")
	return sqlDb
}
