package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	errs "github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/processor"
	"github.com/xatta-trone/words-combinator/utils"
)

type WordGroupService struct {
	db *sqlx.DB
}

type WordGroupServiceInterface interface {
	ProcessWordGroupData(wg model.WordGroupModel)
}

type WordStruct struct {
	Word string
}

func NewWordGroupService(db *sqlx.DB) *WordGroupService {
	return &WordGroupService{db: db}
}

// ProcessWordGroupData function received a word group record
// 1. then it reads the csv file
// 2. then it reads the new words column
// 3. combines all the words together
// 4. then it inserts all words one by one into words table
// 5. if the words table returns a id then its a new word => to be processed later
// 5.a if it doesn't give a id back (0) then its a duplicate word and we continue
// 6. then the new words returned from step 5 is inserted into a new slice called insertedWords(struct) and newWords(only word)
// 7. then it fires InsertGroupRelation with the word group id and the []words
// 8. then it updated the word_groups table with the new found unique words form this operation => newWords (InsertNewWordsToWordsGroupTable)
// 9. then it fires a go routine where the new words are sent to be processed ProcessNewWords(insertedWords,word_group_id)
func (wgService *WordGroupService) ProcessWordGroupData(wg model.WordGroupModel) {

	fmt.Println(wg.Id, wg.Name, wg.FileName, wg.Words)

	words := make([]string, 0)

	// check if it is a type of uploaded file or not

	if wg.FileName != nil {
		// get the words from the file
		wordsFromFile, err := ReadCsv(*wg.FileName)

		if err != nil {
			fmt.Println(errs.New(err).ErrorStack())
			fmt.Println(err.(*errs.Error).ErrorStack())
			return
		}

		words = append(words, wordsFromFile...)
	}

	// check if we have some words other than getting from csv or not
	if wg.Words != nil {
		wordsFromDB := ReadWords(*wg.Words)
		words = append(words, wordsFromDB...)

	}

	newWords := make([]string, 0)
	insertedWords := make([]model.WordModel, 0)
	// now check the unique words
	for _, word := range words {

		if word == "" {
			continue
		}

		wordId, err := wgService.InsertIntoWordsTable(word)

		if err != nil {
			utils.Errorf(err)
			return
		}

		if wordId == -1 {
			continue
		}

		if wordId == 0 {
			// duplicate found
			// append into uniqueWords to use it in the words group relation table
			// uniqueWords = append(uniqueWords, word)
			continue
		}
		// appending for further processing from the wordlist table
		insertedWords = append(insertedWords, model.WordModel{Word: word, Id: wordId})
		// append into uniqueWords to use it in the words group relation table
		newWords = append(newWords, word)

	}

	// now with the unique words make the word group relation table
	wgService.InsertGroupRelation(words, wg.Id)
	// insert the unique words into the word_group database
	wgService.InsertNewWordsToWordsGroupTable(newWords, wg.Id)
	// now get the new words meanings
	go wgService.ProcessNewWords(insertedWords, wg.Id)

	// fmt.Println(words, uniqueWords, insertedWords)

}

func (wgService *WordGroupService) ProcessNewWords(newWords []model.WordModel, groupId int64) {
	query_group := `update word_groups set status = ?, updated_at=now() where id=?`
	// update the word groups id status to processing
	_, err := wgService.db.Exec(query_group, enums.WordGroupProcessing, groupId)
	if err != nil {
		utils.Errorf(err)
		fmt.Println(err.(*errs.Error).ErrorStack())
	}

	for _, word := range newWords {
		// processor.ReadTableAndProcessWord(word.Word)
		processedWordData := processor.ProcessSingleWordData(wgService.db, word)

		// save to the database
		processor.SaveProcessedDataToWordTable(wgService.db, word.Word,word.Id, processedWordData)
	}

	// update the word groups id status to complete
	_, err = database.Gdb.Exec(query_group, enums.WordGroupComplete, groupId)
	if err != nil {
		utils.Errorf(err)
		fmt.Println(err.(*errs.Error).ErrorStack())
	}

}

func (wgService *WordGroupService) InsertNewWordsToWordsGroupTable(newWords []string, groupId int64) {

	query_group := `update word_groups set new_words = ?, updated_at=now() where id=?`

	wordsToInsert := ""

	if len(newWords) > 0 {
		wordsToInsert = strings.Join(newWords, ",")
	}

	res, err := wgService.db.Exec(query_group, wordsToInsert, groupId)
	if err != nil {
		utils.Errorf(err)
		fmt.Println(err.(*errs.Error).ErrorStack())
	}

	rowsAffected, _ := res.RowsAffected()

	if int(rowsAffected) != 1 {
		fmt.Printf("total rows affected %d do not match total relation data %d \n", rowsAffected, 1)
	}

	utils.PrintG(fmt.Sprintf("Unique words inserts into word_groups table with id %d \n", groupId))

}

func (wgService *WordGroupService) InsertGroupRelation(words []string, groupId int64) {

	wordMap := []model.WordModel{}

	fmt.Println("inside InsertGroupRelation")
	// fmt.Println(words)

	query, param, err := sqlx.In("SELECT id, word FROM words WHERE word in (?)", words)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.(*errs.Error).ErrorStack())
	}

	// this will pull words into the slice wordMap
	err = wgService.db.Select(&wordMap, query, param...)

	if err != nil {
		utils.Errorf(err)
		fmt.Println(err.(*errs.Error).ErrorStack())
	}

	// fmt.Println(wordMap)

	// make the word groups
	wordGroupRelations := make([]model.WordGroupRelationModel, 0)

	for _, word := range wordMap {
		wordGroupRelations = append(wordGroupRelations, model.WordGroupRelationModel{WordId: word.Id, WordGroupId: groupId})

		// fmt.Println(word)

	}

	// fmt.Println(wordGroupRelations)

	query_group := `INSERT INTO word_group_relation(word_id,word_group_id,created_at) VALUES (:word_id,:word_group_id,now())`
	res, err := wgService.db.NamedExec(query_group, wordGroupRelations)
	if err != nil {
		utils.Errorf(err)
		fmt.Println(err.(*errs.Error).ErrorStack())
	}

	rowsAffected, _ := res.RowsAffected()

	if int(rowsAffected) != len(wordGroupRelations) {
		fmt.Printf("total rows affected %d do not match total relation data %d \n", rowsAffected, len(wordGroupRelations))
	}

	utils.PrintG(fmt.Sprintf("All words inserted for group id %d \n", groupId))

}

func (wgService *WordGroupService) InsertIntoWordsTable(word string) (int64, error) {
	const insertQuery = "INSERT INTO words(word,created_at) SELECT :word, now() WHERE NOT EXISTS (SELECT word FROM words WHERE word = :word);"
	w := WordStruct{Word: word}

	res, err := wgService.db.NamedExec(insertQuery, w)

	if err != nil {
		fmt.Println(err.(*errs.Error).ErrorStack())
		return -1, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err.(*errs.Error).ErrorStack())
		return -1, err
	}

	return id, nil

}

func ReadCsv(filePath string) ([]string, error) {

	words := make([]string, 0)

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err.(*errs.Error).ErrorStack())
		return nil, err
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	for {
		data, err := csvReader.Read() //read line by line

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err.(*errs.Error).ErrorStack())
			return nil, err
		}

		// sanitize the word
		processedWord := utils.ProcessWord(data[0])

		words = append(words, processedWord...)

	}

	return words, nil
}

func ReadWords(wordString string) []string {

	words := make([]string, 0)

	tempWords := strings.Split(strings.ReplaceAll(wordString, "\r\n", "\n"), "\n")

	for _, str := range tempWords {
		wordGroups := strings.Split(str, ",")
		words = append(words, wordGroups...)
	}

	return words

}
