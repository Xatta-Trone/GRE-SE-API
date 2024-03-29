package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xatta-trone/words-combinator/enums"
)

type ListsIndexReqStruct struct {
	ID        int    `form:"id,default=0" json:"id"`
	Query     string `form:"query" json:"query"`
	OrderBy   string `form:"order_by,default=id" json:"order_by" `
	Order     string `form:"order,default=desc" json:"order" `
	SaveOrder string `form:"save_order,default=desc" json:"save_order" `
	Page      int    `form:"page,default=1" json:"page"`
	PerPage   int    `form:"per_page,default=20" json:"per_page"`
	UserId    uint64 `json:"user_id"`
	Count     int64  `form:"count" json:"count"`
	Filter    string `form:"filter,default=all" json:"filter"`
}

func (c ListsIndexReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		// validation.Field(&c.ID, validation.Required),
		validation.Field(&c.OrderBy, validation.Required),
		validation.Field(&c.Order, validation.Required, validation.In("desc", "asc")),
		validation.Field(&c.SaveOrder, validation.Required, validation.In("desc", "asc")),
		validation.Field(&c.Filter, validation.Required, validation.In(enums.ListFilterAll, enums.ListFilterSaved, enums.ListFilterCrated)),
		validation.Field(&c.Page, validation.Required),
		validation.Field(&c.PerPage, validation.Required),
	)
}

func ListsIndexRequest(c *gin.Context) (*ListsIndexReqStruct, error) {
	var req ListsIndexReqStruct
	err := c.ShouldBind(&req)
	if err != nil {
		return &req, err
	}

	err = req.Validate()

	// set the order dir
	// if req.Order == "1" {
	// 	req.OrderDir = "asc"
	// }

	if err != nil {
		return &req, err
	}

	return &req, nil

}
