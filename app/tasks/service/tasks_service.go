package service

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/copier"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/merchantusers"
	"streetbox.id/app/tasks"
	"streetbox.id/app/taskshomevisit"
	"streetbox.id/app/tasksnonreglog"
	"streetbox.id/app/tasksnonregular"
	"streetbox.id/app/tasksreglog"
	"streetbox.id/app/tasksregular"
	"streetbox.id/app/taskstracking"
	"streetbox.id/app/tasksvisitlog"
	"streetbox.id/app/trxvisitsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TasksService ..
type TasksService struct {
	MerchantUsersRepo merchantusers.RepoInterface
	RegularRepo       tasksregular.RepoInterface
	HomevisitRepo     taskshomevisit.RepoInterface
	NonRegularRepo    tasksnonregular.RepoInterface
	TasksRepo         tasks.RepoInterface
	RegLogRepo        tasksreglog.RepoInterface
	VisitLogRepo      tasksvisitlog.RepoInterface
	NonRegLogRepo     tasksnonreglog.RepoInterface
	TrackingRepo      taskstracking.RepoInterface
	VisitSalesRepo    homevisitsales.RepoInterface
	TrxVisitRepo      trxvisitsales.RepoInterface
}

// New ..
func New(merchantUsersRepo merchantusers.RepoInterface,
	regRepo tasksregular.RepoInterface,
	homeRepo taskshomevisit.RepoInterface,
	nonRegRepo tasksnonregular.RepoInterface,
	tasksRepo tasks.RepoInterface,
	regLogRepo tasksreglog.RepoInterface,
	visitLogRepo tasksvisitlog.RepoInterface,
	nonRegLogRepo tasksnonreglog.RepoInterface,
	trackingRepo taskstracking.RepoInterface,
	visitSalesRepo homevisitsales.RepoInterface,
	trxVisitRepo trxvisitsales.RepoInterface) tasks.ServiceInterface {
	return &TasksService{merchantUsersRepo,
		regRepo, homeRepo, nonRegRepo,
		tasksRepo, regLogRepo, visitLogRepo,
		nonRegLogRepo, trackingRepo, visitSalesRepo, trxVisitRepo}
}

// CreateTasksRegular ..
func (s *TasksService) CreateTasksRegular(
	req *model.ReqCreateTasksRegular) (*entity.Tasks, error) {
	tasks := new(entity.Tasks)
	merchantUserID := s.MerchantUsersRepo.GetByUsersID(req.UsersID).ID
	tasks.MerchantUsersID = merchantUserID
	tasks.Types = util.TasksRegular
	tasks.Status = 1
	if db, err := s.TasksRepo.Create(tasks); err == nil {
		tasksRegular := new(entity.TasksRegular)
		tasksRegular.TrxParkingSpaceSalesID = req.TrxSalesID
		tasksRegular.TasksID = tasks.ID
		tasksRegular.ScheduleDate = req.ScheduleDate
		if err := s.RegularRepo.Create(db, tasksRegular); err == nil {
			db.Commit()
			res := s.TasksRepo.FindOne(tasks.ID)
			return res, nil
		}
		db.Rollback()
	}
	return nil, errors.New("Failed")
}

// IsTasksRegularAssigned check if tasks regular exist
func (s *TasksService) IsTasksRegularAssigned(
	salesID int64, scheduleDate time.Time, foodtruckID int64) bool {
	return s.RegularRepo.IsAssigned(salesID, scheduleDate, foodtruckID)
}

// MyTasksRegByUsersID get by foodtruck
func (s *TasksService) MyTasksRegByUsersID(usersID int64) *[]model.ResMyTasksReg {
	data := make([]model.ResMyTasksReg, 0)
	regular := s.RegularRepo.MyTasks(usersID)
	homevisit := s.HomevisitRepo.MyTasks(usersID)
	if len(*regular) > 0 {
		for _, v := range *regular {
			data = append(data, v)
		}
	}
	if len(*homevisit) > 0 {
		for _, v := range *homevisit {
			data = append(data, v)
		}
	}
	return &data
}

// CreateTasksHomevisit ..
func (s *TasksService) CreateTasksHomevisit(
	req *model.ReqCreateTasksHomevisit) (*entity.Tasks, error) {
	tasks := new(entity.Tasks)
	merchantUserID := s.MerchantUsersRepo.GetByUsersID(req.UsersID).ID
	tasks.MerchantUsersID = merchantUserID
	tasks.Types = util.TasksHomeVisit
	tasks.Status = 1
	if db, err := s.TasksRepo.Create(tasks); err == nil {
		tasksHomevisit := new(entity.TasksHomevisit)
		tasksHomevisit.TrxHomevisitSalesID = req.TrxHomevisitSalesID
		tasksHomevisit.TasksID = tasks.ID
		if err := s.HomevisitRepo.Create(db, tasksHomevisit); err == nil {
			db.Commit()
			res := s.TasksRepo.FindOne(tasks.ID)
			return res, nil
		}
	}
	return nil, errors.New("Failed")
}

// IsTasksHomevisitAssigned ..
func (s *TasksService) IsTasksHomevisitAssigned(trxVisitSalesID, foodtruckID int64) bool {
	return s.HomevisitRepo.IsAssigned(trxVisitSalesID)
}

// UpdateTasksStatus ..
func (s *TasksService) UpdateTasksStatus(tasksID int64, status int) error {
	tasks := new(entity.Tasks)
	tasks.Status = status
	var err error
	if err = s.TasksRepo.UpdateByID(tasks, tasksID); err == nil {
		if status >= 3 {
			s.TrackingRepo.DeleteByTasksID(tasksID)
		}
	}
	return err
}

// CreateRegLog ..
func (s *TasksService) CreateRegLog(
	req *model.ReqTasksRegLog, activity int) (*entity.TasksRegularLog, error) {
	data := new(entity.TasksRegularLog)
	copier.Copy(&data, req)
	data.Activity = activity
	if err := s.RegLogRepo.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}

// CreateVisitLog ..
func (s *TasksService) CreateVisitLog(
	req *model.ReqTasksVisitLog, activity int) (*entity.TasksHomevisitLog, error) {
	data := new(entity.TasksHomevisitLog)
	copier.Copy(&data, req)
	data.Activity = activity
	if err := s.VisitLogRepo.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}

// CreateTasksNonRegular ..
func (s *TasksService) CreateTasksNonRegular(usersID int64) (*entity.Tasks, error) {
	tasks := new(entity.Tasks)
	merchantUserID := s.MerchantUsersRepo.GetByUsersID(usersID).ID
	tasks.MerchantUsersID = merchantUserID
	tasks.Types = util.TasksNonRegular
	tasks.Status = 2
	if db, err := s.TasksRepo.Create(tasks); err == nil {
		tasksNonReg := new(entity.TasksNonregular)
		tasksNonReg.TasksID = tasks.ID
		if err := s.NonRegularRepo.Create(db, tasksNonReg); err == nil {
			db.Commit()
			res := s.TasksRepo.FindOne(tasks.ID)
			return res, nil

		}
		db.Rollback()
	}
	return nil, errors.New("Failed")
}

// CreateNonRegLog ..
func (s *TasksService) CreateNonRegLog(
	req *model.ReqTasksNonRegLog, activity int) (*entity.TasksNonregularLog, error) {
	data := new(entity.TasksNonregularLog)
	copier.Copy(&data, req)
	nonRegular := s.NonRegularRepo.GetByTasksID(req.TasksID)
	data.TasksNonregularID = nonRegular.ID
	data.Activity = activity
	if err := s.NonRegLogRepo.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}

// CreateTracking ..
func (s *TasksService) CreateTracking(
	req *model.ReqCreateTasksTracking) (*model.ResTasksTracking, error) {
	data := new(entity.TasksTracking)
	res := new(model.ResTasksTracking)
	copier.Copy(&data, req)
	data.LogTime = time.Now()
	if err := s.TrackingRepo.Create(data); err != nil {
		return nil, err
	}
	copier.Copy(&res, data)
	res.Status = s.GetTasksByID(req.TasksID).Status
	return res, nil
}

// GetTasksByID ..
func (s *TasksService) GetTasksByID(id int64) *entity.Tasks {
	return s.TasksRepo.FindOne(id)
}

// GetTrackingByTasksID ..
func (s *TasksService) GetTrackingByTasksID(id int64) *model.ResTasksTracking {
	return s.TrackingRepo.GetTrackingByID(id)
}

// MyTasksNonRegByUsersID ..
func (s *TasksService) MyTasksNonRegByUsersID(usersID int64) *model.ResMyTasksNonReg {
	return s.NonRegularRepo.GetTasksByUsersID(usersID)
}

// GetTasksNonRegByUsersID ..
func (s *TasksService) GetTasksNonRegByUsersID(id int64) *entity.Tasks {
	return s.TasksRepo.FindNonRegByUsersID(id)
}

// GetAllByMerchantID ..
func (s *TasksService) GetAllByMerchantID(id int64) *[]model.ResMyTasksReg {
	data := make([]model.ResMyTasksReg, 0)
	regular := s.RegularRepo.MyTasksByMerchantID(id)
	homevisit := s.HomevisitRepo.MyTasksByMerchantID(id)
	nonReg := s.NonRegularRepo.MyTasksByMerchantID(id)
	if len(*regular) > 0 {
		for _, v := range *regular {
			data = append(data, v)
		}
	}
	if len(*homevisit) > 0 {
		for _, v := range *homevisit {
			data = append(data, v)
		}
	}
	if len(*nonReg) > 0 {
		for _, v := range *nonReg {
			data = append(data, v)
		}
	}
	return &data
}

// RegularStatusCompletedByUsersID ..
func (s *TasksService) RegularStatusCompletedByUsersID(id int64) {
	merchantUsersID := s.MerchantUsersRepo.GetByUsersID(id).ID
	tasksList := s.TasksRepo.FindUncompletedRegByMerchantUsersID(merchantUsersID)
	if len(*tasksList) > 0 {
		for _, v := range *tasksList {
			if err := s.UpdateTasksStatus(v.ID, 4); err == nil {
				log.Printf("INFO: Auto Completed Tasks Regular with TasksID %d", v.ID)
				s.TrackingRepo.DeleteByTasksID(v.ID)
			}
		}
	}
	log.Printf("INFO: All Tasks Already Completed")
}

// VisitStatusCompletedByUsersID ..
func (s *TasksService) VisitStatusCompletedByUsersID(id int64) {
	merchantUsersID := s.MerchantUsersRepo.GetByUsersID(id).ID
	tasksList := s.TasksRepo.FindUncompletedVisitByMerchantUsersID(merchantUsersID)
	if len(*tasksList) > 0 {
		for _, v := range *tasksList {
			if err := s.UpdateTasksStatus(v.ID, 4); err == nil {
				data := entity.TrxHomevisitSales{Status: util.TrxVisitStatusClosed}
				// update status to closed
				s.VisitSalesRepo.UpdateByID(&data, v.TrxHomevisitSalesID)
				log.Printf("INFO: Auto Completed Tasks Home Visit with TasksID %d", v.ID)
				s.TrackingRepo.DeleteByTasksID(v.ID)
			}
		}
	}
	log.Printf("INFO: All Tasks Already Completed")
}

// NonRegStatusCompletedByUsersID ..
func (s *TasksService) NonRegStatusCompletedByUsersID(id int64) {
	merchantUsersID := s.MerchantUsersRepo.GetByUsersID(id).ID
	tasksList := s.TasksRepo.FindUncompletedNonRegByMerchantUsersID(merchantUsersID)
	if len(*tasksList) > 0 {
		for _, v := range *tasksList {
			if err := s.UpdateTasksStatus(v.ID, 4); err == nil {
				log.Printf("INFO: Auto Completed Tasks Non Regular with TasksID %d", v.ID)
				s.TrackingRepo.DeleteByTasksID(v.ID)
			}
		}
	}
	log.Printf("INFO: All Tasks Already Completed")
}

// ClosedTrxHomevisitSales after checkout
func (s *TasksService) ClosedTrxHomevisitSales(tasksHomevisitID int64) error {
	qry := entity.TasksHomevisit{ID: tasksHomevisitID}
	tasksVisit := s.HomevisitRepo.FindOne(&qry)
	log.Printf("INFO: tasksVisit -> %+v", tasksVisit)
	updateTrx := entity.TrxHomevisitSales{Status: util.TrxVisitStatusClosed}
	return s.TrxVisitRepo.UpdateByID(&updateTrx, tasksVisit.TrxHomevisitSalesID)
}
