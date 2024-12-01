package model

import (
    "testing"
    m "go-crud-example/internal/model"
)

func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    m.User
        wantErr bool
    }{
        {
            name: "valid user",
            user: m.User{
                Name: "John Doe",
                Age:  25,
            },
            wantErr: false,
        },
        {
            name: "invalid name - too short",
            user: m.User{
                Name: "J",
                Age:  25,
            },
            wantErr: true,
        },
        {
            name: "invalid age - too high",
            user: m.User{
                Name: "John Doe",
                Age:  151,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := tt.user.Validate(); (err != nil) != tt.wantErr {
                t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
