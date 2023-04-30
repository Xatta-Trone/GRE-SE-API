package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

// repository
type WordRepositoryInterface interface {
	FindAll(req *requests.WordIndexReqStruct) ([]model.WordModel, error)
	FindAllByListId(req *requests.WordIndexByListIdReqStruct) ([]model.WordModel, error)
	Create(word string, wordData model.WordDataModel, isReviewed int) (model.WordModel, error)
	FindOne(id int) (model.WordModel, error)
	DeleteOne(id int) (bool, error)
	UpdateById(id int, data model.WordDataModel, isReviewed int) (bool, error)
}

type WordRepository struct {
	Db *sqlx.DB
}

func NewWordRepository(db *sqlx.DB) *WordRepository {
	return &WordRepository{Db: db}
}

func (rep *WordRepository) FindAll(r *requests.WordIndexReqStruct) ([]model.WordModel, error) {

	words := []model.WordModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"word": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	order := r.OrderDir // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// query := fmt.Sprintf("SELECT id,word,word_data,is_reviewed,created_at,updated_at FROM words where word like ? ORDER BY `id` %s LIMIT ?", order)

	searchString := "FROM words where word like :word and id > :id"

	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT id,word,is_reviewed,created_at,updated_at %s order by id %s limit :limit offset :offset",searchString, order)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return words, err
	}
	err = nstmt.Select(&words, queryMap)

	// err := rep.Db.Select(&words, query, "%"+r.Query+"%", r.PerPage)

	if err != nil {
		utils.Errorf(err)
		return words, err
	}

	// get the counts
	queryCount := fmt.Sprintf("SELECT count(id) as count %s limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	r.Count = count.Count

	return words, nil

}

func (rep *WordRepository) FindOne(id int) (model.WordModel, error) {

	word := model.WordModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT id,word,word_data,is_reviewed,created_at,updated_at FROM words where id=:id"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return word, err
	}
	err = nstmt.Get(&word, queryMap)

	if err != nil {
		utils.Errorf(err)
		return word, err
	}

	return word, nil

}

func (rep *WordRepository) DeleteOne(id int) (bool, error) {

	query := fmt.Sprintf("Delete FROM words where id=%d", id)

	res, err := rep.Db.Exec(query)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
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

func (rep *WordRepository) UpdateById(id int, wordData model.WordDataModel, isReviewed int) (bool, error) {

	// marshal the data for inserting
	data, err := json.Marshal(wordData)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	queryMap := map[string]interface{}{"id": id, "word_data": string(data), "is_reviewed": isReviewed}

	res, err := rep.Db.NamedExec("Update words set word_data=:word_data,is_reviewed=:is_reviewed, updated_at=now() where id=:id", queryMap)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
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

func (rep *WordRepository) Create(word string, wordData model.WordDataModel, isReviewed int) (model.WordModel, error) {

	newRecord := model.WordModel{}

	// marshal the data for inserting
	data, err := json.Marshal(wordData)

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	queryMap := map[string]interface{}{"word": word, "word_data": string(data), "is_reviewed": isReviewed}

	res, err := rep.Db.NamedExec("Insert into words(word,word_data,is_reviewed,created_at,updated_at) values(:word,:word_data,:is_reviewed,now(),now())", queryMap)

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	if lastId == 0 {
		return newRecord, fmt.Errorf("there was a problem with the insertion. last id: %d", lastId)
	}

	newRecord, err = rep.FindOne(int(lastId))

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	return newRecord, nil

}

func (rep *WordRepository) FindAllByListId(r *requests.WordIndexByListIdReqStruct) ([]model.WordModel, error) {

	words := []model.WordModel{}

	queryMap := map[string]interface{}{"word": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "list_id": r.ListId}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153

	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT id,word,word_data,is_reviewed,created_at,updated_at FROM words where id IN (SELECT word_id FROM list_word_relation WHERE list_id = :list_id) order by id %s limit :limit offset :offset", order)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return words, err
	}
	err = nstmt.Select(&words, queryMap)

	// err := rep.Db.Select(&words, query, "%"+r.Query+"%", r.PerPage)

	if err != nil {
		utils.Errorf(err)
		return words, err
	}

	return words, nil

}
