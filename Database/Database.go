package Database

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/ini.v1"
	"os"
)

func ConnectToDB(host string, port int, user string,
	password string, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return db, nil
}

func ConnectToDBViaConfig(configPath string) (*sql.DB, error) {
	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	dbSection := cfg.Section("database")
	host := dbSection.Key("host").String()
	port, _ := dbSection.Key("port").Int()
	user := dbSection.Key("user").String()
	passwd := dbSection.Key("password").String()
	dbname := dbSection.Key("db_name").String()

	return ConnectToDB(host, port, user, passwd, dbname)
}

func ConnectDBViaGORM(configPath string) (*gorm.DB, error) {
	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	dbSection := cfg.Section("database")
	host := dbSection.Key("host").String()
	port, _ := dbSection.Key("port").Int()
	user := dbSection.Key("user").String()
	passwd := dbSection.Key("password").String()
	dbname := dbSection.Key("db_name").String()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbname)

	return gorm.Open("postgres", psqlInfo)
}
