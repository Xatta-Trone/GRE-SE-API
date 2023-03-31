package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
)

type WordController struct {
	wordRepository repository.WordRepositoryInterface
}

func NewWordController(wordRepository repository.WordRepositoryInterface) *WordController {
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

func (ctl *WordController) WordSave(c *gin.Context) {

	// validation request
	req, errs := requests.WordCreateRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.All()})
		return
	}

	wordData := model.WordDataModel{
		Word:            req.Word,
		PartsOfSpeeches: req.WordData,
	}

	// get the data
	word, err := ctl.wordRepository.Create(req.Word, wordData, req.IsReviewed)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": word,
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

func (ctl *WordController) UpdateById(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	fmt.Println(idx)

	// check if record exists

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

	if word.Id != int64(idx) {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Id mismatch"})
		return
	}

	// validation request
	req, errs := requests.WordUpdateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.All()})
		return
	}

	fmt.Println(req.IsReviewed)
	fmt.Println(req.WordData)

	wordData := model.WordDataModel{
		Word:            word.Word,
		PartsOfSpeeches: req.WordData,
	}

	// get the data
	ok, err := ctl.wordRepository.UpdateById(idx, wordData, req.IsReviewed)

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
