package processor

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
)

func ProcessListMetaRecord(db *sqlx.DB, listMeta model.ListMetaModel) {
	// update the list meta table
	UpdateListMetaRecordStatus(db, listMeta.Id, enums.ListMetaStatusParsing)

	// now check the type of word to be processed...URL or word

	if listMeta.Url != nil {
		// fire url processor
		fmt.Println(listMeta.Url)
	}

	if listMeta.Words != nil {
		// fire words processor
		fmt.Println(listMeta.Words)
	}

}

func UpdateListMetaRecordStatus(db *sqlx.DB, id uint64, status int) {

	queryMap := map[string]interface{}{"id": id, "status": status, "updated_at": time.Now().UTC()}

	db.NamedExec("Update list_meta set status=:status,updated_at=:updated_at where id=:id", queryMap)

}

func GetWordsFromListMetaRecord(words string) []string {
	var processedWords []string

	// split by new line
	// tempData := strings
	byNewLine := strings.Split(words, "\n")

	for _, value := range byNewLine {
		//    split by comma
		tempWords := strings.Split(value, ",")

		for _, v := range tempWords {
			fmt.Println(v)

		}
	}

	return processedWords
}