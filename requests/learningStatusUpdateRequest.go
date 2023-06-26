package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xatta-trone/words-combinator/enums"
)

type LearningStatusUpdateRequestStruct struct {
	ListId        uint64 `json:"list_id" form:"list_id" `
	WordId        uint64 `json:"word_id" form:"word_id"`
	LearningState int    `json:"learning_state" form:"learning_state,default=0"`
	UserId        uint64 `json:"user_id,omitempty" form:"user_id"`
}

func (c LearningStatusUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ListId, validation.Required),
		validation.Field(&c.WordId, validation.Required),
		validation.Field(&c.LearningState, validation.Required, validation.In(enums.LearningStatusDoNotKnow, enums.LearningStatusLearning, enums.LearningStatusMastered)),
	)
}

func LearningStatusUpdateRequest(c *gin.Context) (LearningStatusUpdateRequestStruct, error) {
	var req LearningStatusUpdateRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		return req, err
	}

	err = req.Validate()

	if err != nil {
		return req, err
	}

	return req, nil

}
