package requests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"github.com/xatta-trone/words-combinator/model"
)

type WordUpdateReqStruct struct {
	IsReviewed int           `json:"is_reviewed" validate:"required|integer"`
	WordData   []model.Combined `json:"word_data" validate:"required"`
}

func WordUpdateRequest(c *gin.Context) (*WordUpdateReqStruct, validate.Errors) {
	var req WordUpdateReqStruct
	err := c.ShouldBindJSON(&req)

	if err != nil {
		return &req, nil
	}

	v := validate.Struct(req)

	fmt.Println(v.Validate())

	if !v.Validate() {
		return &req, v.Errors
	}

	return &req, nil

}
