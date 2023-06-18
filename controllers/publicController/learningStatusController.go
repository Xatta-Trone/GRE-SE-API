package publicController

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type LearningStatusController struct {
	repo repository.LearningStatusInterface
}

func NewLearningStatusController(repo repository.LearningStatusInterface) *LearningStatusController {
	return &LearningStatusController{
		repo: repo,
	}
}

func (ctl *LearningStatusController) Update(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(userId)

	// validation request
	req, errs := requests.LearningStatusUpdateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId

	// now create the record
	RowsAffected, err := ctl.repo.Create(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(RowsAffected)

	c.Status(http.StatusNoContent)

}

func (ctl *LearningStatusController) Delete(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(userId)

	// validation request
	req, errs := requests.LearningStatusDeleteRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId

	// now create the record
	RowsAffected, err := ctl.repo.Delete(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(RowsAffected)

	c.Status(http.StatusNoContent)

}
