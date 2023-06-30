package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type PendingWordsInterface interface {
	Index(req *requests.PendingWordIndexRequestStruct) ([]model.PendingWordModel, error)
	Delete(id int, word string) (bool, error)
	Update(r requests.PendingWordsUpdateRequestStruct) (bool, error)
}

type PendingWordsRepository struct {
	Db *sqlx.DB
}

func NewPendingWordsRepository(db *sqlx.DB) *PendingWordsRepository {
	return &PendingWordsRepository{Db: db}
}

func (rep *PendingWordsRepository) Index(r *requests.PendingWordIndexRequestStruct) ([]model.PendingWordModel, error) {

	models := []model.PendingWordModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"word": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	searchString := "FROM pending_words where word like :word"

	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * %s limit :limit offset :offset", searchString)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return models, err
	}
	err = nstmt.Select(&models, queryMap)

	if err != nil {
		utils.Errorf(err)
		return models, err
	}

	// get the counts
	queryCount := fmt.Sprintf("SELECT count(word) as count %s limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	r.Count = count.Count

	return models, nil

}

func (rep *PendingWordsRepository) Delete(id int, word string) (bool, error) {

	queryMap := map[string]interface{}{"id": id, "word": word}

	res, err := rep.Db.NamedExec("DELETE FROM `pending_words` WHERE list_id=:id and word=:word", queryMap)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	_, err = res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	return true, nil

}

func (rep *PendingWordsRepository) Update(r requests.PendingWordsUpdateRequestStruct) (bool, error) {

	queryMap := map[string]interface{}{"list_id": r.ListId, "word": r.Word}

	res, err := rep.Db.NamedExec("Update `pending_words` set approved=1 WHERE list_id=:list_id and word=:word", queryMap)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	_, err = res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	return true, nil
}