package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type NotificationIndexReqStruct struct {
	ID       int    `form:"id,default=0" json:"id"`
	Query    string `form:"query" json:"query"`
	OrderBy  string `form:"order_by,default=id" json:"order_by" `
	OrderDir string `form:"order_dir,default=desc" json:"order_dir" `
	Page     int    `form:"page,default=1" json:"page"`
	PerPage  int    `form:"per_page,default=20" json:"per_page"`
	Filter   string `form:"filter,default=all" json:"filter"`
	UserId   uint64 `json:"user_id,omitempty"`
	Count    int64  `form:"count" json:"count"`
}

func (c NotificationIndexReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OrderBy, validation.Required),
		validation.Field(&c.OrderDir, validation.Required),
		validation.Field(&c.Page, validation.Required),
		validation.Field(&c.PerPage, validation.Required),
	)
}

func NotificationIndexRequest(c *gin.Context) (*NotificationIndexReqStruct, error) {
	var req NotificationIndexReqStruct
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
