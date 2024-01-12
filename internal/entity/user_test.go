package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Pedro", "pedro@mont.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, user.Name, "Pedro")
	assert.Equal(t, user.Email, "pedro@mont.com")
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("Pedro", "pedro@mont.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ComparePassword("123456"))
	assert.False(t, user.ComparePassword("1234567343"))
	assert.NotEqual(t, "123456", user.Password)
}
