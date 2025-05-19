package usecase

import (
	"errors"
	"gatherly-app/models"
	"gatherly-app/models/dto"
	"gatherly-app/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	CreateUser(input dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(id int) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	UpdateUser(id int, input dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id int) error
}

type userUsecase struct {
	repo repositories.UserRepository
}

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) CreateUser(input dto.CreateUserRequest) (*dto.UserResponse, error) {
	if _, err := uc.repo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("email already exists")
	}
	if _, err := uc.repo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     input.Name,
		Age:      input.Age,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	err = uc.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Age:      user.Age,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (uc *userUsecase) GetUserByID(id int) (*dto.UserResponse, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Age:      user.Age,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (uc *userUsecase) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			ID:       u.ID,
			Name:     u.Name,
			Age:      u.Age,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		})
	}
	return res, nil
}

func (uc *userUsecase) UpdateUser(id int, input dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != "" && input.Username != user.Username {
		if _, err := uc.repo.FindByUsername(input.Username); err == nil {
			return nil, errors.New("username already taken")
		}
	}

	if input.Email != "" && input.Email != user.Email {
		if _, err := uc.repo.FindByEmail(input.Email); err == nil {
			return nil, errors.New("email already taken")
		}
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Age != 0 {
		user.Age = input.Age
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	err = uc.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Age:      user.Age,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (uc *userUsecase) DeleteUser(id int) error {
	return uc.repo.Delete(id)
}
