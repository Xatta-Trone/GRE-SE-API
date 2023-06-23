package publicController

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		if listCountModel.WordCount != nil {
			wordCountMap[listCountModel.ListId] = *listCountModel.WordCount
		} else {
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

	// check if username is given
	if req.UserName != "" {
		// get the user
		user, err := ctl.userRepo.FindOneByUserName(req.UserName)
		if err == nil {
			req.UserId = user.ID
		}
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

	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No folder found."})
		return
	}

	userId, _ := utils.GetUserId(c)

	// get the data
	list, err := ctl.repository.FindOne(listId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check permissions and visibility
	if list.Visibility != enums.ListVisibilityPublic && userId != list.UserId {
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

	// now attach the user in the list meta
	user, _ := ctl.userRepo.FindOne(int(list.UserId))

	user2 := model.UserModel{Name: user.Name, UserName: user.UserName, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}

	list.User = &user2

	c.JSON(200, gin.H{
		"list_meta": list,
		"words":     words,
		"meta":      req,
	})

}

func (ctl *ListsController) FindWords(c *gin.Context) {

	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No list found."})
		return
	}

	userId, _ := utils.GetUserId(c)

	// get the data
	list, err := ctl.repository.FindOne(listId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check permissions and visibility
	if list.Visibility != enums.ListVisibilityPublic && userId != list.UserId {
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

	// word ids
	wordIdx := []string{}

	if req.WordIds != "" {

		for _, wordId := range strings.Split(req.WordIds, ",") {
			wordIdx = append(wordIdx, wordId)
		}

	}

	// check word id
	if len(wordIdx) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": fmt.Errorf("please provide word ids of this list")})
		return
	}

	// get the data
	words, err := ctl.wordRepo.FindWordsById(wordIdx, req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// now attach the user in the list meta
	user, _ := ctl.userRepo.FindOne(int(list.UserId))

	list.User = &user

	c.JSON(200, gin.H{
		"list_meta": list,
		"words":     words,
		"meta":      req,
	})

}

func (ctl *ListsController) Update(c *gin.Context) {

	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No folder found."})
		return
	}

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "No user found."})
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(listId)

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
	req.UserId = userId

	fmt.Println(req)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// check if name changed, then update the slug, otherwise keep the original one
	if req.Name == list.Name {
		// slug not changed, keep the original one
		req.Slug = list.Slug
	}

	// update the data
	ok, err := ctl.repository.Update(listId, req)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"updated": true,
	})

}

func (ctl *ListsController) Delete(c *gin.Context) {

	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No folder found."})
		return
	}

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(listId)

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

	// get the data
	list, err := ctl.repository.FindOne(listId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check permissions and visibility
	if userId == list.UserId {
		// user is the owner of the list
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
	} else {
		ok, err := ctl.repository.DeleteFromSavedList(userId, listId)

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

func (ctl *ListsController) DeleteWordInList(c *gin.Context) {

	// validate the given slug
	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No list found."})
		return
	}

	wordIdTemp := c.Query("word_id")

	if wordIdTemp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "missing param word id"})
		return
	}

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(listId)

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

func (ctl *ListsController) AddWordsInList(c *gin.Context) {

	// get the list id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No list found."})
		return
	}

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(listId)

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

	// request validation
	req, err := requests.ListWordsUpdateRequest(c)

	fmt.Println(req)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err})
		return
	}

	// check if its just word undo
	if req.WordId > 0 {
		err := ctl.listService.InsertListWordRelation(int64(req.WordId), int64(listId))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
			return
		}
		// success response
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	// check if it is words
	if req.Words != "" {
		// process the words
		words := ctl.listService.GetWordsFromListMetaRecord(req.Words)

		if len(words) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "No words found"})
			return
		}

		// send to the processor
		go ctl.listService.ProcessWordsOfSingleGroup(words, int64(listId))
		// success response
		c.JSON(http.StatusNoContent, gin.H{})
		return

	}

	c.JSON(http.StatusBadRequest, gin.H{
		"errors": "could not process the words",
	})

}

func (ctl *ListsController) FoldersInList(c *gin.Context) {

	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No folder found."})
		return
	}

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "No user found."})
		return
	}

	// get the data
	list, err := ctl.repository.FindOne(listId)

	if err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errors": "No list found."})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if list.UserId != userId {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"errors": "You are not the owner of this set."})
		return
	}

	folderIds := []uint64{}

	folders, err := ctl.repository.FoldersByListId(list.Id, userId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"folders": folderIds,
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	for _, v := range folders {
		folderIds = append(folderIds, v.FolderId)
	}

	c.JSON(http.StatusOK, gin.H{
		"folders": folderIds,
	})

}

func (ctl *ListsController) ToggleFolder(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	// get the folder id
	listId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// list id
	folderId := utils.ParseQueryToUint64(c, "folder_id")
	fmt.Println(folderId)

	ok, err := ctl.repository.ToggleFolder(folderId, listId)

	fmt.Println(ok, err)

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "List toggled.",
	})

}
