package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseParamToUint64(c *gin.Context, param string) (uint64,error) {
	var Uint64Value uint64
	stringValue := c.Param(param)
	var err error = nil

	if stringValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("missing param %s.", param)})
		return Uint64Value,errors.New(fmt.Sprintf("missing param %s.", param))
	}

	Uint64Value, err = strconv.ParseUint(stringValue, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("could not parse %s", param)})
		return Uint64Value,errors.New(fmt.Sprintf("could not parse %s", param))
	}

	return Uint64Value,err
}

func ParseQueryToUint64(c *gin.Context, param string) uint64 {
	var Uint64Value uint64
	stringValue := c.Query(param)

	if stringValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("missing query %s.", param)})
		return Uint64Value
	}

	Uint64Value, err := strconv.ParseUint(stringValue, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("could not parse %s", param)})
		return Uint64Value
	}

	return Uint64Value
}

func ParseQueryString(c *gin.Context, query string) string {

	return c.Query(query)

}
