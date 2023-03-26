package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
)

func WordIndex(c *gin.Context) {

	// validation request
	req, err := requests.WordsIndexRequest(c)

	if err != nil {
		// validation request
		c.JSON(http.StatusUnprocessableEntity, err.All())
		return
	}

	// get the data
	wordData := services.WordsIndex(req)

	c.JSON(200, gin.H{
		"data": wordData,
		// "errors": err.Error(),
	})
}
