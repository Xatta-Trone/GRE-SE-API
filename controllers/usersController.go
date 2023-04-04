package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
)

type UsersController struct {
	repository repository.UserRepositoryInterface
}

func NewUsersController(userRepo repository.UserRepositoryInterface) *UsersController {
	return &UsersController{repository: userRepo}
}

func (ctl *UsersController) Index(c *gin.Context) {

	// validation request
	req, errs := requests.UsersIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// get the data
	word, err := ctl.repository.Index(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": word,
		"meta": req,
	})

}

func (ctl *UsersController) FindOne(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	word, err := ctl.repository.FindOne(idx)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": word,
	})

}

func (ctl *UsersController) Delete(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	word, err := ctl.repository.Delete(idx)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": word,
	})

}

func (ctl *UsersController) Create(c *gin.Context) {

	// validation request
	req, errs := requests.UsersCreateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// save the record
	user, err := ctl.repository.Create(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

func (ctl *UsersController) Update(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// check if record exists

	// get the data
	model, err := ctl.repository.FindOne(idx)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if int64(model.ID) != int64(idx) {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Id mismatch"})
		return
	}

	// validation request
	req, errs := requests.UsersUpdateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.Error()})
		return
	}


	// get the data
	ok, err := ctl.repository.Update(idx, req)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"updated": true})

}
