package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type ListRepositoryInterface interface {
	Create(req *requests.ListsCreateRequestStruct) (model.ListMetaModel, error)
	Index(req *requests.ListsIndexReqStruct) ([]model.ListModel, error)
	Update(id uint64, req *requests.ListsUpdateRequestStruct) (bool, error)
	FindOneBySlug(slug string) (model.ListModel, error)
	DeleteFromListMeta(listMetaId uint64) (bool, error)
	Delete(listMetaId uint64) (bool, error)
	DeleteWordInList(wordId, listId uint64) (bool, error)
	ListsByFolderId(req *requests.FolderListIndexReqStruct) ([]model.ListModel, error)
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

func (rep *ListRepository) ListsByFolderId(r *requests.FolderListIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "folder_id": r.FolderId}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * FROM lists where (name like :query and user_id = :user_id) and id in (select list_id from folder_list_relation where folder_id = :folder_id ) order by id %s limit :limit offset :offset", order)

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

func (rep *ListRepository) FindOneBySlug(slug string) (model.ListModel, error) {

	modelx := model.ListModel{}

	queryMap := map[string]interface{}{"slug": slug}

	query := "SELECT * FROM lists where slug=:slug"

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

func (rep *ListRepository) Update(id uint64, req *requests.ListsUpdateRequestStruct) (bool, error) {

	slug := rep.GenerateUniqueListSlug(req.Name, id)

	queryMap := map[string]interface{}{"id": id, "name": req.Name, "slug": slug, "visibility": req.Visibility, "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Update lists set name=:name,slug=:slug,visibility=:visibility,updated_at=:updated_at where id=:id", queryMap)

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

func (rep *ListRepository) GenerateUniqueListSlug(title string, id uint64) string {

	slug := slug.Make(title)
	// now check the slug

	row := rep.Db.QueryRow("SELECT Count(id) FROM lists WHERE slug like ? and id != ?", fmt.Sprintf("%%%s-%%", slug), id)
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

func (rep *ListRepository) DeleteFromListMeta(listMetaId uint64) (bool, error) {

	queryMap := map[string]interface{}{"id": listMetaId}

	query := "Delete FROM list_meta where id=:id"

	res, err := rep.Db.NamedExec(query, queryMap)

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

func (rep *ListRepository) Delete(listId uint64) (bool, error) {

	queryMap := map[string]interface{}{"id": listId}

	query := "Delete FROM lists where id=:id"

	res, err := rep.Db.NamedExec(query, queryMap)

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

func (rep *ListRepository) DeleteWordInList(wordId, listId uint64) (bool, error) {

	queryMap := map[string]interface{}{"list_id": listId, "word_id": wordId}

	query := "Delete FROM list_word_relation where list_id=:list_id and word_id=:word_id "

	res, err := rep.Db.NamedExec(query, queryMap)

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
