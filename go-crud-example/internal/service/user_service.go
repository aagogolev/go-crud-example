package service

import (
    "database/sql"
    "errors"
    "fmt"
    "go-crud-example/internal/model"
    "go-crud-example/internal/repository"
)

type UserService interface {
    GetUsers() ([]model.User, error)
    GetUser(id string) (*model.User, error)
    CreateUser(user *model.User) error
    UpdateUser(user *model.User) error
    DeleteUser(id string) error
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{
        repo: repo,
    }
}

func (s *userService) GetUsers() ([]model.User, error) {
    users, err := s.repo.GetAll()
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }
    return users, nil
}

func (s *userService) GetUser(id string) (*model.User, error) {
    user, err := s.repo.GetByID(id)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}

func (s *userService) CreateUser(user *model.User) error {
    if err := user.Validate(); err != nil {
        return fmt.Errorf("validation error: %w", err)
    }

    if err := s.repo.Create(user); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}

func (s *userService) UpdateUser(user *model.User) error {
    if err := user.Validate(); err != nil {
        return fmt.Errorf("validation error: %w", err)
    }

    if err := s.repo.Update(user); err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("user not found: %w", err)
        }
        return fmt.Errorf("failed to update user: %w", err)
    }
    return nil
}

func (s *userService) DeleteUser(id string) error {
    if err := s.repo.Delete(id); err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("user not found: %w", err)
        }
        return fmt.Errorf("failed to delete user: %w", err)
    }
    return nil
}

// Определение пользовательских ошибок
var (
    ErrEmptyName  = errors.New("name cannot be empty")
    ErrInvalidAge = errors.New("age must be positive")
)
