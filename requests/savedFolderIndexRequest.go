package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SavedFolderIndexReqStruct struct {
	ID      int    `form:"id,default=0" json:"id"`
	Query   string `form:"query" json:"query"`
	OrderBy string `form:"order_by,default=id" json:"order_by" `
	Order   string `form:"order,default=desc" json:"order" `
	Page    int    `form:"page,default=1" json:"page"`
	PerPage int    `form:"per_page,default=20" json:"per_page"`
	Total   int    `json:"total"`
	UserId  uint64 `json:"user_id,omitempty"`
}

func (c SavedFolderIndexReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OrderBy, validation.Required),
		validation.Field(&c.Order, validation.Required, validation.In("desc", "asc")),
		validation.Field(&c.Page, validation.Required),
		validation.Field(&c.PerPage, validation.Required),
	)
}

func SavedFolderIndexRequest(c *gin.Context) (*SavedFolderIndexReqStruct, error) {
	var req SavedFolderIndexReqStruct
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
