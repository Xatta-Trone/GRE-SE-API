package requests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

type WordIndexReqStruct struct {
	ID    int    `form:"id" json:"id" validate:"required" `
	Query string `form:"query" json:"query" validate:"-"`
}

func WordsIndexRequest(c *gin.Context) (WordIndexReqStruct, validate.Errors) {
	var req WordIndexReqStruct
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
