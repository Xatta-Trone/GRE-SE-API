package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type FolderRepositoryInterface interface {
	Index(req *requests.FolderIndexReqStruct) ([]model.FolderModel, error)
	PublicIndex(req *requests.PublicFolderIndexReqStruct) ([]model.FolderModel, error)
	SavedFolders(req *requests.SavedFolderIndexReqStruct) ([]model.FolderModel, error)
	AdminIndex(req *requests.FolderIndexReqStruct) ([]model.FolderModel,model.CountModel, error)
	Create(req *requests.FolderCreateRequestStruct) (model.FolderModel, error)
	SaveFolder(userId, folderId uint64) (bool, error)
	FindOne(id uint64) (model.FolderModel, error)
	Update(id uint64, req *requests.FolderUpdateRequestStruct) (bool, error)
	Delete(folderId uint64, deleteLists bool) (bool, error)
	DeleteSavedFolder(userId, folderId uint64) (bool, error)
	ToggleList(folderId, listId uint64) (bool, error)
	GetCount(ids []uint64) ([]model.FolderListRelationModel, error)
}
type FolderRepository struct {
	Db *sqlx.DB
}

func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{Db: db}
}

func (rep *FolderRepository) Index(r *requests.FolderIndexReqStruct) ([]model.FolderModel, error) {

	models := []model.FolderModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId,  "public_visibility":enums.FolderVisibilityPublic, "filter": r.Filter }

	fmt.Println(r.UserId)

	// select the filter 
	filterQuery := ""
	filterCreatedSql := "(saved_folders.folder_id = folders.id AND saved_folders.user_id = :user_id AND folders.user_id = :user_id)"
	filterSavedSql := "(saved_folders.folder_id = folders.id AND saved_folders.user_id = :user_id AND folders.user_id != :user_id AND folders.visibility = :public_visibility)"

	if (r.Filter == enums.FolderFilterAll) {
		filterQuery = filterCreatedSql + " OR " + filterSavedSql
	}

	if (r.Filter == enums.FolderFilterCrated) {
		filterQuery = filterCreatedSql
	}

	if (r.Filter == enums.FolderFilterSaved) {
		filterQuery = filterSavedSql
	}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * FROM folders where id IN (SELECT saved_folders.folder_id FROM saved_folders INNER JOIN folders ON %s order by saved_folders.created_at %s) and name like :query limit :limit offset :offset",filterQuery, order)

	searchStringCount := fmt.Sprintf("FROM folders where id IN (SELECT saved_folders.folder_id FROM saved_folders INNER JOIN folders ON %s order by saved_folders.created_at %s) and name like :query",filterQuery, order)

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
	queryCount := fmt.Sprintf("SELECT count(folders.id) as count %s limit 1", searchStringCount)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)
	r.Count = count.Count

	return models, nil

}

func (rep *FolderRepository) PublicIndex(r *requests.PublicFolderIndexReqStruct) ([]model.FolderModel, error) {

	models := []model.FolderModel{}
	count := model.CountModel{}


	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "order": r.Order, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "visibility": enums.FolderVisibilityPublic}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	searchString := "FROM folders INNER JOIN folder_list_relation ON folder_list_relation.folder_id = folders.id where folders.name like :query and folders.visibility=:visibility"
	searchStringCount := "FROM folders where folders.name like :query and folders.visibility=:visibility"

	query := fmt.Sprintf("SELECT folders.*, COUNT(folder_list_relation.list_id) AS lists_count %s GROUP BY folders.id order by %s %s limit :limit offset :offset",searchString, queryMap["orderby"], order)

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
	queryCount := fmt.Sprintf("SELECT count(folders.id) as count %s limit 1", searchStringCount)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)
	r.Count = count.Count

	return models, nil

}

func (rep *FolderRepository) SavedFolders(r *requests.SavedFolderIndexReqStruct) ([]model.FolderModel, error) {

	models := []model.FolderModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "order": r.Order, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "visibility": enums.FolderVisibilityPublic}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT folders.*, COUNT(folder_list_relation.list_id) AS lists_count FROM folders INNER JOIN folder_list_relation ON folder_list_relation.folder_id = folders.id where folders.id IN (select folder_id from saved_folders where user_id=:user_id) and folders.name like :query GROUP BY folders.id order by %s %s limit :limit offset :offset",queryMap["orderby"], order)

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

func (rep *FolderRepository) AdminIndex(r *requests.FolderIndexReqStruct) ([]model.FolderModel, model.CountModel, error) {

	models := []model.FolderModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId}

	fmt.Println(r.UserId)

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	searchString := "FROM folders INNER JOIN users on folders.user_id = users.id where folders.name like :query or users.name like :query or users.email like :query or users.username like :query"


	query := fmt.Sprintf("SELECT folders.* %s order by folders.id %s limit :limit offset :offset",searchString, order)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return models,count, err
	}
	err = nstmt.Select(&models, queryMap)

	if err != nil {
		utils.Errorf(err)
		return models,count, err
	}

	// get the counts
	queryCount := fmt.Sprintf("SELECT count(folders.id) as count %s limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	return models,count, nil

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

	if lastId != 0 {
		// create the folder list relation
		queryMapForListFolderRelation := map[string]interface{}{"folder_id": lastId, "user_id": req.UserId, "created_at": time.Now().UTC()}
		// insert into saved folders
		_, err = rep.Db.NamedExec("Insert into saved_folders(user_id,folder_id,created_at) values(:user_id,:folder_id,:created_at)", queryMapForListFolderRelation)
		if err != nil {
			utils.Errorf(err)
			utils.PrintR("there was an error creating list folder relation \n")

		}

	}

	newRecord, err = rep.FindOne(uint64(lastId))

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	// now run a process to process these words

	return newRecord, nil

}

func (rep *FolderRepository) SaveFolder(userId, folderId uint64) (bool, error) {


	queryMap := map[string]interface{}{ "user_id": userId, "folder_id": folderId, "created_at": time.Now().UTC()}

	_, err := rep.Db.NamedExec("Insert ignore into saved_folders(user_id,folder_id,created_at) values(:user_id,:folder_id,:created_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}


	return true, nil

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

func (rep *FolderRepository) DeleteSavedFolder(userId, folderId uint64) (bool, error) {

	queryMap := map[string]interface{}{"folder_id": folderId, "user_id": userId}

	// now delete the folder
	query := "Delete FROM saved_folders where folder_id=:folder_id and user_id=:user_id"

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

	fmt.Println(modelx, err)

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


func (rep *FolderRepository) GetCount(ids []uint64) ([]model.FolderListRelationModel, error) {


	models := []model.FolderListRelationModel{}

	query, args, err := sqlx.In("SELECT folder_list_relation.folder_id, count(folder_list_relation.list_id) as list_count FROM folder_list_relation where folder_id in (?) group by folder_id", ids)

	if err != nil {
		utils.Errorf(err)
		return models, err
	}

	query = rep.Db.Rebind(query)

	err = rep.Db.Select(&models, query, args...)

	if err != nil {
		utils.Errorf(err)
		return models, err
	}

	return models, nil

}