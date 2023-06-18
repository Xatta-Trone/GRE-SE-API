package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/gin-gonic/gin"
)

type WordIndexByListIdReqStruct struct {
	ID      int    `form:"id,default=0" json:"id"`
	Query   string `form:"query" json:"query"`
	OrderBy string `form:"order_by,default=desc" json:"order_by" `
	Page    int    `form:"page,default=1" json:"page"`
	PerPage int    `form:"per_page,default=50" json:"per_page"`
	Total   int    `json:"total"`
	ListId  uint64 `json:"list_id"`
	WordIds string `form:"word_ids" json:"word_ids,omitempty"`
}

func (c WordIndexByListIdReqStruct) Validate() error {
	return validation.ValidateStruct(&c,
		// validation.Field(&c.ID, validation.Required),
		validation.Field(&c.OrderBy, validation.Required, validation.In("desc", "asc")),
		validation.Field(&c.Page, validation.Required),
		validation.Field(&c.PerPage, validation.Required),
	)
}

func WordIndexByListIdRequest(c *gin.Context) (*WordIndexByListIdReqStruct, error) {
	var req WordIndexByListIdReqStruct
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
