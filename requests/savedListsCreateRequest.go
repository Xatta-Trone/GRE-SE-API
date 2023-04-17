package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SavedListsCreateRequestStruct struct {
	ListId uint64 `json:"list_id" form:"list_id" `
	UserId uint64 `json:"user_id" form:"user_id"`
	Scope  string `json:"scope" form:"scope,default=user"` // either admin or user; userId required when scope is admin
}

func (c SavedListsCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ListId, validation.Required),
		validation.Field(&c.Scope, validation.Required),
		validation.Field(&c.UserId, validation.When(c.Scope == "admin", validation.Required)),
	)
}

func SavedListsCreateRequest(c *gin.Context) (*SavedListsCreateRequestStruct, error) {
	var req *SavedListsCreateRequestStruct
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


