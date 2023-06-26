package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LearningStatusDeleteRequestStruct struct {
	ListId uint64 `json:"list_id" form:"list_id" `
	UserId uint64 `json:"user_id,omitempty" form:"user_id"`
}

func (c LearningStatusDeleteRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ListId, validation.Required),
	)
}

func LearningStatusDeleteRequest(c *gin.Context) (LearningStatusDeleteRequestStruct, error) {
	var req LearningStatusDeleteRequestStruct
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
