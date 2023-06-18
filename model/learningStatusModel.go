package model

type LearningStatusModel struct {
	WordId        int64  `db:"word_id" json:"word_id,omitempty"`
	UserId        uint64 `db:"user_id" json:"user_id,omitempty"`
	ListId        uint64 `db:"list_id" json:"list_id,omitempty"`
	LearningState int    `db:"learning_state" json:"learning_state,omitempty"`
}
