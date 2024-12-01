package handler

import (
    "bytes"
    "encoding/json"
    "github.com/gorilla/mux"
    h "go-crud-example/internal/handler"
    "go-crud-example/internal/model"
    "go-crud-example/internal/service"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

type mockService struct {
    users map[string]model.User
}

type testHandler struct {
    *h.UserHandler
    service *mockService
}

func newMockService() *mockService {
    return &mockService{
        users: make(map[string]model.User),
    }
}

func (m *mockService) GetUsers() ([]model.User, error) {
    users := make([]model.User, 0, len(m.users))
    for _, user := range m.users {
        users = append(users, user)
    }
    return users, nil
}

func (m *mockService) GetUser(id string) (*model.User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, service.ErrUserNotFound
    }
    return &user, nil
}

func (m *mockService) CreateUser(user *model.User) error {
    user.ID = "1" // Фиксированный ID для тестов
    m.users[user.ID] = *user
    return nil
}

func (m *mockService) UpdateUser(user *model.User) error {
    if _, exists := m.users[user.ID]; !exists {
        return service.ErrUserNotFound
    }
    m.users[user.ID] = *user
    return nil
}

func (m *mockService) DeleteUser(id string) error {
    if _, exists := m.users[id]; !exists {
        return service.ErrUserNotFound
    }
    delete(m.users, id)
    return nil
}

func setupTestHandler() *testHandler {
    logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
    mockSvc := newMockService()
    return &testHandler{
        UserHandler: h.NewUserHandler(mockSvc, logger),
        service:     mockSvc,
    }
}

func TestUserHandler_CreateUser(t *testing.T) {
    th := setupTestHandler()
    router := mux.NewRouter()
    th.RegisterRoutes(router)

    tests := []struct {
        name       string
        user       model.User
        wantStatus int
    }{
        {
            name: "valid user",
            user: model.User{
                Name: "John Doe",
                Age:  25,
            },
            wantStatus: http.StatusOK,
        },
        {
            name: "invalid user",
            user: model.User{
                Name: "",
                Age:  -1,
            },
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.user)
            req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.wantStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.wantStatus)
            }
        })
    }
}

func TestUserHandler_GetUsers(t *testing.T) {
    th := setupTestHandler()
    router := mux.NewRouter()
    th.RegisterRoutes(router)

    req := httptest.NewRequest("GET", "/users", nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    var users []model.User
    if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
        t.Errorf("Failed to decode response body: %v", err)
    }
}

func TestUserHandler_GetUser(t *testing.T) {
    th := setupTestHandler()
    router := mux.NewRouter()
    th.RegisterRoutes(router)

    // Создаем тестового пользователя
    testUser := model.User{
        ID:   "1",
        Name: "Test User",
        Age:  30,
    }
    th.service.users[testUser.ID] = testUser

    tests := []struct {
        name       string
        userID     string
        wantStatus int
    }{
        {
            name:       "existing user",
            userID:     "1",
            wantStatus: http.StatusOK,
        },
        {
            name:       "non-existing user",
            userID:     "999",
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.wantStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.wantStatus)
            }
        })
    }
}

func TestUserHandler_UpdateUser(t *testing.T) {
    th := setupTestHandler()
    router := mux.NewRouter()
    th.RegisterRoutes(router)

    // Создаем тестового пользователя
    testUser := model.User{
        ID:   "1",
        Name: "Test User",
        Age:  30,
    }
    th.service.users[testUser.ID] = testUser

    tests := []struct {
        name       string
        userID     string
        updateUser model.User
        wantStatus int
    }{
        {
            name:   "valid update",
            userID: "1",
            updateUser: model.User{
                Name: "Updated User",
                Age:  35,
            },
            wantStatus: http.StatusOK,
        },
        {
            name:   "non-existing user",
            userID: "999",
            updateUser: model.User{
                Name: "Updated User",
                Age:  35,
            },
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.updateUser.ID = tt.userID
            body, _ := json.Marshal(tt.updateUser)
            req := httptest.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.wantStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.wantStatus)
            }
        })
    }
}

func TestUserHandler_DeleteUser(t *testing.T) {
    th := setupTestHandler()
    router := mux.NewRouter()
    th.RegisterRoutes(router)

    // Создаем тестового пользователя
    testUser := model.User{
        ID:   "1",
        Name: "Test User",
        Age:  30,
    }
    th.service.users[testUser.ID] = testUser

    tests := []struct {
        name       string
        userID     string
        wantStatus int
    }{
        {
            name:       "existing user",
            userID:     "1",
            wantStatus: http.StatusNoContent,
        },
        {
            name:       "non-existing user",
            userID:     "999",
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.wantStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.wantStatus)
            }
        })
    }
}
