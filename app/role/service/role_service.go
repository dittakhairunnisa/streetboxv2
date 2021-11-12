package service

import (
	"github.com/jinzhu/copier"
	"streetbox.id/app/role"
	IRole "streetbox.id/app/role"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RoleService ...
type RoleService struct {
	Repo role.RepoInterface
}

// New ...
func New(s IRole.RepoInterface) role.ServiceInterface {
	return &RoleService{s}
}

// Create ...
func (s *RoleService) Create(req model.ReqRoleCreate) (*entity.Role, error) {
	role := new(entity.Role)
	copier.Copy(&role, &req)
	if err := s.Repo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

// SearchByName ...
func (s *RoleService) SearchByName(name string) *entity.Role {
	return s.Repo.FindByName(name)
}

// DeleteByID ...
func (s *RoleService) DeleteByID(id int64) error {
	return s.Repo.DeleteByID(id)
}

// GetAll ...
func (s *RoleService) GetAll() *[]entity.Role {
	return s.Repo.GetAll()
}

// GetAllExclude ...
func (s *RoleService) GetAllExclude() *[]entity.Role {
	return s.Repo.GetAllExclude()
}
