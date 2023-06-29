package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type NotificationInterface interface {
	Index(req *requests.NotificationIndexReqStruct) ([]model.NotificationModel, error)
	Update(userId uint64)
}

type NotificationRepository struct {
	Db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{Db: db}
}

func (rep *NotificationRepository) Index(r *requests.NotificationIndexReqStruct) ([]model.NotificationModel, error) {

	models := []model.NotificationModel{}
	count := model.CountModel{}

	queryMap := map[string]interface{}{"query": "%" + r.Query + "%", "user_id": r.UserId, "orderby": r.OrderBy, "limit": r.PerPage, "offset": (r.Page - 1) * r.PerPage}

	order := r.OrderDir // problem with order by https://github.com/jmoiron/sqlx/issues/153

	searchString := "FROM notifications where user_id=:user_id"

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
	queryCount := fmt.Sprintf("SELECT count(id) as count %s and read_at is NULL limit 1", searchString)
	nstmt1, _ := rep.Db.PrepareNamed(queryCount)
	_ = nstmt1.Get(&count, queryMap)

	r.Count = count.Count

	return models, nil

}

func (rep *NotificationRepository) Update(userId uint64) {

	queryMap := map[string]interface{}{"user_id": userId, "read_at": time.Now().UTC()}

	rep.Db.NamedExec("Update notifications set read_at=:read_at where user_id=:user_id", queryMap)

}
