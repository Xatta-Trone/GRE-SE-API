package publicController

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
)

type ListsController struct {
	repository  repository.ListRepositoryInterface
	listService services.ListProcessorServiceInterface
	wordRepo    repository.WordRepositoryInterface
}

func NewListsController(repository repository.ListRepositoryInterface, listService services.ListProcessorServiceInterface, wordRepo repository.WordRepositoryInterface) *ListsController {
	return &ListsController{
		repository:  repository,
		listService: listService,
		wordRepo:    wordRepo,
	}
}

func (ctl *ListsController) ListsByUserId(c *gin.Context) {
	// get all the lists by user id

	c.JSON(http.StatusOK, gin.H{
		"data": "user",
	})
}

func (ctl *ListsController) Index(c *gin.Context) {

	// get the user id
	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	fmt.Println(userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse the user id"})
		return
	}

	// validation request
	req, errs := requests.ListsIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// attach the user id to the request
	req.UserId = userId

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

func (ctl *ListsController) Create(c *gin.Context) {
	// get the user id
	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse the user id"})
		return
	}

	// request validation
	req, err := requests.ListsCreateRequest(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err})
		return
	}

	// set the user id
	req.UserId = userId

	// now create the record
	listMeta, err := ctl.repository.Create(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// now process the record
	go ctl.listService.ProcessListMetaRecord(listMeta)

	c.JSON(http.StatusCreated, gin.H{
		"data":    listMeta,
		"message": "Your list has been created. You will get a notification after processing the list shortly.",
	})
}

func (ctl *ListsController) FindOne(c *gin.Context) {

	// validate the given slug
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param slug"})
		return
	}

	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse the user id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOneBySlug(slug)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check permissions and visibility
	if userId != list.UserId && enums.ListVisibilityPublic != list.Visibility {
		c.JSON(http.StatusForbidden, gin.H{"errors": "The list either not public or deleted."})
		return
	}

	// all good now get the word data
	// validation request
	req, errs := requests.WordIndexByListIdRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}
	// set the list id
	req.ListId = list.Id

	// get the data
	words, err := ctl.wordRepo.FindAllByListId(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"list_meta": list,
		"words":     words,
		"meta":      req,
	})

}

func (ctl *ListsController) Update(c *gin.Context) {

	// validate the given slug
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param slug"})
		return
	}

	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse the user id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOneBySlug(slug)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check permissions and visibility
	if userId != list.UserId {
		c.JSON(http.StatusForbidden, gin.H{"errors": "Unauthorized."})
		return
	}

	// all good now get the word data
	// validation request
	req, errs := requests.ListsUpdateRequest(c)

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	if req.Name != list.Name {
		// update the data
		ok, err := ctl.repository.Update(list.Id, req)
		if err != nil || !ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	}

	c.JSON(http.StatusNoContent, gin.H{
		"updated": true,
	})

}

func (ctl *ListsController) Delete(c *gin.Context) {

	// validate the given slug
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param slug"})
		return
	}

	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse the user id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOneBySlug(slug)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check permissions and visibility
	if userId != list.UserId {
		c.JSON(http.StatusForbidden, gin.H{"errors": "Unauthorized."})
		return
	}

	// check if it belongs to list meta or not
	if list.ListMetaId != nil {
		ok, err := ctl.repository.DeleteFromListMeta(*list.ListMetaId)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
			return
		}

		if err != nil || !ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	} else {
		ok, err := ctl.repository.Delete(list.Id)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
			return
		}

		if err != nil || !ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	}

	c.JSON(http.StatusNoContent, gin.H{
		"deleted": true,
	})

}
