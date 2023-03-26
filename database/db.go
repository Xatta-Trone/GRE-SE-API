package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Gdb *sqlx.DB

func InitializeDB() *sqlx.DB {

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))

	// Gdb = db
	if err != nil {
		log.Fatalln(err)
	}

	// defer db.Close()

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to db!")
	return db

}
