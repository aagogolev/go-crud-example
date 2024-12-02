package handler

import (
    "bytes"
    "encoding/json"
    "go-crud-example/internal/handler"
    "go-crud-example/internal/model"
    "go-crud-example/internal/repository"
    "go-crud-example/internal/service"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gorilla/mux"
    "log"
)

type mockUserService struct {
    users map[string]model.User
}

func (m *mockUserService) CreateUser(user *model.User) error {
    if user.Name == "" || user.Age == 0 {
        return service.ErrInvalidUser
    }
    user.ID = "1"
    m.users[user.ID] = *user
    return nil
}

func (m *mockUserService) GetUsers() ([]model.User, error) {
    users := make([]model.User, 0, len(m.users))
    for _, user := range m.users {
        users = append(users, user)
    }
    return users, nil
}

func (m *mockUserService) GetUser(id string) (*model.User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, repository.ErrUserNotFound
    }
    return &user, nil
}

func (m *mockUserService) UpdateUser(user *model.User) error {
    if _, exists := m.users[user.ID]; !exists {
        return repository.ErrUserNotFound
    }
    m.users[user.ID] = *user
    return nil
}

func (m *mockUserService) DeleteUser(id string) error {
    if _, exists := m.users[id]; !exists {
        return repository.ErrUserNotFound
    }
    delete(m.users, id)
    return nil
}

func setupTest() (*handler.UserHandler, *mockUserService) {
    mockService := &mockUserService{
        users: make(map[string]model.User),
    }
    logger := log.New(log.Writer(), "TEST: ", log.LstdFlags)
    return handler.NewUserHandler(mockService, logger), mockService
}

func TestUserHandler_CreateUser(t *testing.T) {
    tests := []struct {
        name       string
        user      model.User
        wantCode  int
    }{
        {
            name: "valid user",
            user: model.User{Name: "John", Age: 30},
            wantCode: http.StatusOK,
        },
        {
            name: "invalid user",
            user: model.User{},
            wantCode: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h, _ := setupTest()

            body, _ := json.Marshal(tt.user)
            req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
            w := httptest.NewRecorder()

            h.CreateUser(w, req)

            if w.Code != tt.wantCode {
                t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.wantCode)
            }
        })
    }
}

func TestUserHandler_GetUsers(t *testing.T) {
    h, mockService := setupTest()
    mockService.users["1"] = model.User{ID: "1", Name: "John", Age: 30}

    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()

    h.GetUsers(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
    }
}

func TestUserHandler_GetUser(t *testing.T) {
    tests := []struct {
        name     string
        userID   string
        wantCode int
    }{
        {
            name:     "existing user",
            userID:   "1",
            wantCode: http.StatusOK,
        },
        {
            name:     "non-existing user",
            userID:   "999",
            wantCode: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h, mockService := setupTest()
            mockService.users["1"] = model.User{ID: "1", Name: "John", Age: 30}

            req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
            req = mux.SetURLVars(req, map[string]string{"id": tt.userID})
            w := httptest.NewRecorder()

            h.GetUser(w, req)

            if w.Code != tt.wantCode {
                t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.wantCode)
            }
        })
    }
}

func TestUserHandler_UpdateUser(t *testing.T) {
    tests := []struct {
        name     string
        userID   string
        user     model.User
        wantCode int
    }{
        {
            name:     "valid update",
            userID:   "1",
            user:     model.User{ID: "1", Name: "John Updated", Age: 31},
            wantCode: http.StatusOK,
        },
        {
            name:     "non-existing user",
            userID:   "999",
            user:     model.User{ID: "999", Name: "Non Existing", Age: 25},
            wantCode: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h, mockService := setupTest()
            mockService.users["1"] = model.User{ID: "1", Name: "John", Age: 30}

            body, _ := json.Marshal(tt.user)
            req := httptest.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(body))
            req = mux.SetURLVars(req, map[string]string{"id": tt.userID})
            w := httptest.NewRecorder()

            h.UpdateUser(w, req)

            if w.Code != tt.wantCode {
                t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.wantCode)
            }
        })
    }
}

func TestUserHandler_DeleteUser(t *testing.T) {
    tests := []struct {
        name     string
        userID   string
        wantCode int
    }{
        {
            name:     "existing user",
            userID:   "1",
            wantCode: http.StatusNoContent,
        },
        {
            name:     "non-existing user",
            userID:   "999",
            wantCode: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h, mockService := setupTest()
            mockService.users["1"] = model.User{ID: "1", Name: "John", Age: 30}

            req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
            req = mux.SetURLVars(req, map[string]string{"id": tt.userID})
            w := httptest.NewRecorder()

            h.DeleteUser(w, req)

            if w.Code != tt.wantCode {
                t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.wantCode)
            }
        })
    }
}
