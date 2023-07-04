package processor

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/scrapper"
	"github.com/xatta-trone/words-combinator/utils"
)

func ProcessPendingWords(db *sqlx.DB) {

	// word words to process
	pendingWords := []model.PendingWordModel{}

	err := db.Select(&pendingWords, "SELECT * FROM `pending_words` WHERE `tried` = 0 and `approved` = 1")

	if err != nil {
		utils.Errorf(err)
		utils.PrintR(err.Error())
		return
	}

	if len(pendingWords) == 0 {
		utils.PrintR("No words found")
		return
	}

	fmt.Println(len(pendingWords))

	for _, pendingWord := range pendingWords {

		// check in wordlist

		// created in word list
		wordListData, err := InsertIntoWordListTable(db, pendingWord.Word)

		if err != nil {
			UpdateAsTried(pendingWord, db)
			utils.PrintR("error in insert")
			return
		}

		// now get the data from the internet

		// var wg *sync.WaitGroup

		// wg.Add(5)

		// go scrapper.GetGoogleResultAndSaveWithWG(db, wordListData, wg)
		// go scrapper.GetWikiResultAndSaveWithWg(db, wordListData, wg)
		// go scrapper.GetThesaurusResultAndSaveWithWg(db, wordListData, wg)
		// go scrapper.GetWordsResultAndSaveWithWg(db, wordListData, wg)
		// go scrapper.GetMWResultAndSaveWithWg(db, wordListData, wg)

		scrapper.GetGoogleResultAndSave(db, wordListData)
		scrapper.GetWikiResultAndSave(db, wordListData)
		scrapper.GetThesaurusResultAndSave(db, wordListData)
		scrapper.GetWordsResultAndSave(db, wordListData)
		scrapper.GetMWResultAndSave(db, wordListData)

		// wg.Wait()

		// now get the data again
		updatedRawData, err := GetFromWordListTable(db, pendingWord.Word)
		if err != nil {
			UpdateAsTried(pendingWord, db)
			utils.PrintR("error in updated data")
			continue
		}

		fmt.Println(updatedRawData.Google)

		// now process the data
		processedData, err := ProcessWordData(db, updatedRawData)

		if err != nil {
			utils.Errorf(err)
			utils.PrintR(err.Error())
			UpdateAsTried(pendingWord, db)
			SetTried(db, updatedRawData.ID)
			continue
		}

		// save new word
		err = SaveNewWordData(db, updatedRawData.Word, int64(updatedRawData.ID), processedData)

		if err != nil {
			utils.Errorf(err)
			utils.PrintR(err.Error())
			UpdateAsTried(pendingWord, db)
			SetTried(db, updatedRawData.ID)
			continue
		}

		utils.PrintG(fmt.Sprintf("processed word id %d %s ", updatedRawData.ID, updatedRawData.Word))

		// all done.. now save the new data with the relation and delete the record
		wordsTableId := GetWordId(db, wordListData.Word)

		//
		if wordsTableId != 0 {
			// save to word list relation
			InsertListWordRelation(int64(wordsTableId), int64(pendingWord.ListId), db)
			// delete from the pending list
			DeletePendingWords(pendingWord,db)
			// db.Exec("Delete from `pending_words` where word=?, list_id=?", pendingWord.Word, pendingWord.ListId)

		}

	}

}

func UpdateAsTried(data model.PendingWordModel, db *sqlx.DB) {
	tx, err := db.Begin()
	if err != nil {
		utils.Errorf(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update pending_words set tried=1 where word=? and list_id=?", data.Word, data.ListId)

	if err != nil {
		utils.Errorf(err)
	}

	err = tx.Commit()

	if err != nil {
		utils.Errorf(err)
	}
}

func DeletePendingWords(data model.PendingWordModel, db *sqlx.DB) {
	tx, err := db.Begin()
	if err != nil {
		utils.Errorf(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Delete from pending_words where word=? and list_id=?", data.Word, data.ListId)

	if err != nil {
		utils.Errorf(err)
	}

	err = tx.Commit()

	if err != nil {
		utils.Errorf(err)
	}
}

func GetWordId(db *sqlx.DB, word string) uint64 {
	var Id uint64
	db.Get(&Id, "SELECT id FROM words where word=? LIMIT 1", word)

	return Id
}

// :keep
func InsertListWordRelation(wordId, listId int64, db *sqlx.DB) error {

	queryMap := map[string]interface{}{"word_id": wordId, "list_id": listId, "created_at": time.Now().UTC()}

	res, err := db.NamedExec("Insert ignore into list_word_relation(word_id,list_id,created_at) values(:word_id,:list_id,:created_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return err
	}

	lastId, err := res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
		return err
	}

	if lastId == 0 {
		return fmt.Errorf("there was a problem with the insertion. rows affected: %d", lastId)
	}

	return nil

}
