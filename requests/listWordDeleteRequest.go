package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ListWordDeleteRequestStruct struct {
	ListId     uint64 `json:"list_id" form:"list_id"`
	WordId     uint64 `json:"word_id" form:"word_id"`
}

func (c ListWordDeleteRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ListId, validation.Required),
		validation.Field(&c.WordId, validation.Required),
	)
}

func ListWordDeleteRequest(c *gin.Context) (*ListWordDeleteRequestStruct, error) {
	var req ListWordDeleteRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		return &req, err
	}

	err = req.Validate()

	if err != nil {
		return &req, err
	}

	return &req, nil

}
