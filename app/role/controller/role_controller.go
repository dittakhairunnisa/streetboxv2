package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"streetbox.id/app/role"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RoleController ...
type RoleController struct {
	Service role.ServiceInterface
}

// Create godoc
// @Summary Create new Role (permission = superadmin)
// @Id CreateRole
// @Tags Master Role
// @Security Token
// @Param body body model.ReqRoleCreate true "role"
// @Success 200 {object} entity.Role "data: entity.Role, message: "New Role has created" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Create Role Failed" "
// @Router /role [post]
func (ct *RoleController) Create(c *gin.Context) {
	req := model.ReqRoleCreate{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		err  error
		data *entity.Role
	)
	if data, err = ct.Service.Create(req); err != nil {
		model.ResponseError(c, "Create Role Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseCreated(c, gin.H{"data": data, "message": "New Role has created"})
	return
}

// GetAll godoc
// @Summary Get All Role (permission = superadmin)
// @Id GetAllRole
// @Tags Master Role
// @Security Token
// @Success 200 {object} []entity.Role "data: []entity.Role "
// @Router /role [get]
func (ct *RoleController) GetAll(c *gin.Context) {
	model.ResponseJSON(c, ct.Service.GetAll())
}

// GetAllExclude godoc
// @Summary Get All Role Exclude Foodtruck (permission = superadmin)
// @Id GetAllExclude
// @Tags Master Role
// @Security Token
// @Success 200 {object} []entity.Role "data: []entity.Role "
// @Router /role/exclude [get]
func (ct *RoleController) GetAllExclude(c *gin.Context) {
	model.ResponseJSON(c, ct.Service.GetAllExclude())
}
