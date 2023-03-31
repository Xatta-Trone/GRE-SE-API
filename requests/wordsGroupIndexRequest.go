package requests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

type WordsGroupIndexRequestStruct struct {
	ID      int    `form:"id,default=0" json:"id" validate:"integer" `
	Query   string `form:"query" json:"query" validate:"-"`
	OrderBy string `form:"order_by,default=desc" json:"order_by" validate:"in:asc,desc" `
	Page    int    `form:"page,default=1" json:"page" validate:"required|integer" `
	PerPage int    `form:"per_page,default=20" json:"per_page" validate:"required|integer" `
}

func WordsGroupIndexRequest(c *gin.Context) (WordsGroupIndexRequestStruct, validate.Errors) {
	var req WordsGroupIndexRequestStruct
	err := c.ShouldBindQuery(&req)

	if err != nil {
		return req, nil
	}

	v := validate.Struct(req)

	fmt.Println(v.Validate())

	if !v.Validate() {
		return req, v.Errors
	}

	return req, nil

}
