package publicController

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type FolderController struct {
	repository     repository.FolderRepositoryInterface
	listRepository repository.ListRepositoryInterface
}

func NewFolderController(repository repository.FolderRepositoryInterface, listRepository repository.ListRepositoryInterface) *FolderController {
	return &FolderController{
		repository:     repository,
		listRepository: listRepository,
	}

}

func (ctl *FolderController) Index(c *gin.Context) {
	userId := utils.GetUserId(c)

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

	c.JSON(200, gin.H{
		"data": folders,
		"meta": req,
	})
}

func (ctl *FolderController) Create(c *gin.Context) {
	userId := utils.GetUserId(c)

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

func (ctl *FolderController) FindOne(c *gin.Context) {

	userId := utils.GetUserId(c)

	fmt.Println(userId)

	// get the folder id
	folderId := utils.ParseParamToUint64(c, "id")

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
	if userId != folder.UserId && enums.FolderVisibilityPublic != folder.Visibility {
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

	c.JSON(200, gin.H{
		"lists":  lists,
		"folder": folder,
		"meta":   req,
	})

}

func (ctl *FolderController) Update(c *gin.Context) {
	userId := utils.GetUserId(c)

	fmt.Println(userId)

	// get the folder id
	folderId := utils.ParseParamToUint64(c, "id")

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

	if req.Name != folder.Name {
		// update the data
		ok, err := ctl.repository.Update(folder.Id, req)
		if err != nil || !ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
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

	userId := utils.GetUserId(c)

	fmt.Println(userId)

	// get the folder id
	folderId := utils.ParseParamToUint64(c, "id")

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

func (ctl *FolderController) ToggleList(c *gin.Context) {

	userId := utils.GetUserId(c)

	fmt.Println(userId)

	// get the folder id
	folderId := utils.ParseParamToUint64(c, "id")
	fmt.Println(folderId)

	// list id
	listId := utils.ParseQueryToUint64(c, "list_id")
	fmt.Println(listId)

	ok, err := ctl.repository.ToggleList(folderId, listId)

	fmt.Println(ok, err)

	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List toggled.",
	})

}
