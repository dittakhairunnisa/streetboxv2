package controller

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"streetbox.id/app/appsetting"
	"streetbox.id/model"
	"streetbox.id/util"
)

// AppSettingController ...
type AppSettingController struct {
	AppSettingSvc appsetting.ServiceInterface
}

// GetByKey godoc
// @Summary Get App Setting Value (permission = all)
// @Id GetByKey
// @Tags AppSetting
// @Security Token
// @Param key query string false " " default("nearby_radius")
// @Success 200 {object} model.AppSetting "data: model.AppSetting"
// @Router /appsetting/get-by-key/:key [GET]
func (r *AppSettingController) GetByKey(c *gin.Context) {
	key := c.Param("key")
	result := r.AppSettingSvc.GetByKey(key)

	model.ResponseJSON(c, result)
	return
}

// UpdateByKey godoc
// @Summary Update App Setting Value (permission = superadmin)
// @Id UpdateByKey
// @Tags AppSetting
// @Security Token
// @Param key query string false " " default("nearby_radius")
// @Param req body model.ReqUpdateAppSettingByKey true "Update AppSetting"
// @Success 200 {object} model.AppSetting "data: model.AppSetting"
// @Router /appsetting/update-by-key/:key [POST]
func (r *AppSettingController) UpdateByKey(c *gin.Context) {
	req := model.ReqUpdateAppSettingByKey{}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "superadmin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Super Admin Can See Menu"})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	key := c.Param("key")
	result := r.AppSettingSvc.UpdateByKey(key, &req)

	model.ResponseJSON(c, result)
	return
}