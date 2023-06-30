package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type PendingWordsUpdateRequestStruct struct {
	Word   string `json:"word" form:"word" `
	ListId int `json:"list_id" form:"list_id" `
}

func (c PendingWordsUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Word, validation.Required),
		validation.Field(&c.ListId, validation.Required),
	)
}

func PendingWordsUpdateRequest(c *gin.Context) (PendingWordsUpdateRequestStruct, error) {
	var req PendingWordsUpdateRequestStruct
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
