package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CouponCreateRequestStruct struct {
	Coupon  string `json:"coupon" form:"coupon" `
	Type    string `json:"type" form:"type,default=one_time" `
	MaxUse  int    `json:"max_use" form:"max_use,default=0" `
	Months  int    `json:"months" form:"months" `
	Expires string `json:"expires" form:"expires" `
}

func (c CouponCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Coupon, validation.Required),
		validation.Field(&c.Type, validation.Required),
		// validation.Field(&c.MaxUse, validation.Required),
		validation.Field(&c.Months, validation.Required),
	)
}

func CouponCreateRequest(c *gin.Context) (CouponCreateRequestStruct, error) {
	var req CouponCreateRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		return req, err
	}

	err = req.Validate()

	if err != nil {
		return req, err
	}

	return req, nil

}
