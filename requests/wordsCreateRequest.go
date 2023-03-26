package requests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"github.com/xatta-trone/words-combinator/model"
)

type WordCreateReqStruct struct {
	Word       string           `json:"word" validate:"required|string"`
	WordData   []model.Combined `json:"word_data" validate:"required"`
	IsReviewed int              `form:"is_reviewed,default=0" json:"is_reviewed" validate:"required|integer"`
}

func WordCreateRequest(c *gin.Context) (*WordCreateReqStruct, validate.Errors) {
	var req WordCreateReqStruct
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
