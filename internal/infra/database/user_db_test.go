package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pedrohenmonteiro/golang-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, err := openDBUserAndCreateTable()
	assert.NoError(t, err)

	defer dropTableUserAndCloseDB(db)

	user, err := entity.NewUser("Pedro", "p@m.com", "123456")
	assert.NoError(t, err)

	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.QueryRow("select * from users where id = ?", user.ID).Scan(&userFound.ID, &userFound.Name, &userFound.Email, &userFound.Password)

	assert.Nil(t, err)

	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)

}

func TestFindByEmail(t *testing.T) {
	db, err := openDBUserAndCreateTable()
	assert.NoError(t, err)

	defer dropTableUserAndCloseDB(db)

	user, err := entity.NewUser("Pedro", "p@m.com", "123456")
	assert.NoError(t, err)
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)

	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)

}

//auxiliares

func openDBUserAndCreateTable() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT,
			email TEXT,
			password TEXT
		)
	`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dropTableUserAndCloseDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}
