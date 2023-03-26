package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/requests"
)

type WordController struct {
	wordRepository model.WordRepository
}

func NewWordController(wordRepository model.WordRepository) *WordController {
	return &WordController{
		wordRepository: wordRepository,
	}

}

func (ctl *WordController) WordIndex(c *gin.Context) {

	// validation request
	req, errs := requests.WordsIndexRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.All()})
		return
	}

	// get the data
	word, err := ctl.wordRepository.FindAll(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": word,
		"meta": req,
	})

}

func (ctl *WordController) WordById(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	word, err := ctl.wordRepository.FindOne(idx)

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

func (ctl *WordController) DeleteById(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	ok, err := ctl.wordRepository.DeleteOne(idx)

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
