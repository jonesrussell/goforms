package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetPassword(t *testing.T) {
	user := &User{}
	err := user.SetPassword("securepassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.HashedPassword)
}

func TestCheckPassword(t *testing.T) {
	user := &User{}
	err := user.SetPassword("securepassword")
	assert.NoError(t, err)

	// Test correct password
	assert.True(t, user.CheckPassword("securepassword"))

	// Test incorrect password
	assert.False(t, user.CheckPassword("wrongpassword"))
}
