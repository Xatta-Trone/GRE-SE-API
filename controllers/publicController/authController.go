package publicController

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/requests"
	"github.com/xatta-trone/words-combinator/services"
	"github.com/xatta-trone/words-combinator/utils"
)

type AuthController struct {
	userRepo    repository.UserRepositoryInterface
	couponRepo  repository.CouponInterface
	authService services.AuthInterface
}

func NewAuthController(userRepo repository.UserRepositoryInterface, authService services.AuthInterface, couponRepo repository.CouponInterface) *AuthController {
	return &AuthController{
		userRepo:    userRepo,
		authService: authService,
		couponRepo:  couponRepo,
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

		// check coupon
		fmt.Println(req.Promo)
		if req.Promo != "" {
			ctl.UpgradeUserWithCoupon(user, req.Promo)
		}
	}

	// record found, now issue a token
	token, exp, err := ctl.authService.GenerateTokenFromEmail(user)

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
	} else {
		cookieDomain = fmt.Sprintf(".%s", cookieDomain)
	}

	c.SetCookie("grese_token", token, ttlValue, "/", cookieDomain, false, true)
	// set a cookie domain to localhost too...for development
	// c.SetCookie("grese_token", token, ttlValue, "/", "localhost", false, true)

	user2, _ := ctl.userRepo.FindOneByEmail(req.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"exp":   exp,
		"user":  user2,
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

func (ctl *AuthController) Upgrade(c *gin.Context) {

	// find the user
	email := c.GetString("email")

	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthenticated."})
		return
	}

	// record found, now issue a token
	user, err := ctl.userRepo.FindOneByEmail(email)

	if err == sql.ErrNoRows {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		utils.Errorf(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(user.ExpiresOn)

	// validation request
	req, errs := requests.CouponValidateRequest(c)

	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	// upgrade user
	e := ctl.UpgradeUserWithCoupon(user, req.Coupon)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": e.Error()})
		return
	}

	user, _ = ctl.userRepo.FindOneByEmail(email)

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})

}

func (ctl *AuthController) UpgradeUserWithCoupon(user model.UserModel, coupon string) error {
	var err error

	fmt.Println(user.ExpiresOn)

	// check already premium user
	// check if premium expired
	today := time.Now().UTC()
	if user.ExpiresOn != nil {
		fmt.Println(today.After(*user.ExpiresOn))
		if today.Before(*user.ExpiresOn) {
			err = fmt.Errorf("you already have a ongoing subscription")
			return err
		}
	}

	couponData, err := ctl.couponRepo.FindByCoupon(coupon)

	if err == sql.ErrNoRows {
		utils.Errorf(err)
		err = fmt.Errorf("coupon expired or not found")
		return err
	}

	// check coupon expiry

	if couponData.Expires != nil && today.After(*couponData.Expires) {
		err = fmt.Errorf("coupon expired or not found")
		return err
	}

	if couponData.Type == "one_time" && couponData.UserId != nil {
		err = fmt.Errorf("coupon already used")
		return err
	}

	// check coupon max use
	if couponData.Type == "multiple" && (couponData.MaxUse-couponData.Used) <= 0 {
		err = fmt.Errorf("coupon already used to maximum limit")
		return err
	}

	expires := time.Now().UTC().AddDate(0, couponData.Months, 0)
	ok, err := ctl.userRepo.UpdateExpiresOn(expires, user.ID)

	if err != nil {
		utils.Errorf(err)
		return err
	}

	if !ok {
		err = fmt.Errorf("there was a problem updating profile")
		return err
	}

	// ctl.couponRepo.UpdateUserId(couponData.ID, user.ID)
	ctl.couponRepo.UpdateCouponStatus(couponData, user.ID)
	return nil

}

func (ctl *AuthController) PurchaseSuccess(c *gin.Context) {

	// find the user
	email := c.GetString("email")

	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthenticated."})
		return
	}

	// record found, now issue a token
	user, err := ctl.userRepo.FindOneByEmail(email)

	if err == sql.ErrNoRows {
		utils.Errorf(err)
		c.JSON(http.StatusNotFound, gin.H{"errors": "No record found."})
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		utils.Errorf(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	fmt.Println(user.ExpiresOn)

	// check already premium user
	// check if premium expired
	today := time.Now().UTC()
	if user.ExpiresOn != nil {
		fmt.Println(today.After(*user.ExpiresOn))
		if today.Before(*user.ExpiresOn) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "you already have a ongoing subscription"})
			return
		}
	}

	expires := time.Now().UTC().AddDate(0, 5, 0)
	ok, err := ctl.userRepo.UpdateExpiresOn(expires, user.ID)

	if err != nil {
		utils.Errorf(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "there was a problem updating profile."})
		return
	}

	user, _ = ctl.userRepo.FindOneByEmail(email)

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})

}
