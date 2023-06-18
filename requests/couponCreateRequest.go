package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CouponCreateRequestStruct struct {
	Coupon string `json:"coupon" form:"coupon" `
	MaxUse int    `json:"max_use" form:"max_use" `
	Months int    `json:"months" form:"months" `
}

func (c CouponCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Coupon, validation.Required),
		validation.Field(&c.MaxUse, validation.Required),
	)

	return nil
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
