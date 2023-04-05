package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/xatta-trone/words-combinator/enums"
)

type ListsCreateRequestStruct struct {
	Name       string `json:"name" form:"name" `
	Url        *string `json:"url" form:"url" `
	Words      *string `json:"words" form:"words" `
	Visibility int    `json:"visibility" form:"visibility,default=1" `
	UserId     uint64 `json:"user_id" form:"user_id"`
}

func (c ListsCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Url, validation.When(c.Words == nil, validation.Required), is.URL),
		validation.Field(&c.Words, validation.When(c.Url == nil, validation.Required)),
		validation.Field(&c.Visibility, validation.Required, validation.In(enums.ListVisibilityMe, enums.ListVisibilityPublic)),
	)
}

func ListsCreateRequest(c *gin.Context) (*ListsCreateRequestStruct, error) {
	var req *ListsCreateRequestStruct
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
