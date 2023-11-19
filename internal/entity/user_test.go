package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Heitor Saraiva","espetinho@gmail.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Email)
	assert.Equal(t, user.Name, "Heitor Saraiva")
	assert.Equal(t, user.Email, "espetinho@gmail.com")
}

func TestUserCheckPassword(t *testing.T) {
	user, err := NewUser("Heitor Saraiva","espetinho@gmail.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.CheckPassword("123456"))
	assert.False(t, user.CheckPassword("1234567"))
	assert.NotEqual(t,"123456",user.Password)
	
}