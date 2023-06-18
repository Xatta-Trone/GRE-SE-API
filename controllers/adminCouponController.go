package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
)

type AdminCouponController struct {
	repository repository.CouponInterface
}

func NewAdminCouponController(repository repository.CouponInterface) *AdminCouponController {
	return &AdminCouponController{
		repository: repository,
	}

}

func (ctl *AdminCouponController) Index(c *gin.Context) {

	// validation request
	req, errs := requests.CouponIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// // get the data
	coupons, err := ctl.repository.Index(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": coupons,
		"meta": req,
	})

}

func (ctl *AdminCouponController) Create(c *gin.Context) {

	// validation request
	req, errs := requests.CouponCreateRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}
	// check if coupon is empty then auto generate

	// check if coupon exists
	_, errs = ctl.repository.Create(req, req.Coupon)
	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": 0})

}

func (ctl *AdminCouponController) Delete(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	ok, err := ctl.repository.Delete(idx)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"deleted": true})

}

// func generateToken(n int) string {

// 	if n == 0 {
// 		n = 8
// 	}

// 	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }
