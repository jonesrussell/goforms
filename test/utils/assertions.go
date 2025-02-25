package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertJSONResponse asserts common JSON response properties
func AssertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()
	assert.Equal(t, expectedStatus, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
}

// AssertErrorResponse asserts error response properties
func AssertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	t.Helper()
	AssertJSONResponse(t, rec, expectedStatus)
	var response map[string]interface{}
	err := ParseJSONResponse(rec, &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	if expectedError != "" {
		assert.Equal(t, expectedError, response["error"])
	}
}

// AssertSuccessResponse asserts success response properties
func AssertSuccessResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()
	AssertJSONResponse(t, rec, expectedStatus)
	var response map[string]interface{}
	err := ParseJSONResponse(rec, &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")
}
