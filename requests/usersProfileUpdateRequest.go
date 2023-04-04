package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UsersProfileUpdateRequestStruct struct {
	UserName string `json:"username" form:"username" `
}

func (c UsersProfileUpdateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UserName, validation.Required),
	)
}

func UsersProfileUpdateRequest(c *gin.Context) (*UsersProfileUpdateRequestStruct, error) {
	var req UsersProfileUpdateRequestStruct
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
