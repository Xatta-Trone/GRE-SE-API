package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FolderIndexReqStruct struct {
	ID       int    `form:"id,default=0" json:"id"`
	Query    string `form:"query" json:"query"`
	OrderBy  string `form:"order_by,default=id" json:"order_by" `
	Order    string `form:"order,default=0" json:"order" `
	OrderDir string `form:"order_dir,default=desc" json:"order_dir" `
	Page     int    `form:"page,default=1" json:"page"`
	PerPage  int    `form:"per_page,default=20" json:"per_page"`
	Count    int64  `form:"count" json:"count"`
	UserId   uint64 `json:"user_id,omitempty"`
}

func (c FolderIndexReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		// validation.Field(&c.ID, validation.Required),
		validation.Field(&c.OrderBy, validation.Required),
		validation.Field(&c.OrderDir, validation.Required),
		validation.Field(&c.Page, validation.Required),
		validation.Field(&c.PerPage, validation.Required),
	)
}

func FolderIndexRequest(c *gin.Context) (*FolderIndexReqStruct, error) {
	var req FolderIndexReqStruct
	err := c.ShouldBind(&req)
	if err != nil {
		return &req, err
	}

	err = req.Validate()

	// set the order dir 
	if req.Order == "1" {
		req.OrderDir = "asc"
	}

	if err != nil {
		return &req, err
	}

	return &req, nil

}
