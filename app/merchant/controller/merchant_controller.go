package controller

import (
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"streetbox.id/app/fcm"
	"streetbox.id/app/merchant"
	"streetbox.id/app/trx"
	"streetbox.id/app/user"
	"streetbox.id/cfg"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// MerchantController ...
type MerchantController struct {
	MerchantSvc merchant.ServiceInterface
	UserSvc     user.ServiceInterface
	TrxSvc      trx.ServiceInterface
	FcmSvc      fcm.ServiceInterface
}

// CreateMerchantCategory godoc
// @Summary Create Merchant Category (permission = superadmin)
// @Id CreateMerchantCategory
// @Tags Merchant Category
// @Security Token
// @Param merchantcategory body entity.MerchantCategory true "category field mandatory"
// @Success 200 {object} entity.MerchantCategory "{ "data": Model }"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Create category failed"
// @Router /merchant/category [post]
func (r *MerchantController) CreateCategory(c *gin.Context) {
	var req entity.MerchantCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if err := r.MerchantSvc.CreateCategory(&req); err != nil {
		model.ResponseError(c, "Create category failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, req)
	return
}

// GetAllMerchantCategory godoc
// @Summary Get All Merchant Category (permission = superadmin)
// @Id GetAllMerchantCategory
// @Tags Merchant Category
// @Security Token
// @Success 200 {object} []entity.MerchantCategory "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry categories failed"
// @Router /merchant/category [get]
func (r *MerchantController) GetAllCategory(c *gin.Context) {
	var (
		cats []entity.MerchantCategory
		err  error
	)
	if cats, err = r.MerchantSvc.GetAllCategory(); err != nil {
		model.ResponseError(c, "Inquiry categories failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, cats)
	return
}

// UpdateMerchantCategory godoc
// @Summary Update Merchant Category (permission = superadmin)
// @Id UpdateMerchantCategory
// @Tags Merchant Category
// @Security Token
// @Param merchantcategory body entity.MerchantCategory true "all fields mandatory"
// @Success 200 {object} entity.MerchantCategory "{ "data": Model }"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Update category failed"
// @Router /merchant/category [put]
func (r *MerchantController) UpdateCategory(c *gin.Context) {
	var req entity.MerchantCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if err := r.MerchantSvc.UpdateCategory(&req); err != nil {
		model.ResponseError(c, "Update category failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, req)
	return
}

// CreateMerchantCategory godoc
// @Summary Delete Merchant Category (permission = superadmin)
// @Id DeleteMerchantCategory
// @Tags Merchant Category
// @Security Token
// @Param id path string true "category's id"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 500 {object} model.ResponseErrors "Delete category failed"
// @Router /merchant/:id/category [delete]
func (r *MerchantController) DeleteCategory(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	if err := r.MerchantSvc.DeleteCategory(id); err != nil {
		model.ResponseError(c, "Delete category failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// CreateMerchant godoc
// @Summary Create Merchant (permission = admin)
// @Id CreateMerchant
// @Tags Merchant
// @Security Token
// @Param merchant body model.ReqCreateMerchant true "all fields mandatory"
// @Success 200 {object} entity.Merchant "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant [post]
func (r *MerchantController) CreateMerchant(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Merchant Admin Can Create Merchant"})
		return
	}
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant != nil && merchant.ID > 0 {
		model.ResponseJSON(c, gin.H{"message": "Sorry, You Can Only Have One Merchant"})
		return
	}
	req := model.ReqCreateMerchant{}
	var (
		data *entity.Merchant
		err  error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if data, err = r.MerchantSvc.CreateMerchant(&req, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Create Merchant Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// XenditGenerateSubAccount godoc
// @Summary Xendit Generate Sub-Account Merchant (permission = admin)
// @Id XenditGenerateSubAccount
// @Tags Merchant
// @Security Token
// @Param merchant body model.ReqXenditGenerateSubAccount true "all fields mandatory"
// @Success 200 {object} entity.Merchant "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/xendit-generate-subaccount [post]
func (r *MerchantController) XenditGenerateSubAccount(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "superadmin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Super Admin Can XenditGenerateSubAccount"})
		return
	}
	req := model.ReqXenditGenerateSubAccount{}
	var (
		data *entity.Merchant
		err  error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if data, err = r.MerchantSvc.XenditGenerateSubAccount(&req, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Xendit Generate Sub-Account Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// GetAllFoodtruck godoc
// @Summary Get All Foodtruck (permission = admin)
// @Id GetAllFoodtruck
// @Tags Merchant
// @Security Token
// @Success 200 {object} []model.ResGetFoodtruckTasks "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/foodtruck/all [get]
func (r *MerchantController) GetAllFoodtruck(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	if merchantID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	// data := r.MerchantSvc.GetAllFoodtruck(merchantID)
	data := r.MerchantSvc.GetFoodtruckTasks(jwtModel.UserID)
	model.ResponseJSON(c, data)
	return
}

// CreateFoodtruck godoc
// @Summary Create Foodtruck (permission = admin)
// @Id CreateFoodtruck
// @Tags Merchant
// @Security Token
// @Param req body model.ReqCreateFoodtruck true "username mandatory"
// @Success 200 {object} entity.Users "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/foodtruck [post]
func (r *MerchantController) CreateFoodtruck(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Merchant Admin Can Add Foodtruck"})
		return
	}
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	req := model.ReqCreateFoodtruck{}
	if merchantID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	var (
		data *entity.Users
		err  error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if r.UserSvc.GetUserByUserName(req.UserName) != nil {
		model.ResponseError(c, "Username Exist", http.StatusUnprocessableEntity)
		return
	}
	if data, err = r.MerchantSvc.CreateFoodtruck(&req, merchantID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// CreateMenu godoc
// @Summary Create Menu (permission = admin)
// @Id CreateMenu
// @Tags Merchant
// @Security Token
// @Param req body model.ReqCreateMerchantMenu true " "
// @Success 200 {object} entity.MerchantMenu "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/menu [post]
func (r *MerchantController) CreateMenu(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Merchant Admin Can Add Menu"})
		return
	}
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	req := model.ReqCreateMerchantMenu{}

	if merchantID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	var (
		data *entity.MerchantMenu
		err  error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	if data, err = r.MerchantSvc.CreateMenu(&req, merchantID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// UpdateMenu godoc
// @Summary Update Merchant Menu (permission = merchant)
// @Id UpdateMenu
// @Tags Merchant
// @Security Token
// @Param id path integer true "Merchant Menu ID"
// @Param req body model.ReqUpdateMerchantMenu true " "
// @Success 200 {object} entity.MerchantMenu "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/menu/{id} [put]
func (r *MerchantController) UpdateMenu(c *gin.Context) {
	req := model.ReqUpdateMerchantMenu{}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Merchant Admin Can Edit Menu"})
		return
	}
	id := util.ParamIDToInt64(c.Param("id"))
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.MerchantMenu
		err  error
	)
	if data, err = r.MerchantSvc.UpdateMenu(&req, merchantID, id); err != nil {
		model.ResponseError(c,
			"Failed to Update Merchant", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// DeleteMenu godoc
// @Summary Delete Merchant Menu (permission = merchant)
// @Id DeleteMenu
// @Tags Merchant
// @Security Token
// @Param id path integer true "Merchant Menu ID"
// @Success 200 {object} model.ResponseSuccess "message: "Menu Successfully Deleted" "
// @Failure 500 {object} model.ResponseErrors  "message: "Failed to Delete Menu""
// @Router /merchant/{id}/menu/delete [delete]
func (r *MerchantController) DeleteMenu(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Merchant Admin Can Delete Menu"})
		return
	}
	id := util.ParamIDToInt64(c.Param("id"))

	var err error
	if err = r.MerchantSvc.DeleteMenu(id); err != nil {
		model.ResponseError(c,
			"Failed to Delete Menu", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Menu Successfully Deleted"})
	return
}

// ListPaginateMenu godoc
// @Summary List Paginate Menu (permission = admin)
// @Id ListPaginateMenu
// @Tags Merchant
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc/ id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /merchant/menu/all [get]
func (r *MerchantController) ListPaginateMenu(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Merchant Admin Can See Menu"})
		return
	}
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sort := util.SortedBy(c.QueryArray("sort"))
	data := r.MerchantSvc.GetMenuPagination(merchantID, limit, page, sort)

	model.ResponsePagination(c, data)
	return
}

// UploadMenu godoc
// @Summary Upload Image Menu (permission = admin)
// @Id UploadMenu
// @Tags Merchant
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "Merchant Menu ID"
// @Param image formData file true "image menu"
// @Success 201 {object} model.ResponseSuccess "message: "Upload Menu Merchant Success" "
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Form/Upload Image Menu error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Image Meta data failed" "
// @Router /merchant/upload-menu/{id} [put]
func (r *MerchantController) UploadMenu(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	id := util.ParamIDToInt64(c.Param("id"))
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(
			c,
			gin.H{"message": "Sorry, Admin Access Only"},
			http.StatusUnprocessableEntity)
		return
	}
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	path := cfg.Config.Path.Image
	image, err := c.FormFile("image")
	if err != nil {
		model.ResponseError(c, "Form/Upload Image Menu error", http.StatusBadRequest)
		return
	}
	filename := util.GeneratedUUID(filepath.Base(image.Filename))
	pathImg := path + filename
	// Upload logo
	if err := c.SaveUploadedFile(image, pathImg); err != nil {
		model.ResponseError(c, "Upload Menu image error", http.StatusBadRequest)
		return
	}
	merchantMenu := r.MerchantSvc.GetMenuByID(merchant.ID, id)
	if merchantMenu.Photo != "" {
		os.Remove(path + merchantMenu.Photo)
	}
	// Save filename
	if err := r.MerchantSvc.UploadMenu(filename, merchant.ID, id); err != nil {
		model.ResponseError(c, "Upload Menu failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Upload Menu Merchant Success"})
	return
}

// DeleteImageMenu godoc
// @Summary Delete Image Menu(permission = admin)
// @Id DeleteImageMenu
// @Tags Merchant
// @Security Token
// @Produce json
// @Param id path integer true "Merchant Menu ID"
// @Success 201 {object} model.ResponseSuccess "message: "Delete Image Merchant Menu Success" "
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Delete Image Menu error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Delete Image Meta data failed" "
// @Router /merchant/delete-image-menu/{id} [put]
func (r *MerchantController) DeleteImageMenu(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	id := util.ParamIDToInt64(c.Param("id"))
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(
			c,
			gin.H{"message": "Sorry, Admin Access Only"},
			http.StatusUnprocessableEntity)
		return
	}
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	path := cfg.Config.Path.Image
	merchantMenu := r.MerchantSvc.GetMenuByID(merchant.ID, id)
	if merchantMenu.Photo != "" {
		os.Remove(path + merchantMenu.Photo)

		if err := r.MerchantSvc.RemoveImageMenu(merchantMenu, id); err != nil {
			model.ResponseError(
				c,
				gin.H{"message": "Image Already Deleted or Problem When Trying to Delete Image"},
				http.StatusUnprocessableEntity)
			return
		}
	}
	model.ResponseJSON(c, gin.H{"message": "Image Menu Successfully Deleted!"})
	return
}

// GetInfo godoc
// @Summary Get Info Merchant (permission = merchant)
// @Id GetInfoMerchant
// @Tags Merchant
// @Security Token
// @Success 200 {object} model.Merchant "data: model.Merchant"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/info [get]
func (r *MerchantController) GetInfo(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	data := r.MerchantSvc.GetInfo(jwtModel.UserID)
	model.ResponseJSON(c, data)
	return
}

// UpdateMerchant godoc
// @Summary Update Merchant (permission = merchant)
// @Id UpdateMerchant
// @Tags Merchant
// @Security Token
// @Param req body model.ReqUpdateMerchant true "Update Merchant"
// @Success 200 {object} entity.Merchant "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant [put]
func (r *MerchantController) UpdateMerchant(c *gin.Context) {
	req := model.ReqUpdateMerchant{}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	id := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.Merchant
		err  error
	)
	if data, err = r.MerchantSvc.Update(&req, id); err != nil {
		model.ResponseError(c,
			"Failed to Update Merchant", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// GetAll godoc
// @Summary Get All Merchant (permission = superadmin)
// @Id GetAllMerchant
// @Tags Merchant
// @Security Token
// @Success 200 {object} []entity.Merchant "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/all [get]
func (r *MerchantController) GetAll(c *gin.Context) {
	model.ResponseJSON(c, r.MerchantSvc.GetAll())
	return
}

// DeleteByID godoc
// @Summary Delete Merchant by ID (permission = superadmin)
// @Id DeleteByIDMerchant
// @Tags Merchant
// @Security Token
// @Param id path integer true "MerchantID"
// @Success 200 {object} model.ResponseSuccess "{ "message": "Success" }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/{id} [delete]
func (r *MerchantController) DeleteByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.MerchantSvc.DeleteByMerchantID(id, jwtModel.UserID); err != nil {
		model.ResponseError(c,
			"Failed to Delete Merchant", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// GetByID godoc
// @Summary Get Merchant by ID (permission = superadmin)
// @Id GetByIDMerchant
// @Tags Merchant
// @Security Token
// @Param id path integer true "merchantID"
// @Success 200 {object} entity.Merchant "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/info/{id} [get]
func (r *MerchantController) GetByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, r.MerchantSvc.GetByID(id))
	return
}

// UpdateFoodtruck godoc
// @Summary Update Foodtruck Info by Admin (permission = admin)
// @Id UpdateFoodtruck
// @Tags Merchant
// @Security Token
// @Param id path integer true "foodtruckID"
// @Param req body model.ReqUserUpdate true "Update Foodtruck"
// @Success 200 {object} entity.Users "data: entity.Users"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/foodtruck/{id}/update [put]
func (r *MerchantController) UpdateFoodtruck(c *gin.Context) {
	foodtruckID := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantIDFoodTruck := r.MerchantSvc.GetInfo(foodtruckID).ID
	merchantIDAdmin := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	if merchantIDAdmin != merchantIDFoodTruck {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Foodtruck Isn't Belong to You"})
		return
	}
	req := model.ReqUserUpdate{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.Users
		err  error
	)
	if data, err = r.MerchantSvc.UpdateFoodtruck(&req, foodtruckID); err != nil {
		model.ResponseError(c,
			"Failed to Update Merchant", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// ResetPasswordFoodTruck godoc
// @Summary Reset Password Footruck by Admin (permission = admin)
// @Id ResetPasswordFoodTruck
// @Tags Merchant
// @Security Token
// @Param id path integer true "foodtruckID"
// @Param req body model.ReqChangePassword true "Reset Password Foodtruck"
// @Success 201 {object} model.ResponseSuccess "message: "Reset Food Truck Password Success!" "
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Client Errors" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Server Errors" "
// @Router /merchant/foodtruck/{id}/resetpassword [put]
func (r *MerchantController) ResetPasswordFoodTruck(c *gin.Context) {
	foodTruckID := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantIDFoodTruck := r.MerchantSvc.GetInfo(foodTruckID).ID
	merchantIDAdmin := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	if merchantIDAdmin != merchantIDFoodTruck {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Foodtruck Isn't Belong to you"})
		return
	}
	req := model.ReqChangePassword{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	if err := r.UserSvc.ChangePassword(req.Password, foodTruckID); err != nil {
		model.ResponseError(c,
			"Failed to Change Password", http.StatusInternalServerError)
		return
	}

	model.ResponseJSON(c, gin.H{
		"message": "Change Food Truck Password Success!",
	})
	return
}

// UploadLogo godoc
// @Summary Upload Logo (permission = admin)
// @Id UploadLogo
// @Tags Merchant
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param logo formData file true "image logo"
// @Success 201 {object} model.ResponseSuccess "message: "Upload Logo Merchant Success" "
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Form/Upoad image error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Image Meta data failed" "
// @Router /merchant/upload-logo [put]
func (r *MerchantController) UploadLogo(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(
			c,
			gin.H{"message": "Sorry, Admin Access Only"},
			http.StatusUnprocessableEntity)
		return
	}
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	path := cfg.Config.Path.Image
	image, err := c.FormFile("logo")
	if err != nil {
		model.ResponseError(c, "Form image error", http.StatusBadRequest)
		return
	}
	filename := util.GeneratedUUID(filepath.Base(image.Filename))
	pathImg := path + filename
	// Upload logo
	if err := c.SaveUploadedFile(image, pathImg); err != nil {
		model.ResponseError(c, "Upload image error", http.StatusBadRequest)
		return
	}
	if merchant.Logo != "" {
		os.Remove(path + merchant.Logo)
	}
	// Save filename
	if err := r.MerchantSvc.UploadLogo(filename, merchant.ID); err != nil {
		model.ResponseError(c, "Upload Logo failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Upload Logo Merchant Success"})
	return
}

// UploadBanner godoc
// @Summary Upload Banner (permission = admin)
// @Id UploadBanner
// @Tags Merchant
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param banner formData file true "image banner"
// @Success 201 {object} model.ResponseSuccess "message: "Upload Banner Merchant Success" "
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Form/Upoad image error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Image Meta data failed" "
// @Router /merchant/upload-banner [put]
func (r *MerchantController) UploadBanner(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(
			c,
			gin.H{"message": "Sorry, Admin Role Only"},
			http.StatusUnprocessableEntity)
		return
	}
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	path := cfg.Config.Path.Image
	image, err := c.FormFile("banner")
	if err != nil {
		model.ResponseError(c, "Form image error", http.StatusBadRequest)
		return
	}
	filename := util.GeneratedUUID(filepath.Base(image.Filename))
	pathImg := path + filename
	// Upload banner
	if err := c.SaveUploadedFile(image, pathImg); err != nil {
		model.ResponseError(c, "Upload image error", http.StatusBadRequest)
		return
	}
	if merchant.Banner != "" {
		os.Remove(path + merchant.Banner)
	}
	// Save filename
	if err := r.MerchantSvc.UploadBanner(filename, merchant.ID); err != nil {
		model.ResponseError(c, "Upload Banner failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Upload Banner Merchant Success"})
	return
}

// GetFoodtruckByID godoc
// @Summary Get Foodtruck by ID (permission = admin)
// @Id GetFoodtruckByID
// @Tags Merchant
// @Security Token
// @Param id path integer true "foodtruckID"
// @Success 200 {object} entity.Users "data: entity.Users"
// @Router /merchant/info/{id}/foodtruck [get]
func (r *MerchantController) GetFoodtruckByID(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(
			c,
			gin.H{"message": "Sorry, Admin Role Only"},
			http.StatusUnprocessableEntity)
		return
	}
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, r.MerchantSvc.GetFoodtruckByID(id))
	return
}

// DeleteFoodtruckByID godoc
// @Summary Delete Foodtruck by ID (permission = admin)
// @Id DeleteFoodtruckByID
// @Tags Merchant
// @Security Token
// @Param id path integer true "FoodtruckID"
// @Success 200 {object} model.ResponseSuccess "{ "message": "Success" }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /merchant/{id}/foodtruck [delete]
func (r *MerchantController) DeleteFoodtruckByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.UserSvc.DeleteByID(id, jwtModel.UserID); err != nil {
		model.ResponseError(c,
			"Failed to Delete Foodtruck", http.StatusInternalServerError)
		return
	}
	if err := r.MerchantSvc.DeleteFoodTruckByID(id); err != nil {
		model.ResponseError(c,
			"Failed to Delete Foodtruck", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// ListMenu godoc
// @Summary Merchant List Menu (permission = merchant)
// @Id ListMenu
// @Tags Merchant
// @Security Token
// @Param filter query string false "all, nearby, visit" default(all)
// @Success 200 {object} []entity.MerchantMenu ""
// @Router /merchant/list/menu [get]
func (r *MerchantController) ListMenu(c *gin.Context) {
	var nearby, visit bool
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "foodtruck" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Foodtruck Can See the Menu"})
		return
	}
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	filter := c.DefaultQuery("filter", "all")
	if filter == "all" {
		nearby, visit = true, true
	} else if filter == "nearby" {
		nearby = true
	} else if filter == "visit" {
		visit = true
	}
	menus := r.MerchantSvc.GetMenuList(merchantID, nearby, visit)
	var (
		updatedMenus       []model.ResMerchantMenuList
		updatedMenu        model.ResMerchantMenuList
		price              float64
		priceAfterDiscount float64
	)
	for _, value := range *menus {
		price = float64(value.Price)
		priceAfterDiscount = math.Floor(price * float64(value.Discount/100))
		copier.Copy(&updatedMenu, value)
		updatedMenu.Price = price
		updatedMenu.PriceAfterDiscount = priceAfterDiscount
		updatedMenus = append(updatedMenus, updatedMenu)
	}

	model.ResponseJSON(c, updatedMenus)
	return
}

// GetMenuByID godoc
// @Summary Get Menu By ID (permission = merchant)
// @Id GetMenuByID
// @Tags Merchant
// @Security Token
// @Param id path integer true "Merchant Menu ID"
// @Success 200 {object} entity.MerchantMenu ""
// @Router /merchant/menu/info/{id} [get]
func (r *MerchantController) GetMenuByID(c *gin.Context) {
	menuID := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "admin" {
		model.ResponseJSON(c, gin.H{"message": "Sorry, Only Foodtruck Can See the Menu"})
		return
	}
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	menu := r.MerchantSvc.GetMenuByID(merchantID, menuID)
	var (
		price              float64
		priceAfterDiscount float64
	)

	price = float64(menu.Price)
	priceAfterDiscount = math.Floor(price * float64(menu.Discount/100))

	updatedMenu := new(model.ResMerchantMenuList)
	copier.Copy(&updatedMenu, menu)
	updatedMenu.Price = price
	updatedMenu.PriceAfterDiscount = priceAfterDiscount

	model.ResponseJSON(c, updatedMenu)
	return
}

// CountFoodtruck godoc
// @Summary Count Foodtruck (permission = admin)
// @Id CountFoodtruck
// @Tags Merchant
// @Security Token
// @Success 200 {object} model.ResCountFoodtruck "data: model.ResCountFoodtruck"
// @Router /merchant/count/foodtruck [GET]
func (r *MerchantController) CountFoodtruck(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	Foodtrucks := r.MerchantSvc.CountFoodtruckByMerchantID(merchantID)

	model.ResponseJSON(c, model.ResCountFoodtruck{Foodtrucks: Foodtrucks})
	return
}

// GetTaxSetting godoc
// @Summary Get Merchant Tax Setting (permission = merchant)
// @Id GetTaxSetting
// @Tags Merchant
// @Security Token
// @Success 200 {object} entity.MerchantTax "data: entity.MerchantTax"
// @Router /merchant/taxsetting/menu [GET]
func (r *MerchantController) GetTaxSetting(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "foodtruck" && jwtModel.RoleName != "admin" {
		model.ResponseError(c, gin.H{"message": "Sorry, Foodtruck Role Only"}, http.StatusUnprocessableEntity)
		return
	}

	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}

	model.ResponseJSON(c, r.MerchantSvc.GetTax(merchant.ID))
	return
}

// SetMerchantTax godoc
// @Summary Set Merchant Tax (permission = admin)
// @Id MerchantTax
// @Tags Merchant
// @Security Token
// @Param req body model.MerchantTax true " "
// @Success 200 {object} entity.MerchantTax "data: entity.MerchantTax"
// @Router /merchant/tax/menu [POST]
func (r *MerchantController) SetMerchantTax(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(c, gin.H{"message": "Sorry, Admin Role Only"}, http.StatusUnprocessableEntity)
		return
	}

	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}
	req := model.MerchantTax{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	checkTax := r.MerchantSvc.GetTax(merchant.ID)

	var (
		data *entity.MerchantTax
		err  error
	)
	if checkTax == nil {
		if data, err = r.MerchantSvc.CreateTax(&req, merchant.ID); err != nil {
			model.ResponseError(c, "Problem When Create Tax", http.StatusBadRequest)
			return
		}
	} else {
		if data, err = r.MerchantSvc.UpdateTax(&req, checkTax.ID, merchant.ID); err != nil {
			model.ResponseError(c, "Problem When Update Tax", http.StatusBadRequest)
			return
		}
	}
	model.ResponseJSON(c, data)
	return
}

// ImportCSV godoc
// @Summary Upload CSV (permission = admin)
// @Id ImportCSV
// @Tags Merchant
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param csv formData file true "csv file"
// @Param imageZip formData file false "image zip"
// @Success 200 {object} []entity.MerchantMenu "{ "data": Model }"
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Form/Upoad image error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Image Meta data failed" "
// @Router /merchant/csv/uploadmenu [POST]
func (r *MerchantController) ImportCSV(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "admin" {
		model.ResponseError(c, gin.H{"message": "Sorry, Admin Role Only"}, http.StatusUnprocessableEntity)
		return
	}

	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant.ID == 0 {
		model.ResponseJSON(c, gin.H{"message": "Please Create Merchant First"})
		return
	}

	var (
		err      error
		pathCsv  string
		records  [][]string
		datas    []entity.MerchantMenu
		data     *entity.MerchantMenu
		req      model.ReqCreateMerchantMenu
		price    int64
		qty      int
		discount float32
	)
	imageZip, _ := c.FormFile("imageZip")

	if imageZip != nil {
		uuidFilenames, filenames, err := util.Unzip(imageZip, cfg.Config.Path.Image, c)
		if strings.Contains(imageZip.Filename, ".zip") == false {
			model.ResponseError(c, "File Must Be zip", http.StatusBadRequest)
			return
		}
		if err != nil {
			model.ResponseError(c, "Unzip File Error", http.StatusBadRequest)
			return
		}

		records, err = util.ReadCsv(c, "csv")
		if err != nil {
			model.ResponseError(c, err.Error(), http.StatusBadRequest)
			return
		}

		records = records[1:]
		//uuidFilename := reflect.ValueOf(uuidFilenames)
		var indexImage = 1

		for _, value := range records {
			price = util.ParamIDToInt64(value[3])
			qty = util.ParamIDToInt(value[5])
			if value[4] == "" {
				discount = 0
			} else {
				discount = util.ParamToFloat32(value[4])
			}
			req = model.ReqCreateMerchantMenu{
				Name:        value[1],
				Description: value[2],
				Price:       price,
				Qty:         qty,
				IsActive:    true,
				IsNearby:    true,
				IsVisit:     true,
				Discount:    discount,
			}
			//check filename in csv are same as in extracted zip file images
			if util.FindStringInSlice(filenames, filenames[0]+value[0]) {
				if data, err = r.MerchantSvc.CreateMenuWithImage(&req, uuidFilenames[indexImage], merchant.ID); err != nil {
					model.ResponseError(c, err.Error(), http.StatusInternalServerError)
					return
				}
				datas = append(datas, *data)
				indexImage++
				continue
			}

			if data, err = r.MerchantSvc.CreateMenu(&req, merchant.ID); err != nil {
				model.ResponseError(c, err.Error(), http.StatusInternalServerError)
				return
			}
			datas = append(datas, *data)
		}
		os.Remove(pathCsv)
		model.ResponseJSON(c, datas)
		return
	}

	records, err = util.ReadCsv(c, "csv")
	if err != nil {
		model.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	records = records[1:]

	for _, value := range records {
		price = util.ParamIDToInt64(value[3])
		qty = util.ParamIDToInt(value[5])
		if value[4] == "" {
			discount = 0
		} else {
			discount = util.ParamToFloat32(value[4])
		}
		req = model.ReqCreateMerchantMenu{
			Name:        value[1],
			Description: value[2],
			Price:       price,
			Qty:         qty,
			IsActive:    true,
			IsNearby:    true,
			IsVisit:     true,
			Discount:    discount,
		}

		if data, err = r.MerchantSvc.CreateMenu(&req, merchant.ID); err != nil {
			model.ResponseError(c, err.Error(), http.StatusInternalServerError)
			return
		}
		datas = append(datas, *data)
	}
	os.Remove(pathCsv)
	model.ResponseJSON(c, datas)
	return
}

// DownloadMenuTemplateCSV ...
func (r *MerchantController) DownloadMenuTemplateCSV(c *gin.Context) {
	fileName := "menutemplate.csv"
	targetPath := filepath.Join(cfg.Config.Path.Csv, fileName)
	c.Header("Content-Disposition", "attachment;filename="+fileName)
	c.File(targetPath)
}

// GetPosTransaction godoc
// @Summary Get Order POS Transaction
// @Id GetOrderTransaction
// @Tags Merchant
// @Security Token
// @Param startDate query string false "date DD/MM/YYYY"
// @Param endDate query string false "date DD/MM/YYYY"
// @Param keyword query string false "keyword like trx number and order number"
// @Success 200 {object} model.ResTrxOrderList "data: model.ResTrxOrderList"
// @Router /merchant/pos/gettransaction [get]
func (r *MerchantController) GetPosTransaction(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// admin only
	if jwtModel.RoleName != "foodtruck" {
		model.ResponseError(
			c,
			gin.H{"message": "Sorry, Food Truck Role Only"},
			http.StatusUnprocessableEntity)
		return
	}
	var (
		startDates string
		endDates   string
	)
	if c.DefaultQuery("startDate", "") != "" {
		startDate, _ := time.Parse("02/01/2006", c.DefaultQuery("startDate", ""))
		startDates = startDate.Format("2006-01-02")
	} else {
		startDates = ""
	}

	if c.DefaultQuery("endDate", "") != "" {
		endDate, _ := time.Parse("02/01/2006", c.DefaultQuery("endDate", ""))
		endDates = endDate.Format("2006-01-02")
	} else {
		endDates = ""
	}
	merchantUsersID := r.MerchantSvc.GetMerchantUsersByUsersID(jwtModel.UserID).ID
	model.ResponseJSON(c, r.TrxSvc.GetOrderTrx(merchantUsersID, startDates, endDates, c.DefaultQuery("keyword", "")))
	return
}

// RegistrationToken godoc
// @Summary Send Registration Token FCM, Its should be done to archive push notification (permission = merchant)
// @Id RegistrationToken
// @Tags Merchant
// @Security Token
// @Param token path string true "fcm token from client SDK"
// @Success 200 {object} model.ResponseSuccess "data: Registration Token Successfully"
// @Router /merchant/registration-token/{token} [post]
func (r *MerchantController) RegistrationToken(c *gin.Context) {
	regisToken := c.Param("token")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.MerchantSvc.RegistrationToken(regisToken, jwtModel.UserID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(c, "Registration Token Successfully")
	return
}

// RemoveLogo godoc
// @Summary Remove Logo (permission = admin)
// @Id RemoveLogo
// @Tags Merchant
// @Security Token
// @Param filename path string true "logo filename"
// @Success 200 {object} model.ResponseSuccess "data: Logo Removed Successfully"
// @Router /merchant/remove-logo/{filename} [put]
func (r *MerchantController) RemoveLogo(c *gin.Context) {
	filename := c.Param("filename")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant == nil {
		model.ResponseError(c, "User Doesn't have merchant", http.StatusUnprocessableEntity)
		return
	}
	if err := r.MerchantSvc.RemoveImage(filename, "logo", merchant.ID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	path := cfg.Config.Path.Image
	if err := os.Remove(path + filename); err != nil {
		log.Printf("INFO: Image %s Not Found", filename)
	}
	model.ResponseJSON(c, "Logo Removed Successfully")
	return
}

// RemoveBanner godoc
// @Summary Remove Banner (permission = admin)
// @Id RemoveBanner
// @Tags Merchant
// @Security Token
// @Param filename path string true "banner filename"
// @Success 200 {object} model.ResponseSuccess "data: Banner Removed Successfully"
// @Router /merchant/remove-banner/{filename} [put]
func (r *MerchantController) RemoveBanner(c *gin.Context) {
	filename := c.Param("filename")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant == nil {
		model.ResponseError(c, "User Doesn't have merchant", http.StatusUnprocessableEntity)
		return
	}
	if err := r.MerchantSvc.RemoveImage(filename, "banner", merchant.ID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	path := cfg.Config.Path.Image
	if err := os.Remove(path + filename); err != nil {
		log.Printf("INFO: Image %s Not Found", filename)
	}
	model.ResponseJSON(c, "Banner Removed Successfully")
	return
}
