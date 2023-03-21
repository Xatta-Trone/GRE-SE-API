package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/processor"
	"github.com/xatta-trone/words-combinator/utils"
)

// for named exec
type WordStruct struct {
	Word string
}

// const insertQuery = "Insert ignore into words(word,created_at) values(?,now()); "
const insertQuery = "INSERT INTO words(word,created_at) SELECT :word, now() WHERE NOT EXISTS (SELECT word FROM words WHERE word = :word); "

func ReadAndImportNamedCsv(fileName, groupName string) {

	fmt.Printf("Reading %s file with the name %s \n", fileName, groupName)

	// create the group name
	result, err := database.Gdb.Exec("Insert into word_groups(name,created_at) values(?,now())", groupName)

	if err != nil {
		log.Fatal(err)
	}

	groupId, _ := result.LastInsertId()

	fmt.Println("Inserted word group into table with id", groupId)

	f, err := os.Open(fileName)

	if err != nil {
		log.Fatal("Could not open the file", fileName, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	for {
		data, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(data)
		processedWord := utils.ProcessWord(data[0])

		if len(processedWord) > 1 {
			for _, w := range processedWord {
				w := WordStruct{Word: w}

				res, err := database.Gdb.NamedExec(insertQuery, w)

				if err != nil {
					log.Fatal(err)
				}

				id, err := res.LastInsertId()
				if err != nil {
					log.Fatal("last insert", err)
				}

				count, err := res.RowsAffected()
				if err != nil {
					log.Fatal("rows affected", err)
				}

				fmt.Println("last inserted id", id, "rows affected", count)

				// now insert into group relation

				insertGroupRelation(id, groupId)
				processNewWord(w.Word)

			}

		} else {
			w := WordStruct{Word: processedWord[0]}

			res, err := database.Gdb.NamedExec(insertQuery, w)

			if err != nil {
				log.Fatal(err)
			}

			id, err := res.LastInsertId()
			if err != nil {
				log.Fatal("last insert", err)
			}

			count, err := res.RowsAffected()
			if err != nil {
				log.Fatal("rows affected", err)
			}

			fmt.Println("last inserted id", id, "rows affected", count)
			// now insert into group relation

			insertGroupRelation(id, groupId)
			processNewWord(w.Word)

		}

		// res, err := database.Gdb.Exec("Insert IGNORE into words(word) values(?)", processedWord[0])

		// if err != nil {
		// 	log.Fatal(err)
		// }

		// id, _ := res.LastInsertId()
		// // now insert into group relation

		// fmt.Println(id)

		// insertGroupRelation(id, groupId)

	}

}

func insertGroupRelation(wordId, groupId int64) {
	if wordId != 0 {
		_, err := database.Gdb.Exec("Insert into word_group_relation(word_id,word_group_id,created_at) values(?,?,now())", wordId, groupId)

		if err != nil {
			log.Fatal("grp", err)
		}

		fmt.Printf("Inserted word %d group %d", wordId, groupId)

	} else {
		fmt.Println("===duplicate found===")
	}

}

func processNewWord(word string) {
	// now process the word
	processor.ReadTableAndProcessWord(word)
}

func readRemoteFile() {

	// truncate the table first
	// _, err := database.Gdb.Exec("TRUNCATE TABLE wordlist")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	res, err := http.Get("https://raw.githubusercontent.com/Xatta-Trone/gre-words-collection/main/word-list/combined.csv")

	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	reader := csv.NewReader(res.Body)

	totalRows := 0

	for {
		data, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		_, err = database.Gdb.Exec("INSERT IGNORE INTO wordlist(word) values (?)", data[0])

		if err != nil {
			fmt.Println(err)
		}

		totalRows++

		// fmt.Println(data[0])
	}

	fmt.Println("Total rows inserted", totalRows)

}
