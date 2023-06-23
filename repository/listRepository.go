package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type ListRepositoryInterface interface {
	Create(req *requests.ListsCreateRequestStruct) (model.ListMetaModel, error)
	SaveListItem(req *requests.SavedListsCreateRequestStruct) (bool, error)
	Index(req *requests.ListsIndexReqStruct) ([]model.ListModel, error)
	PublicIndex(req *requests.PublicListsIndexReqStruct) ([]model.ListModel, error)
	SavedLists(req *requests.SavedListsIndexReqStruct) ([]model.ListModel, error)
	AdminIndex(req *requests.ListsIndexReqStruct) ([]model.ListModel, error)
	Update(id uint64, req *requests.ListsUpdateRequestStruct) (bool, error)
	FindOneBySlug(slug string) (model.ListModel, error)
	FindOne(id uint64) (model.ListModel, error)
	DeleteFromListMeta(listMetaId uint64) (bool, error)
	Delete(listMetaId uint64) (bool, error)
	DeleteFromSavedList(userId, listId uint64) (bool, error)
	DeleteWordInList(wordId, listId uint64) (bool, error)
	ListsByFolderId(req *requests.FolderListIndexReqStruct) ([]model.ListModel, error)
	GetCount(ids []uint64) ([]model.ListWordModel, error)
	FoldersByListId(listId, userId uint64) ([]model.FolderListRelationModel, error)
}

type ListRepository struct {
	Db *sqlx.DB
}

func NewListRepository(db *sqlx.DB) *ListRepository {
	return &ListRepository{Db: db}
}

func (rep *ListRepository) Index(r *requests.ListsIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "public_visibility": enums.ListVisibilityPublic}

	// select the filter
	filterQuery := ""
	filterCreatedSql := "(saved_lists.list_id = lists.id AND saved_lists.user_id = :user_id AND lists.user_id = :user_id)"
	filterSavedSql := "(saved_lists.list_id = lists.id AND saved_lists.user_id = :user_id AND lists.user_id != :user_id AND lists.visibility = :public_visibility)"

	if r.Filter == enums.ListFilterAll {
		filterQuery = filterCreatedSql + " OR " + filterSavedSql
	}

	if r.Filter == enums.ListFilterCrated {
		filterQuery = filterCreatedSql
	}

	if r.Filter == enums.ListFilterSaved {
		filterQuery = filterSavedSql
	}

	order := r.Order         // problem with order by https://github.com/jmoiron/sqlx/issues/153
	saveOrder := r.SaveOrder // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * FROM lists where id IN (SELECT saved_lists.list_id FROM saved_lists INNER JOIN lists ON %s order by saved_lists.created_at %s) and name like :query order by %s %s limit :limit offset :offset", filterQuery, saveOrder, r.OrderBy, order)

	searchStringCount := fmt.Sprintf("FROM lists where id IN (SELECT saved_lists.list_id FROM saved_lists INNER JOIN lists ON %s order by saved_lists.created_at %s) and name like :query", filterQuery, order)

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
	queryCount := fmt.Sprintf("SELECT count(lists.id) as count %s limit 1", searchStringCount)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)
	r.Count = count.Count

	return models, nil

}

func (rep *ListRepository) PublicIndex(r *requests.PublicListsIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": "lists." + r.OrderBy, "order": r.Order, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "visibility": enums.ListVisibilityPublic}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153

	// if there is user id and user name present then filter by the user
	userFilter := ""

	if r.UserName != "" && r.UserId != 0 {
		userFilter = "and lists.user_id=:user_id"
	}

	searchString := fmt.Sprintf("FROM lists INNER JOIN list_word_relation ON list_word_relation.list_id = lists.id where lists.name like :query and lists.visibility=:visibility %s", userFilter)
	searchStringCount := fmt.Sprintf("FROM lists where lists.name like :query and lists.visibility=:visibility %s", userFilter)

	query := fmt.Sprintf("SELECT lists.*, COUNT(list_word_relation.word_id) AS word_count %s GROUP BY lists.id order by %s %s limit :limit offset :offset", searchString, queryMap["orderby"], order)

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
	queryCount := fmt.Sprintf("SELECT count(lists.id) as count %s limit 1", searchStringCount)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)
	r.Count = count.Count

	return models, nil

}

func (rep *ListRepository) SavedLists(r *requests.SavedListsIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}

	fmt.Println(*r)

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": "lists." + r.OrderBy, "order": r.Order, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "visibility": enums.ListVisibilityPublic}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153

	query := fmt.Sprintf("SELECT lists.*, COUNT(list_word_relation.word_id) AS word_count FROM lists INNER JOIN list_word_relation ON list_word_relation.list_id = lists.id where lists.id IN (select saved_lists.list_id from saved_lists where saved_lists.user_id = :user_id) and lists.name like :query GROUP BY lists.id order by %s %s limit :limit offset :offset", queryMap["orderby"], order)

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

func (rep *ListRepository) AdminIndex(r *requests.ListsIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153

	searchString := "FROM lists INNER JOIN users on lists.user_id = users.id where lists.name like :query or users.name like :query or users.email like :query or users.username like :query"
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT lists.* %s order by lists.id %s limit :limit offset :offset", searchString, order)

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
	queryCount := fmt.Sprintf("SELECT count(lists.id) as count %s limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	r.Count = count.Count

	return models, nil

}

func (rep *ListRepository) ListsByFolderId(r *requests.FolderListIndexReqStruct) ([]model.ListModel, error) {

	models := []model.ListModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "order": r.Order, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage, "user_id": r.UserId, "folder_id": r.FolderId}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * FROM lists where (name like :query) and id in (select list_id from folder_list_relation where folder_id = :folder_id ) order by %s %s limit :limit offset :offset", r.Order, order)

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

	queryMap := map[string]interface{}{"name": req.Name, "url": req.Url, "words": req.Words, "visibility": req.Visibility, "user_id": req.UserId, "folder_id": req.FolderId, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Insert into list_meta(name,url,words,visibility,user_id,folder_id,created_at,updated_at) values(:name,nullif(:url,\"\"),nullif(:words,\"\"),:visibility,:user_id,nullif(:folder_id,0),created_at,:updated_at)", queryMap)

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

	newRecord, err = rep.FindOneListMeta(uint64(lastId))

	if err != nil {
		utils.Errorf(err)
		return newRecord, err
	}

	// now run a process to process these words

	return newRecord, nil

}

func (rep *ListRepository) SaveListItem(req *requests.SavedListsCreateRequestStruct) (bool, error) {

	queryMap := map[string]interface{}{"user_id": req.UserId, "list_id": req.ListId, "created_at": time.Now().UTC()}

	_, err := rep.Db.NamedExec("Insert ignore into saved_lists(user_id,list_id,created_at) values(:user_id,:list_id, :created_at)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return false, err
	}

	return true, nil

}

func (rep *ListRepository) FindOneListMeta(id uint64) (model.ListMetaModel, error) {

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

func (rep *ListRepository) FindOne(id uint64) (model.ListModel, error) {

	modelx := model.ListModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT * FROM lists where id=:id"

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

	wordCounts, _ := rep.GetCount([]uint64{modelx.Id})

	for _, wordCount := range wordCounts {
		if wordCount.ListId == modelx.Id {
			modelx.WordCount = wordCount.WordCount
		}
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

	wordCounts, _ := rep.GetCount([]uint64{modelx.Id})

	for _, wordCount := range wordCounts {
		if wordCount.ListId == modelx.Id {
			modelx.WordCount = wordCount.WordCount
		}
	}

	return modelx, nil

}

func (rep *ListRepository) Update(id uint64, req *requests.ListsUpdateRequestStruct) (bool, error) {

	// check if name changed
	if req.Slug == "" {
		req.Slug = rep.GenerateUniqueListSlug(req.Name, id)
	}

	queryMap := map[string]interface{}{"id": id, "name": req.Name, "slug": req.Slug, "visibility": req.Visibility, "updated_at": time.Now().UTC()}

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
	timestampedSlug := fmt.Sprintf("%s-%d", slug, time.Now().UnixMilli())

	// fmt.Println(slug, fmt.Sprintf("%s-%%", slug), totalCount)

	if err != nil {
		// just add the timestamp and return
		return timestampedSlug
	}

	if totalCount > 0 {

		newSlug := fmt.Sprintf("%s-%d", slug, totalCount+1)
		// check if this slug exists
		row := rep.Db.QueryRow("SELECT Count(id) FROM lists WHERE slug like ? and id != ?", fmt.Sprintf("%%%s-%%", newSlug), id)
		var totalCount int
		err := row.Scan(&totalCount)

		// fmt.Println(slug, fmt.Sprintf("%s-%%", slug), totalCount)

		if err != sql.ErrNoRows {
			// we can return the new slug
			return newSlug
		}

		return timestampedSlug

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

func (rep *ListRepository) DeleteFromSavedList(userId, listId uint64) (bool, error) {

	queryMap := map[string]interface{}{"list_id": listId, "user_id": userId}

	query := "Delete FROM saved_lists where list_id=:list_id and user_id=:user_id"

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

func (rep *ListRepository) GetCount(ids []uint64) ([]model.ListWordModel, error) {

	models := []model.ListWordModel{}

	query, args, err := sqlx.In("SELECT list_word_relation.list_id, count(list_word_relation.word_id) as word_count FROM list_word_relation where list_id in (?) group by list_id", ids)

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

func (rep *ListRepository) FoldersByListId(listId, userId uint64) ([]model.FolderListRelationModel, error) {
	models := []model.FolderListRelationModel{}

	queryMap := map[string]interface{}{"user_id": userId, "list_id": listId}

	// get folders created by this user
	// select the filter
	fmt.Println(listId)

	// I am using named execution to make it more clear
	query := "SELECT folder_id FROM folder_list_relation where list_id = :list_id "

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
