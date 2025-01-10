package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/pkg/response"
	"github.com/wuwen/hello-go/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// @Summary     Register user
// @Description Register a new user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body     service.RegisterRequest true "User info"
// @Success     200  {object} response.Response{data=model.User}
// @Failure     400  {object} response.Response
// @Failure     500  {object} response.Response
// @Router      /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		switch err {
		case service.ErrUserExist:
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, user)
}

// @Summary     Login user
// @Description Login with username and password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body     service.LoginRequest true "Login info"
// @Success     200  {object} response.Response{data=service.LoginResponse}
// @Failure     400  {object} response.Response
// @Failure     401  {object} response.Response
// @Router      /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svc.Login(&req)
	if err != nil {
		switch err {
		case service.ErrInvalidAuth:
			response.Error(c, http.StatusUnauthorized, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, resp)
}

// @Summary     Update user
// @Description Update user info
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path     int                  true "User ID"
// @Param       user body     service.UpdateRequest true "User info"
// @Success     200  {object} response.Response{data=model.User}
// @Failure     400  {object} response.Response
// @Failure     404  {object} response.Response
// @Failure     500  {object} response.Response
// @Security    BearerAuth
// @Router      /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req service.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.svc.UpdateUser(uint(id), &req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, user)
}

// @Summary     Update user role
// @Description Update user role
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id      path     int  true "User ID"
// @Param       role_id path     int  true "Role ID"
// @Success     200     {object} response.Response{data=model.User}
// @Failure     400     {object} response.Response
// @Failure     404     {object} response.Response
// @Failure     500     {object} response.Response
// @Security    BearerAuth
// @Router      /users/{id}/roles/{role_id} [put]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	err = h.svc.UpdateUserRole(uint(userID), uint(roleID))
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		case service.ErrRoleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, nil)
}

// UpdateUserRole 更新用户角色
// @Summary     Update user role
// @Description Update user role
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path     int  true "User ID"
// @Param       role body     UpdateRoleRequest true "Role info"
// @Success     200  {object} response.Response
// @Failure     400  {object} response.Response
// @Failure     404  {object} response.Response
// @Failure     500  {object} response.Response
// @Security    BearerAuth
// @Router      /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.svc.UpdateUserRole(uint(userID), req.RoleID); err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.Error(c, http.StatusNotFound, "user not found")
		case service.ErrRoleNotFound:
			response.Error(c, http.StatusNotFound, "role not found")
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, nil)
}
