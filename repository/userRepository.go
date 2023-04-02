package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
)

type UserRepositoryInterface interface {
	Index(req requests.UsersIndexReqStruct) ([]model.UserModel, error)
	FindOne(id int) (model.UserModel, error)
	Delete(id int) (bool, error)
	Create(req *requests.UsersCreateRequestStruct) (model.UserModel, error)
	Update(id int, req *requests.UsersUpdateRequestStruct) (bool, error)
}

type UserRepository struct {
	Db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (rep *UserRepository) Index(r requests.UsersIndexReqStruct) ([]model.UserModel, error) {

	models := []model.UserModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	order := r.OrderBy // problem with order by https://github.com/jmoiron/sqlx/issues/153
	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT id,name,email,username,created_at FROM users where name like :query or email like :query or username like :query  and id > :id order by id %s limit :limit offset :offset", order)

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		fmt.Println(err)
		return models, err
	}
	err = nstmt.Select(&models, queryMap)

	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return modelx, err
	}
	err = nstmt.Get(&modelx, queryMap)

	if err != nil {
		fmt.Println(err)
		return modelx, err
	}

	return modelx, nil

}

func (rep *UserRepository) Delete(id int) (bool, error) {

	queryMap := map[string]interface{}{"id": id}

	query := "Delete FROM users where id=:id"

	res, err := rep.Db.NamedExec(query, queryMap)

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

func (rep *UserRepository) Create(req *requests.UsersCreateRequestStruct) (model.UserModel, error) {

	var newRecord model.UserModel

	username := services.GeneRateUserName(*rep.Db)

	queryMap := map[string]interface{}{"name": req.Name, "email": req.Email, "username": username, "created_at": time.Now().UTC(), "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Insert into  users(name,email,username, created_at,updated_at) values(:name,:email,:username,now(),now())", queryMap)

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

func (rep *UserRepository) Update(id int, req *requests.UsersUpdateRequestStruct) (bool, error) {

	queryMap := map[string]interface{}{"id": id, "name": req.Name, "email": req.Email, "username": req.UserName, "updated_at": time.Now().UTC()}

	res, err := rep.Db.NamedExec("Update users set name=:name,email=:email,username=:username,updated_at=:updated_at where id=:id", queryMap)

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
