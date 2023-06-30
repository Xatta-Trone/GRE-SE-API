package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/processor"
	"github.com/xatta-trone/words-combinator/scrapper"
	"github.com/xatta-trone/words-combinator/utils"
)

type ListProcessorService struct {
	db *sqlx.DB
}

func NewListProcessorService(db *sqlx.DB) *ListProcessorService {
	return &ListProcessorService{db: db}

}

type ListProcessorServiceInterface interface {
	ProcessListMetaRecord(listMeta model.ListMetaModel)
	InsertListWordRelation(wordId, listId int64) error
	ProcessWordsOfSingleGroup(words []string, listId int64)
	GetWordsFromListMetaRecord(words string) []string
}

func (listService *ListProcessorService) ProcessListMetaRecord(listMeta model.ListMetaModel) {
	// update the list meta table
	listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusParsing)
	listService.SendNotification(listMeta, enums.ListMetaStatusParsing, "")

	// now check the type of word to be processed...URL or word

	// steps
	// 1. get the words slice either from array or url parser
	// 2. for each word
	// 2.1 check if word exists in the words table
	// 2.2 if yes then map to list words relation table
	// 2.3 if no then insert into words and also into word list and make the words_list table relation
	// 2.4 run a function to get the word data form internet, then process the data then finally insert the processed data into words table

	// 1. get the words slice either from array or url parser

	if listMeta.Words != nil {
		var words []string
		// fire words processor
		fmt.Println(*listMeta.Words)
		processedWordStruct := listService.GetWordsFromListMetaRecord(*listMeta.Words)
		words = append(words, processedWordStruct...)

		var folderId uint64

		if listMeta.FolderId != nil {
			folderId = *listMeta.FolderId
		}

		_, unprocessedWords := listService.createListAndSaveWordsAndFolder(listMeta, listMeta.Name, listMeta.Visibility, listMeta.UserId, words, folderId)
		fmt.Println(unprocessedWords)

		if len(unprocessedWords) == 0 {
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusComplete)
			listService.SendNotification(listMeta, enums.ListMetaStatusComplete, "")
		} else {
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
			additionalText := fmt.Sprintf("we could not process these words %s", strings.Join(unprocessedWords, ","))
			listService.SendNotification(listMeta, enums.ListMetaStatusError, additionalText)
		}

	}

	if listMeta.Url != nil {
		// fire url processor
		fmt.Println(*listMeta.Url)
		listService.ProcessWordsFromUrl(listMeta)
		return

	}

	// // crate list record from list meta record
	// listId, err := listService.CreateListRecordFromListMeta(listMeta, listMeta.Name, 0)

	// if err != nil {
	// 	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
	// 	return
	// }

	// // now follow the steps
	// // listService.ProcessWordsOfSingleGroup(words, listId)

	// // now add to saved lists
	// listService.AddToSavedList(uint64(listId), listMeta.UserId)

	// UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusComplete)

	// fmt.Println(words)
	utils.PrintG("Processing complete")

}

func (listService *ListProcessorService) createListAndSaveWordsAndFolder(listMeta model.ListMetaModel, listName string, listVisibility int, userId uint64, words []string, folderId uint64) (bool, []string) {
	allOk := true
	unsuccessfulWords := []string{}
	// now create the list
	listId, err := listService.CreateList(userId, listMeta.Id, listName, enums.ListVisibilityMe, folderId)

	if err != nil {
		// todo: add notification err
		listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
		return false, unsuccessfulWords
	}

	// insert words into list
	_, errWords := listService.InsertWordsIntoList(uint64(listId), words)

	if len(errWords) > 0 {
		unsuccessfulWords = append(unsuccessfulWords, errWords...)
		// todo: insert a notification with missing words
		listService.ProcessUnSuccessfulWords(listMeta, errWords, uint64(listId))
	}

	// update list visibility
	listService.UpdateListVisibility(uint64(listId), listVisibility)

	listService.AddToSavedList(uint64(listId), userId)

	return allOk, unsuccessfulWords

}

// Step 1.1
func (listService *ListProcessorService) CreateList(userId uint64, listMetaId uint64, name string, visibility int, folderId uint64) (int64, error) {

	if visibility == 0 {
		// set as private
		visibility = enums.ListVisibilityMe
	}

	// generate slug
	slug := listService.GenerateUniqueListSlug(name)

	queryMap := map[string]interface{}{"name": name, "slug": slug, "list_meta_id": listMetaId, "visibility": visibility, "user_id": userId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	// create the lists record
	res, err := listService.db.NamedExec("Insert into lists(name,slug,list_meta_id,visibility,user_id,created_at,updated_at) values(:name,:slug,NullIf(:list_meta_id,0),:visibility,:user_id,:created_at,:updated_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	listId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	if listId == 0 {
		return -1, fmt.Errorf("there was a problem with the insertion. last id: %d", listId)
	}

	// now insert the folder

	var folderIdToInsert uint64

	if folderId != 0 {
		folderIdToInsert = folderId
	}

	if folderIdToInsert != 0 {
		// create the folder list relation
		queryMapForListFolderRelation := map[string]interface{}{"list_id": listId, "folder_id": folderIdToInsert, "user_id": userId, "created_at": time.Now().UTC()}
		_, err = listService.db.NamedExec("Insert ignore into folder_list_relation(folder_id,list_id) values(:folder_id,:list_id)", queryMapForListFolderRelation)
		if err != nil {
			utils.Errorf(err)
			utils.PrintR("there was an error creating list folder relation \n")

		}
		// insert into saved lists
		_, err = listService.db.NamedExec("Insert ignore into saved_lists(user_id,list_id,created_at) values(:user_id,:list_id,:created_at)", queryMapForListFolderRelation)
		if err != nil {
			utils.Errorf(err)
			utils.PrintR("there was an error creating list folder relation \n")

		}

	}

	return listId, nil

}

// Step 1.2
func (listService *ListProcessorService) InsertWordsIntoList(listId uint64, words []string) (bool, []string) {
	succeed := true
	unsuccessWords := []string{}

	// regex
	var IsLetter = regexp.MustCompile(`^[a-zA-Z\s-]+$`).MatchString
	// 2. for each word
	// 2.1 check if word exists in the words table
	// 2.2 if yes then map to list words relation table
	// 2.3 if no then insert into words and also into word list and make the words_list table relation
	// 2.4 run a function to get the word data form internet, then process the data then finally insert the processed data into words table

	for _, word := range words {
		if word == "" {
			// double check safety
			continue
		}

		// check if it contains letter space and -
		if !IsLetter(word) {
			unsuccessWords = append(unsuccessWords, word)
			continue
		}

		// 2.1 check if word exists in the table
		wordId, err := listService.CheckWordsTableForWord(word)

		if err != nil && err != sql.ErrNoRows {
			utils.Errorf(err)
			utils.PrintR(fmt.Sprintf("Could not find the word %s in words table", word))
		}

		// if word exists then insert the relation
		if wordId != 0 {
			_ = listService.InsertListWordRelation(int64(wordId), int64(listId))
			// go to the next word
			continue
		}

		// if the word not exists
		if err == sql.ErrNoRows {
			unsuccessWords = append(unsuccessWords, word)
			continue
			// first check in word-data table
			// wordListId := listService.CheckWordListTableForWord(word)

			// if wordListId == 0 {
			// 	// todo: dump to words to scrap
			// 	unsuccessWords = append(unsuccessWords, word)
			// 	continue
			// }

			// if wordListId != 0 {
			// 	// process the word from the wordlist table and save to words table
			// 	wordDataRaw, err := processor.CheckWordListTable(listService.db, word)

			// 	if err != nil {
			// 		utils.Errorf(err)
			// 		unsuccessWords = append(unsuccessWords, word)
			// 		continue
			// 	}

			// 	wordProcessedData, err := processor.ProcessWordData(listService.db, wordDataRaw)

			// 	if err != nil {
			// 		utils.Errorf(err)
			// 		// unsuccessWords = append(unsuccessWords, word)
			// 		continue
			// 	}

			// 	// save to words table
			// 	wordModel, err := listService.InsertIntoWordsTableWithData(word, wordProcessedData)

			// 	if err != nil {
			// 		utils.Errorf(err)
			// 		unsuccessWords = append(unsuccessWords, word)
			// 		continue
			// 	}

			// 	// insert the words relation
			// 	_ = listService.InsertListWordRelation(int64(wordModel.Id), int64(listId))
			// 	// go to the next word
			// 	continue
			// }

		} else {
			unsuccessWords = append(unsuccessWords, word)
			continue
		}

	}

	return succeed, unsuccessWords
}

func (listService *ListProcessorService) ProcessWordsFromUrl(listMeta model.ListMetaModel) {
	fmt.Println("Processing form url ", *listMeta.Url)
	// at first check if its a valid url
	// next determine which one is it
	if strings.Contains(*listMeta.Url, "vocabulary.com") {
		listService.ProcessVocabularyWords(listMeta)
	}
	if strings.Contains(*listMeta.Url, "quizlet.com") {
		listService.ProcessQuizletWords(listMeta)
	}
	if strings.Contains(*listMeta.Url, "memrise.com") {
		listService.ProcessMemriseWords(listMeta)
	}

}

func (listService *ListProcessorService) ProcessMemriseWords(listMeta model.ListMetaModel) {

	memriseSet, err := scrapper.GetMemriseSets(*listMeta.Url)

	fmt.Println(memriseSet.Urls)

	if err != nil {
		utils.Errorf(err)
		listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusURLError)
		listService.SendNotification(listMeta, enums.ListMetaStatusURLError, "")
		return
	}

	// create the folder
	// crate folder record from list meta record
	folderId, err := listService.CreateFolderFromListMeta(listMeta, memriseSet.Title)

	if err != nil {
		utils.Errorf(err)
		listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
		listService.SendNotification(listMeta, enums.ListMetaStatusError, "")
		return
	}

	// now create a saved folder record
	listService.AddToSavedFolders(uint64(folderId), listMeta.UserId)

	for i, url := range memriseSet.Urls {
		words, err := scrapper.ScrapMemrise(url)

		if err != nil {
			utils.Errorf(err)
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusURLError)
			listService.SendNotification(listMeta, enums.ListMetaStatusURLError, "")
			return
		}

		fmt.Println(words)

		if len(words) == 0 {
			utils.PrintR("ProcessMemriseWords No word found ")
			return
		}

		// list title
		title := fmt.Sprintf("%s-Group-%d", memriseSet.Title, i+1)
		title = utils.NormalizeString(title)

		// crate list record from list meta record
		_, unprocessedWords := listService.createListAndSaveWordsAndFolder(listMeta, title, listMeta.Visibility, listMeta.UserId, words, uint64(folderId))
		fmt.Println(unprocessedWords)

		if len(unprocessedWords) > 0 {
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)

			additionalText := fmt.Sprintf("we could not process these words %s", strings.Join(unprocessedWords, ","))
			listService.SendNotification(listMeta, enums.ListMetaStatusError, additionalText)
			return
		}

		// listId, err := listService.CreateListRecordFromListMeta(listMeta, title, uint64(folderId))

		// if err != nil {
		// 	utils.Errorf(err)
		// 	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
		// 	return
		// }

		// now follow the steps
		// listService.ProcessWordsOfSingleGroup(words, listId)
		// now add to saved lists
		// listService.AddToSavedList(uint64(listId), listMeta.UserId)

	}
	listService.UpdateFolderVisibility(uint64(folderId), listMeta.Visibility)
	listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusComplete)
	listService.SendNotification(listMeta, enums.ListMetaStatusComplete, "")
	utils.PrintG("Processing complete")

}

func (listService *ListProcessorService) ProcessQuizletWords(listMeta model.ListMetaModel) {

	// check if it is a folder
	if strings.Contains(*listMeta.Url, "folders") && strings.Contains(*listMeta.Url, "sets") {
		urls, setTitle, err := scrapper.GetQuizletUrlMaps(*listMeta.Url)

		fmt.Println(setTitle)
		fmt.Println(urls)

		if err != nil {
			utils.Errorf(err)
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusURLError)
			listService.SendNotification(listMeta, enums.ListMetaStatusURLError, "")
			return
		}

		// create the folder
		// crate folder record from list meta record
		folderId, err := listService.CreateFolderFromListMeta(listMeta, setTitle)

		if err != nil {
			utils.Errorf(err)
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
			listService.SendNotification(listMeta, enums.ListMetaStatusError, "")
			return
		}

		// now create a saved folder record
		listService.AddToSavedFolders(uint64(folderId), listMeta.UserId)

		for _, set := range urls {
			words, title, err := scrapper.ScrapQuizlet(set.Url)

			if err != nil {
				utils.Errorf(err)
				listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusURLError)
				listService.SendNotification(listMeta, enums.ListMetaStatusURLError, "")
				return
			}
			if title == "" {
				title = "Unknown title"
			}
			title = utils.NormalizeString(title)

			fmt.Println(words)
			fmt.Println(title)
			utils.Errorf(err)

			if len(words) == 0 {
				utils.PrintR("ProcessQuizletWords No word found ")
				return
			}

			// // crate list record from list meta record
			// listId, err := listService.CreateListRecordFromListMeta(listMeta, title, uint64(folderId))

			// if err != nil {
			// 	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
			// 	return
			// }

			// // now follow the steps
			// listService.ProcessWordsOfSingleGroup(words, listId)
			// // now add to saved lists
			// listService.AddToSavedList(uint64(listId), listMeta.UserId)

			// crate list record from list meta record
			_, unprocessedWords := listService.createListAndSaveWordsAndFolder(listMeta, title, listMeta.Visibility, listMeta.UserId, words, uint64(folderId))
			fmt.Println(unprocessedWords)

			if len(unprocessedWords) > 0 {
				listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
				additionalText := fmt.Sprintf("we could not process these words %s", strings.Join(unprocessedWords, ","))
				listService.SendNotification(listMeta, enums.ListMetaStatusError, additionalText)
				return

			}

		}
		listService.UpdateFolderVisibility(uint64(folderId), listMeta.Visibility)

	} else {

		words, title, err := scrapper.ScrapQuizlet(*listMeta.Url)

		if err != nil {
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusURLError)
			listService.SendNotification(listMeta, enums.ListMetaStatusURLError, "")
			return
		}

		title = utils.NormalizeString(title)

		fmt.Println(words)
		fmt.Println(title)
		utils.Errorf(err)

		if len(words) == 0 {
			utils.PrintR("ProcessQuizletWords No word found ")
			return
		}

		// crate list record from list meta record
		// listId, err := listService.CreateListRecordFromListMeta(listMeta, title, 0)

		// if err != nil {
		// 	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
		// 	return
		// }

		// // now follow the steps
		// listService.ProcessWordsOfSingleGroup(words, listId)

		// // now add to saved lists
		// listService.AddToSavedList(uint64(listId), listMeta.UserId)

		_, unprocessedWords := listService.createListAndSaveWordsAndFolder(listMeta, title, listMeta.Visibility, listMeta.UserId, words, 0)
		fmt.Println(unprocessedWords)

		if len(unprocessedWords) > 0 {
			listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
			additionalText := fmt.Sprintf("we could not process these words %s", strings.Join(unprocessedWords, ","))
			listService.SendNotification(listMeta, enums.ListMetaStatusError, additionalText)
			return
		}

	}

	listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusComplete)
	listService.SendNotification(listMeta, enums.ListMetaStatusComplete, "")
	utils.PrintG("Processing complete")

}

func (listService *ListProcessorService) ProcessVocabularyWords(listMeta model.ListMetaModel) {
	words, title, err := scrapper.ScrapVocabulary(*listMeta.Url)

	if err != nil {
		listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusURLError)
		listService.SendNotification(listMeta, enums.ListMetaStatusURLError, "")
		return
	}

	fmt.Println(words)
	fmt.Println(title)
	utils.Errorf(err)

	if len(words) == 0 {
		utils.PrintR("ProcessVocabularyWords No word found ")
		return
	}

	if title == "" {
		title = "Unknown title"
	}

	title = utils.NormalizeString(title)

	// crate list record from list meta record
	// listId, err := listService.CreateListRecordFromListMeta(listMeta, title, 0)
	// if err != nil {
	// 	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
	// 	return
	// }

	// now follow the steps
	// listService.ProcessWordsOfSingleGroup(words, listId)

	// add to saved list
	// listService.AddToSavedList(uint64(listId), listMeta.UserId)

	// crate list record from list meta record
	_, unprocessedWords := listService.createListAndSaveWordsAndFolder(listMeta, title, listMeta.Visibility, listMeta.UserId, words, 0)
	fmt.Println(unprocessedWords)

	if len(unprocessedWords) > 0 {
		listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusError)
		additionalText := fmt.Sprintf("we could not process these words %s", strings.Join(unprocessedWords, ","))
		listService.SendNotification(listMeta, enums.ListMetaStatusError, additionalText)
		utils.PrintR("Processing error ")
		return
	}

	// update the list meta status
	listService.UpdateListMetaRecordStatus(listMeta, enums.ListMetaStatusComplete)
	listService.SendNotification(listMeta, enums.ListMetaStatusComplete, "")
	utils.PrintG("Processing complete")

}

func (listService *ListProcessorService) ProcessWordsOfSingleGroup(words []string, listId int64) {

	// regex
	var IsLetter = regexp.MustCompile(`^[a-zA-Z\s-]+$`).MatchString

	// 2. for each word
	// 2.1 check if word exists in the words table
	// 2.2 if yes then map to list words relation table
	// 2.3 if no then insert into words and also into word list and make the words_list table relation
	// 2.4 run a function to get the word data form internet, then process the data then finally insert the processed data into words table

	for _, word := range words {
		if word == "" {
			continue
		}
		// check if it contains letter space and -
		if !IsLetter(word) {
			continue
		}

		// 2.1 check if word exists in the table
		checkId, err := listService.CheckWordsTableForWord(word)

		if err != nil && err != sql.ErrNoRows {
			fmt.Println("err in List processor service CheckWordsTableForWord func ==")
			utils.Errorf(err)
			continue
		}

		if checkId != 0 {
			// 2.2 if yes then map to list words relation table
			_ = listService.InsertListWordRelation(int64(checkId), listId)
			fmt.Printf("inserted word relation with word id %d list id %d\n", checkId, listId)
			continue
		}

		if err == sql.ErrNoRows {
			// 2.3 if no then insert into words and also into word list and make the words_list table relation
			// lets insert the word into the words table
			wordModel, err := listService.InsertIntoWordsTable(word)
			if err != nil {
				fmt.Println("err in List processor service InsertIntoWordsTable func ==")
				utils.Errorf(err)
				continue
			}

			// now insert into word list table and get the processed data
			processedWordData := processor.ProcessSingleWordData(listService.db, wordModel)

			if len(processedWordData) > 0 {
				// now update the words table with this data
				processor.SaveProcessedDataToWordTable(listService.db, wordModel.Word, wordModel.Id, processedWordData)
				// now insert the relation
				_ = listService.InsertListWordRelation(wordModel.Id, listId)
				fmt.Printf("inserted into word list and inserted word relation with word id %d list id %d\n", checkId, listId)

			} else {
				fmt.Printf("could not insert word relation with word id %d list id %d\n", checkId, listId)
			}

			fmt.Println(checkId)
		}

	}

}

func (listService *ListProcessorService) InsertIntoWordsTable(word string) (model.WordModel, error) {
	var model model.WordModel

	queryMap := map[string]interface{}{"word": word, "created_at": time.Now().UTC()}

	res, err := listService.db.NamedExec("Insert into words(word,created_at) values(:word,:created_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return model, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return model, err
	}

	if lastId == 0 {
		return model, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	query := "SELECT * FROM `words` WHERE `id` = ?;"

	result := listService.db.QueryRowx(query, lastId)

	err = result.StructScan(&model)

	if err != nil {
		utils.Errorf(err)
		return model, err
	}

	return model, nil
}

// :keep
func (listService *ListProcessorService) InsertIntoWordsTableWithData(word string, wordData []model.Combined) (model.WordModel, error) {
	var modelToExport model.WordModel

	// insert the word into the words table
	wodDataModel := model.WordDataModel{
		Word:            word,
		PartsOfSpeeches: wordData,
	}

	data, _ := json.Marshal(wodDataModel)

	queryMap := map[string]interface{}{"word": word, "wordJsonData": data, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := listService.db.NamedExec("Insert into words(word,word_data,created_at,updated_at) values(:word,:wordJsonData,:created_at,:updated_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return modelToExport, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return modelToExport, err
	}

	if lastId == 0 {
		return modelToExport, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	query := "SELECT * FROM `words` WHERE `id` = ?;"

	result := listService.db.QueryRowx(query, lastId)

	err = result.StructScan(&modelToExport)

	if err != nil {
		utils.Errorf(err)
		return modelToExport, err
	}

	return modelToExport, nil
}

// :keep
func (listService *ListProcessorService) InsertListWordRelation(wordId, listId int64) error {

	queryMap := map[string]interface{}{"word_id": wordId, "list_id": listId, "created_at": time.Now().UTC()}

	res, err := listService.db.NamedExec("Insert ignore into list_word_relation(word_id,list_id,created_at) values(:word_id,:list_id,:created_at)", queryMap)

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

// :keep
func (listService *ListProcessorService) CheckWordsTableForWord(word string) (uint64, error) {

	query := "SELECT `id` FROM `words` WHERE `word` = ?;"

	result := listService.db.QueryRowx(query, word)
	var ResultId uint64

	err := result.Scan(&ResultId)

	if err != nil {

		utils.Errorf(err)
		return ResultId, err
	}

	return ResultId, err

}

// :keep
func (listService *ListProcessorService) CheckWordListTableForWord(word string) uint64 {

	query := "SELECT `id` FROM `wordlist` WHERE `word` = ?;"

	result := listService.db.QueryRowx(query, word)
	var ResultId uint64

	err := result.Scan(&ResultId)

	if err != nil {
		utils.Errorf(err)
		return ResultId
	}

	return ResultId

}

// :keep
func (listService *ListProcessorService) UpdateListMetaRecordStatus(listMeta model.ListMetaModel, status int) {

	queryMap := map[string]interface{}{"id": listMeta.Id, "status": status, "updated_at": time.Now().UTC()}

	listService.db.NamedExec("Update list_meta set status=:status,updated_at=:updated_at where id=:id", queryMap)

	// for notifications
	// listService.SendNotification(listMeta, status, "")

}

// :keep
func (listService *ListProcessorService) SendNotification(listMeta model.ListMetaModel, status int, additionalText string) {

	notificationText := enums.GetListMetaStatusText(status)

	fmt.Println("notificationText")
	fmt.Println(notificationText)

	if notificationText != "" {

		finalNotificationText := ""

		if listMeta.Url == nil {
			finalNotificationText += fmt.Sprintf("::%d:: %s ::Status:: %s %s", listMeta.Id, listMeta.Name, notificationText, additionalText)
		}

		if status == enums.ListMetaStatusURLError {
			finalNotificationText += fmt.Sprintf("::%d:: %s ::Status:: %s %s", listMeta.Id, listMeta.Name, notificationText, additionalText)
		}

		data := model.NotificationModel{
			Id:        ulid.Make().String(),
			Content:   finalNotificationText,
			UserId:    listMeta.UserId,
			CreatedAt: time.Now().UTC(),
			URL:       "",
		}

		listService.db.NamedExec("INSERT INTO `notifications`(`id`, `content`, `user_id`, `url`, `created_at`) VALUES (:id,:content,:user_id,:url,:created_at)", data)

	}

}

// :keep
func (listService *ListProcessorService) UpdateListVisibility(listId uint64, visibility int) {
	queryMap := map[string]interface{}{"id": listId, "visibility": visibility, "updated_at": time.Now().UTC()}
	listService.db.NamedExec("Update lists set visibility=:visibility,updated_at=:updated_at where id=:id", queryMap)
}

// :keep
func (listService *ListProcessorService) UpdateFolderVisibility(folderId uint64, visibility int) {
	queryMap := map[string]interface{}{"id": folderId, "visibility": visibility, "updated_at": time.Now().UTC()}
	listService.db.NamedExec("Update folders set visibility=:visibility,updated_at=:updated_at where id=:id", queryMap)
}

func (listService *ListProcessorService) GetWordsFromListMetaRecord(words string) []string {
	var processedWords []string

	// trim white spaces, then split by new line
	tempData := strings.TrimSpace(words)
	splitBySpace := strings.Split(tempData, " ")

	for _, wordSplitBySpace := range splitBySpace {
		splitByNewLine := strings.Split(wordSplitBySpace, "\n")

		for _, value := range splitByNewLine {
			//    split by comma
			tempWords := strings.Split(value, ",")

			for _, word := range tempWords {

				// now run through the processor to remove extra characters
				words := utils.ProcessWord(word)
				processedWords = append(processedWords, words...)

			}
		}

	}

	return processedWords
}

// here the title is different because it could be parsed from the other sources like quizlet or so...
func (listService *ListProcessorService) CreateListRecordFromListMeta(listMeta model.ListMetaModel, title string, folderId uint64) (int64, error) {

	slug := listService.GenerateUniqueListSlug(title)

	var folderIdToInsert uint64
	var ListIdToInsert uint64

	if listMeta.FolderId != nil {
		folderIdToInsert = *listMeta.FolderId
	}

	if folderId != 0 {
		folderIdToInsert = folderId
	}

	ListIdToInsert = listMeta.Id

	queryMap := map[string]interface{}{"name": title, "slug": slug, "list_meta_id": ListIdToInsert, "visibility": listMeta.Visibility, "user_id": listMeta.UserId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	// create the lists record
	res, err := listService.db.NamedExec("Insert into lists(name,slug,list_meta_id,visibility,user_id,created_at,updated_at) values(:name,:slug,NullIf(:list_meta_id,0),:visibility,:user_id,:created_at,:updated_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	if lastId == 0 {
		return -1, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	if folderIdToInsert != 0 {
		// create the folder list relation
		queryMapForListFolderRelation := map[string]interface{}{"list_id": lastId, "folder_id": folderIdToInsert, "user_id": listMeta.UserId, "created_at": time.Now().UTC()}
		_, err = listService.db.NamedExec("Insert ignore into folder_list_relation(folder_id,list_id) values(:folder_id,:list_id)", queryMapForListFolderRelation)
		if err != nil {
			utils.Errorf(err)
			utils.PrintR("there was an error creating list folder relation \n")

		}
		// insert into saved lists
		_, err = listService.db.NamedExec("Insert ignore into saved_lists(user_id,list_id,created_at) values(:user_id,:list_id,:created_at)", queryMapForListFolderRelation)
		if err != nil {
			utils.Errorf(err)
			utils.PrintR("there was an error creating list folder relation \n")

		}

	}

	return lastId, nil

}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// CreateFolderFromListMeta makes a folder from the given list meta model.
// this is used for quizlet and memrise
func (listService *ListProcessorService) CreateFolderFromListMeta(listMeta model.ListMetaModel, title string) (int64, error) {

	slug := listService.GenerateUniqueFolderSlug(title)

	queryMap := map[string]interface{}{"name": title, "slug": slug, "list_meta_id": listMeta.Id, "visibility": enums.FolderVisibilityMe, "user_id": listMeta.UserId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := listService.db.NamedExec("Insert into folders(name,slug,list_meta_id,visibility,user_id,created_at,updated_at) values(:name,:slug,:list_meta_id,:visibility,:user_id,:created_at,:updated_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	if lastId == 0 {
		return -1, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	if lastId != 0 {
		// create the folder list relation
		queryMapForListFolderRelation := map[string]interface{}{"folder_id": lastId, "user_id": listMeta.UserId, "created_at": time.Now().UTC()}
		// insert into saved folders
		_, err = listService.db.NamedExec("Insert into saved_folders(user_id,folder_id,created_at) values(:user_id,:folder_id,:created_at)", queryMapForListFolderRelation)
		if err != nil {
			utils.Errorf(err)
			utils.PrintR("there was an error creating list folder relation \n")

		}

	}

	return lastId, nil

}

func (listService *ListProcessorService) GenerateUniqueListSlug(title string) string {

	slug := fmt.Sprintf("%s-%d", slug.Make(title), time.Now().UnixMilli())
	// now check the slug

	row := listService.db.QueryRow("SELECT Count(id) FROM lists WHERE slug like ?", fmt.Sprintf("%%%s-%%", slug))
	var totalCount int
	err := row.Scan(&totalCount)

	// fmt.Println(slug, fmt.Sprintf("%s-%%", slug), totalCount)

	if err != nil {
		// just add the timestamp and return
		return fmt.Sprintf("%s-%d", slug, time.Now().UnixMilli())
	}

	if totalCount > 0 {
		return fmt.Sprintf("%s-%d", slug, totalCount+1)

	}

	return slug
}

func (listService *ListProcessorService) GenerateUniqueFolderSlug(title string) string {

	slug := fmt.Sprintf("%s-%d", slug.Make(title), time.Now().UnixMilli())
	// now check the slug

	row := listService.db.QueryRow("SELECT Count(id) FROM folders WHERE slug like ?", fmt.Sprintf("%%%s-%%", slug))
	var totalCount int
	err := row.Scan(&totalCount)

	// fmt.Println(slug, fmt.Sprintf("%s-%%", slug), totalCount)

	if err != nil {
		// just add the timestamp and return
		return fmt.Sprintf("%s-%d", slug, time.Now().UnixMilli())
	}

	if totalCount > 0 {
		return fmt.Sprintf("%s-%d", slug, totalCount+1)

	}

	return slug
}

func (listService *ListProcessorService) AddToSavedList(listId, userId uint64) {

	queryMap := map[string]interface{}{"user_id": userId, "list_id": listId, "created_at": time.Now().UTC()}

	_, err := listService.db.NamedExec("Insert ignore into saved_lists(user_id,list_id,created_at) values(:user_id,:list_id,:created_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
	}
}

func (listService *ListProcessorService) AddToSavedFolders(folderId, userId uint64) {

	queryMap := map[string]interface{}{"user_id": userId, "folder_id": folderId, "created_at": time.Now().UTC()}

	_, err := listService.db.NamedExec("Insert ignore into saved_folders(user_id,folder_id,created_at) values(:user_id,:folder_id,:created_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
	}
}

func (listService *ListProcessorService) ProcessUnSuccessfulWords(listMeta model.ListMetaModel, words []string, listId uint64) {

	wordList := []model.PendingWordModel{}

	for _, word := range words {
		w := model.PendingWordModel{
			Word:   word,
			ListId: listId,
		}

		wordList = append(wordList, w)
	}

	_, err := listService.db.NamedExec(`INSERT INTO pending_words (word,list_id) VALUES (:word,:list_id)`, wordList)

	if err != nil {
		utils.Errorf(err)
		return
	}

	// send notification

	url := os.Getenv("DISCORD_WEBHOOK_URL")

	if url == "" {
		url = "https://discord.com/api/webhooks/1123787115050844181/roFKRIy_iZ6SWhfNHEtue4rbVixP1X_PKBRcJPl5N73DCkwnyCgSSMeBO733ZcQG2hgr"
	}

	client, err := webhook.NewWithURL(url)

	if err != nil {
		utils.Errorf(err)
		return
	}

	msg := fmt.Sprintf("%d words needs to be checked, for list id %d \n :: words :: \n %s", len(words), listId, strings.Join(words, ", "))

	_, err = client.CreateMessage(discord.WebhookMessageCreate{
		Content: msg,
	})

	if err != nil {
		utils.Errorf(err)
		return
	}

}
