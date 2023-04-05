package publicController

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
)

type ListsController struct {
	repository repository.ListRepositoryInterface
}

func NewListsController(repository repository.ListRepositoryInterface) *ListsController {
	return &ListsController{
		repository: repository,
	}
}

func (ctl *ListsController) ListsByUserId(c *gin.Context) {
	// get all the lists by user id

	c.JSON(http.StatusOK, gin.H{
		"data": "user",
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

	c.JSON(http.StatusCreated, gin.H{
		"data":    listMeta,
		"message": "Your list has been created. You will get a notification after processing the list shortly.",
	})
}
