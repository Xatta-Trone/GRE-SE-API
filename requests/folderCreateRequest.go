package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xatta-trone/words-combinator/enums"
)

type FolderCreateRequestStruct struct {
	Name       string  `json:"name" form:"name" `
	Visibility int     `json:"visibility" form:"visibility,default=1" `
	UserId     uint64  `json:"user_id" form:"user_id"`
}

func (c FolderCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Visibility, validation.Required, validation.In(enums.FolderVisibilityMe, enums.FolderVisibilityPublic)),
	)
}

func FolderCreateRequest(c *gin.Context) (*FolderCreateRequestStruct, error) {
	var req *FolderCreateRequestStruct
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


