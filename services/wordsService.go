package services

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
)

type WordService struct {
	Db *sqlx.DB
}

var wordService WordService

func NewWordService(db *sqlx.DB) {
	wordService.Db = db
}

type WordResponseStruct struct {
	Word     string
	WordData model.CombinedWithWord `db:"word_data"`
}

func WordsIndex(r requests.WordIndexReqStruct) []WordResponseStruct {
	fmt.Println(r.ID, r.Query)

	// fetch at most 10 place names
	var words []WordResponseStruct

	// queryString := ``

	// fmt.Println(queryString)

	err := wordService.Db.Select(&words, "SELECT word,word_data FROM words where word like ? limit ?", "%"+r.Query+"%", 10)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(words)

	return words
}
