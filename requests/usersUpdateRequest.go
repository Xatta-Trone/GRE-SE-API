package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UsersUpdateRequestStruct struct {
	Name     string `json:"name" form:"name" `
	Email    string `json:"email" form:"email" `
	UserName string `json:"username" form:"username" `
}

func (c UsersUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.UserName, validation.Required),
	)
}

func UsersUpdateRequest(c *gin.Context) (*UsersUpdateRequestStruct, error) {
	var req UsersUpdateRequestStruct
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
