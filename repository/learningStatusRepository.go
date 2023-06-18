package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type LearningStatusInterface interface {
	Create(req requests.LearningStatusUpdateRequestStruct) (int64, error)
	Delete(listId, userId uint64) (int64, error)
	FindWordIdsByListId(listId uint64) ([]int64, error)
	FindLearningStatusByListId(listId, userId uint64) ([]model.LearningStatusModel, error)
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

func (rep *LearningStatusRepository) Delete(listId, userId uint64) (int64, error) {

	queryMap := map[string]interface{}{"user_id": userId, "list_id": listId}

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

func (rep *LearningStatusRepository) FindWordIdsByListId(listId uint64) ([]int64, error) {
	listIds := []model.ListWordModel{}
	queryMap := map[string]interface{}{"list_id": listId}

	query := "SELECT word_id FROM list_word_relation where list_id=:list_id order by word_id asc"

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return nil, err
	}
	err = nstmt.Select(&listIds, queryMap)

	if err != nil {
		utils.Errorf(err)
		return nil, err
	}

	// now iterate over the list
	idx := []int64{}

	for _, record := range listIds {
		idx = append(idx, int64(record.WordId))
	}

	return idx, nil

}

func (rep *LearningStatusRepository) FindLearningStatusByListId(listId, userId uint64) ([]model.LearningStatusModel, error) {
	learningStatuses := []model.LearningStatusModel{}
	queryMap := map[string]interface{}{"list_id": listId, "user_id": userId}

	query := fmt.Sprintf("SELECT word_id, learning_state FROM learning_status where list_id=:list_id and user_id=:user_id order by word_id asc")

	nstmt, err := rep.Db.PrepareNamed(query)

	if err != nil {
		utils.Errorf(err)
		return nil, err
	}
	err = nstmt.Select(&learningStatuses, queryMap)

	if err != nil {
		utils.Errorf(err)
		return nil, err
	}

	return learningStatuses, nil

}
