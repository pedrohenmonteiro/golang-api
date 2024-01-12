package database

import (
	"database/sql"

	"github.com/pedrohenmonteiro/golang-api/internal/entity"
)

type User struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *entity.User) error {
	stmt, err := u.DB.Prepare("insert into users(id, name, email, password) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID.String(), user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	stmt, err := u.DB.Prepare("select * from users where email = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user entity.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
