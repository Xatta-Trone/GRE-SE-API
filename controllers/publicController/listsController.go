package publicController

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
	"github.com/xatta-trone/words-combinator/utils"
)

type ListsController struct {
	repository  repository.ListRepositoryInterface
	listService services.ListProcessorServiceInterface
	wordRepo    repository.WordRepositoryInterface
	userRepo    repository.UserRepositoryInterface
}

func NewListsController(
	repository repository.ListRepositoryInterface,
	listService services.ListProcessorServiceInterface,
	wordRepo repository.WordRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *ListsController {
	return &ListsController{
		repository:  repository,
		listService: listService,
		wordRepo:    wordRepo,
		userRepo:    userRepo,
	}
}

func (ctl *ListsController) ListsByUserId(c *gin.Context) {
	// get all the lists by user id

	c.JSON(http.StatusOK, gin.H{
		"data": "user",
	})
}

func (ctl *ListsController) Index(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	// validation request
	req, errs := requests.ListsIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// attach the user id to the request
	req.UserId = userId

	fmt.Println(req)

	// get the data
	lists, err := ctl.repository.Index(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// make a temporary variable to copy the result then export via gin
	listsToExport := make([]model.ListModel, 0)
	userIds := []uint64{}
	listIds := []uint64{}
	usersMap := make(map[uint64]model.UserModel)
	wordCountMap := make(map[uint64]int)

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
		listIds = append(listIds, list.Id)
	}

	// get the users
	users, _ := ctl.userRepo.In(userIds, "id", "username")
	// map the users to user map to avoid second level iteration
	for _, user := range users {
		usersMap[user.ID] = user
	}

	// get the users
	listsCount, _ := ctl.repository.GetCount(listIds)
	// map the users to user map to avoid second level iteration
	for _, listCountModel := range listsCount {
		if (listCountModel.WordCount != nil) {
			wordCountMap[listCountModel.ListId] = *listCountModel.WordCount
		}else {
			wordCountMap[listCountModel.ListId] = 0
		}
		
	}

	// now attach the users to the folders result
	for _, list := range lists {
		user := usersMap[list.UserId]
		wordCount := wordCountMap[list.Id]
		f := model.ListModel(list)
		f.User = &user
		f.WordCount = &wordCount

		listsToExport = append(listsToExport, f)
	}

	c.JSON(200, gin.H{
		"data": listsToExport,
		"meta": req,
	})

}

func (ctl *ListsController) PublicLists(c *gin.Context) {

	// validation request
	req, errs := requests.PublicListsIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// get the data
	lists, err := ctl.repository.PublicIndex(req)

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
	users, _ := ctl.userRepo.In(userIds, "id", "username")
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

func (ctl *ListsController) SavedLists(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	// validation request
	req, errs := requests.SavedListsIndexRequest(c)

	// fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId

	// get the data
	lists, err := ctl.repository.SavedLists(req)

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
	users, _ := ctl.userRepo.In(userIds, "id", "username")
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

func (ctl *ListsController) Create(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
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

func (ctl *ListsController) SaveListItem(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}
	// request validation
	req, err := requests.SavedListsCreateRequest(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err})
		return
	}

	// set the user id
	req.UserId = userId

	// now create the record
	_, err = ctl.repository.SaveListItem(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Your list has been saved.",
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
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

func (ctl *ListsController) DeleteSavedList(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	listId, err := utils.ParseParamToUint64(c, "list_id")

	if err != nil {
		return
	}

	ok, err := ctl.repository.DeleteFromSavedList(userId,listId)

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

func (ctl *ListsController) DeleteWordInList(c *gin.Context) {

	// validate the given slug
	slug := c.Param("slug")
	wordIdTemp := c.Query("word_id")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param slug"})
		return
	}
	if wordIdTemp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param word id"})
		return
	}

	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "user id not found"})
		return
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
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
