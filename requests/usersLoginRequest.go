package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UsersLoginRequestStruct struct {
	Email string `json:"email" form:"email" `
	Token string `json:"token" form:"token" `
}

func (c UsersLoginRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Token, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}

func UsersLoginRequest(c *gin.Context) (*UsersLoginRequestStruct, error) {
	var req *UsersLoginRequestStruct
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
