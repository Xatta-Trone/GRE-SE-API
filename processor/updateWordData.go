package processor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/scrapper"
	"github.com/xatta-trone/words-combinator/utils"
)

func ProcessSingleWordData(db *sqlx.DB, word model.WordModel) []model.Combined {

	var result []model.Combined

	// check if the there is data in the wordlist table
	data, err := CheckWordListTable(db, word.Word)

	fmt.Println("===inside ProcessSingleWordData ==")
	// fmt.Println(data, err)

	if err == sql.ErrNoRows {
		// if no then insert into wordlist and get the wordlist data
		wordListModel, err := InsertIntoWordListTable(db, word.Word)
		if err != nil {
			utils.Errorf(err)
			utils.Errorf(err)
			return result
		}
		// then process the data
		scrapper.ScrapWordById(db, wordListModel)
		// then process this result again
		return ProcessSingleWordData(db, word)

	}

	if err != nil {
		utils.Errorf(err)
		utils.Errorf(err)
		return result
	}

	if data.ID != 0 {
		// if yes process the data
		processedData, _ := ProcessWordData(db, data)

		result = append(result, processedData...)

	}

	return result

}

func InsertIntoWordListTable(db *sqlx.DB, word string) (model.Result, error) {
	var model model.Result

	_, err := db.Exec("INSERT IGNORE INTO wordlist(word,created_at,updated_at) values (?,now(),now())", word)

	if err != nil {
		utils.Errorf(err)
		utils.Errorf(err)
		return model, err
	}

	query := "SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus`,`mw` FROM `wordlist` WHERE `word` = ?;"

	result := db.QueryRowx(query, word)

	err = result.StructScan(&model)

	if err != nil {
		utils.Errorf(err)
		utils.Errorf(err)
		return model, err
	}

	return model, nil

}

func GetFromWordListTable(db *sqlx.DB, word string) (model.Result, error) {
	var model model.Result
	query := "SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus`,`mw` FROM `wordlist` WHERE `word` = ?;"

	result := db.QueryRowx(query, word)

	err := result.StructScan(&model)

	if err != nil {
		utils.Errorf(err)
		utils.Errorf(err)
		return model, err
	}

	return model, nil

}

func CheckWordListTable(db *sqlx.DB, word string) (model.Result, error) {

	query := "SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus`,`mw` FROM `wordlist` WHERE `word` = ?;"

	result := db.QueryRowx(query, word)

	var model model.Result

	err := result.StructScan(&model)

	if err != nil {
		// utils.Errorf(err)
		utils.Errorf(err)
		return model, err
	}

	return model, err
}

func GetUnProcessedWordDataById(db *sqlx.DB, wordId uint64) (model.Result, error) {

	query := "SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus`,`mw` FROM `wordlist` WHERE `id` = ?;"

	result := db.QueryRowx(query, wordId)

	var model model.Result

	err := result.StructScan(&model)

	if err != nil {
		// utils.Errorf(err)
		utils.Errorf(err)
		return model, err
	}

	return model, err
}

func ProcessWordData(db *sqlx.DB, wordlist model.Result) ([]model.Combined, error) {
	fmt.Println("==inside ProcessWordData==", wordlist.ID)
	// init final processed model
	finalResult := []model.Combined{}

	fmt.Println(wordlist.Google.MainWord)

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
	// fmt.Println("wiki")
	// fmt.Println(wordlist.Wiki)

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
					// fmt.Printf("checking for %s in wiki : %t\n", synonym, check)

					if check {
						// add to gre list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					} else {
						// add to normal list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
							finalResult[posIndex].SynonymsN = append(finalResult[posIndex].SynonymsN, synonym)
						}
					}

				}
			}

		}
	}

	// fmt.Println("words api")
	// fmt.Println(wordlist.WordsApi)
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
					// fmt.Printf("checking for %s in words api : %t\n", synonym, check)

					if check {
						// add to gre list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					} else {
						// add to normal list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
							finalResult[posIndex].SynonymsN = append(finalResult[posIndex].SynonymsN, synonym)
						}
					}

				}
			}

		}

	}

	// 4. process thesaurus data
	// fmt.Println("thesaurus api")
	// fmt.Println(wordlist.Thesaurus)
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
					// fmt.Printf("checking for %s in thesaurus : %t\n", synonym, check)

					if check {
						// add to gre list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
							finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
						}
					} else {
						// add to normal list
						if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
							finalResult[posIndex].SynonymsN = append(finalResult[posIndex].SynonymsN, synonym)
						}
					}

				}
			}

		}

	}

	// 5.0 process the MW data
	// fmt.Println("MW api")
	// fmt.Println(wordlist.Mw)
	if len(wordlist.Mw.Data.PartsOfSpeeches) > 0 {
		for _, tPos := range wordlist.Mw.Data.PartsOfSpeeches {
			// get the parts of speech
			pos := utils.GetPos(tPos.PartsOfSpeech)
			// get the index of this parts of speech form the final result data
			posIndex := getPartsOfSpeechIndex(finalResult, pos)

			if posIndex != -1 {
				// loop through each data
				for _, d := range tPos.Data {
					// parts of speech found in the final result
					for _, synonym := range d.Synonyms {
						check := checkForSynonymWord(db, synonym)
						// fmt.Printf("check for word %s %t \n", synonym, check)

						if check {
							// add to gre list
							if checkExists := checkUnique(finalResult[posIndex].SynonymsG, synonym); !checkExists {
								finalResult[posIndex].SynonymsG = append(finalResult[posIndex].SynonymsG, synonym)
							}
						} else {
							// add to normal list
							if checkExists := checkUnique(finalResult[posIndex].SynonymsN, synonym); !checkExists {
								finalResult[posIndex].SynonymsN = append(finalResult[posIndex].SynonymsN, synonym)
							}
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
		// utils.Errorf(err)
		return false
	}

	if err != nil {
		utils.Errorf(err)
		utils.Errorf(err)
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

func SaveProcessedDataToWordTable(db *sqlx.DB, word string, wordId int64, wordData []model.Combined) {

	tx, err := db.Begin()
	if err != nil {
		utils.Errorf(err)
	}
	// insert the word into the words table
	wodDataModel := model.WordDataModel{
		Word:            word,
		PartsOfSpeeches: wordData,
	}
	data, err := json.Marshal(wodDataModel)
	if err != nil {
		utils.Errorf(err)
	}
	_, err = tx.Exec("Update words set word_data=? ,updated_at=now() where id=?", string(data), wordId)

	if err != nil {
		utils.Errorf(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set is_all_parsed=1, in_words=1,tried=1 where word=?", word)

	if err != nil {
		utils.Errorf(err)
	}

	err = tx.Commit()

	if err != nil {
		utils.Errorf(err)
	}

}

func SaveNewWordData(db *sqlx.DB, word string, wordId int64, wordData []model.Combined) error {

	tx, err := db.Begin()
	if err != nil {
		utils.Errorf(err)
		return err
	}
	// insert the word into the words table
	wodDataModel := model.WordDataModel{
		Word:            word,
		PartsOfSpeeches: wordData,
	}
	data, err := json.Marshal(wodDataModel)
	if err != nil {
		utils.Errorf(err)
		return err
	}
	_, err = tx.Exec("INSERT INTO `words`(`word`, `word_data`, `created_at`, `updated_at`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE word_data=?", word, string(data), time.Now().UTC(), time.Now().UTC(),string(data),)

	if err != nil {
		utils.Errorf(err)
		return err
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set is_all_parsed=1, in_words=1,tried=1 where word=?", word)

	if err != nil {
		utils.Errorf(err)
		return err
	}

	err = tx.Commit()

	if err != nil {
		utils.Errorf(err)
		return err
	}

	return nil

}

func UpdateNeedAttentionFlag(db *sqlx.DB, wordListData model.Result) {
	tx, err := db.Begin()
	if err != nil {
		utils.Errorf(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set needs_attention=1 where id=?", wordListData.ID)

	if err != nil {
		utils.Errorf(err)
	}

	err = tx.Commit()

	if err != nil {
		utils.Errorf(err)
	}

}

func SetTried(db *sqlx.DB, wordId uint64) {
	tx, err := db.Begin()
	if err != nil {
		utils.Errorf(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set tried=1 where id=?", wordId)

	if err != nil {
		utils.Errorf(err)
	}

	err = tx.Commit()

	if err != nil {
		utils.Errorf(err)
	}

}
