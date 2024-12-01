package service

import (
    "database/sql"
    "testing"
    "go-crud-example/internal/model"
    svc "go-crud-example/internal/service"
)
type mockRepository struct {
    users map[string]model.User
}

func newMockRepository() *mockRepository {
    return &mockRepository{
        users: make(map[string]model.User),
    }
}

func (m *mockRepository) GetAll() ([]model.User, error) {
    users := make([]model.User, 0, len(m.users))
    for _, user := range m.users {
        users = append(users, user)
    }
    return users, nil
}

func (m *mockRepository) GetByID(id string) (*model.User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, sql.ErrNoRows
    }
    return &user, nil
}

func (m *mockRepository) Create(user *model.User) error {
    user.ID = "1" // Для тестов используем фиксированный ID
    m.users[user.ID] = *user
    return nil
}

func (m *mockRepository) Update(user *model.User) error {
    if _, exists := m.users[user.ID]; !exists {
        return sql.ErrNoRows
    }
    m.users[user.ID] = *user
    return nil
}

func (m *mockRepository) Delete(id string) error {
    if _, exists := m.users[id]; !exists {
        return sql.ErrNoRows
    }
    delete(m.users, id)
    return nil
}

func TestUserService_CreateUser(t *testing.T) {
    repo := newMockRepository()
    service := svc.NewUserService(repo)

    tests := []struct {
        name    string
        user    model.User
        wantErr bool
    }{
        {
            name: "valid user",
            user: model.User{
                Name: "John Doe",
                Age:  25,
            },
            wantErr: false,
        },
        {
            name: "invalid user",
            user: model.User{
                Name: "",
                Age:  -1,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := service.CreateUser(&tt.user)
            if (err != nil) != tt.wantErr {
                t.Errorf("UserService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
