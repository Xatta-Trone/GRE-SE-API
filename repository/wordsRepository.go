package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
)

type WordRepository struct {
	Db *sqlx.DB
}

func NewWordRepository(db *sqlx.DB) *WordRepository {
	return &WordRepository{Db: db}
}

func (rep *WordRepository) FindAll(r requests.WordIndexReqStruct) ([]model.WordModel, error) {

	words := []model.WordModel{}

	queryMap := map[string]interface{}{"word": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// query := fmt.Sprintf("SELECT id,word,word_data,is_reviewed,created_at,updated_at FROM words where word like ? ORDER BY `id` %s LIMIT ?", order)

	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT id,word,word_data,is_reviewed,created_at,updated_at FROM words where word like :word and id > :id order by id %s limit :limit offset :offset", order)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		fmt.Println(err)
		return words, err
	}
	err = nstmt.Select(&words, queryMap)

	// err := rep.Db.Select(&words, query, "%"+r.Query+"%", r.PerPage)

	if err != nil {
		fmt.Println(err)
		return words, err
	}

	return words, nil

}

func (rep *WordRepository) FindOne(id int) (model.WordModel, error) {

	word := model.WordModel{}

	queryMap := map[string]interface{}{"id": id}

	query := fmt.Sprintf("SELECT id,word,word_data,is_reviewed,created_at,updated_at FROM words where id=:id")

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
