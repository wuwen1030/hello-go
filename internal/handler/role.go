package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/pkg/response"
	"github.com/wuwen/hello-go/internal/service"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// @Summary     Create role
// @Description Create a new role
// @Tags        roles
// @Accept      json
// @Produce     json
// @Param       role body     service.CreateRoleRequest true "Role info"
// @Success     200  {object} response.Response{data=model.Role}
// @Failure     400  {object} response.Response
// @Failure     500  {object} response.Response
// @Security    BearerAuth
// @Router      /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.svc.Create(&req)
	if err != nil {
		switch err {
		case service.ErrRoleExist:
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, role)
}

// @Summary     Get role
// @Description Get role by ID
// @Tags        roles
// @Accept      json
// @Produce     json
// @Param       id   path     int true "Role ID"
// @Success     200  {object} response.Response{data=model.Role}
// @Failure     404  {object} response.Response
// @Failure     500  {object} response.Response
// @Security    BearerAuth
// @Router      /roles/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	role, err := h.svc.Get(uint(id))
	if err != nil {
		switch err {
		case service.ErrRoleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, role)
}

// @Summary     Update role
// @Description Update role by ID
// @Tags        roles
// @Accept      json
// @Produce     json
// @Param       id   path     int                  true "Role ID"
// @Param       role body     service.UpdateRoleRequest true "Role info"
// @Success     200  {object} response.Response{data=model.Role}
// @Failure     400  {object} response.Response
// @Failure     404  {object} response.Response
// @Failure     500  {object} response.Response
// @Security    BearerAuth
// @Router      /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.svc.Update(uint(id), &req)
	if err != nil {
		switch err {
		case service.ErrRoleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, role)
}

// @Summary     Delete role
// @Description Delete role by ID
// @Tags        roles
// @Accept      json
// @Produce     json
// @Param       id   path     int true "Role ID"
// @Success     200  {object} response.Response
// @Failure     404  {object} response.Response
// @Failure     500  {object} response.Response
// @Security    BearerAuth
// @Router      /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	err = h.svc.Delete(uint(id))
	if err != nil {
		switch err {
		case service.ErrRoleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, nil)
}
