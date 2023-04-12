package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type FolderRepositoryInterface interface {
	Index(req *requests.FolderIndexReqStruct) ([]model.FolderModel, error)
	Create(req *requests.FolderCreateRequestStruct) (model.FolderModel, error)
	FindOne(id uint64) (model.FolderModel, error)
	Update(id uint64, req *requests.FolderUpdateRequestStruct) (bool, error)
	Delete(folderId uint64, deleteLists bool) (bool, error)
	ToggleList(folderId, listId uint64) (bool, error)
}
type FolderRepository struct {
	Db *sqlx.DB
}

func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{Db: db}
}

func (rep *FolderRepository) Index(r *requests.FolderIndexReqStruct) ([]model.FolderModel, error) {

	models := []model.FolderModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * FROM folders where name like :query and user_id = :user_id order by id %s limit :limit offset :offset", order)

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

func (rep *FolderRepository) Create(req *requests.FolderCreateRequestStruct) (model.FolderModel, error) {

	var newRecord model.FolderModel

	slug := rep.GenerateUniqueFolderSlug(req.Name, 0)

	queryMap := map[string]interface{}{"name": req.Name, "slug": slug, "visibility": req.Visibility, "user_id": req.UserId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Insert into folders(name,slug,visibility,user_id,created_at,updated_at) values(:name,:slug,:visibility,:user_id,:created_at,:updated_at)", queryMap)

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

	newRecord, err = rep.FindOne(uint64(lastId))

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	// now run a process to process these words

	return newRecord, nil

}

func (rep *FolderRepository) FindOne(id uint64) (model.FolderModel, error) {

	modelx := model.FolderModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT * FROM folders where id=:id"

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

func (rep *FolderRepository) Update(id uint64, req *requests.FolderUpdateRequestStruct) (bool, error) {

	slug := rep.GenerateUniqueFolderSlug(req.Name, id)

	queryMap := map[string]interface{}{"id": id, "name": req.Name, "slug": slug, "visibility": req.Visibility, "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Update folders set name=:name,slug=:slug,visibility=:visibility,updated_at=:updated_at where id=:id", queryMap)

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

func (rep *FolderRepository) GenerateUniqueFolderSlug(title string, id uint64) string {

	slug := slug.Make(title)
	// now check the slug

	row := rep.Db.QueryRow("SELECT slug FROM folders WHERE slug like ? and id != ? order by id desc limit 1", fmt.Sprintf("%%%s-%%", slug), id)
	var lastSlug string
	err := row.Scan(&lastSlug)

	if err == sql.ErrNoRows {
		return fmt.Sprintf("%s-%d", slug, 1)
	}

	if err != nil {
		utils.Errorf(err)
		// just add the timestamp and return
		return fmt.Sprintf("%s-%d", slug, time.Now().UnixMilli())
	}

	if lastSlug == "" {
		return fmt.Sprintf("%s-%d", slug, 1)

	}

	// get the serial number form the slug
	slugParts := strings.Split(lastSlug, "-")

	number, err := strconv.ParseInt(slugParts[len(slugParts)-1], 10, 64)

	if err != nil {
		utils.Errorf(err)
		// just add the timestamp and return
		return fmt.Sprintf("%s-%d", slug, time.Now().UnixMilli())
	}

	return fmt.Sprintf("%s-%d", slug, number+1)

}

func (rep *FolderRepository) Delete(folderId uint64, deleteLists bool) (bool, error) {

	queryMap := map[string]interface{}{"id": folderId, "delete_lists": deleteLists}

	// get list ids

	if deleteLists {
		var listIds []uint64
		err := rep.Db.Select(&listIds, "SELECT list_id FROM folder_list_relation where folder_id = ?", folderId)
		if err != nil {
			utils.Errorf(err)
			return false, err
		}
		// delete lists
		// now delete the folder
		query := "Delete FROM folder where id in ?"

		_, err = rep.Db.Exec(query, listIds)

		if err != nil {
			utils.Errorf(err)
			return false, err
		}

	}

	// now delete the folder
	query := "Delete FROM folders where id=:id"

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

func (rep *FolderRepository) ToggleList(folderId, listId uint64) (bool, error) {
	// check if exists
	modelx := model.FolderListRelationModel{}

	queryMap := map[string]interface{}{"folder_id": folderId, "list_id": int64(listId)}

	query := "SELECT folder_id,list_id FROM folder_list_relation where folder_id=:folder_id and list_id=:list_id"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}
	err = nstmt.Get(&modelx, queryMap)

	fmt.Println(modelx,err)

	if err == sql.ErrNoRows {
		// insert the record
		_, err := rep.Db.NamedExec("Insert into folder_list_relation(folder_id,list_id) values(:folder_id,:list_id)", queryMap)

		if err != nil {
			utils.Errorf(err)
			return false, err
		}

		return true, nil

	} else {
		// delete the record
		// now delete the folder
		query := "Delete FROM folder_list_relation where folder_id=:folder_id and list_id=:list_id"

		_, err := rep.Db.NamedExec(query, queryMap)

		if err != nil {
			utils.Errorf(err)
			return false, err
		}

		return true, nil
	}

}
