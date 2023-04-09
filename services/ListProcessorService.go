package services

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
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
}

func (listService *ListProcessorService) ProcessListMetaRecord(listMeta model.ListMetaModel) {
	// update the list meta table
	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusParsing)

	// now check the type of word to be processed...URL or word

	// steps
	// 1. get the words slice either from array or url parser
	// 2. for each word
	// 2.1 check if word exists in the words table
	// 2.2 if yes then map to list words relation table
	// 2.3 if no then insert into words and also into word list and make the words_list table relation
	// 2.4 run a function to get the word data form internet, then process the data then finally insert the processed data into words table

	// 1. get the words slice either from array or url parser
	var words []string

	if listMeta.Words != nil {
		// fire words processor
		fmt.Println(*listMeta.Words)
		processedWordStruct := GetWordsFromListMetaRecord(*listMeta.Words)
		words = append(words, processedWordStruct...)
	}

	if listMeta.Url != nil {
		// fire url processor
		fmt.Println(*listMeta.Url)
		listService.ProcessWordsFromUrl(listMeta)
		return

	}

	// crate list record from list meta record
	listId, err := listService.CreateListRecordFromListMeta(listMeta, listMeta.Name)

	if err != nil {
		UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
		return
	}

	// now follow the steps
	listService.ProcessWordsOfSingleGroup(words, listId)

	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusComplete)

	fmt.Println(words)

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

	urls, err := scrapper.GetMemriseSets(*listMeta.Url)

	fmt.Println(urls.Urls)

	if err != nil {
		UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusURLError)
		return
	}

	for i, url := range urls.Urls {
		words, err := scrapper.ScrapMemrise(url)

		if err != nil {
			UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusURLError)
			return
		}

		fmt.Println(words)
		fmt.Println(err)

		if len(words) == 0 {
			utils.PrintR("ProcessMemriseWords No word found ")
			return
		}

		// list title
		title := fmt.Sprintf("%s-Group-%d", urls.Title, i+1)

		// crate list record from list meta record
		listId, err := listService.CreateListRecordFromListMeta(listMeta, title)

		if err != nil {
			UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
			return
		}

		// now follow the steps
		listService.ProcessWordsOfSingleGroup(words, listId)

	}

	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusComplete)

}

func (listService *ListProcessorService) ProcessQuizletWords(listMeta model.ListMetaModel) {

	// check if it is a folder
	if strings.Contains(*listMeta.Url, "folders") && strings.Contains(*listMeta.Url, "sets") {
		urls, err := scrapper.GetQuizletUrlMaps(*listMeta.Url)

		if err != nil {
			UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusURLError)
			return
		}

		for _, set := range urls {
			words, title, err := scrapper.ScrapQuizlet(set.Url)

			if err != nil {
				UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusURLError)
				return
			}

			fmt.Println(words)
			fmt.Println(title)
			fmt.Println(err)

			if len(words) == 0 {
				utils.PrintR("ProcessQuizletWords No word found ")
				return
			}

			// crate list record from list meta record
			listId, err := listService.CreateListRecordFromListMeta(listMeta, title)

			if err != nil {
				UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
				return
			}

			// now follow the steps
			listService.ProcessWordsOfSingleGroup(words, listId)

		}

	} else {

		words, title, err := scrapper.ScrapQuizlet(*listMeta.Url)

		if err != nil {
			UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusURLError)
			return
		}

		fmt.Println(words)
		fmt.Println(title)
		fmt.Println(err)

		if len(words) == 0 {
			utils.PrintR("ProcessQuizletWords No word found ")
			return
		}

		// crate list record from list meta record
		listId, err := listService.CreateListRecordFromListMeta(listMeta, title)

		if err != nil {
			UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
			return
		}

		// now follow the steps
		listService.ProcessWordsOfSingleGroup(words, listId)

	}

	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusComplete)

}

func (listService *ListProcessorService) ProcessVocabularyWords(listMeta model.ListMetaModel) {
	words, title, err := scrapper.ScrapVocabulary(*listMeta.Url)

	if err != nil {
		UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusURLError)
		return
	}

	fmt.Println(words)
	fmt.Println(title)
	fmt.Println(err)

	if len(words) == 0 {
		utils.PrintR("ProcessVocabularyWords No word found ")
		return
	}

	// crate list record from list meta record
	listId, err := listService.CreateListRecordFromListMeta(listMeta, title)

	if err != nil {
		UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusError)
		return
	}

	// now follow the steps
	listService.ProcessWordsOfSingleGroup(words, listId)

	UpdateListMetaRecordStatus(listService.db, listMeta.Id, enums.ListMetaStatusComplete)

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
			fmt.Println(err)
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
				fmt.Println(err)
				continue
			}

			// now insert into word list table and get the processed data
			processedWordData := processor.ProcessSingleWordData(listService.db, wordModel)

			if len(processedWordData) > 0 {
				// now update the words table with this data
				processor.SaveProcessedDataToWordTable(listService.db, wordModel, processedWordData)
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
		fmt.Println(err)
		return model, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return model, err
	}

	if lastId == 0 {
		return model, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	query := "SELECT * FROM `words` WHERE `id` = ?;"

	result := listService.db.QueryRowx(query, lastId)

	err = result.StructScan(&model)

	if err != nil {
		fmt.Println(err)
		return model, err
	}

	return model, nil
}

func (listService *ListProcessorService) InsertListWordRelation(wordId, listId int64) error {

	queryMap := map[string]interface{}{"word_id": wordId, "list_id": listId, "created_at": time.Now().UTC()}

	res, err := listService.db.NamedExec("Insert into list_word_relation(word_id,list_id,created_at) values(:word_id,:list_id,:created_at)", queryMap)

	if err != nil {
		fmt.Println(err)
		return err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if lastId == 0 {
		return fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	return nil

}

func (listService *ListProcessorService) CheckWordsTableForWord(word string) (uint64, error) {

	query := "SELECT `id` FROM `words` WHERE `word` = ?;"

	result := listService.db.QueryRowx(query, word)
	var ResultId uint64

	err := result.Scan(&ResultId)

	if err != nil {

		fmt.Println(err)
		return ResultId, err
	}

	return ResultId, err

}

func UpdateListMetaRecordStatus(db *sqlx.DB, id uint64, status int) {

	queryMap := map[string]interface{}{"id": id, "status": status, "updated_at": time.Now().UTC()}

	db.NamedExec("Update list_meta set status=:status,updated_at=:updated_at where id=:id", queryMap)

}

func GetWordsFromListMetaRecord(words string) []string {
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

func (listService *ListProcessorService) GenerateUniqueListSlug(title string) string {

	slug := slug.Make(title)
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

	return fmt.Sprintf("%s-%d", slug, 0)
}

// here the title is different because it could be parsed from the other sources like quizlet or so...
func (listService *ListProcessorService) CreateListRecordFromListMeta(listMeta model.ListMetaModel, title string) (int64, error) {

	slug := listService.GenerateUniqueListSlug(title)

	queryMap := map[string]interface{}{"name": title, "slug": slug, "list_meta_id": listMeta.Id, "visibility": listMeta.Visibility, "user_id": listMeta.UserId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := listService.db.NamedExec("Insert into lists(name,slug,list_meta_id,visibility,user_id,created_at,updated_at) values(:name,:slug,:list_meta_id,:visibility,:user_id,:created_at,:updated_at)", queryMap)

	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	if lastId == 0 {
		return -1, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	return lastId, nil

}
