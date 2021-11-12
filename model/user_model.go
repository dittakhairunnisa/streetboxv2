package model

import "time"

// ReqUserCreate create  user superadmin/admin
type ReqUserCreate struct {
	UserName string `json:"userName" binding:"required"`
	RoleID   int    `json:"roleId" binding:"required"`
}

// ReqChangePassword ...
type ReqChangePassword struct {
	Password string `json:"password" binding:"required"`
}

// ReqResetPassword ...
type ReqResetPassword struct {
	Username string `json:"userName" binding:"required"`
}

// ReqUserConsumerCreate ... create consumer user
type ReqUserConsumerCreate struct {
	UserName string
	Password string
	RoleName string
}

// ReqUserUpdate ...
type ReqUserUpdate struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
	PlatNo  string `json:"platNo,omitempty"`
}

// ReqEndUserUpdateWithImage ...
type ReqEndUserUpdateWithImage struct {
	Name         string `json:"name,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Address      string `json:"address,omitempty"`
	Email        string `json:"email,omitempty"`
	PhotoProfile string `json:"photoProfile,omitempty"`
}

// ReqEndUserUpdate ...
type ReqEndUserUpdate struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
	Email   string `json:"email,omitempty"`
}

// ReqUserLogin ...
type ReqUserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ReqUserLoginGoogle ..
type ReqUserLoginGoogle struct {
	IDToken        string `json:"idToken" binding:"required"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profilePicture"`
}

// ReqRoleCreate ...
type ReqRoleCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// ReqUserRoleCreate ...
type ReqUserRoleCreate struct {
	ID     int64
	UserID int64
	RoleID int64
}

// ResUserMerchant Response Get User Merchant
type ResUserMerchant struct {
	UserName   string `json:"username"`
	Name       string `json:"merchantName"`
	MerchantID int64  `json:"merchantId"`
	UsersID    int64  `json:"usersId"`
}

// ResUserAll ...
type ResUserAll struct {
	ID                int64      `json:"id"`
	UserName          string     `json:"userName"`
	Name              string     `json:"name"`
	Phone             string     `json:"phone"`
	Address           string     `json:"address"`
	RoleID            int64      `json:"roleId"`
	RoleName          string     `json:"roleName"`
	PlatNo            string     `json:"platNo"`
	ProfilePicture    string     `json:"profilePicture"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	RegistrationToken string     `json:"-"`
}
