package processor

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func ProcessSingleWordData(db *sqlx.DB, word model.WordModel) {

	// check if the there is data in the wordlist table

	data, err := CheckWordListTable(db, word)

	if err != nil {
		fmt.Println(err)
	}

	// if yes process the data
	processedData, _ := ProcessWordData(db, data)

	// if no then insert into wordlist and get the wordlist data

	// then process the data

	// save to the database
	SaveProcessedDataToWordTable(db, word, processedData)

}

func CheckWordListTable(db *sqlx.DB, word model.WordModel) (model.Result, error) {

	query := "SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus` FROM `wordlist` WHERE `word` = ?;"

	result := db.QueryRowx(query, word.Word)

	var model model.Result

	err := result.StructScan(&model)

	if err != nil {
		fmt.Println(err)
		return model, err
	}

	return model, err
}

func ProcessWordData(db *sqlx.DB, wordlist model.Result) ([]model.Combined, error) {
	// init final processed model
	finalResult := []model.Combined{}

	// check if main word exists
	if wordlist.Google.MainWord == "" {
		// update the database to review manually
		UpdateNeedAttentionFlag(db, wordlist)
		return finalResult, errors.New("main word is empty")
	}

	// total parts of speeches for this word
	fmt.Printf("Word %s, partsOfSpeech %d \n", wordlist.Word, len(wordlist.Google.PartsOfSpeeches))

	if len(wordlist.Google.PartsOfSpeeches) == 0 {
		// update the database manually
		UpdateNeedAttentionFlag(db, wordlist)
		return finalResult, errors.New("no parts of speech found")
	}

	// 1. Process the google result
	for _, pos := range wordlist.Google.PartsOfSpeeches {

		singleResult := model.Combined{}

		//  get the parts of speech
		thisPos := utils.GetPos(pos.PartsOfSpeech)

		if thisPos == "" {
			// parts of speech does not match
			continue
		}

		singleResult.PartsOfSpeech = thisPos

		// add the definitions
		tmpDefinitions := []string{}
		tmpExamples := []string{}
		tmpSynonymsGre := []string{}    // synonyms exists in the gre db
		tmpSynonymsNormal := []string{} // synonyms do not exists in the gre db

		// now loop through definitions
		for _, definition := range pos.Definitions {
			// add the definition data
			tmpDefinitions = append(tmpDefinitions, definition.Definition)
			tmpExamples = append(tmpExamples, definition.Example)

			// now check for the synonyms
			for _, synonym := range definition.Synonyms {
				check := checkForSynonymWord(db, synonym)

				if check {
					// add to gre list
					if checkExists := checkUnique(tmpSynonymsGre, synonym); !checkExists {
						tmpSynonymsGre = append(tmpSynonymsGre, synonym)
					}
				} else {
					// add to normal list
					if checkExists := checkUnique(tmpSynonymsNormal, synonym); !checkExists {
						tmpSynonymsNormal = append(tmpSynonymsNormal, synonym)
					}
				}

			}

		}

		// now add the synonyms example definitions to the single res result
		singleResult.Definitions = tmpDefinitions
		singleResult.Examples = tmpExamples
		singleResult.SynonymsG = tmpSynonymsGre
		singleResult.SynonymsN = tmpSynonymsNormal

		// now append to the final result
		finalResult = append(finalResult, singleResult)
	}

	// 2. Process the wikipedia result
	for _, wikiPos := range wordlist.Wiki.PartsOfSpeeches {
		if len(wikiPos.Synonyms) > 0 {
			// get the parts of speech
			pos := utils.GetPos(wikiPos.PartsOfSpeech)
			// get the index of this parts of speech form the final result data
			posIndex := getPartsOfSpeechIndex(finalResult, pos)

			if posIndex != -1 {
				// parts of speech found in the final result
				for _, synonym := range wikiPos.Synonyms {
					check := checkForSynonymWord(db, synonym)

					if check {
						// add to gre list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					} else {
						// add to normal list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					}

				}
			}

		}
	}

	// 3.process words api result
	for _, wAPos := range wordlist.WordsApi.Results {

		if len(wAPos.Synonyms) > 0 {
			// get the parts of speech
			pos := utils.GetPos(wAPos.PartOfSpeech)
			// get the index of this parts of speech form the final result data
			posIndex := getPartsOfSpeechIndex(finalResult, pos)

			if posIndex != -1 {
				// parts of speech found in the final result
				for _, synonym := range wAPos.Synonyms {
					check := checkForSynonymWord(db, synonym)

					if check {
						// add to gre list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					} else {
						// add to normal list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					}

				}
			}

		}

	}

	// 4. process thesaurus data

	if len(wordlist.Thesaurus.Data.Synonyms) > 0 {
		for _, tPos := range wordlist.Thesaurus.Data.Synonyms {
			// get the parts of speech
			pos := utils.GetPos(tPos.PartsOfSpeech)
			// get the index of this parts of speech form the final result data
			posIndex := getPartsOfSpeechIndex(finalResult, pos)

			if posIndex != -1 {
				// parts of speech found in the final result
				for _, synonym := range tPos.Synonym {
					check := checkForSynonymWord(db, synonym)

					if check {
						// add to gre list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					} else {
						// add to normal list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					}

				}
			}

		}

	}

	return finalResult, nil

}

// checks if word exists in the wordlist table
func checkForSynonymWord(db *sqlx.DB, synonym string) bool {

	var result int
	err := db.Get(&result, "SELECT id FROM wordlist WHERE word=? LIMIT 1", synonym)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		log.Printf("error: %s\n", err)
		return false
	}

	return true
}

func checkUnique(data []string, word string) bool {
	for _, val := range data {
		if val == word {
			return true
		}
	}
	return false

}

func getPartsOfSpeechIndex(results []model.Combined, pos string) int {
	for i, res := range results {
		if res.PartsOfSpeech == pos {
			return i
		}

	}
	return -1
}

func SaveProcessedDataToWordTable(db *sqlx.DB, word model.WordModel, wordData []model.Combined) {

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	// insert the word into the words table
	data, err := json.Marshal(wordData)
	if err != nil {
		fmt.Println(err)
	}
	_, err = tx.Exec("Update words set word_data=? ,updated_at=now() where id=?", string(data), word.Id)

	if err != nil {
		fmt.Println(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set is_all_parsed=1 where word=?", word.Word)

	if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println(err)
	}

}

func UpdateNeedAttentionFlag(db *sqlx.DB, wordListData model.Result) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set needs_attention=1 where id=?", wordListData.ID)

	if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println(err)
	}

}
