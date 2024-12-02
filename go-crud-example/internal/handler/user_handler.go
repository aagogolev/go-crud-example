package handler

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "go-crud-example/internal/model"
    "go-crud-example/internal/service"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "log"
    "net/http"
    "time"
)

type UserHandler struct {
    service service.UserService
    logger  *log.Logger
}

func NewUserHandler(service service.UserService, logger *log.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/users", h.CreateUser).Methods("POST")
    router.HandleFunc("/users", h.GetUsers).Methods("GET")
    router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
    router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
    router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

    // Metrics endpoint
    router.Handle("/metrics", promhttp.Handler())

    // Health check endpoint
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(map[string]string{
            "status":    "OK",
            "timestamp": time.Now().Format(time.RFC3339),
        }); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
    }).Methods("GET")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateUser(&user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.service.GetUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    user, err := h.service.GetUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.ID = id
    if err := h.service.UpdateUser(&user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.service.DeleteUser(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
