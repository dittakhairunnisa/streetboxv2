package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"streetbox.id/app/logactivitymerchant"
	"streetbox.id/app/merchant"
	"streetbox.id/model"
	"streetbox.id/util"
)

// LogMerchantController ..
type LogMerchantController struct {
	Svc         logactivitymerchant.ServiceInterface
	MerchantSvc merchant.ServiceInterface
}

// GetAll godoc
// @Summary Get All Log Merchant Pagination (permission = admin)
// @Id GetAllLogMerchant
// @Tags Log Merchant
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /log-merchant [get]
func (s *LogMerchantController) GetAll(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sorted := util.SortedBy(c.QueryArray("sort"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := s.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant == nil {
		model.ResponseError(c, "Merchant Not Found", http.StatusUnprocessableEntity)
		return
	}
	data := s.Svc.GetAll(limit, page, sorted, merchant.ID)
	model.ResponsePagination(c, data)
	return
}
