package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/processor"
	"github.com/xatta-trone/words-combinator/utils"
)

func ProcessWordGroupData(wg model.WordGroupModel) {

	fmt.Println(wg.Id, wg.Name, wg.FileName, wg.Words)

	words := make([]string, 0)

	// check if it is a type of uploaded file or not

	if wg.FileName != nil {
		// get the words from the file
		wordsFromFile, err := ReadCsv(*wg.FileName)

		if err != nil {
			fmt.Printf("err %s \n", err)
			return
		}

		words = append(words, wordsFromFile...)
	}

	// check if we have some words or not
	if wg.Words != nil {
		wordsFromDB := ReadWords(*wg.Words)
		words = append(words, wordsFromDB...)

	}

	newWords := make([]string, 0)
	insertedWords := make([]InsertedWordStruct, 0)
	// now check the unique words
	for _, word := range words {

		wordId, err := InsertIntoWordsTable(word)

		if err != nil {
			fmt.Println(err)
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
		insertedWords = append(insertedWords, InsertedWordStruct{Word: word, Id: wordId})
		// append into uniqueWords to use it in the words group relation table
		newWords = append(newWords, word)

	}

	// now with the unique words make the word group relation table
	InsertGroupRelation(words, wg.Id)
	// insert the unique words into the word_group database
	InsertNewWordsToWordsGroupTable(newWords, wg.Id)
	// now get the new words meanings
	go ProcessNewWords(insertedWords, wg.Id)

	// fmt.Println(words, uniqueWords, insertedWords)

}

type WordStruct struct {
	Word string
}

type InsertedWordStruct struct {
	Word string `db:"word"`
	Id   int64  `db:"id"`
}

type WordGroupRelationMap struct {
	WordId      int64 `db:"word_id"`
	WordGroupId int64 `db:"word_group_id"`
}

func ProcessNewWords(newWords []InsertedWordStruct, groupId int64) {
	query_group := `update word_groups set status = ?, updated_at=now() where id=?`
	// update the word groups id status to processing
	_, err := database.Gdb.Exec(query_group, enums.WordGroupProcessing, groupId)
	if err != nil {
		fmt.Println(err)
	}

	for _, word := range newWords {
		processor.ReadTableAndProcessWord(word.Word)
	}

	// update the word groups id status to complete
	_, err = database.Gdb.Exec(query_group, enums.WordGroupComplete, groupId)
	if err != nil {
		fmt.Println(err)
	}

}

func InsertNewWordsToWordsGroupTable(newWords []string, groupId int64) {

	query_group := `update word_groups set new_words = ?, updated_at=now() where id=?`

	wordsToInsert := ""

	if len(newWords) > 0 {
		wordsToInsert = strings.Join(newWords, ",")
	}

	res, err := database.Gdb.Exec(query_group, wordsToInsert, groupId)
	if err != nil {
		fmt.Println(err)
	}

	rowsAffected, _ := res.RowsAffected()

	if int(rowsAffected) != 1 {
		fmt.Printf("total rows affected %d do not match total relation data %d \n", rowsAffected, 1)
	}

	utils.PrintG(fmt.Sprintf("Unique words inserts into word_groups table with id %d \n", groupId))

}

func InsertGroupRelation(words []string, groupId int64) {

	wordMap := []InsertedWordStruct{}

	fmt.Println("inside InsertGroupRelation")
	// fmt.Println(words)

	query, param, err := sqlx.In("SELECT id, word FROM words WHERE word in (?)", words)
	if err != nil {
		log.Fatal(err)
	}

	// this will pull words into the slice wordMap
	err = database.Gdb.Select(&wordMap, query, param...)

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(wordMap)

	// make the word groups
	wordGroupRelations := make([]WordGroupRelationMap, 0)

	for _, word := range wordMap {
		wordGroupRelations = append(wordGroupRelations, WordGroupRelationMap{WordId: word.Id, WordGroupId: groupId})

		// fmt.Println(word)

	}

	// fmt.Println(wordGroupRelations)

	query_group := `INSERT INTO word_group_relation(word_id,word_group_id,created_at) VALUES (:word_id,:word_group_id,now())`
	res, err := database.Gdb.NamedExec(query_group, wordGroupRelations)
	if err != nil {
		fmt.Println(err)
	}

	rowsAffected, _ := res.RowsAffected()

	if int(rowsAffected) != len(wordGroupRelations) {
		fmt.Printf("total rows affected %d do not match total relation data %d \n", rowsAffected, len(wordGroupRelations))
	}

	utils.PrintG(fmt.Sprintf("All words inserted for group id %d \n", groupId))

}

func InsertIntoWordsTable(word string) (int64, error) {
	const insertQuery = "INSERT INTO words(word,created_at) SELECT :word, now() WHERE NOT EXISTS (SELECT word FROM words WHERE word = :word);"
	w := WordStruct{Word: word}

	res, err := database.Gdb.NamedExec(insertQuery, w)

	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		// log.Fatal("last insert", err)
		return -1, err
	}

	return id, nil

}

func ReadCsv(filePath string) ([]string, error) {

	words := make([]string, 0)

	f, err := os.Open(filePath)

	if err != nil {
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
