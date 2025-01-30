package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jonesrussell/goforms/internal/domain/common"
)

func TestSetPassword(t *testing.T) {
	user := &common.User{}
	err := user.SetPassword("securepassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.HashedPassword)
}

func TestCheckPassword(t *testing.T) {
	user := &common.User{}
	err := user.SetPassword("securepassword")
	assert.NoError(t, err)

	// Test correct password
	assert.True(t, user.CheckPassword("securepassword"))

	// Test incorrect password
	assert.False(t, user.CheckPassword("wrongpassword"))
}
