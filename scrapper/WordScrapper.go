package scrapper

import (
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
)

// ScrapWordById accepts the model of wordlist then scraps the data from the internet and saves into the db
func ScrapWordById(db *sqlx.DB, word model.Result) {
	// get the google data
	GetGoogleResultAndSave(db, word)
	// get the wiki result
	GetWikiResultAndSave(db, word)
	// get the thesaurus result
	GetThesaurusResultAndSave(db, word)
	// get the words api result
	GetWordsResultAndSave(db, word)
	// get the mw result
	GetMWResultAndSave(db,word)

}
