package model

type User struct {
    ID   string `json:"id"`
    Name string `json:"name" validate:"required,min=2,max=100"`
    Age  int    `json:"age" validate:"required,gte=0,lte=150"`
}
