package service

import (
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/model"
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/repository/couchdb"
)


type UserService struct {
	repo couchdb.UserRepository
}

func NewUserService(repo couchdb.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (userService *UserService) CreateUsers(users []model.User) error {
	for _, user := range users {
		if err := userService.repo.Create([]model.User{user}); err != nil {
			return err 
		}
	}
	return nil
}


func (userService *UserService) UpdateUsers(id string ,users []model.User) error {
	for _, user := range users {
		if err := userService.repo.Update([]model.User{user}); err != nil {
			return err 
		}
	}
	return nil
}

func (userService *UserService) DeleteUsers(id string) error {
	
	if err := userService.repo.Delete(id); err != nil {
		return err 
	}
	return nil
}

func (s *UserService) GetUserById(id string) (map[string]*model.User, error) {
	user, err := s.repo.Read(id)
	if err != nil {
		return nil, err
	}
	return map[string]*model.User{id: user}, nil
}

func (userService *UserService) GetAllUsers() ([]model.User, error) {
	return userService.repo.ReadAll()
}
