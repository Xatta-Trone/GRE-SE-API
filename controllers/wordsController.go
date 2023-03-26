package controllers

import (
	"fmt"
	"net/http"

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
