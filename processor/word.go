package processor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func ReadTableAndProcessWord(word string) {

	fmt.Println("getting word result for ", word)

	rs := database.Gdb.QueryRowx("SELECT `id`, `word`, `google`, `wiki`, `words_api`,`thesaurus` FROM `wordlist` WHERE `word` = ?;", word)

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
	fmt.Println(r.Google.MainWord == "")

	if r.Google.MainWord == "" {
		// this data needs to reviewed manually

		UpdateDBNeedsAttention(r.Word)
		return
	}

	// total parts of speeches
	fmt.Println("total parts of speeches", len(r.Google.PartsOfSpeeches))

	finalResult := []model.Combined{}

	if len(r.Google.PartsOfSpeeches) > 0 {

		for _, pos := range r.Google.PartsOfSpeeches {
			var SingleRes model.Combined

			poss := utils.GetPos(pos.PartsOfSpeech)

			fmt.Println("POS ", poss)
			if poss == "" {
				fmt.Println("No pos found")

			}

			SingleRes.PartsOfSpeech = poss

			// add the definitions
			tmpDef := []string{}
			tmpEx := []string{}
			tmpSynG := []string{} // synonyms exists in the gre db
			tmpSynN := []string{} // synonyms do not exists in the gre db

			for _, d := range pos.Definitions {
				tmpDef = append(tmpDef, d.Definition)
				tmpEx = append(tmpEx, d.Example)
				// fmt.Println(d.Synonyms)

				// check if synonyms are found in the db

				for _, s := range d.Synonyms {
					v := checkSynExists(s)
					// fmt.Printf("synonym %s status %v \n", s, v)

					if v == 1 {
						// synonyms exists in the db
						// check if its already added
						// if not added then add
						if ok := checkExists(tmpSynG, s); !ok {
							tmpSynG = append(tmpSynG, s)
						}

					} else {

						if ok := checkExists(tmpSynN, s); !ok {
							tmpSynN = append(tmpSynN, s)
						}
					}
				}

			}

			SingleRes.Definitions = tmpDef
			SingleRes.Examples = tmpEx
			SingleRes.SynonymsG = tmpSynG
			SingleRes.SynonymsN = tmpSynN

			finalResult = append(finalResult, SingleRes)
		}

		// data, _ := json.MarshalIndent(finalResult, "", "\t")

		// fmt.Println(string(data))
	}

	// process wiki result
	// var wikis model.Wiki

	// json.Unmarshal(r.Wiki,wikis)
	// fmt.Println("wiki", r.Wiki)

	for _, wikiPos := range r.Wiki.PartsOfSpeeches {
		// fmt.Println(wikiPos)
		// fmt.Println()

		if len(wikiPos.Synonyms) > 0 {
			poss := utils.GetPos(wikiPos.PartsOfSpeech)
			i := getPosIndex(finalResult, poss)

			if i != -1 {
				// now iterate over the synonyms
				for _, s := range wikiPos.Synonyms {
					v := checkSynExists(s)
					// fmt.Printf("synonym %s status %v \n", s, v)

					if v == 1 {
						// synonyms exists in the db
						// check if its already added
						// if not added then add
						if ok := checkExists(finalResult[i].SynonymsG, s); !ok {
							finalResult[i].SynonymsG = append(finalResult[i].SynonymsG, s)
						}

					} else {

						if ok := checkExists(finalResult[i].SynonymsN, s); !ok {
							finalResult[i].SynonymsN = append(finalResult[i].SynonymsN, s)
						}
					}

				}
			}

			// fmt.Println("index of pos in wiki", i, wikiPos.PartsOfSpeech)

		}

	}

	// data, _ := json.MarshalIndent(finalResult, "", "\t")

	// fmt.Println(string(data))

	// process words api
	// fmt.Println(r.WordsApi)

	for _, pos := range r.WordsApi.Results {
		// get the pos
		poss := utils.GetPos(pos.PartOfSpeech)
		i := getPosIndex(finalResult, poss)

		if i != -1 {
			// now iterate over the synonyms
			for _, s := range pos.Synonyms {
				v := checkSynExists(s)
				// fmt.Printf("synonym %s status %v \n", s, v)

				if v == 1 {
					// synonyms exists in the db
					// check if its already added
					// if not added then add
					if ok := checkExists(finalResult[i].SynonymsG, s); !ok {
						finalResult[i].SynonymsG = append(finalResult[i].SynonymsG, s)
					}

				} else {

					if ok := checkExists(finalResult[i].SynonymsN, s); !ok {
						finalResult[i].SynonymsN = append(finalResult[i].SynonymsN, s)
					}
				}
			}
		}

	}

	// data, _ := json.MarshalIndent(finalResult, "", "\t")

	// fmt.Println(string(data))

	// add thesaurus

	for _, pos := range r.Thesaurus.Data.Synonyms {
		// get the pos
		poss := utils.GetPos(pos.PartsOfSpeech)
		i := getPosIndex(finalResult, poss)

		if i != -1 {
			// now iterate over the synonyms
			for _, s := range pos.Synonym {
				v := checkSynExists(s)
				// fmt.Printf("synonym %s status %v \n", s, v)

				if v == 1 {
					// synonyms exists in the db
					// check if its already added
					// if not added then add
					if ok := checkExists(finalResult[i].SynonymsG, s); !ok {
						finalResult[i].SynonymsG = append(finalResult[i].SynonymsG, s)
					}

				} else {

					if ok := checkExists(finalResult[i].SynonymsN, s); !ok {
						finalResult[i].SynonymsN = append(finalResult[i].SynonymsN, s)
					}
				}
			}
		}

	}

	// data, _ := json.MarshalIndent(finalResult, "", "\t")

	// fmt.Println(string(data))

	finalData := model.CombinedWithWord{
		Word:            word,
		PartsOfSpeeches: finalResult,
	}
	finalData.Word = word
	finalData.PartsOfSpeeches = finalResult

	// fmt.Println(finalData)

	SaveToDB(finalData)

}

func ReturnUnique(words []string) []string {
	keyDict := make(map[string]bool)
	uniqueWords := []string{}

	for _, val := range words {

		_, ok := keyDict[val]
		if !ok {
			keyDict[val] = true
			uniqueWords = append(uniqueWords, val)
		}

	}
	return uniqueWords
}

func checkExists(data []string, word string) bool {
	for _, val := range data {
		if val == word {
			return true
		}
	}
	return false
}

func checkSynExists(synonym string) int {

	var result int
	err := database.Gdb.Get(&result, "SELECT id FROM wordlist WHERE word=? LIMIT 1", synonym)

	if err == sql.ErrNoRows {
		return 0
	}

	if err == nil {
		return 1
	}

	log.Printf("error: %s\n", err)

	return 0

}

func getPosIndex(results []model.Combined, pos string) int {

	for i, res := range results {
		if res.PartsOfSpeech == pos {
			return i
		}

	}

	return -1

}

func SaveToDB(wordData model.CombinedWithWord) {
	tx, err := database.Gdb.Begin()
	if err != nil {
		fmt.Println(err)
	}
	// insert the word into the words table
	data, err := json.Marshal(wordData)
	if err != nil {
		fmt.Println(err)
	}
	_, err = tx.Exec("Update words set word_data=? where word=?", string(data), wordData.Word)

	if err != nil {
		fmt.Println(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set is_all_parsed=1 where word=?", wordData.Word)

	if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println(err)
	}

}

func UpdateDBNeedsAttention(word string) {
	tx, err := database.Gdb.Begin()
	if err != nil {
		fmt.Println(err)
	}
	// update the wordlist table
	_, err = tx.Exec("Update wordlist set needs_attention=1 where word=?", word)

	if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println(err)
	}

}
