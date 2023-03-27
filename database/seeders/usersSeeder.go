package seeders

import (
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
)

func UserSeed(db *sqlx.DB) {
	const Total = 10000
	color.Yellow("=== Inside the users seeder ===")
	color.Yellow(fmt.Sprintf("=== Seeding %d users ===", Total))

	// batch import per 1000 records
	totalRun := int(math.Ceil(Total / 1000))

	for i := 0; i < totalRun; i++ {

		users := []model.UserModel{}

		for i := 1; i <= 1000; i++ {
			user := model.UserModel{
				Name:  faker.Name(),
				Email: faker.Email(options.WithGenerateUniqueValues(true)),
				// Email:    fmt.Sprintf("test-email%d@example.com", i),
				UserName:  faker.Username(options.WithGenerateUniqueValues(true)),
				CreatedAt: time.Now().UTC().Local(),
			}

			// fmt.Println(user)

			users = append(users, user)

		}

		query := `INSERT Ignore INTO users (name,email,username,created_at) VALUES (:name,:email,:username,:created_at)`
		res, err := db.NamedExec(query, users)
		if err != nil {
			fmt.Println("err", err)
		}

		totalImported, err := res.RowsAffected()

		if err != nil {
			fmt.Println("err", err)
		}

		color.Green(fmt.Sprintf("Imported users: %d", int(totalImported)))

	}

	color.Green(fmt.Sprintf("Total Imported users: %d", int(Total)))

}
