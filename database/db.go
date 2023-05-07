package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Gdb *sqlx.DB

func InitializeDB() *sqlx.DB {
	var count int64

	for {
		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))

		// Gdb = db
		if err != nil {
			log.Println("DB not yet ready....")
			count++
		} else {
			pingErr := db.Ping()

			if pingErr != nil {
				log.Fatal(pingErr)
			}
			fmt.Println("Connected to db!")
			return db
		}

		if count > 10 {
			log.Fatalln(err)
			return nil
		}

		log.Println("Backing off for 2s")
		time.Sleep(2 * time.Second)
		continue

	}

	// db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))

	// // Gdb = db
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // defer db.Close()

	// pingErr := db.Ping()

	// if pingErr != nil {
	// 	log.Fatal(pingErr)
	// }
	// fmt.Println("Connected to db!")
	// return db

}
