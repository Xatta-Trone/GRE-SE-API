package requests

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/xatta-trone/words-combinator/enums"
)

type ListsCreateRequestStruct struct {
	Name       string `json:"name" form:"name" `
	Url        string `json:"url" form:"url" `
	Words      string `json:"words" form:"words" `
	Visibility int    `json:"visibility" form:"visibility,default=1" `
	UserId     uint64 `json:"user_id" form:"user_id"`
	Scope      string `json:"scope" form:"scope,default=user"` // either admin or user; userId required when scope is admin
}

func (c ListsCreateRequestStruct) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Scope, validation.Required),
		validation.Field(&c.Url, validation.When(c.Words == "", validation.Required, is.URL, validation.By(checkUrl(c.Url)))),
		validation.Field(&c.Words, validation.When(c.Url == "", validation.Required)),
		validation.Field(&c.UserId, validation.When(c.Scope == "admin", validation.Required)),
		validation.Field(&c.Visibility, validation.Required, validation.In(enums.ListVisibilityMe, enums.ListVisibilityPublic)),
	)
}

func ListsCreateRequest(c *gin.Context) (*ListsCreateRequestStruct, error) {
	var req *ListsCreateRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		return req, err
	}

	err = req.Validate()

	if err != nil {
		return req, err
	}

	return req, nil

}

func checkUrl(url string) validation.RuleFunc {
	return func(value interface{}) error {
		fmt.Println(url)
		if strings.Contains(url, "vocabulary.com") || strings.Contains(url, "quizlet.com") || strings.Contains(url, "memrise.com") {
			return nil
		}
		return errors.New("vocabulary.com | quizlet.com | memrise.com are only allowed")
	}
}
