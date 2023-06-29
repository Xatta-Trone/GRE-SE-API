package publicController

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/utils"
)

type NotificationController struct {
	notificationRepo repository.NotificationInterface
	userRepo         repository.UserRepositoryInterface
}

func NewNotificationController(repo repository.NotificationInterface, userRepo repository.UserRepositoryInterface) *NotificationController {
	return &NotificationController{
		notificationRepo: repo,
		userRepo:         userRepo,
	}
}

func (ctl *NotificationController) Index(c *gin.Context) {

	// find the user
	userId, err := utils.GetUserId(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		return
	}

	req, errs := requests.NotificationIndexRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	req.UserId = userId

	notifications, err := ctl.notificationRepo.Index(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": notifications,
		"meta": req,
	})

}


func (ctl *NotificationController) MarkAsRead(c *gin.Context) {

	// find the user
	userId, err := utils.GetUserId(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		return
	}


	ctl.notificationRepo.Update(userId)

	c.JSON(204, gin.H{
		"data": "ok",
	})

}
