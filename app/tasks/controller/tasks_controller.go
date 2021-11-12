package controller

import (
	"fmt"
	"net/http"

	"streetbox.id/app/logactivitymerchant"
	"streetbox.id/app/merchant"
	"streetbox.id/app/tasks"
	"streetbox.id/app/user"
	"streetbox.id/util"

	"streetbox.id/entity"

	"github.com/gin-gonic/gin"
	"streetbox.id/model"
)

// TasksController ..
type TasksController struct {
	TasksSvc    tasks.ServiceInterface
	MerchantSvc merchant.ServiceInterface
	LogSvc      logactivitymerchant.ServiceInterface
	UserSvc     user.ServiceInterface
}

// CreateTasksRegular godoc
// @Summary Create New Task Regular FoodTruck (permission = admin)
// @Id CreateTasksRegular
// @Tags Tasks
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param trxSalesId formData integer true " "
// @Param usersId formData integer true " "
// @Param scheduleDate formData string true " "
// @Success 200 {object} entity.Tasks "data: entity.Tasks"
// @Router /tasks/regular [post]
func (r *TasksController) CreateTasksRegular(c *gin.Context) {
	req := model.ReqCreateTasksRegular{}
	if err := c.ShouldBind(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// Check is admin and Foodtruck at same merchant
	merchantIDAdmin := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	merchantIDOps := r.MerchantSvc.GetInfo(req.UsersID).ID
	if merchantIDAdmin != merchantIDOps {
		model.ResponseError(c, "Foodtruck is not belong to your merchant",
			http.StatusUnprocessableEntity)
		return
	}
	// Check if Tasks Regular Exist
	if r.TasksSvc.IsTasksRegularAssigned(req.TrxSalesID, req.ScheduleDate, req.UsersID) {
		model.ResponseError(c, "Tasks Already Assigned",
			http.StatusUnprocessableEntity)
		return
	}
	data := new(entity.Tasks)
	var err error
	// Auto Completed Tasks with status < 4
	r.TasksSvc.RegularStatusCompletedByUsersID(req.UsersID)
	if data, err = r.TasksSvc.CreateTasksRegular(&req); err != nil {
		model.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}
	foodtruck := r.UserSvc.GetUserByID(req.UsersID)
	msg := fmt.Sprintf("Anda Assign Tasks Regular Pada %s", foodtruck.PlatNo)
	r.LogSvc.Create(merchantIDOps, msg)
	model.ResponseJSON(c, data)
	return
}

// CreateTasksHomevisit godoc
// @Summary Create New Task Home Visit FoodTruck (permission = merchant)
// @Id CreateTasksHomevisit
// @Tags Tasks
// @Security Token
// @Param task body model.ReqCreateTasksHomevisit true "all fields mandatory"
// @Success 200 {object} entity.Tasks "data: entity.Tasks"
// @Router /tasks/homevisit [post]
func (r *TasksController) CreateTasksHomevisit(c *gin.Context) {
	req := model.ReqCreateTasksHomevisit{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// Check is admin and Foodtruck at same merchant
	merchantIDAdmin := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	merchantIDOps := r.MerchantSvc.GetInfo(req.UsersID).ID
	if merchantIDAdmin != merchantIDOps {
		model.ResponseError(c, "Foodtruck is not belong to your merchant",
			http.StatusUnprocessableEntity)
		return
	}
	// Check if Tasks Regular Exist
	if r.TasksSvc.IsTasksHomevisitAssigned(
		req.TrxHomevisitSalesID, req.UsersID) {
		model.ResponseError(c, "Tasks Already Assigned",
			http.StatusUnprocessableEntity)
		return
	}
	data := new(entity.Tasks)
	var err error
	// Auto Completed Tasks with status < 4
	r.TasksSvc.VisitStatusCompletedByUsersID(req.UsersID)
	if data, err = r.TasksSvc.CreateTasksHomevisit(&req); err != nil {
		model.ResponseError(c, "Failed Create Tasks", http.StatusInternalServerError)
		return
	}
	foodtruck := r.UserSvc.GetUserByID(req.UsersID)
	msg := fmt.Sprintf("Anda Assign Tasks Homevisit Pada %s", foodtruck.PlatNo)
	r.LogSvc.Create(merchantIDOps, msg)
	model.ResponseJSON(c, data)
	return
}

// MyTaskRegular godoc
// @Summary Get Tasks Regular/Homevisit List Foodtruck (permission = merchant)
// @Id MyTaskRegular
// @Tags Tasks
// @Security Token
// @Success 200 {object} []model.ResMyTasksReg "data: model.ResMyTasksReg"
// @Router /tasks/regular/list [get]
func (r *TasksController) MyTaskRegular(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	data := r.TasksSvc.MyTasksRegByUsersID(jwtModel.UserID)
	model.ResponseJSON(c, data)
	return
}

// MyTaskNonRegular godoc
// @Summary Get Free Tasks List (permission = merchant)
// @Id MyTaskNonRegular
// @Tags Tasks
// @Security Token
// @Success 200 {object} model.ResMyTasksNonReg "data: model.ResMyTasksNonReg"
// @Router /tasks/nonregular/list [get]
func (r *TasksController) MyTaskNonRegular(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	data := r.TasksSvc.MyTasksNonRegByUsersID(jwtModel.UserID)
	model.ResponseJSON(c, data)
	return
}

// NonRegStatus godoc
// @Summary Get Free Tasks Status (permission = merchant)
// @Id NonRegStatus
// @Tags Tasks
// @Security Token
// @Success 200 {object} entity.Tasks "data: entity.Tasks"
// @Router /tasks/nonregular/status [get]
func (r *TasksController) NonRegStatus(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	tasks := r.TasksSvc.GetTasksNonRegByUsersID(jwtModel.UserID)
	if tasks != nil {
		model.ResponseJSON(c, tasks)
		return
	}
	model.ResponseJSON(c, "")
	return
}

// ChangeToOpen godoc
// @Summary Change Status Tasks from 2 to 1 (1=open,2=ongoing) (permission = merchant)
// @Id ChangeToOpen
// @Tags Tasks
// @Security Token
// @Param id path integer true "TasksID"
// @Success 200 {object} model.ResponseSuccess "data: "Success Update Tasks Status" "
// @Router /tasks/undo/{id} [put]
func (r *TasksController) ChangeToOpen(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	// Update Tasks Status -> 1 (open)
	if err := r.TasksSvc.UpdateTasksStatus(id, 1); err != nil {
		model.ResponseError(c, "Failed To Update Tasks Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	msg := fmt.Sprintf("%s Membatalkan Perjalanan Task", foodtruck.PlatNo)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, "Success Update Tasks Status")
	return
}

// ShiftIn godoc
// @Summary Create Task Log Shift In (permission = merchant)
// @Id ShiftIn
// @Tags Tasks
// @Security Token
// @Success 200 {object} entity.MerchantUsersShift "data: entity.MerchantUsersShift"
// @Router /tasks/shift-in [post]
func (r *TasksController) ShiftIn(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	var (
		data *entity.MerchantUsersShift
		err  error
	)
	if data, err = r.MerchantSvc.CreateShift(jwtModel.UserID, "IN"); err != nil {
		model.ResponseError(c, "Failed To Create Shift-In", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// ShiftOut godoc
// @Summary Create Task Log Shift Out (permission = merchant)
// @Id ShiftOut
// @Tags Tasks
// @Security Token
// @Success 200 {object} entity.MerchantUsersShift "{ "data": Model }"
// @Failure 300 {object} model.ResponseErrors "Redirect"
// @Failure 400 {object} model.ResponseErrors "Client Errors"
// @Failure 500 {object} model.ResponseErrors "Server Errors"
// @Router /tasks/shift-out [post]
func (r *TasksController) ShiftOut(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	var (
		data *entity.MerchantUsersShift
		err  error
	)
	if data, err = r.MerchantSvc.CreateShift(jwtModel.UserID, "OUT"); err != nil {
		model.ResponseError(c, "Failed To Create Shift Out", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// ShiftInStatus godoc
// @Summary Check Status Shift-In (permission = merchant)
// @Id ShiftInStatus
// @Tags Tasks
// @Security Token
// @Success 200 {object} model.ResponseSuccess "shiftIn : bool "
// @Router /tasks/shift-in/status [get]
func (r *TasksController) ShiftInStatus(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	model.ResponseJSON(c,
		gin.H{"shiftIn": r.MerchantSvc.IsUsersShiftIn(jwtModel.UserID)})
	return
}

// CreateRegularCheckIn godoc
// @Summary Create Tasks Regular Log CheckIn (permission = merchant)
// @Id CreateRegularCheckIn
// @Tags Tasks
// @Security Token
// @Param req body model.ReqTasksRegLog true "All Mandatory"
// @Success 200 {object} entity.TasksRegularLog "data : entity.TasksRegularLog"
// @Router /tasks/regular/check-in [post]
func (r *TasksController) CreateRegularCheckIn(c *gin.Context) {
	req := model.ReqTasksRegLog{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.TasksRegularLog
		err  error
	)
	if data, err = r.TasksSvc.CreateRegLog(&req, 1); err != nil {
		model.ResponseError(c, "Failed To Create Tasks Reguler Log", http.StatusInternalServerError)
		return
	}
	// Update Tasks Status -> 3 (arrived)
	if err := r.TasksSvc.UpdateTasksStatus(req.TasksID, 3); err != nil {
		model.ResponseError(c, "Failed To Update Task Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	msg := fmt.Sprintf("%s Check-In Task Regular Di Parking Space %s",
		foodtruck.PlatNo, req.ParkingSpaceName)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, data)
	return
}

// CreateRegularCheckOut godoc
// @Summary Create Tasks Regular Log CheckOut (permission = merchant)
// @Id CreateRegularCheckOut
// @Tags Tasks
// @Security Token
// @Param req body model.ReqTasksRegLog true "All Mandatory"
// @Success 200 {object} entity.TasksRegularLog "data : entity.TasksRegularLog"
// @Router /tasks/regular/check-out [post]
func (r *TasksController) CreateRegularCheckOut(c *gin.Context) {
	req := model.ReqTasksRegLog{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.TasksRegularLog
		err  error
	)
	if data, err = r.TasksSvc.CreateRegLog(&req, 2); err != nil {
		model.ResponseError(c, "Failed To Create Tasks Reguler Log", http.StatusInternalServerError)
		return
	}
	// Update Tasks Status -> 4 (completed)
	if err := r.TasksSvc.UpdateTasksStatus(req.TasksID, 4); err != nil {
		model.ResponseError(c, "Failed To Update Task Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	msg := fmt.Sprintf("%s Check-Out Task Regular Di Parking Space %s",
		foodtruck.PlatNo, req.ParkingSpaceName)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, data)
	return
}

// CreateHomeVisitCheckIn godoc
// @Summary Create Tasks HomeVisit Log CheckIn (permission = merchant)
// @Id CreateHomeVisitCheckIn
// @Tags Tasks
// @Security Token
// @Param req body model.ReqTasksVisitLog true "All Mandatory"
// @Success 200 {object} entity.TasksRegularLog "data : entity.TasksRegularLog"
// @Router /tasks/homevisit/check-in [post]
func (r *TasksController) CreateHomeVisitCheckIn(c *gin.Context) {
	req := model.ReqTasksVisitLog{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.TasksHomevisitLog
		err  error
	)
	if data, err = r.TasksSvc.CreateVisitLog(&req, 1); err != nil {
		model.ResponseError(c, "Failed To Create Tasks Reguler Log", http.StatusInternalServerError)
		return
	}
	// Update Tasks Status -> 3 (arrived)
	if err := r.TasksSvc.UpdateTasksStatus(req.TasksID, 3); err != nil {
		model.ResponseError(c, "Failed To Update Task Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	msg := fmt.Sprintf("%s Check-In Task Home Visit Pada Pelanggan %s",
		foodtruck.PlatNo, req.CustomerName)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, data)
	return
}

// CreateHomeVisitCheckOut godoc
// @Summary Create Tasks HomeVisit Log CheckOut (permission = merchant)
// @Id CreateHomeVisitCheckOut
// @Tags Tasks
// @Security Token
// @Param req body model.ReqTasksVisitLog true "All Mandatory"
// @Success 200 {object} entity.TasksRegularLog "data : entity.TasksRegularLog"
// @Router /tasks/homevisit/check-out [post]
func (r *TasksController) CreateHomeVisitCheckOut(c *gin.Context) {
	req := model.ReqTasksVisitLog{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.TasksHomevisitLog
		err  error
	)
	if data, err = r.TasksSvc.CreateVisitLog(&req, 2); err != nil {
		model.ResponseError(c, "Failed To Create Tasks Reguler Log", http.StatusInternalServerError)
		return
	}
	// Update Tasks Status -> 4 (completed)
	if err := r.TasksSvc.UpdateTasksStatus(req.TasksID, 4); err != nil {
		model.ResponseError(c, "Failed To Update Task Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	msg := fmt.Sprintf("%s Check-Out Task Home Visit Pada Pelanggan %s",
		foodtruck.PlatNo, req.CustomerName)
	r.LogSvc.Create(merchantID, msg)
	if err := r.TasksSvc.ClosedTrxHomevisitSales(req.TasksHomevisitID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// CreateNonRegCheckIn godoc
// @Summary Create Tasks Non Regular Log CheckIn (permission = merchant)
// @Id CreateNonRegCheckIn
// @Tags Tasks
// @Security Token
// @Param req body model.ReqTasksNonRegLog true "All Mandatory"
// @Success 200 {object} entity.Tasks "data : entity.Tasks"
// @Router /tasks/nonregular/check-in [post]
func (r *TasksController) CreateNonRegCheckIn(c *gin.Context) {
	req := model.ReqTasksNonRegLog{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		tasks *entity.TasksNonregularLog
		err   error
	)
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if tasks, err = r.TasksSvc.CreateNonRegLog(&req, 1); err != nil {
		model.ResponseError(c, "Failed To Create Tasks Non Regular", http.StatusInternalServerError)
		return
	}
	// Update Tasks Status -> 3 (arrived)
	if err := r.TasksSvc.UpdateTasksStatus(req.TasksID, 3); err != nil {
		model.ResponseError(c, "Failed To Update Task Status", http.StatusInternalServerError)
		return
	}
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	msg := fmt.Sprintf("%s Check-In Free Task Di %s",
		foodtruck.PlatNo, req.Address)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, tasks)
	return
}

// CreateNonRegCheckOut godoc
// @Summary Create Tasks Non Regular Log CheckOut (permission = merchant)
// @Id CreateNonRegCheckOut
// @Tags Tasks
// @Security Token
// @Param req body model.ReqTasksNonRegLog true "All Mandatory"
// @Success 200 {object} entity.TasksNonregularLog "data : entity.TasksNonregularLog"
// @Router /tasks/nonregular/check-out [post]
func (r *TasksController) CreateNonRegCheckOut(c *gin.Context) {
	req := model.ReqTasksNonRegLog{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.TasksNonregularLog
		err  error
	)
	if data, err = r.TasksSvc.CreateNonRegLog(&req, 2); err != nil {
		model.ResponseError(c, "Failed To Create Tasks Reguler Log", http.StatusInternalServerError)
		return
	}
	// Update Tasks Status -> 4 (completed)
	if err := r.TasksSvc.UpdateTasksStatus(req.TasksID, 4); err != nil {
		model.ResponseError(c, "Failed To Update Task Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	msg := fmt.Sprintf("%s Check-Out Free Task Di %s",
		foodtruck.PlatNo, req.Address)
	r.LogSvc.Create(merchantID, msg)
	r.TasksSvc.NonRegStatusCompletedByUsersID(jwtModel.UserID)
	model.ResponseJSON(c, data)
	return
}

// CreateTracking godoc
// @Summary Create Tasks Non Regular Log CheckOut (permission = merchant)
// @Id CreateTracking
// @Tags Tasks
// @Security Token
// @Param req body model.ReqCreateTasksTracking true "All Mandatory"
// @Success 200 {object} model.ResTasksTracking "data : model.ResTasksTracking"
// @Router /tasks/tracking [post]
func (r *TasksController) CreateTracking(c *gin.Context) {
	req := model.ReqCreateTasksTracking{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	var (
		data *model.ResTasksTracking
		err  error
	)
	if data, err = r.TasksSvc.CreateTracking(&req); err != nil {
		model.ResponseError(c, "Failed to Create Tracking", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// GetTracking godoc
// @Summary Create Tasks Non Regular Log CheckOut (permission = merchant)
// @Id GetTracking
// @Tags Tasks
// @Security Token
// @Param id path integer true "TasksID"
// @Success 200 {object} model.ResTasksTracking "data : model.ResTasksTracking"
// @Router /tasks/tracking/{id}/foodtruck [get]
func (r *TasksController) GetTracking(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, r.TasksSvc.GetTrackingByTasksID(id))
	return
}

// ChangeToOngoing godoc
// @Summary Change Status Tasks to Ongoing (permission = merchant)
// @Id ChangeToOngoing
// @Tags Tasks
// @Security Token
// @Param id path integer true "TasksID"
// @Success 200 {object} model.ResponseSuccess "data: "Success Update Tasks Status" "
// @Router /tasks/ongoing/{id} [put]
func (r *TasksController) ChangeToOngoing(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	// Update Tasks Status -> 2 (ongoing)
	if err := r.TasksSvc.UpdateTasksStatus(id, 2); err != nil {
		model.ResponseError(c, "Failed To Update Tasks Status", http.StatusInternalServerError)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	msg := fmt.Sprintf("%s Sedang Dalam Perjalanan Task", foodtruck.PlatNo)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, "Success Update Tasks Status")
	return
}

// CreateTasksNonRegular godoc
// @Summary Create New Task NonRegular FoodTruck (permission = merchant)
// @Id CreateTasksNonRegular
// @Tags Tasks
// @Security Token
// @Success 200 {object} entity.Tasks "data: entity.Tasks"
// @Router /tasks/nonregular [post]
func (r *TasksController) CreateTasksNonRegular(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	data := new(entity.Tasks)
	var err error
	// Auto Completed Tasks with status < 4
	r.TasksSvc.NonRegStatusCompletedByUsersID(jwtModel.UserID)
	if data, err = r.TasksSvc.CreateTasksNonRegular(jwtModel.UserID); err != nil {
		model.ResponseError(c, "Failed Create Tasks", http.StatusInternalServerError)
		return
	}
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	foodtruck := r.UserSvc.GetUserByID(jwtModel.UserID)
	msg := fmt.Sprintf("%s Memulai Free Task", foodtruck.PlatNo)
	r.LogSvc.Create(merchantID, msg)
	model.ResponseJSON(c, data)
	return
}

// GetAll godoc
// @Summary Get All Tasks Uncompleted (permission = admin)
// @Id GetAllTasks
// @Tags Tasks
// @Security Token
// @Success 200 {object} []model.ResMyTasksReg "data: []model.ResMyTasksReg"
// @Router /tasks/all [get]
func (r *TasksController) GetAll(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	data := r.TasksSvc.GetAllByMerchantID(merchantID)
	model.ResponseJSON(c, data)
}
