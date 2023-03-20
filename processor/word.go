package processor

import (
	"fmt"
	"log"

	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
)

func ReadTableAndProcessWord(word string) {

	fmt.Println("getting word result for ", word)

	rs := database.Gdb.QueryRowx("SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus`, `ninja` FROM `wordlist` WHERE `word` = ?;", word)

	var r model.Result

	if rs.Err() != nil {
		log.Fatal(rs.Err())
	}

	err := rs.StructScan(&r)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r.ID)
	fmt.Println(r.Word)
	fmt.Println(r.Google.PartsOfSpeeches)
	// fmt.Println(r.Wiki)
	fmt.Println(r.WordsApi)
	// fmt.Println(r.Thesaurus)
	// fmt.Println(r.Ninja)

}