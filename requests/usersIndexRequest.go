package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UsersIndexReqStruct struct {
	ID      int    `form:"id,default=0" json:"id"`
	Query   string `form:"query" json:"query"`
	OrderBy string `form:"order_by,default=desc" json:"order_by" `
	Page    int    `form:"page,default=1" json:"page"`
	PerPage int    `form:"per_page,default=20" json:"per_page"`
}

func (c UsersIndexReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ID, validation.Required),
		validation.Field(&c.OrderBy, validation.Required, validation.In("desc", "asc")),
		validation.Field(&c.Page, validation.Required),
		validation.Field(&c.PerPage, validation.Required),
	)
}

func UsersIndexRequest(c *gin.Context) (UsersIndexReqStruct, error) {
	var req UsersIndexReqStruct
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
