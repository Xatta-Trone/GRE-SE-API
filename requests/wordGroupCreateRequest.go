package requests

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type WordGroupCreateReqStruct struct {
	Name     string                `json:"name" form:"name" `
	File     *multipart.FileHeader `json:"file" form:"file" `
	Words    *string               `json:"words" form:"words"`
	FileName string                `json:"file_name" db:"file_name"`
}

func (c WordGroupCreateReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Words, validation.When(c.File == nil, validation.Required)),
		validation.Field(&c.File, validation.When(c.Words == nil, validation.Required)),
	)
}

func WordGroupCreateRequest(c *gin.Context) (*WordGroupCreateReqStruct, error) {
	var req WordGroupCreateReqStruct
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
