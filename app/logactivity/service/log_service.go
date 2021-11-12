package service

import (
	"streetbox.id/app/logactivity"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// LogService ..
type LogService struct {
	Repo logactivity.RepoInterface
}

// New ..
func New(repo logactivity.RepoInterface) logactivity.ServiceInterface {
	return &LogService{repo}
}

// GetAll ..
func (s *LogService) GetAll(limit int, page int, sort []string) model.Pagination {
	data, count, offset := s.Repo.GetAllPagination(limit, page, sort)
	totalPages := util.TotalPages(count, limit)
	model := model.Pagination{
		Data:         data,
		Offset:       offset,
		TotalPages:   totalPages,
		Limit:        limit,
		NextPage:     util.NextPage(page, totalPages),
		Page:         page,
		PrevPage:     util.PrevPage(page),
		TotalRecords: count,
	}
	return model
}

// GetList ...
func (s *LogService) GetList() *[]entity.LogActivity {
	return s.Repo.GetList()
}
