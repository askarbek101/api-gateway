package handlers

import (
	"api-gateway/internal/models/requests"
	"api-gateway/internal/models/responses"
	"api-gateway/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService services.IUserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles user creation requests
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body requests.CreateUserRequest true "User information"
// @Success 201 {object} responses.UserResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req requests.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, responses.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(201, user)
}

// GetUser handles user retrieval requests
// @Summary Get a user by ID
// @Description Get detailed information about a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security Bearer
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Error: "invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(404, responses.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, user)
}

// UpdateUser handles user update requests
// @Summary Update a user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body requests.UpdateUserRequest true "User information"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security Bearer
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Error: "invalid user ID"})
		return
	}

	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, responses.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, user)
}

// DeleteUser handles user deletion requests
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} responses.MessageResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security Bearer
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Error: "invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(400, responses.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(204)
}

// ListUsers handles user listing requests
// @Summary List all users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {array} responses.UserResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, users)
}
