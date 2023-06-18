package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type LearningStatusInterface interface {
	Create(req requests.LearningStatusUpdateRequestStruct) (int64, error)
	Delete(req requests.LearningStatusDeleteRequestStruct) (int64, error)
}

type LearningStatusRepository struct {
	Db *sqlx.DB
}

func NewLearningStatusRepository(db *sqlx.DB) *LearningStatusRepository {
	return &LearningStatusRepository{Db: db}
}

func (rep *LearningStatusRepository) Create(req requests.LearningStatusUpdateRequestStruct) (int64, error) {

	queryMap := map[string]interface{}{"user_id": req.UserId, "list_id": req.ListId, "word_id": req.WordId, "learning_state": req.LearningState}

	res, err := rep.Db.NamedExec("Insert into  learning_status(user_id,list_id,word_id,learning_state) values(:user_id,:list_id,:word_id,:learning_state) ON DUPLICATE KEY UPDATE learning_state=:learning_state", queryMap)

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	RowsAffected, err := res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}


	return RowsAffected, nil

}

func (rep *LearningStatusRepository) Delete(req requests.LearningStatusDeleteRequestStruct) (int64, error) {

	queryMap := map[string]interface{}{"user_id": req.UserId, "list_id": req.ListId,}

	res, err := rep.Db.NamedExec("DELETE FROM `learning_status` WHERE user_id=:user_id and list_id=:list_id", queryMap)

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}

	RowsAffected, err := res.RowsAffected()

	if err != nil {
		utils.Errorf(err)
		return -1, err
	}


	return RowsAffected, nil

}
