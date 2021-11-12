package r

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"streetbox.id/app/user"
	"streetbox.id/cfg"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// UserController ...
type UserController struct {
	UserService user.ServiceInterface
}

// Login godoc
// @Summary Login user
// @Id Login
// @Tags Authorization
// @Param CLIENT_ID header string true " " Default(streetbox-mobile-merchant)
// @Param login body model.ReqUserLogin true "all fields mandatory"
// @Success 200 {object} model.ResponseSuccess "token: "exampletokenresponse" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 401 {object} model.ResponseErrors "code: 401, message: "username or password not valid, please try again" "
// @Router /login [post]
func (r *UserController) Login(c *gin.Context) {
	req := model.ReqUserLogin{}
	clientID := c.GetHeader("CLIENT_ID")
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var token string
	if token = r.UserService.Login(req, clientID); token == "" {
		model.ResponseUnauthorized(c, "username or password not valid, please try again")
		return
	}
	if token == "invalid foodtruck role" {
		model.ResponseUnauthorized(c, "Users Doesn't Have Foodtruck Role")
		return
	}
	model.ResponseJSON(c, gin.H{"token": token})
	return
}

// LoginGoogle godoc
// @Summary Login End User Using Google
// @Id LoginGoogle
// @Tags Authorization
// @Param login body model.ReqUserLoginGoogle true "Google"
// @Success 200 {object} model.ResponseSuccess "message: "exampletokenresponse" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 401 {object} model.ResponseErrors "code: 401, message: "username or password not valid, please try again" "
// @Router /login/google [post]
func (r *UserController) LoginGoogle(c *gin.Context) {
	req := model.ReqUserLoginGoogle{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var token string
	if token = r.UserService.LoginGoogle(&req); token == "" {
		model.ResponseUnauthorized(c, "Google Account Not Valid, please try again")
		return
	}
	model.ResponseJSON(c, gin.H{"token": token})
	return
}

// CreateUser godoc
// @Summary Create User Superadmin or Admin (permission = superadmin)
// @Id CreateUser
// @Tags Master User
// @Security Token
// @Param user body model.ReqUserCreate true "new user"
// @Success 201 {object} entity.Users "data: entity.Users, message: "Create User Success" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Create User Failed" "
// @Router /user [post]
func (r *UserController) CreateUser(c *gin.Context) {
	req := model.ReqUserCreate{}
	user := new(entity.Users)
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if user = r.UserService.CreateUser(req, jwtModel.UserID); user.ID == 0 {
		model.ResponseError(c, "Create User Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseCreated(c, gin.H{"data": user, "message": "Create User Success"})
	return
}

// GetAll godoc
// @Summary Get All user (permission = superadmin)
// @Id GetAllUser
// @Tags Master User
// @Security Token
// @Param filter query string false " "
// @Success 200 {object} []entity.Users "data: []entity.Users"
// @Router /user [get]
func (r *UserController) GetAll(c *gin.Context) {
	filter := c.DefaultQuery("filter", "")
	model.ResponseJSON(c, r.UserService.GetAllUser(filter))
	return
}

// UpdateUser godoc
// @Summary Update user (permission = all)
// @Id UpdateUser
// @Tags Master User
// @Security Token
// @Param user body model.ReqUserUpdate true "update user"
// @Success 200 {object} entity.Users "data: entity.Users"
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update User Failed" "
// @Router /user/update [put]
func (r *UserController) UpdateUser(c *gin.Context) {
	req := model.ReqUserUpdate{}
	user := new(entity.Users)
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if user, err = r.UserService.UpdateUser(req, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Update User Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, user)
	return
}

// GetUser godoc
// @Summary Get user profile (permission = all)
// @Id GetUser
// @Tags Master User
// @Security Token
// @Success 200 {object} model.ResUserAll "data:model.ResUserAll"
// @Router /user/info [get]
func (r *UserController) GetUser(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	model.ResponseJSON(c, r.UserService.GetUserByID(jwtModel.UserID))
	return
}

// ForgotPassword godoc
// @Summary Send Email Forgot Password
// @Id ForgotPassword
// @Tags Authorization
// @Param user body model.ReqResetPassword true "username"
// @Success 200 {object} model.ResponseSuccess "message: "Send Email Forgot Password Success" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Username not found" "
// @Router /forgotpassword [post]
func (r *UserController) ForgotPassword(c *gin.Context) {
	req := model.ReqResetPassword{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	if isValid := r.UserService.SendEmailResetPassword(req.Username); isValid == false {
		model.ResponseError(c, "Username not found", http.StatusNotAcceptable)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Send Email Forgot Password Success"})
	return
}

// ChangePassword godoc
// @Summary Change password user (permission = all)
// @Id ChangePassword
// @Tags Master User
// @Security Token
// @Param user body model.ReqChangePassword true "change password user"
// @Success 200 {object} model.ResponseSuccess "message: "Change Password Success" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Change Password User Failed" "
// @Router /user/changepassword [put]
func (r *UserController) ChangePassword(c *gin.Context) {
	req := model.ReqChangePassword{}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if err := r.UserService.ChangePassword(req.Password, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Change Password User Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Change Password Success"})
	return
}

// ResetPassword godoc
// @Summary Reset password user (permission  = all)
// @Id ResetPassword
// @Tags Authorization
// @Param user body model.ReqChangePassword true "change password user"
// @Success 200 {object} model.ResponseSuccess "message: "Reset Password Success" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Reset Password Failed" "
// @Router /resetpassword [put]
func (r *UserController) ResetPassword(c *gin.Context) {
	req := model.ReqChangePassword{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	token := util.ExtractToken(c.Request)
	if err := r.UserService.ResetForgotPassword(token, req.Password); err != nil {
		model.ResponseError(c, "Reset Password Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Reset Password Success"})
	return
}

// GetUserMerchant godoc
// @Summary Get All User Merchant Admin (permission = superadmin)
// @Id GetUserMerchant
// @Tags Master User
// @Security Token
// @Success 200 {object} []model.ResUserMerchant "data: []model.ResUserMerchant"
// @Router /user/merchant [get]
func (r *UserController) GetUserMerchant(c *gin.Context) {
	model.ResponseJSON(c, r.UserService.GetUserAdmin())
	return
}

// DeleteByID godoc
// @Summary Delete Parking Space by ID (permission = superadmin)
// @Id DeleteByIDUser
// @Tags Master User
// @Security Token
// @Param id path integer true "userID"
// @Success 200 {object} model.ResponseSuccess "message: "Success""
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Failed Delete User" "
// @Router /user/{id}/delete [delete]
func (r *UserController) DeleteByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.UserService.DeleteByID(id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Failed Delete User", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// UpdateUserRole godoc
// @Summary Update user role (permission = all)
// @Id UpdateUserRole
// @Tags Master User
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "usersID"
// @Param roleID formData integer true "Role ID"
// @Success 200 {object} model.ResponseSuccess "message: "Success" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Failed Update Role User" "
// @Router /user/role/{id}/update [put]
func (r *UserController) UpdateUserRole(c *gin.Context) {
	usersID := util.ParamIDToInt64(c.Param("id"))
	roleID := util.ParamIDToInt64(c.PostForm("roleID"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.UserService.UpdateRole(usersID, roleID, jwtModel.UserID); err != nil {
		model.ResponseError(c,
			"Failed Update Role User", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// CheckToken godoc
// @Summary Check Token Forget Password
// @Id CheckToken
// @Tags Authorization
// @Param token query string true "token reset"
// @Router /check [get]
func (r *UserController) CheckToken(c *gin.Context) {
	token := c.Query("token")
	webAddress := cfg.Config.Web.Host
	webPort := cfg.Config.Web.Port
	backoffice := fmt.Sprintf("%s:%s", webAddress, webPort)
	if _, err := r.UserService.CheckJwt(token); err != nil {
		forgotPage := backoffice + "/forgotpassword?error=true"
		c.Redirect(http.StatusSeeOther, forgotPage)
		return
	}
	resetPswdPage := backoffice + "/resetpassword/" + token
	c.Redirect(http.StatusSeeOther, resetPswdPage)
	return
}

// GetUserByID godoc
// @Summary Get user profile by ID (permission = superadmin)
// @Id GetUserByID
// @Tags Master User
// @Security Token
// @Param id path integer true "userID"
// @Success 200 {object} entity.Users "data: entity.Users"
// @Router /user/info/{id} [get]
func (r *UserController) GetUserByID(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if jwtModel.RoleName != "superadmin" {
		model.ResponseError(
			c,
			"Sorry, SuperAdmin Access Only",
			http.StatusUnprocessableEntity)
		return
	}
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, r.UserService.GetUserByID(id))
	return
}

// CreateAddress godoc
// @Summary Create User Address (permission = consumer)
// @Id CreateAddress
// @Tags Master User
// @Security Token
// @Param useraddress body entity.UsersAddress true "id and primary fields are not mandatory"
// @Success 200 {object} entity.UsersAddress "{ "data": Model }"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Create User Address Failed"
// @Router /user/address [post]
func (r *UserController) CreateAddress(c *gin.Context) {
	var req entity.UsersAddress
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	req.UserID = jwtModel.UserID
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if err := r.UserService.CreateAddress(&req); err != nil {
		model.ResponseError(c, "Create User Address Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, req)
	return
}

// GetPrimaryAddressByUserID godoc
// @Summary Get User' Primary Address (permission = consumer)
// @Id GetPrimaryAddressByUserID
// @Tags Master User
// @Security Token
// @Success 200 {object} entity.UsersAddress "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry User Address Failed"
// @Router /user/address/primary [get]
func (r *UserController) GetPrimaryAddressByUserID(c *gin.Context) {
	var addrs entity.UsersAddress
	var err error
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if addrs, err = r.UserService.GetPrimaryAddressByUserID(jwtModel.UserID); err != nil {
		model.ResponseError(c, "Inquiry User Address Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, addrs)
	return
}

// GetAddressByUserID godoc
// @Summary Get User' Addresses (permission = consumer)
// @Id GetAddressByUserID
// @Tags Master User
// @Security Token
// @Success 200 {object} []entity.UsersAddress "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry User Address Failed"
// @Router /user/address [get]
func (r *UserController) GetAddressByUserID(c *gin.Context) {
	var addrs []entity.UsersAddress
	var err error
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if addrs, err = r.UserService.GetAddressByUserID(jwtModel.UserID); err != nil {
		model.ResponseError(c, "Inquiry User Address Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, addrs)
	return
}

// UpdateUserAddress godoc
// @Summary Update User Address (permission = consumer)
// @Id UpdateUserAddress
// @Tags Master User
// @Security Token
// @Param useraddress body entity.UsersAddress true "all fields are not mandatory"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 401 {object} model.ResponseErrors "Invalid authorization"
// @Failure 500 {object} model.ResponseErrors "Create User Address Failed"
// @Router /user/address [put]
func (r *UserController) UpdateAddress(c *gin.Context) {
	var req entity.UsersAddress
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	if jwtModel.UserID != req.UserID {
		model.ResponseError(c, "Invalid authorization", http.StatusUnauthorized)
		return
	}
	if err := r.UserService.UpdateAddress(req); err != nil {
		model.ResponseError(c, "Update User Address Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// DeleteUserAddress godoc
// @Summary Delete User Address (permission = consumer)
// @Id DeleteUserAddress
// @Tags Master User
// @Security Token
// @Param id path string true "User address' id"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 500 {object} model.ResponseErrors "Delete user address failed"
// @Router /user/:id/address [delete]
func (r *UserController) DeleteAddress(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.UserService.DeleteAddress(id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Delete User Address Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// SwitchUserAddress godoc
// @Summary Switch Primary User Address (permission = consumer)
// @Id SwitchUserAddress
// @Tags Master User
// @Security Token
// @Param id path string true "User address' id"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 500 {object} model.ResponseErrors "Switch Primary user address failed"
// @Router /user/address/:id/switch [put]
func (r *UserController) SwitchAddress(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.UserService.SwitchAddress(id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Switch Primary User Address Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// UpdateUserRadius godoc
// @Summary Update User Radius (permission = superadmin)
// @Id UpdateUserRadius
// @Tags Master User
// @Security Token
// @Param rad path string true "Radius in km"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 500 {object} model.ResponseErrors "Update Users Radius Configuration Failed"
// @Router /user/config/:rad [put]
func (r *UserController) UpdateRadius(c *gin.Context) {
	rad := util.ParamIDToInt(c.Param("rad"))
	if err := r.UserService.UpdateRadius(rad); err != nil {
		model.ResponseError(c, "Update Users Radius Configuration Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// GetUserConfiguration godoc
// @Summary Get User' Configuration (permission = consumer)
// @Id GetUserConfiguration
// @Tags Master User
// @Security Token
// @Success 200 {object} entity.UsersConfig "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry User Configuration Failed"
// @Router /user/config [get]
func (r *UserController) GetConfig(c *gin.Context) {
	var cfg entity.UsersConfig
	var err error
	if cfg, err = r.UserService.GetConfig(); err != nil {
		model.ResponseError(c, "Inquiry User Configuration Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, cfg)
	return
}