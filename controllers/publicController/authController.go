package publicController

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
)

type AuthController struct {
	userRepo    repository.UserRepositoryInterface
	authService services.AuthInterface
}

func NewAuthController(userRepo repository.UserRepositoryInterface, authService services.AuthInterface) *AuthController {
	return &AuthController{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (ctl *AuthController) Register(c *gin.Context) {

	// validation request it is also in the admin user create request
	req, errs := requests.UsersCreateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// save the record
	user, err := ctl.userRepo.Create(req)

	if err != nil && strings.Contains(err.Error(), "key 'users.email'") {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The email has already been taken."})
		return
	}

	if err != nil && strings.Contains(err.Error(), "key 'users.username'") {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The username has already been taken."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})

}

func (ctl *AuthController) Login(c *gin.Context) {
	user := model.UserModel{}

	// validation request it is also in the admin user create request
	req, errs := requests.UsersCreateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// find the user
	user, err := ctl.userRepo.FindOneByEmail(req.Email)

	if err == sql.ErrNoRows {
		// save the record
		user, err = ctl.userRepo.Create(req)

		if err != nil && strings.Contains(err.Error(), "key 'users.email'") {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "The email has already been taken."})
			return
		}

		if err != nil && strings.Contains(err.Error(), "key 'users.username'") {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "The username has already been taken."})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	}

	// record found, now issue a token
	token, err := ctl.authService.GenerateTokenFromEmail(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// set the cookie
	ttl := os.Getenv("AUTH_TTL")
	ttlValue, _ := strconv.Atoi(ttl)
	cookieDomain := os.Getenv("COOKIE_URL")

	if cookieDomain == "" {
		cookieDomain = "localhost"
	}

	c.SetCookie("grese_token", token, ttlValue, "/", cookieDomain, false, true)
	// set a cookie domain to localhost too...for development
	c.SetCookie("grese_token", token, ttlValue, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})

}

func (ctl *AuthController) Me(c *gin.Context) {

	// find the user
	email := c.GetString("email")

	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthenticated."})
		return
	}

	// record found, now issue a token
	user, err := ctl.userRepo.FindOneByEmail(email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})

}

func (ctl *AuthController) Update(c *gin.Context) {

	// find the user
	email := c.GetString("email")

	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthenticated."})
		return
	}

	// record found, now issue a token
	user, err := ctl.userRepo.FindOneByEmail(email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// validation request
	req, errs := requests.UsersProfileUpdateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs.Error()})
		return
	}

	// update profile
	ok, err := ctl.userRepo.UpdateUserName(user.ID, req.UserName)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "There was a problem updating the record. NR"})
		return
	}

	if err != nil && strings.Contains(err.Error(), "key 'users.username'") {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The username has already been taken."})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"data": ok,
	})

}

func (ctl *AuthController) Logout(c *gin.Context) {
	// set the cookie
	cookieDomain := os.Getenv("COOKIE_URL")

	if cookieDomain == "" {
		cookieDomain = "localhost"
	}

	c.SetCookie("grese_token", "deleted", -1, "/", cookieDomain, false, true)

	c.JSON(http.StatusNoContent, gin.H{
		"token": "deleted",
	})
	
}