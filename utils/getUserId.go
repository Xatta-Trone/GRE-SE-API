package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (uint64,error) {
	var userId uint64

	// get the user id
	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "user id not found"})
		return userId,errors.New("user id not found")
	}

	userId, errs := strconv.ParseUint(userIdString, 10, 64)

	if errs != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "could not parse the user id"})
		return userId,errors.New("could not parse the user id")
	}

	return userId,nil

}
