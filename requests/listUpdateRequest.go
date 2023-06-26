package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xatta-trone/words-combinator/enums"
)

type ListsUpdateRequestStruct struct {
	Name       string `json:"name" form:"name" `
	Slug       string `json:"slug" form:"slug" `
	Visibility int    `json:"visibility" form:"visibility,default=1" `
	UserId     uint64 `json:"user_id" form:"user_id"`
}

func (c ListsUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Visibility, validation.Required, validation.In(enums.ListVisibilityMe, enums.ListVisibilityPublic)),
	)
}

func ListsUpdateRequest(c *gin.Context) (*ListsUpdateRequestStruct, error) {
	var req ListsUpdateRequestStruct
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
