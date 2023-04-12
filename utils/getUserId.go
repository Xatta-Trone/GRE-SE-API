package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) uint64 {
	var userId uint64

	// get the user id
	userIdString := c.GetString("user_id")

	if userIdString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return userId
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse the user id"})
		return userId
	}

	return userId

}
