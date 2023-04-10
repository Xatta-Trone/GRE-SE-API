package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type ListRepositoryInterface interface {
	Create(req *requests.ListsCreateRequestStruct) (model.ListMetaModel, error)
	Index(req *requests.ListsIndexReqStruct) ([]model.ListModel, error)
}

type ListRepository struct {
	Db *sqlx.DB
}

func NewListRepository(db *sqlx.DB) *ListRepository {
	return &ListRepository{Db: db}
}

func (rep *ListRepository) Index(r *requests.ListsIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT id,name,slug,visibility,list_meta_id,status,created_at,updated_at FROM lists where name like :query and user_id = :user_id order by id %s limit :limit offset :offset", order)

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


	return models, nil

}

func (rep *ListRepository) Create(req *requests.ListsCreateRequestStruct) (model.ListMetaModel, error) {

	var newRecord model.ListMetaModel

	queryMap := map[string]interface{}{"name": req.Name, "url": req.Url, "words": req.Words, "visibility": req.Visibility, "user_id": req.UserId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Insert into list_meta(name,url,words,visibility,user_id,created_at,updated_at) values(:name,nullif(:url,\"\"),nullif(:words,\"\"),:visibility,:user_id,created_at,:updated_at)", queryMap)

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

	// now run a process to process these words

	return newRecord, nil

}

func (rep *ListRepository) FindOne(id int) (model.ListMetaModel, error) {

	modelx := model.ListMetaModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT * FROM list_meta where id=:id"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return modelx, err
	}
	err = nstmt.Get(&modelx, queryMap)

	if err != nil {
		utils.Errorf(err)
		return modelx, err
	}

	return modelx, nil

}
