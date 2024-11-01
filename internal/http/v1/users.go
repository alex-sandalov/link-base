package v1

import (
	"link-base/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type userSignUpRequest struct {
	Email        string `json:"email" binding:"required,email,min=2,max=64"`
	Password     string `json:"password" binding:"required,max=64"`
	ReferralCode string `json:"referral_code"`
}

type userSignInRequest struct {
	Email    string `json:"email" binding:"required,email,min=2,max=64"`
	Password string `json:"password" binding:"required,max=64"`
}

type refreshRequest struct {
	Token string `json:"token" binding:"required"`
}

type referralCreateRequest struct {
	TTL string `json:"ttl" binding:"required"`
}

type sendEmailRequest struct {
	Email string `json:"email" binding:"required,email,min=2,max=64"`
}

func (h *Handler) initUsersRouter(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)

		referral := users.Group("", h.userIdentity)
		{
			referral.GET("/referral", h.getReferrals)
			referral.POST("/create-code", h.createCode)
			referral.POST("/send-email")
		}

	}
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpRequest true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/sign-up [post]
func (h *Handler) userSignUp(c *gin.Context) {
	var inp userSignUpRequest
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.service.User.SignUp(c.Request.Context(), service.SignUpInput{
		Email:        inp.Email,
		Password:     inp.Password,
		ReferralCode: inp.ReferralCode,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary User SignIn
// @Tags users-auth
// @Description user sign in
// @ModuleID userSignInRequest
// @Accept  json
// @Produce  json
// @Param input body userSignInRequest true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/sign-in [post]
func (h *Handler) userSignIn(c *gin.Context) {
	var inp userSignInRequest
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.service.User.SignIn(c.Request.Context(), service.SignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshRequest true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/auth/refresh [post]
func (h *Handler) userRefresh(c *gin.Context) {
	var inp refreshRequest
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.User.RefreshTokens(c.Request.Context(), inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary User Referrals
// @Security UsersAuth
// @Tags users-referral
// @Description get user referral
// @Accept  json
// @Produce  json
// @Success 200 {array} uuid.UUID
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/referral [get]
func (h *Handler) getReferrals(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.service.Referral.FindReferralByUserID(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Create Referral Code
// @Security UsersAuth
// @Tags users-referral
// @Description create referral code for the current user
// @ModuleID createCode
// @Accept  json
// @Produce  json
// @Param input body referralCreateRequest true "Create referral code request"
// @Success 200 {string} string "referral code"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/create-code [post]
func (h *Handler) createCode(c *gin.Context) {
	var inp referralCreateRequest
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ttl, err := time.ParseDuration(inp.TTL)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.service.Referral.CreateCode(c.Request.Context(), service.ReferralInput{
		UserId: id,
		TTL:    ttl,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Send Email
// @Security UsersAuth
// @Tags users-referral
// @Description send email for the current user
// @ModuleID sendEmail
// @Accept  json
// @Produce  json
// @Param input body sendEmailRequest true "Send email request"
// @Success 200
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/send-email [post]
func (h *Handler) sendEmail(c *gin.Context) {
	var inp sendEmailRequest
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.service.Referral.SendEmail(c.Request.Context(), id, inp.Email)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
