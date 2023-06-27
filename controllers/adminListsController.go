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
	"github.com/xatta-trone/words-combinator/services"
	"github.com/xatta-trone/words-combinator/utils"
)

type AdminListsController struct {
	repository  repository.ListRepositoryInterface
	listService services.ListProcessorServiceInterface
	wordRepo    repository.WordRepositoryInterface
	userRepo    repository.UserRepositoryInterface
}

func NewListsController(repository repository.ListRepositoryInterface, listService services.ListProcessorServiceInterface, wordRepo repository.WordRepositoryInterface, userRepo repository.UserRepositoryInterface) *AdminListsController {
	return &AdminListsController{
		repository:  repository,
		listService: listService,
		wordRepo:    wordRepo,
		userRepo:    userRepo,
	}
}

func (ctl *AdminListsController) Index(c *gin.Context) {

	fmt.Println("param", c.Param("scope"), c.Query("scope"))

	// validation request
	req, errs := requests.AdminListsIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// get the data
	lists, err := ctl.repository.AdminIndex(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// make a temporary variable to copy the result then export via gin
	listsToExport := make([]model.ListModel, 0)
	userIds := []uint64{}
	usersMap := make(map[uint64]model.UserModel)

	// check if len is zero
	if len(lists) == 0 {
		// send empty response
		c.JSON(200, gin.H{
			"data": listsToExport,
			"meta": req,
		})
		return
	}

	for _, list := range lists {
		userIds = append(userIds, list.UserId)
	}

	// get the users
	users, _ := ctl.userRepo.In(userIds)
	// map the users to user map to avoid second level iteration
	for _, user := range users {
		usersMap[user.ID] = user
	}

	// now attach the users to the folders result
	for _, list := range lists {
		user := usersMap[list.UserId]
		f := model.ListModel(list)
		f.User = &user

		listsToExport = append(listsToExport, f)

	}

	c.JSON(200, gin.H{
		"data": listsToExport,
		"meta": req,
	})

}

func (ctl *AdminListsController) Create(c *gin.Context) {
	// request validation
	req, err := requests.ListsCreateRequest(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err})
		return
	}

	fmt.Println(req)

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

func (ctl *AdminListsController) FindOne(c *gin.Context) {

	// get the folder id
	Id, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(Id)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
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
		"list":  list,
		"words": words,
		"meta":  req,
	})

}

func (ctl *AdminListsController) Update(c *gin.Context) {

	// get the list id
	Id, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(Id)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
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

	ok, err := ctl.repository.Update(list.Id, req)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	updatedData, _ := ctl.repository.FindOne(Id)

	c.JSON(http.StatusOK, gin.H{
		"data": updatedData,
	})

}

func (ctl *AdminListsController) Delete(c *gin.Context) {

	// get the list id
	Id, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(Id)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
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

func (ctl *AdminListsController) DeleteWordInList(c *gin.Context) {

	// get the list id
	Id, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	wordIdTemp := c.Query("word_id")
	if wordIdTemp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param word id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(Id)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	wordId, _ := strconv.ParseUint(wordIdTemp, 10, 64)

	// delete the record
	ok, err := ctl.repository.DeleteWordInList(wordId, list.Id)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"deleted": true,
	})

}
