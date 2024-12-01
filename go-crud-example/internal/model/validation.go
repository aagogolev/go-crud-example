package model

import (
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (u *User) Validate() error {
    return validate.Struct(u)
}
