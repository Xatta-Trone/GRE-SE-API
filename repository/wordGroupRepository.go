package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
)

type WordGroupInterface interface {
	Create(req *requests.WordGroupCreateReqStruct) (model.WordGroupModel, error)
	FindOne(id int) (model.WordGroupModel, error)
	DeleteOne(id int) (bool, error)
}

type WordGroupRepository struct {
	Db *sqlx.DB
}

func NewWordGroupRepository(db *sqlx.DB) *WordGroupRepository {
	return &WordGroupRepository{Db: db}
}

func (rep *WordGroupRepository) Create(req *requests.WordGroupCreateReqStruct) (model.WordGroupModel, error) {

	var newRecord model.WordGroupModel

	queryMap := map[string]interface{}{"name": req.Name, "words": req.Words, "status": enums.WordGroupUploaded, "file_name": req.FileName}

	res, err := rep.Db.NamedExec("Insert into  word_groups(name,words,status,file_name, created_at,updated_at) values(:name,:words,:status,:file_name,now(),now())", queryMap)

	if err != nil {
		fmt.Println(err)
		return newRecord, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return newRecord, err
	}

	if lastId == 0 {
		return newRecord, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	newRecord, err = rep.FindOne(int(lastId))

	if err != nil {
		fmt.Println(err)
		return newRecord, err
	}

	return newRecord, nil

}

func (rep *WordGroupRepository) FindOne(id int) (model.WordGroupModel, error) {

	word := model.WordGroupModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT * FROM word_groups where id=:id"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		fmt.Println(err)
		return word, err
	}
	err = nstmt.Get(&word, queryMap)

	if err != nil {
		fmt.Println(err)
		return word, err
	}

	return word, nil

}

func (rep *WordGroupRepository) DeleteOne(id int) (bool, error) {
	word := model.WordGroupModel{}

	word, err := rep.FindOne(id)

	if err != nil {
		return false, err
	}

	// delete the file
	if word.FileName != nil {
		err := os.Remove(*word.FileName)
		if err != nil {
			return false, err
		}
	}
	// delete the record
	query := fmt.Sprintf("Delete FROM word_groups where id=%d", id)

	res, err := rep.Db.Exec(query)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if rows == 0 {
		return false, sql.ErrNoRows
	}

	if rows != 1 {

		return false, fmt.Errorf("number of rows affected %d", rows)
	}

	return true, nil

}
