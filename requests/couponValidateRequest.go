package requests

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CouponValidateRequestStruct struct {
	Coupon string `json:"coupon" form:"coupon" `
}

func (c CouponValidateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Coupon, validation.Required),
	)

}

func CouponValidateRequest(c *gin.Context) (CouponValidateRequestStruct, error) {
	var req CouponValidateRequestStruct
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
