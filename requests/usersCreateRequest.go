package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UsersCreateRequestStruct struct {
	Name  string `json:"name" form:"name" `
	Email string `json:"email" form:"email" `
}

func (c UsersCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
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
