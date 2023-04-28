package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type AdminFolderController struct {
	repository     repository.FolderRepositoryInterface
	listRepository repository.ListRepositoryInterface
	userRepository repository.UserRepositoryInterface
}

func NewAdminFolderController(repository repository.FolderRepositoryInterface, listRepository repository.ListRepositoryInterface, userRepository repository.UserRepositoryInterface) *AdminFolderController {
	return &AdminFolderController{
		repository:     repository,
		listRepository: listRepository,
		userRepository: userRepository,
	}

}

func (ctl *AdminFolderController) Index(c *gin.Context) {

	req, errs := requests.FolderIndexRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// get the folder records
	folders,count, err := ctl.repository.AdminIndex(req)

	req.Count = count.Count

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// make a temporary variable to copy the result then export via gin
	foldersToExport := make([]model.FolderModel, 0)
	userIds := []uint64{}
	usersMap := make(map[uint64]model.UserModel)

	// check if folders len is zero
	if len(folders) == 0 {
		// send empty response
		c.JSON(200, gin.H{
			"data": foldersToExport,
			"meta": req,
		})
		return
	}

	for _, folder := range folders {
		userIds = append(userIds, folder.UserId)
	}
	// get the users
	users, _ := ctl.userRepository.In(userIds)
	// map the users to user map to avoid second level iteration
	for _, user := range users {
		usersMap[user.ID] = user
	}

	// now attach the users to the folders result
	for _, folder := range folders {
		user := usersMap[folder.UserId]
		f := model.FolderModel(folder)
		f.User = &user

		foldersToExport = append(foldersToExport, f)

	}

	// fmt.Println(folders)

	c.JSON(200, gin.H{
		"data": foldersToExport,
		"meta": req,
	})
}

func (ctl *AdminFolderController) FindOne(c *gin.Context) {

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// get the data
	folder, err := ctl.repository.FindOne(folderId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// now get the lists associated with this folder
	// validation request
	req, errs := requests.FolderListIndexRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = folder.UserId
	req.FolderId = folder.Id

	lists, err := ctl.listRepository.ListsByFolderId(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"lists":  lists,
		"folder": folder,
		"meta":   req,
	})

}

func (ctl *AdminFolderController) Create(c *gin.Context) {

	// request validation
	req, err := requests.FolderCreateRequest(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err})
		return
	}

	// now create the record
	folder, err := ctl.repository.Create(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    folder,
		"message": "Folder created.",
	})
}

func (ctl *AdminFolderController) Update(c *gin.Context) {

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// get the data
	folder, err := ctl.repository.FindOne(folderId)

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
	req, errs := requests.FolderUpdateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	ok, err := ctl.repository.Update(folder.Id, req)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// send the new data 
	newData,_ := ctl.repository.FindOne(folderId)

	c.JSON(http.StatusOK, gin.H{
		"data": newData,
	})

}

func (ctl *AdminFolderController) Delete(c *gin.Context) {
	// determine if the lists should be deleted or not
	var deleteLists bool = false
	delete := utils.ParseQueryString(c, "delete_lists")

	if delete == "1" {
		deleteLists = true
	}

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// get the data
	folder, err := ctl.repository.FindOne(folderId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check if it belongs to list meta or not
	if folder.ListMetaId != nil {
		ok, err := ctl.listRepository.DeleteFromListMeta(*folder.ListMetaId)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
			return
		}

		if err != nil || !ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	} else {
		ok, err := ctl.repository.Delete(folderId, deleteLists)

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

func (ctl *AdminFolderController) ToggleList(c *gin.Context) {


	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// list id
	listId := utils.ParseQueryToUint64(c, "list_id")

	ok, err := ctl.repository.ToggleList(folderId, listId)

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List toggled.",
	})

}
