package publicController

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type FolderController struct {
	repository     repository.FolderRepositoryInterface
	listRepository repository.ListRepositoryInterface
	userRepo       repository.UserRepositoryInterface
}

func NewFolderController(
	repository repository.FolderRepositoryInterface,
	listRepository repository.ListRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *FolderController {
	return &FolderController{
		repository:     repository,
		listRepository: listRepository,
		userRepo:       userRepo,
	}

}

func (ctl *FolderController) Index(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	req, errs := requests.FolderIndexRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId

	folders, err := ctl.repository.Index(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// make a temporary variable to copy the result then export via gin
	foldersToExport := make([]model.FolderModel, 0)
	userIds := []uint64{}
	folderIds := []uint64{}
	usersMap := make(map[uint64]model.UserModel)
	listCountMap := make(map[uint64]int)

	// check if len is zero
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
		folderIds = append(folderIds, folder.Id)
	}

	// get the users
	users, _ := ctl.userRepo.In(userIds, "id", "username")
	// map the users to user map to avoid second level iteration
	for _, user := range users {
		usersMap[user.ID] = user
	}

	// get the users
	listsCount, _ := ctl.repository.GetCount(folderIds)
	// map the users to user map to avoid second level iteration
	for _, listCountModel := range listsCount {
		if listCountModel.ListCount != nil {
			listCountMap[listCountModel.FolderId] = *listCountModel.ListCount
		} else {
			listCountMap[listCountModel.FolderId] = 0
		}

	}

	// now attach the users to the folders result
	for _, folder := range folders {
		user := usersMap[folder.UserId]
		listCount := listCountMap[folder.Id]
		f := model.FolderModel(folder)
		f.User = &user
		f.ListsCount = &listCount

		foldersToExport = append(foldersToExport, f)
	}

	c.JSON(200, gin.H{
		"data": foldersToExport,
		"meta": req,
	})
}

func (ctl *FolderController) PublicFolders(c *gin.Context) {

	req, errs := requests.PublicFolderIndexRequest(c)

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

	folders, err := ctl.repository.PublicIndex(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// make a temporary variable to copy the result then export via gin
	foldersToExport := make([]model.FolderModel, 0)
	userIds := []uint64{}
	usersMap := make(map[uint64]model.UserModel)

	// check if len is zero
	if len(folders) == 0 {
		// send empty response
		c.JSON(200, gin.H{
			"data": foldersToExport,
			"meta": req,
		})
		return
	}

	for _, list := range folders {
		userIds = append(userIds, list.UserId)
	}

	// get the users
	users, _ := ctl.userRepo.In(userIds, "id", "username")
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

	c.JSON(200, gin.H{
		"data": foldersToExport,
		"meta": req,
	})
}

func (ctl *FolderController) SavedFolders(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	req, errs := requests.SavedFolderIndexRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId

	folders, err := ctl.repository.SavedFolders(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// make a temporary variable to copy the result then export via gin
	foldersToExport := make([]model.FolderModel, 0)
	userIds := []uint64{}
	usersMap := make(map[uint64]model.UserModel)

	// check if len is zero
	if len(folders) == 0 {
		// send empty response
		c.JSON(200, gin.H{
			"data": foldersToExport,
			"meta": req,
		})
		return
	}

	for _, list := range folders {
		userIds = append(userIds, list.UserId)
	}

	// get the users
	users, _ := ctl.userRepo.In(userIds, "id", "username")
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

	c.JSON(200, gin.H{
		"data": foldersToExport,
		"meta": req,
	})
}

func (ctl *FolderController) Create(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	// request validation
	req, err := requests.FolderCreateRequest(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err})
		return
	}

	// set the user id
	req.UserId = userId

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

func (ctl *FolderController) SaveFolder(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	folderId, err := utils.ParseParamToUint64(c, "folder_id")

	if err != nil {
		return
	}

	// now create the record
	_, err = ctl.repository.SaveFolder(userId, folderId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Folder Saved.",
	})
}

func (ctl *FolderController) FindOne(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No user found."})
		return
	}

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No folder found."})
		return
	}

	fmt.Println(folderId)

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

	// check permissions and visibility
	if enums.FolderVisibilityPublic != folder.Visibility && userId != folder.UserId {
		c.JSON(http.StatusForbidden, gin.H{"errors": "The folder either not public or deleted."})
		return
	}

	// now get the lists associated with this folder
	// validation request
	req, errs := requests.FolderListIndexRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId
	req.FolderId = folderId

	lists, err := ctl.listRepository.ListsByFolderId(req)

	fmt.Println(lists)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// get the user data of this folder
	user, _ := ctl.userRepo.FindOne(int(folder.UserId))
	folder.User = &user

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
			"lists":  listsToExport,
			"folder": folder,
			"meta":   req,
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
	listsCount, _ := ctl.listRepository.GetCount(listIds)
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

	// get the total folder count

	c.JSON(200, gin.H{
		"lists":  listsToExport,
		"folder": folder,
		"meta":   req,
	})

}

func (ctl *FolderController) Update(c *gin.Context) {
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	fmt.Println(folderId)

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

	// check permissions and visibility
	if userId != folder.UserId {
		c.JSON(http.StatusForbidden, gin.H{"errors": "Unauthorized."})
		return
	}

	// all good now get the word data
	// validation request
	req, errs := requests.FolderUpdateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// check if name changed, then update the slug, otherwise keep the original one
	if req.Name == folder.Name {
		// slug not changed, keep the original one
		req.Slug = folder.Slug
	}

	// update the data
	ok, err := ctl.repository.Update(folderId, req)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"updated": true,
	})

}

func (ctl *FolderController) Delete(c *gin.Context) {
	// determine if the lists should be deleted or not
	var deleteLists bool = false
	delete := utils.ParseQueryString(c, "delete_lists")

	if delete == "1" {
		deleteLists = true
	}

	fmt.Println(delete, deleteLists)

	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	fmt.Println(folderId)

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

	// check permissions and visibility
	if userId != folder.UserId {
		c.JSON(http.StatusForbidden, gin.H{"errors": "Unauthorized."})
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

func (ctl *FolderController) DeleteSaveFolder(c *gin.Context) {
	// determine if the lists should be deleted or not
	var deleteLists bool = false
	delete := utils.ParseQueryString(c, "delete_lists")

	if delete == "1" {
		deleteLists = true
	}
	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	folderId, err := utils.ParseParamToUint64(c, "folder_id")

	if err != nil {
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

	// check permissions and visibility
	if userId == folder.UserId {
		// the user is the owner of the folder
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
	} else {
		// just remove from the relation
		ok, err := ctl.repository.DeleteSavedFolder(userId, folderId)

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

func (ctl *FolderController) ToggleList(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		return
	}

	fmt.Println(userId)

	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

	if err != nil {
		utils.Errorf(err)
		return
	}

	// list id
	listId := utils.ParseQueryToUint64(c, "list_id")
	fmt.Println(listId)

	ok, err := ctl.repository.ToggleList(folderId, listId)

	fmt.Println(ok, err)

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "List toggled.",
	})

}

func (ctl *FolderController) ListsInFolder(c *gin.Context) {

	// validate the given slug
	// get the folder id
	folderId, err := utils.ParseParamToUint64(c, "id")

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
	folder, err := ctl.repository.FindOne(folderId)

	if err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errors": "No result found."})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if folder.UserId != userId {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"errors": "You are not the owner of this folder."})
		return
	}

	listIds := []uint64{}

	lists, err := ctl.repository.ListIdsByFolderId(folder.Id, userId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"lists": listIds,
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	for _, v := range lists {
		listIds = append(listIds, v.ListId)
	}

	c.JSON(http.StatusOK, gin.H{
		"lists": listIds,
	})

}
