package publicController

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/enums"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
	"golang.org/x/exp/slices"
)

type LearningStatusController struct {
	repo repository.LearningStatusInterface
}

func NewLearningStatusController(repo repository.LearningStatusInterface) *LearningStatusController {
	return &LearningStatusController{
		repo: repo,
	}
}

func (ctl *LearningStatusController) Index(c *gin.Context) {

	userId, err := utils.GetUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(userId)

	listId, err := utils.ParseParamToUint64(c, "list_id")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(listId)

	wordIds, err := ctl.repo.FindWordIdsByListId(listId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// now get the learning state ids
	learningStates, err := ctl.repo.FindLearningStatusByListId(listId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// now sort the learning states
	unknownWordIdx := []int64{}
	learningWordIdx := []int64{}
	masteredWordIdx := []int64{}

	unknownWordIdx = append(unknownWordIdx, wordIds...) // initially take all as unknown words

	for _, el := range learningStates {
		switch el.LearningState {
		case enums.LearningStatusMastered:
			// push to mastered list
			masteredWordIdx = append(masteredWordIdx, el.WordId)
			// remove from unknown word list
			// find index
			idx := slices.Index(unknownWordIdx, el.WordId)
			unknownWordIdx = slices.Delete(unknownWordIdx, idx, idx+1)

		case enums.LearningStatusLearning:
			// push to  learningWordIdx
			learningWordIdx = append(learningWordIdx, el.WordId)
			// remove from unknown word list
			// find index
			idx := slices.Index(unknownWordIdx, el.WordId)
			unknownWordIdx = slices.Delete(unknownWordIdx, idx, idx+1)
		}

	}

	// now return the data
	c.JSON(http.StatusOK, gin.H{
		"mastered": masteredWordIdx,
		"learning": learningWordIdx,
		"unknown":  unknownWordIdx,
	})

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

	
	listId, err := utils.ParseParamToUint64(c, "list_id")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// now create the record
	RowsAffected, err := ctl.repo.Delete(listId,userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(RowsAffected)

	c.Status(http.StatusNoContent)

}
