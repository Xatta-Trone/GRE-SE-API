package database

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/database/seeders"
)

func InitSeeder(db *sqlx.DB) {
	seeder := flag.Bool("seed", false, "Run all the seeder")
	flag.Parse()

	color.Yellow("=== Seed the database ===")
	color.HiGreen(fmt.Sprintf("> Should seed ? > %t ", *seeder))

	if *seeder {
		// run all the seeders
		seeders.UserSeed(db)

	}

}
