package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xatta-trone/words-combinator/model"
)

type WordUpdateReqStruct struct {
	IsReviewed int                 `json:"is_reviewed" form:"is_reviewed"`
	WordData   model.WordDataModel `json:"word_data" form:"word_data"`
}

func (c WordUpdateReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		// validation.Field(&c.ID, validation.Required),
		validation.Field(&c.IsReviewed, validation.Required),
		validation.Field(&c.WordData, validation.Required),
	)
}

func WordUpdateRequest(c *gin.Context) (*WordUpdateReqStruct, error) {
	var req WordUpdateReqStruct
	err := c.ShouldBind(&req)

	if err != nil {
		return &req, err
	}

	return &req, nil

}
