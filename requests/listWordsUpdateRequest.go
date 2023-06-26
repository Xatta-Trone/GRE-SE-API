package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ListWordsUpdateRequestStruct struct {
	Words   string `json:"words" form:"words" `
	UserId  uint64 `json:"user_id" form:"user_id"`
	WordId uint64 `json:"word_id" form:"word_id"`
}

func (c ListWordsUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Words, validation.When(c.WordId == 0, validation.Required)),
		validation.Field(&c.WordId, validation.When(c.Words == "", validation.Required)),
	)
}

func ListWordsUpdateRequest(c *gin.Context) (ListWordsUpdateRequestStruct, error) {
	var req ListWordsUpdateRequestStruct
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
