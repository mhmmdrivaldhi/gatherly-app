package controllers

import (
	"gatherly-app/models/dto"
	"gatherly-app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authUC   usecase.AuthenticationUseCase
	validate *validator.Validate
}

func NewAuthController(authUC usecase.AuthenticationUseCase, rg *gin.RouterGroup) {
	ctrl := &AuthController{
		authUC:   authUC,
		validate: validator.New(),
	}

	rg.POST("/login", ctrl.Login)
}

// @Summary User Login
// @Description Authenticates a user and returns an access token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "User credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Unauthorized - Invalid credentials"
// @Router /api/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.authUC.Login(c, request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
