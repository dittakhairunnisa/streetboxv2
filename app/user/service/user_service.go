package s

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"streetbox.id/cfg"
	"streetbox.id/util"

	"streetbox.id/app/logactivity"
	"streetbox.id/app/role"
	"streetbox.id/app/user"
	"streetbox.id/app/useraddress"
	"streetbox.id/app/userauth"
	"streetbox.id/app/userrole"
	"streetbox.id/app/userconfig"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// UserService ...
type UserService struct {
	UserRepo        user.RepoInterface
	RoleRepo        role.RepoInterface
	UserRoleRepo    userrole.RepoInterface
	UserAuthRepo    userauth.RepoInterface
	UserAddressRepo useraddress.RepoInterface
	userConfigRepo  usersconfig.RepoInterface
	LogRepo         logactivity.RepoInterface
}

// New ... init s
func New(
	userRepo user.RepoInterface,
	roleRepo role.RepoInterface,
	userRoleRepo userrole.RepoInterface,
	userAddressRepo useraddress.RepoInterface,
	userAuthRepo userauth.RepoInterface,
	userConfigRepo usersconfig.RepoInterface,
	logRepo logactivity.RepoInterface) user.ServiceInterface {
	return &UserService{
		userRepo, roleRepo,
		userRoleRepo, userAuthRepo, userAddressRepo, userConfigRepo, logRepo}
}

// CreateUser ...
func (s *UserService) CreateUser(
	req model.ReqUserCreate, usersID int64) *entity.Users {
	user := new(entity.Users)
	user.UserName = req.UserName
	username := strings.Split(req.UserName, "@")
	user.Name = username[0]
	password := username[0]

	if trx, err := s.UserRepo.Create(user); err == nil {
		userRole := new(entity.UsersRole)
		userAuth := new(entity.UsersAuth)

		userRole.UsersID = user.ID
		userRole.RoleID = int64(req.RoleID)

		userAuth.UserName = user.UserName
		userAuth.Password = util.HashPassword(password)
		if err := s.UserRoleRepo.Create(trx, userRole); err == nil {
			if err := s.UserAuthRepo.Create(trx, userAuth); err == nil {
				trx.Commit()
				userAdmin := s.UserRepo.FindByID(usersID).UserName
				log := fmt.Sprintf("Add New User %s by %s", user.UserName, userAdmin)
				s.LogRepo.Create(log)
				return user
			}
		}
	}
	return user
}

// Login ...
func (s *UserService) Login(req model.ReqUserLogin, clientID string) string {
	match := s.UserAuthRepo.FindByUsernameAndPassword(req.Username, req.Password)
	if match == false {
		return ""
	}
	user := s.UserRepo.FindByUsername(req.Username)
	// get role name
	roleName := s.UserRoleRepo.GetNameByUserID(user.ID)
	// login pos only
	if clientID == "streetbox-mobile-pos" {
		if roleName != "foodtruck" {
			return "invalid foodtruck role"
		}
	}
	// create token
	return util.CreateToken(user.ID, roleName)
}

// GetUserByUserName ...
func (s *UserService) GetUserByUserName(username string) *entity.Users {
	return s.UserRepo.FindByUsername(username)
}

// ResetPassword ...
func (s *UserService) ResetPassword(newPassword, userName string) error {
	hashPassword := util.HashPassword(newPassword)
	return s.UserAuthRepo.ResetPassword(userName, hashPassword)
}

// SendEmailResetPassword ...
func (s *UserService) SendEmailResetPassword(userName string) bool {
	apiHost := cfg.Config.Api.Host
	apiPort := cfg.Config.Api.Port
	emailCfg := cfg.Config.Smpt.Email
	passwordCfg := cfg.Config.Smpt.Password
	emailHostCfg := cfg.Config.Smpt.Host
	portCfg := cfg.Config.Smpt.Port
	env := cfg.Config.Env
	var urlReset string
	if env == "development" {
		urlReset = fmt.Sprintf("%s:%s", apiHost, apiPort)
	}
	urlReset = fmt.Sprintf("%s", apiHost)
	//
	// generate Token
	user := s.UserRepo.FindByUsername(userName)
	if user == nil {
		return false
	}
	roleName := s.UserRoleRepo.GetNameByUserID(user.ID)
	token := util.CreateTokenReset(user.ID, roleName)
	// init
	auth := smtp.PlainAuth("", emailCfg, passwordCfg, emailHostCfg)
	smtpAddr := fmt.Sprintf("%s:%s", emailHostCfg, portCfg)
	to := []string{userName}
	msg := []byte("To: " + userName + "\r\n" +
		"Subject: Reset Password\r\n" +
		"\r\n" +
		"Untuk mereset password, silahkan klik link ini\r\n" +
		urlReset + "/check?token=" + token)
	err := smtp.SendMail(smtpAddr, auth, emailCfg, to, msg)
	if err != nil {
		return false
	}
	return true
}

// GetAllUser ...
func (s *UserService) GetAllUser(filter string) *[]model.ResUserAll {
	return s.UserRepo.GetAll(filter)
}

// UpdateUser ...
func (s *UserService) UpdateUser(req model.ReqUserUpdate, id int64) (*entity.Users, error) {
	user := new(entity.Users)
	copier.Copy(user, req)
	if err := s.UserRepo.Update(user, id); err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// GetUserByID ...
func (s *UserService) GetUserByID(id int64) *model.ResUserAll {
	return s.UserRepo.FindByID(id)
}

// ChangePassword ...
func (s *UserService) ChangePassword(password string, id int64) error {
	hashPassword := util.HashPassword(password)
	user := s.UserRepo.FindByID(id)
	return s.UserAuthRepo.ResetPassword(user.UserName, hashPassword)
}

// ResetForgotPassword ...
func (s *UserService) ResetForgotPassword(tokenStr string, password string) error {
	token, err := s.CheckJwt(tokenStr)
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	// extract jwt
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if !ok {
			return err
		}
		if err := s.ChangePassword(password, userID); err != nil {
			return err
		}
	}
	return nil

}

// GetUserAdmin ...
func (s *UserService) GetUserAdmin() *[]model.ResUserMerchant {
	return s.UserRepo.GetUserAdmin()
}

// DeleteByID ...
func (s *UserService) DeleteByID(id int64, userID int64) error {
	if trx, err := s.UserRepo.DeleteByID(id); err == nil {
		if err := s.UserRoleRepo.DeleteByID(trx, id); err == nil {
			userName := s.UserRepo.FindByID(id).UserName
			if err := s.UserAuthRepo.DeleteByID(trx, userName); err == nil {
				trx.Commit()
				return nil
			}
		}
	}
	userName := s.UserRepo.FindByID(id).UserName
	userAdmin := s.UserRepo.FindByID(userID).UserName
	msg := fmt.Sprintf("Delete User %s by %s",
		userName, userAdmin)
	s.LogRepo.Create(msg)
	return errors.New("Failed to Delete Users")
}

// UpdateRole ...
func (s *UserService) UpdateRole(usersID int64, roleID int64, userID int64) error {
	userName := s.UserRepo.FindByID(usersID).UserName
	roleName := s.RoleRepo.GetOne(roleID).Name
	userAdmin := s.UserRepo.FindByID(userID).UserName
	msg := fmt.Sprintf("Updated Role with UserName %s and Role %s by %s",
		userName, roleName, userAdmin)
	s.LogRepo.Create(msg)
	return s.UserRoleRepo.Update(usersID, roleID)
}

// CheckJwt ...
func (s *UserService) CheckJwt(tokenParam string) (*jwt.Token, error) {
	jwtKey := cfg.Config.JwtKey
	token, err := jwt.Parse(tokenParam, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// LoginGoogle ..
func (s *UserService) LoginGoogle(req *model.ReqUserLoginGoogle) string {
	token, err := util.ExtractGoogleTokenInfo(req.IDToken)
	data := new(entity.Users)
	if err != nil {
		return ""
	}
	if data = s.UserRepo.FindByUsername(token.Email); data == nil {
		user := new(entity.Users)
		user.UserName = token.Email
		if len(strings.TrimSpace(req.Name)) == 0 {
			username := strings.Split(token.Email, "@")
			user.Name = username[0]
		}
		user.Name = strings.TrimSpace(req.Name)
		if len(req.Phone) != 0 {
			user.Phone = req.Phone
		}
		if len(req.ProfilePicture) != 0 {
			user.ProfilePicture = req.ProfilePicture
		}
		if trx, err := s.UserRepo.Create(user); err == nil {
			userRole := new(entity.UsersRole)

			userRole.UsersID = user.ID
			userRole.RoleID = s.RoleRepo.FindByName("consumer").ID

			if err := s.UserRoleRepo.Create(trx, userRole); err == nil {
				trx.Commit()
				log := fmt.Sprintf("Created New End User %s", user.UserName)
				s.LogRepo.Create(log)
				return util.CreateToken(user.ID, "consumer")
			}
		}

	}

	return util.CreateToken(data.ID, "consumer")
}

func (s *UserService) CreateAddress(addr *entity.UsersAddress) (err error) {
	return s.UserAddressRepo.Create(addr)
}

func (s *UserService) GetPrimaryAddressByUserID(userID int64) (addrs entity.UsersAddress, err error) {
	return s.UserAddressRepo.GetPrimaryByUserID(userID)
}

func (s *UserService) GetAddressByUserID(userID int64) (addrs []entity.UsersAddress, err error) {
	return s.UserAddressRepo.GetByUserID(userID)
}

func (s *UserService) DeleteAddress(id, userID int64) (err error) {
	return s.UserAddressRepo.Delete(id, userID)
}

func (s *UserService) UpdateAddress(addr entity.UsersAddress) (err error) {
	return s.UserAddressRepo.Update(addr)
}

func (s *UserService) SwitchAddress(id, userID int64) (err error) {
	return s.UserAddressRepo.Switch(id, userID)
}

func (s *UserService) UpdateRadius(rad int) (err error) {
	return s.userConfigRepo.UpdateRadius(rad)
}

func (s *UserService) GetConfig() (cfg entity.UsersConfig, err error) {
	return s.userConfigRepo.GetConfig()
}
