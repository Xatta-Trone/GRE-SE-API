package seeders

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
)

func UserSeed(db *sqlx.DB) {
	const Total = 1000
	color.Yellow("=== Inside the users seeder ===")
	color.Yellow(fmt.Sprintf("=== Seeding %d users ===", Total))

	users := []model.UserModel{}

	for i := 1; i <= Total; i++ {
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

	color.Green(fmt.Sprintf("Total imported users: %d", int(totalImported)))

	// fmt.Println(users)

}
