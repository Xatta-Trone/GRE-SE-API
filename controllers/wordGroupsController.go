package controllers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
)

type WordGroupController struct {
	WgRepo    repository.WordGroupInterface
	wgService services.WordGroupServiceInterface
}

func NewWordGroupController(wordGroupRepo repository.WordGroupInterface, wgService services.WordGroupServiceInterface) *WordGroupController {
	return &WordGroupController{
		WgRepo:    wordGroupRepo,
		wgService: wgService,
	}
}

func (ctl *WordGroupController) Index(c *gin.Context) {

	// validation request
	req, errs := requests.WordsGroupIndexRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.All()})
		return
	}

	// get the data
	word, err := ctl.WgRepo.FindAll(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": word,
		"meta": req,
	})

}

func (ctl *WordGroupController) Import(c *gin.Context) {

	// validation request
	req, errs := requests.WordGroupCreateRequest(c)

	fmt.Println(req, errs)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// upload the file if present
	filePath := ""

	if req.File != nil {
		file, header, err := c.Request.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		// upload the file
		fileExt := filepath.Ext(header.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
		now := time.Now()
		filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePath = "uploads/" + filename

		out, err := os.Create("uploads/" + filename)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		// file uploaded

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		}

		req.FileName = filePath

	}

	// save the record
	wordGroup, err := ctl.WgRepo.Create(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// fire the word process service
	go ctl.wgService.ProcessWordGroupData(wordGroup)

	c.JSON(http.StatusCreated, gin.H{
		"data": wordGroup,
	})
}

func (ctl *WordGroupController) FindOne(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	word, err := ctl.WgRepo.FindOne(idx)

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

func (ctl *WordGroupController) DeleteById(c *gin.Context) {

	// validate the given id
	id := c.Param("id")

	idx, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The given id is not valid integer"})
		return
	}

	// get the data
	ok, err := ctl.WgRepo.DeleteOne(idx)

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
