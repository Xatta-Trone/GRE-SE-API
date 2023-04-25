package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
	"github.com/xatta-trone/words-combinator/utils"
)

type UserRepositoryInterface interface {
	Index(req requests.UsersIndexReqStruct) ([]model.UserModel, model.CountModel, error)
	In(ids []uint64, columns ...string) ([]model.UserModel, error)
	FindOne(id int) (model.UserModel, error)
	FindOneByEmail(email string) (model.UserModel, error)
	FindOneByUserName(username string) (model.UserModel, error)
	Delete(id int) (bool, error)
	Create(req *requests.UsersCreateRequestStruct) (model.UserModel, error)
	Update(id int, req *requests.UsersUpdateRequestStruct) (bool, error)
	UpdateUserName(id uint64, username string) (bool, error)
}

type UserRepository struct {
	Db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (rep *UserRepository) Index(r requests.UsersIndexReqStruct) ([]model.UserModel, model.CountModel, error) {

	models := []model.UserModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy,"order": r.Order, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	order := r.OrderDir // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	
	searchString := "FROM users where name like :query or email like :query or username like :query  and id > :id"

	query := fmt.Sprintf("SELECT id,name,email,username,created_at,updated_at %s order by %s %s limit :limit offset :offset", searchString,r.OrderBy, order)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return models, count, err
	}
	err = nstmt.Select(&models, queryMap)

	if err != nil {
		utils.Errorf(err)
		return models, count, err
	}

	// get the counts
	queryCount := fmt.Sprintf("SELECT count(id) as count %s limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	return models, count, nil

}

func (rep *UserRepository) In(ids []uint64, columns ...string) ([]model.UserModel, error) {

	columnsToSelect := "*"

	if len(columns) > 0 {
		columnsToSelect = ""

		for i, col := range columns {
			if i == len(columns)-1 {
				columnsToSelect += col
			} else {
				columnsToSelect += fmt.Sprintf("%s,", col)
			}

		}
	}

	models := []model.UserModel{}

	query, args, err := sqlx.In(fmt.Sprintf("SELECT %s FROM users where id in (?)", columnsToSelect), ids)

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
func (rep *UserRepository) FindOne(id int) (model.UserModel, error) {

	modelx := model.UserModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT id,name,email,username,created_at FROM users where id=:id"

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

func (rep *UserRepository) FindOneByEmail(email string) (model.UserModel, error) {

	modelx := model.UserModel{}

	queryMap := map[string]interface{}{"email": email}

	query := "SELECT id,name,email,username,created_at FROM users where email=:email"

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

func (rep *UserRepository) FindOneByUserName(username string) (model.UserModel, error) {

	modelx := model.UserModel{}

	queryMap := map[string]interface{}{"username": username}

	query := "SELECT id,name,email,username,created_at FROM users where username=:username"

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

func (rep *UserRepository) Delete(id int) (bool, error) {

	queryMap := map[string]interface{}{"id": id}

	query := "Delete FROM users where id=:id"

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

func (rep *UserRepository) Create(req *requests.UsersCreateRequestStruct) (model.UserModel, error) {

	var newRecord model.UserModel

	username := services.GeneRateUserName(*rep.Db)

	queryMap := map[string]interface{}{"name": req.Name, "email": req.Email, "username": username, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Insert into  users(name,email,username, created_at,updated_at) values(:name,:email,:username,:created_at,:updated_at)", queryMap)

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

func (rep *UserRepository) Update(id int, req *requests.UsersUpdateRequestStruct) (bool, error) {

	queryMap := map[string]interface{}{"id": id, "name": req.Name, "email": req.Email, "username": req.UserName, "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Update users set name=:name,email=:email,username=:username,updated_at=:updated_at where id=:id", queryMap)

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

func (rep *UserRepository) UpdateUserName(id uint64, username string) (bool, error) {

	queryMap := map[string]interface{}{"id": id, "username": username, "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Update users set username=:username,updated_at=:updated_at where id=:id", queryMap)

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
