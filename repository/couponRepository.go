package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type CouponInterface interface {
	Index(req *requests.CouponIndexRequestStruct) ([]model.CouponModel, error)
	Create(req requests.CouponCreateRequestStruct, coupon string) (int64, error)
	FindOne(id uint64) (model.CouponModel, error)
	FindByCoupon(coupon string) (model.CouponModel, error)
	Delete(id int) (bool, error)
	UpdateUserId(id, userId uint64) (bool, error)
}

type CouponRepository struct {
	Db *sqlx.DB
}

func NewCouponRepository(db *sqlx.DB) *CouponRepository {
	return &CouponRepository{Db: db}
}

func (rep *CouponRepository) Index(r *requests.CouponIndexRequestStruct) ([]model.CouponModel, error) {

	models := []model.CouponModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"coupon": "%" + r.Query + "%", "id": r.ID, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	order := r.Order // problem with order by https://github.com/jmoiron/sqlx/issues/153

	searchString := "FROM coupons where coupon like :coupon and id > :id"

	// I am using named execution to make it more clear
	query := fmt.Sprintf("SELECT * %s order by id %s limit :limit offset :offset", searchString, order)

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
	queryCount := fmt.Sprintf("SELECT count(id) as count %s limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	r.Count = count.Count

	return models, nil

}

func (rep *CouponRepository) Create(req requests.CouponCreateRequestStruct,coupon string) (int64, error) {

	queryMap := map[string]interface{}{"coupon": coupon, "max_use": req.MaxUse, "expires": "", "months": req.Months}

	if req.Expires != "" {
		expires,_ := time.Parse("2006-01-02", req.Expires)
		queryMap["expires"] = expires
	}

	res, err := rep.Db.NamedExec("Insert into coupons(coupon,max_use,expires,months) values(:coupon,nullif(:max_use,0),nullif(:expires,\"\"), :months)", queryMap)

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	LastInsertId, err := res.LastInsertId()

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	return LastInsertId, nil

}

func (rep *CouponRepository) FindOne(id uint64) (model.CouponModel, error) {

	model := model.CouponModel{}

	queryMap := map[string]interface{}{"id": id}

	query := "SELECT * FROM coupons where id=:id"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return model, err
	}
	err = nstmt.Get(&model, queryMap)

	if err != nil {
		utils.Errorf(err)
		return model, err
	}

	return model, nil

}

func (rep *CouponRepository) FindByCoupon(coupon string) (model.CouponModel, error) {

	model := model.CouponModel{}

	queryMap := map[string]interface{}{"coupon": coupon}

	query := "SELECT * FROM coupons where coupon=:coupon"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return model, err
	}
	err = nstmt.Get(&model, queryMap)

	if err != nil {
		utils.Errorf(err)
		return model, err
	}

	return model, nil

}

func (rep *CouponRepository) Delete(id int) (bool, error) {

	queryMap := map[string]interface{}{"id": id,}

	res, err := rep.Db.NamedExec("DELETE FROM `coupons` WHERE id=:id", queryMap)

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

func (rep *CouponRepository) UpdateUserId(id, userId uint64) (bool, error) {

	queryMap := map[string]interface{}{"id": id,"user_id": userId}

	res, err := rep.Db.NamedExec("Update `coupons` set user_id=:user_id WHERE id=:id", queryMap)

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