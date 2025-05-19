package controllers

import (
	"gatherly-app/models/dto"
	"gatherly-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userUC   usecase.UserUsecase
	validate *validator.Validate
}

// Getter untuk userUC
func (uc *UserController) GetUserUC() usecase.UserUsecase {
	return uc.userUC
}

// Getter untuk validate
func (uc *UserController) GetValidate() *validator.Validate {
	return uc.validate
}

func NewUserController(userUC usecase.UserUsecase, rg *gin.RouterGroup) {
	ctrl := &UserController{
		userUC:   userUC,
		validate: validator.New(),
	}

	rg.POST("/users", ctrl.CreateUser)
	rg.GET("/users/:id", ctrl.GetUserByID)
	rg.GET("/users", ctrl.GetAllUsers)
	rg.PUT("/users/:id", ctrl.UpdateUser)
	rg.DELETE("/users/:id", ctrl.DeleteUser)
}

// @Summary Create a new user
// @Description Adds a new user to the database
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User Data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} string "Invalid request body"
// @Router /api/v1/users [post]
func (ctl *UserController) CreateUser(c *gin.Context) {
	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctl.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctl.userUC.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Get user by ID
// @Description Retrieves a specific user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} string "Invalid user ID"
// @Failure 404 {object} string "User not found"
// @Router /api/v1/users/{id} [get]
func (ctl *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctl.userUC.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get all users
// @Description Retrieves a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/users [get]
func (ctl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctl.userUC.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Update user by ID
// @Description Modifies an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserRequest true "Updated User Data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} string "Invalid request body"
// @Router /api/v1/users/{id} [put]
func (ctl *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctl.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctl.userUC.UpdateUser(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user by ID
// @Description Removes a user from the database
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} string "User deleted successfully"
// @Failure 400 {object} string "Invalid user ID"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/users/{id} [delete]
func (ctl *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = ctl.userUC.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
