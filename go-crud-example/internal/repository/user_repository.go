package repository

import (
    "database/sql"
    "go-crud-example/internal/model"
)

type UserRepository interface {
    GetAll() ([]model.User, error)
    GetByID(id string) (*model.User, error)
    Create(user *model.User) error
    Update(user *model.User) error
    Delete(id string) error
}

type PostgresUserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetAll() ([]model.User, error) {
    rows, err := r.db.Query("SELECT id, name, age FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []model.User{}
    for rows.Next() {
        var u model.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

func (r *PostgresUserRepository) GetByID(id string) (*model.User, error) {
    var u model.User
    err := r.db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).
        Scan(&u.ID, &u.Name, &u.Age)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *PostgresUserRepository) Create(user *model.User) error {
    return r.db.QueryRow(
        "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id",
        user.Name,
        user.Age,
    ).Scan(&user.ID)
}

func (r *PostgresUserRepository) Update(user *model.User) error {
    result, err := r.db.Exec(
        "UPDATE users SET name = $1, age = $2 WHERE id = $3",
        user.Name,
        user.Age,
        user.ID,
    )
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}

func (r *PostgresUserRepository) Delete(id string) error {
    result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}
