package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UsersCreateRequestStruct struct {
	Email string `json:"email" form:"email" `
	Token string `json:"token" form:"token" `
	Name string `json:"name" form:"name" `
}

func (c UsersCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Token, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Name, validation.Required),
	)
}

func UsersCreateRequest(c *gin.Context) (*UsersCreateRequestStruct, error) {
	var req *UsersCreateRequestStruct
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
