package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type CouponInterface interface {
	Index() ([]model.CouponModel, error)
	Create(req requests.CouponCreateRequestStruct, coupon string) (int64, error)
	FindOne(id uint64) (model.CouponModel, error)
	FindByCoupon(coupon string) (model.CouponModel, error)
	Delete(id int) (bool, error)
}

type CouponRepository struct {
	Db *sqlx.DB
}

func NewCouponRepository(db *sqlx.DB) *CouponRepository {
	return &CouponRepository{Db: db}
}

func (rep *CouponRepository) Index() ([]model.CouponModel, error) {

	models := []model.CouponModel{}
	queryMap := map[string]interface{}{}

	query := "select * from coupons order by id desc"

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

func (rep *CouponRepository) Create(req requests.CouponCreateRequestStruct,coupon string) (int64, error) {

	queryMap := map[string]interface{}{"coupon": coupon, "max_use": req.MaxUse, "expires": ""}

	if req.Months > 0 {
		expires := time.Now().UTC().AddDate(0, req.Months, 0)
		queryMap["expires"] = expires
	}

	res, err := rep.Db.NamedExec("Insert into coupons(coupon,max_use,expires) values(:coupon,nullif(:max_use,0),nullif(:expires,\"\"))", queryMap)

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

	query := "SELECT * FROM coupons where coupon=:coupon limit 1"

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
