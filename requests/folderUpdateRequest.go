package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xatta-trone/words-combinator/enums"
)

type FolderUpdateRequestStruct struct {
	Name       string  `json:"name" form:"name" `
	Visibility int     `json:"visibility" form:"visibility,default=1" `
	UserId     uint64  `json:"user_id" form:"user_id"`
}

func (c FolderUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Visibility, validation.Required, validation.In(enums.FolderVisibilityMe, enums.FolderVisibilityPublic)),
	)
}

func FolderUpdateRequest(c *gin.Context) (*FolderUpdateRequestStruct, error) {
	var req *FolderUpdateRequestStruct
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


