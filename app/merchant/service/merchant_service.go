package service

import (
	"bytes"
	bs64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/logactivity"
	"streetbox.id/app/logactivitymerchant"
	"streetbox.id/app/merchant"
	"streetbox.id/app/merchantcategory"
	"streetbox.id/app/merchantmenu"
	"streetbox.id/app/merchanttax"
	"streetbox.id/app/merchantusers"
	"streetbox.id/app/merchantusersshift"
	"streetbox.id/app/role"
	"streetbox.id/app/tasks"
	"streetbox.id/app/user"
	"streetbox.id/app/userauth"
	"streetbox.id/app/userrole"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

var (
	xenditHeader http.Header
	xenditHost   string
)

// MerchantService ...
type MerchantService struct {
	MerchantRepo         merchant.RepoInterface
	MerchantUsersRepo    merchantusers.RepoInterface
	MerchantTaxRepo      merchanttax.RepoInterface
	MerchantCategoryRepo merchantcategory.RepoInterface
	UserRepo             user.RepoInterface
	RoleRepo             role.RepoInterface
	UserRoleRepo         userrole.RepoInterface
	UserAuthRepo         userauth.RepoInterface
	ShiftRepo            merchantusersshift.RepoInterface
	TasksRepo            tasks.RepoInterface
	MerchantMenusRepo    merchantmenu.RepoInterface
	LogMerchantRepo      logactivitymerchant.RepoInterface
	LogActivityRepo      logactivity.RepoInterface
}

// New ...
func New(
	merchantRepo merchant.RepoInterface,
	merchantUsersRepo merchantusers.RepoInterface,
	merchantTaxRepo merchanttax.RepoInterface,
	MerchantCategoryRepo merchantcategory.RepoInterface,
	userRepo user.RepoInterface,
	roleRepo role.RepoInterface,
	userRoleRepo userrole.RepoInterface,
	userAuthRepo userauth.RepoInterface,
	shiftRepo merchantusersshift.RepoInterface,
	tasksRepo tasks.RepoInterface,
	merchantMenusRepo merchantmenu.RepoInterface,
	logMerchantRepo logactivitymerchant.RepoInterface,
	logActivityRepo logactivity.RepoInterface,
	xenditAPIKey, xenditHosts string) merchant.ServiceInterface {

	data := xenditAPIKey + ":"
	xenditAuthHeader := "Basic " + bs64.StdEncoding.EncodeToString([]byte(data))
	xenditHeader = http.Header{
		"Authorization": []string{xenditAuthHeader},
		"Content-Type":  []string{"application/json"},
	}
	xenditHost = xenditHosts
	return &MerchantService{
		MerchantRepo:         merchantRepo,
		MerchantUsersRepo:    merchantUsersRepo,
		MerchantTaxRepo:      merchantTaxRepo,
		MerchantCategoryRepo: MerchantCategoryRepo,
		UserRepo:             userRepo,
		RoleRepo:             roleRepo,
		UserRoleRepo:         userRoleRepo,
		UserAuthRepo:         userAuthRepo,
		ShiftRepo:            shiftRepo,
		TasksRepo:            tasksRepo,
		MerchantMenusRepo:    merchantMenusRepo,
		LogMerchantRepo:      logMerchantRepo,
		LogActivityRepo:      logActivityRepo}
}

func (s *MerchantService) CreateCategory(cat *entity.MerchantCategory) (err error) {
	return s.MerchantCategoryRepo.Create(cat)
}

func (s *MerchantService) GetAllCategory() (cats []entity.MerchantCategory, err error) {
	return s.MerchantCategoryRepo.GetAll()
}

func (s *MerchantService) UpdateCategory(cat *entity.MerchantCategory) (err error) {
	return s.MerchantCategoryRepo.Update(cat)
}

func (s *MerchantService) DeleteCategory(id int64) (err error) {
	return s.MerchantCategoryRepo.Delete(id)
}

// CreateMerchant ...
func (s *MerchantService) CreateMerchant(
	req *model.ReqCreateMerchant, usersID int64) (*entity.Merchant, error) {
	data := new(entity.Merchant)
	copier.Copy(&data, req)
	var (
		db  *gorm.DB
		err error
	)
	if db, err = s.MerchantRepo.Create(data); err != nil {
		db.Rollback()
		return nil, err
	}
	if err = s.MerchantUsersRepo.Create(db, data.ID, usersID); err != nil {
		db.Rollback()
		return nil, err
	}

	reqData := map[string]interface{}{
		"account_email": data.Email,
		"type":          "OWNED",
		"business_profile": map[string]string{
			"business_name": data.Name,
		},
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return data, err
	}

	postReq, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts", xenditHost), bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return data, err
	}
	postReq.Header = xenditHeader
	var respBody []byte
	if respBody, err = util.DoRequest(postReq); err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return data, err
	}
	var respData map[string]string
	if err = json.Unmarshal(respBody, &respData); err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return data, err
	}
	db.Commit()
	upd := new(entity.Merchant)
	upd.XenditID = respData["user_id"]
	if err = s.MerchantRepo.Update(upd, data.ID); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return data, err
	}
	s.MerchantTaxRepo.GetTax(data.ID)
	return data, nil
}

// XenditGenerateSubAccount ...
func (s *MerchantService) XenditGenerateSubAccount(
	req *model.ReqXenditGenerateSubAccount, usersID int64) (*entity.Merchant, error) {
	data := new(entity.Merchant)
	copier.Copy(&data, req)
	var (
		err error
	)

	reqData := map[string]interface{}{
		"account_email": data.Email,
		"type":          "OWNED",
		"business_profile": map[string]string{
			"business_name": data.Name,
		},
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return data, err
	}

	postReq, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts", xenditHost), bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return data, err
	}
	postReq.Header = xenditHeader
	var respBody []byte
	if respBody, err = util.DoRequest(postReq); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return data, err
	}
	var respData map[string]string
	if err = json.Unmarshal(respBody, &respData); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return data, err
	}
	upd := new(entity.Merchant)
	upd.XenditID = respData["user_id"]
	data.XenditID = respData["user_id"]
	if err = s.MerchantRepo.Update(upd, data.ID); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return data, err
	}
	return data, nil
}

// CreateFoodtruck ...
func (s *MerchantService) CreateFoodtruck(
	req *model.ReqCreateFoodtruck, merchantID int64) (*entity.Users, error) {
	user := new(entity.Users)
	copier.Copy(&user, req)
	userName := strings.Split(req.UserName, "@")
	user.Name = userName[0]
	password := userName[0]

	if db, err := s.UserRepo.Create(user); err == nil {
		userRole := new(entity.UsersRole)
		userAuth := new(entity.UsersAuth)

		userRole.UsersID = user.ID
		userRole.RoleID = s.RoleRepo.FindByName("foodtruck").ID

		userAuth.UserName = user.UserName
		userAuth.Password = util.HashPassword(password)
		if err := s.UserRoleRepo.Create(db, userRole); err == nil {
			if err := s.UserAuthRepo.Create(db, userAuth); err == nil {
				if s.MerchantUsersRepo.IsExist(user.ID) {
					return user, errors.New("UsersID already exist in MerchantUsers")
				}
				if err := s.MerchantUsersRepo.Create(
					db, merchantID, user.ID); err == nil {
					db.Commit()
					return user, nil
				}

			}
		}
	}
	return user, errors.New("Email Already Delete, Please Use Another Email")
}

// CreateMenu ...
func (s *MerchantService) CreateMenu(req *model.ReqCreateMerchantMenu, merchantID int64) (*entity.MerchantMenu, error) {
	menu := new(entity.MerchantMenu)
	copier.Copy(&menu, req)
	menu.MerchantID = merchantID
	if err := s.MerchantMenusRepo.CreateMenu(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

// CreateMenuWithImage ...
func (s *MerchantService) CreateMenuWithImage(req *model.ReqCreateMerchantMenu, image string, merchantID int64) (*entity.MerchantMenu, error) {
	menu := new(entity.MerchantMenu)
	copier.Copy(&menu, req)
	menu.MerchantID = merchantID
	menu.Photo = image
	if err := s.MerchantMenusRepo.CreateMenu(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

// CreateTax ...
func (s *MerchantService) CreateTax(req *model.MerchantTax, merchantID int64) (*entity.MerchantTax, error) {
	tax := new(entity.MerchantTax)
	copier.Copy(&tax, req)
	tax.MerchantID = merchantID
	return s.MerchantTaxRepo.Create(tax)
}

// UpdateTax ...
func (s *MerchantService) UpdateTax(req *model.MerchantTax, ID int64, merchantID int64) (*entity.MerchantTax, error) {
	tax := new(entity.MerchantTax)
	copier.Copy(&tax, req)
	tax.MerchantID = merchantID
	tax.ID = ID
	return s.MerchantTaxRepo.Update(tax, merchantID, ID)
}

// UpdateMenu ...
func (s *MerchantService) UpdateMenu(req *model.ReqUpdateMerchantMenu, merchantID int64, ID int64) (*entity.MerchantMenu, error) {
	menu := new(entity.MerchantMenu)
	copier.Copy(&menu, req)
	menu.ID = ID
	menu.MerchantID = merchantID
	if err := s.MerchantMenusRepo.Update(menu, merchantID, ID, false); err != nil {
		return nil, err
	}
	return menu, nil
}

// GetMenuPagination ...
func (s *MerchantService) GetMenuPagination(merchantID int64, limit int, page int, sort []string) model.Pagination {
	data, count, offset := s.MerchantMenusRepo.GetAllMenu(merchantID, limit, page, sort)
	totalPages := util.TotalPages(count, limit)
	model := model.Pagination{
		Data:         data,
		Limit:        limit,
		NextPage:     util.NextPage(page, totalPages),
		Offset:       offset,
		Page:         page,
		PrevPage:     util.PrevPage(page),
		TotalPages:   totalPages,
		TotalRecords: count,
	}
	return model
}

// GetMenuList ...
func (s *MerchantService) GetMenuList(merchantID int64, nearby, visit bool) *[]entity.MerchantMenu {
	return s.MerchantMenusRepo.GetListMenu(merchantID, nearby, visit)
}

// GetMenuByID ...
func (s *MerchantService) GetMenuByID(merchantID int64, ID int64) *entity.MerchantMenu {
	return s.MerchantMenusRepo.GetMenuByID(merchantID, ID)
}

// GetAllFoodtruck ...
func (s *MerchantService) GetAllFoodtruck(merchantID int64) *[]entity.Users {
	return s.MerchantRepo.GetAllFoodtruck(merchantID)
}

// GetInfo ...
func (s *MerchantService) GetInfo(usersID int64) *model.Merchant {
	return s.MerchantRepo.GetInfo(usersID)
}

// IsExist ...
func (s *MerchantService) IsExist(usersID int64) bool {
	return s.MerchantUsersRepo.IsExist(usersID)
}

// CreateShift ...
func (s *MerchantService) CreateShift(
	usersID int64, shift string) (*entity.MerchantUsersShift, error) {
	merchantUsersID := s.MerchantUsersRepo.GetByUsersID(usersID).ID
	data := new(entity.MerchantUsersShift)
	data.Shift = shift
	data.MerchantUsersID = merchantUsersID
	return s.ShiftRepo.Create(data)
}

// Update ...
func (s *MerchantService) Update(
	req *model.ReqUpdateMerchant, id int64) (*entity.Merchant, error) {
	data := new(entity.Merchant)
	copier.Copy(&data, req)
	if err := s.MerchantRepo.Update(data, id); err != nil {
		return nil, err
	}
	return data, nil
}

// GetAll ..
func (s *MerchantService) GetAll() *[]model.Merchant {
	return s.MerchantRepo.GetAll()
}

// DeleteByMerchantID ...
func (s *MerchantService) DeleteByMerchantID(id int64, UserID int64) error {
	var (
		db  *gorm.DB
		err error
	)
	merchant := s.MerchantRepo.GetByID(id).Name
	userMerchants := s.MerchantUsersRepo.GetUserIdsByMerchantID(id)
	userNames := s.UserRepo.FindUsernameByMultipleID(*userMerchants)
	if db, err = s.MerchantRepo.DeleteByID(id); err != nil {
		return err
	}

	if userMerchants != nil {
		if err = s.MerchantUsersRepo.DeleteByMerchantID(db, id); err != nil {
			return err
		}
		if err = s.UserRoleRepo.DeleteByMultipleID(db, *userMerchants); err != nil {
			return err
		}
		if db, err = s.UserRepo.DeleteByMultipleID(db, *userMerchants); err != nil {
			return err
		}
		if userNames != nil {
			if db, err = s.UserAuthRepo.DeleteByMultipleUsername(db, *userNames); err != nil {
				return err
			}
		}
	}
	db.Commit()
	userName := s.UserRepo.FindByID(UserID).UserName
	msg := fmt.Sprintf("Delete Business Merchant %s by %s", merchant, userName)
	s.LogActivityRepo.Create(msg)
	return nil
}

// GetByID ...
func (s *MerchantService) GetByID(id int64) *entity.Merchant {
	return s.MerchantRepo.GetByID(id)
}

// IsUsersShiftIn ...
func (s *MerchantService) IsUsersShiftIn(usersID int64) bool {
	return s.ShiftRepo.IsUsersShiftIn(usersID)
}

// UpdateFoodtruck by Admin
func (s *MerchantService) UpdateFoodtruck(
	req *model.ReqUserUpdate, id int64) (*entity.Users, error) {
	data := new(entity.Users)
	copier.Copy(&data, req)
	if err := s.UserRepo.Update(data, id); err != nil {
		return nil, err
	}
	return data, nil
}

// UploadLogo ..
func (s *MerchantService) UploadLogo(
	logo string, id int64) error {
	data := new(entity.Merchant)
	data.Logo = logo
	if err := s.MerchantRepo.Update(data, id); err != nil {
		return err
	}
	return nil
}

// UploadMenu ..
func (s *MerchantService) UploadMenu(
	photo string, merchantID int64, ID int64) error {
	data := new(entity.MerchantMenu)
	data.Photo = photo
	data.MerchantID = merchantID
	data.ID = ID
	if err := s.MerchantMenusRepo.Update(data, merchantID, ID, true); err != nil {
		return err
	}
	return nil
}

// UploadBanner ..
func (s *MerchantService) UploadBanner(
	banner string, id int64) error {
	data := new(entity.Merchant)
	data.Banner = banner
	if err := s.MerchantRepo.Update(data, id); err != nil {
		return err
	}
	return nil
}

// GetFoodtruckByID ..
func (s *MerchantService) GetFoodtruckByID(id int64) *model.ResUserAll {
	return s.UserRepo.FindByID(id)
}

// DeleteFoodTruckByID first call user_service.DeleteByID
// to exec soft delete at users, users_role and users_auth
// and then exec this method to soft delete merchant_users
func (s *MerchantService) DeleteFoodTruckByID(id int64) error {
	return s.MerchantUsersRepo.DeleteByFoodtruckID(id)
}

// GetFoodtruckTasks ..
func (s *MerchantService) GetFoodtruckTasks(usersID int64) *[]model.ResGetFoodtruckTasks {
	merchantUser := s.MerchantUsersRepo.GetByUsersID(usersID)
	data := make([]model.ResGetFoodtruckTasks, 0)
	foodtrucks := s.GetAllFoodtruck(merchantUser.MerchantID)
	for _, v := range *foodtrucks {
		ft := new(model.ResGetFoodtruckTasks)
		copier.Copy(&ft, v)
		tasks := s.TasksRepo.FindByStatusUsersID(v.ID, 2)
		if tasks != nil {
			ft.TasksID = tasks.ID
			ft.Status = tasks.Status
		}
		data = append(data, *ft)
	}
	return &data
}

// DeleteMenu ..
func (s *MerchantService) DeleteMenu(ID int64) error {
	if err := s.MerchantMenusRepo.Delete(ID); err != nil {
		return err
	}
	return nil
}

// CountFoodtruckByMerchantID ..
func (s *MerchantService) CountFoodtruckByMerchantID(merchantID int64) int {
	return s.MerchantUsersRepo.CountFoodtruckByMerchantID(merchantID)
}

// GetTax ...
func (s *MerchantService) GetTax(merchantID int64) *entity.MerchantTax {
	return s.MerchantTaxRepo.GetTax(merchantID)
}

// RegistrationToken ..
func (s *MerchantService) RegistrationToken(token string, usersID int64) error {
	merchantUsers := s.MerchantUsersRepo.GetByUsersID(usersID)
	data := new(entity.MerchantUsers)
	data.RegistrationToken = token
	return s.MerchantUsersRepo.Update(data, merchantUsers.ID)
}

// GetMerchantUsersByUsersID ..
func (s *MerchantService) GetMerchantUsersByUsersID(id int64) *entity.MerchantUsers {
	return s.MerchantUsersRepo.GetByUsersID(id)
}

// GetMerchantUsersAdminByMerchantID ..
func (s *MerchantService) GetMerchantUsersAdminByMerchantID(id int64) *entity.MerchantUsers {
	return s.MerchantUsersRepo.GetAdminByMerchantID(id)
}

// GetMerchantUsersByID ..
func (s *MerchantService) GetMerchantUsersByID(id int64) *entity.MerchantUsers {
	return s.MerchantUsersRepo.GetOne(id)
}

// RemoveImage method for delete logo or banner
// Types : logo/banner
func (s *MerchantService) RemoveImage(filename, types string, merchantID int64) error {
	return s.MerchantRepo.RemoveImage(types, merchantID)
}

// RemoveImageMenu ...
func (s *MerchantService) RemoveImageMenu(menu *entity.MerchantMenu, ID int64) error {
	return s.MerchantMenusRepo.DeleteImageMenu(menu, ID)
}

// CheckStock ...
func (s *MerchantService) CheckStock(reqProductSales []model.TrxOrderProductSales) error {
	for _, productSales := range reqProductSales {
		if err := s.MerchantMenusRepo.CekStock(productSales.MerchantMenuID, productSales.Qty); !err {
			return errors.New(productSales.Name + " out of stock")
		}
	}
	return nil
}
