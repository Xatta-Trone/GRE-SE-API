package processor

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/utils"
)

func UpdateUnUpdatedWords(db *sqlx.DB) {

	// word id's to process
	wordListIds := []uint64{}

	err := db.Select(&wordListIds, "SELECT id FROM `wordlist` WHERE `tried` = 0 and `in_words` = 0 and google is not null")

	if err != nil {
		utils.Errorf(err)
		utils.PrintR(err.Error())
		return
	}

	if len(wordListIds) == 0 {
		utils.PrintR("No words found")
		return
	}

	// process word by word

	for _, id := range wordListIds {
		utils.PrintS(fmt.Sprintf("processing word id %d ", id))

		// get the word data
		unprocessedWordData, err := GetUnProcessedWordDataById(db, id)

		if err != nil {
			utils.Errorf(err)
			utils.PrintR(err.Error())
			SetTried(db, id)
			continue
		}

		// now process the data
		processedData, err := ProcessWordData(db, unprocessedWordData)

		if err != nil {
			utils.Errorf(err)
			utils.PrintR(err.Error())
			SetTried(db, id)
			continue
		}

		// save new word
		err = SaveNewWordData(db, unprocessedWordData.Word, int64(unprocessedWordData.ID), processedData)

		if err != nil {
			utils.Errorf(err)
			utils.PrintR(err.Error())
			SetTried(db, id)
			continue
		}

		utils.PrintG(fmt.Sprintf("processed word id %d %s ", id, unprocessedWordData.Word))

	}

}
