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

type PendingWordsController struct {
	repository repository.PendingWordsInterface
}

func NewPendingWordsController(repository repository.PendingWordsInterface) *PendingWordsController {
	return &PendingWordsController{
		repository: repository,
	}

}

func (ctl *PendingWordsController) Index(c *gin.Context) {

	// validation request
	req, errs := requests.PendingWordIndexRequest(c)

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

func (ctl *PendingWordsController) Delete(c *gin.Context) {

	// validate the given id
	id := c.Param("list_id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	word := c.Query("word")

	// get the data
	ok, err := ctl.repository.Delete(idx, word)

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

func (ctl *PendingWordsController) Update(c *gin.Context) {

	// validation request
	req, errs := requests.PendingWordsUpdateRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}
	// check if coupon is empty then auto generate

	// check if coupon exists
	_, errs = ctl.repository.Update(req)

	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"id": 0})

}
