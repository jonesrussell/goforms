package handlers

import (
	"net/http"
	"os"
	"testing"

	"log"

	"github.com/joho/godotenv"
	"github.com/jonesrussell/goforms/internal/models"
	"github.com/jonesrussell/goforms/test/fixtures"
	"github.com/jonesrussell/goforms/test/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type SubscriptionTestSuite struct {
	suite.Suite
	testDB  *setup.TestDB
	handler *SubscriptionHandler
	fixture *fixtures.SubscriptionFixture
}

func (s *SubscriptionTestSuite) SetupSuite() {
	var err error

	// Add this to help debug the current working directory
	dir, _ := os.Getwd()
	log.Printf("Current working directory: %s", dir)

	// Option 1: Use relative path from test file location
	err = godotenv.Load("../../.env.test")

	// OR Option 2: Use project root relative path
	// err = godotenv.Load("../../../.env.test")

	require.NoError(s.T(), err, "Error loading .env.test file")

	// Setup test database
	s.testDB, err = setup.NewTestDB()
	require.NoError(s.T(), err)

	// Run migrations
	err = s.testDB.RunMigrations()
	require.NoError(s.T(), err)

	// Setup handler
	logger, _ := zap.NewDevelopment()
	store := models.NewSubscriptionStore(s.testDB.DB)
	s.handler = NewSubscriptionHandler(logger, store)

	// Setup fixture
	s.fixture = fixtures.NewSubscriptionFixture(s.handler.CreateSubscription)
}

func (s *SubscriptionTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Cleanup(true)
		if err != nil {
			s.T().Logf("Failed to cleanup test database: %v", err)
		}
	}
}

func (s *SubscriptionTestSuite) SetupTest() {
	err := s.testDB.ClearData()
	require.NoError(s.T(), err)
}

func (s *SubscriptionTestSuite) TestSubscriptionIntegration() {
	// Test successful subscription
	rec, err := s.fixture.CreateSubscriptionRequest("integration@test.com")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = fixtures.ParseResponse(rec, &response)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "integration@test.com", response["email"])

	// Verify database record
	var exists bool
	err = s.testDB.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM subscriptions WHERE email = ?)", "integration@test.com")
	require.NoError(s.T(), err)
	assert.True(s.T(), exists)

	// Test duplicate subscription
	rec, err = s.fixture.CreateSubscriptionRequest("integration@test.com")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusConflict, rec.Code)

	var errResponse map[string]string
	err = fixtures.ParseResponse(rec, &errResponse)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Email already subscribed", errResponse["error"])

	// Test invalid origin
	rec, err = s.fixture.CreateSubscriptionRequestWithOrigin("new@test.com", "https://invalid-origin.com")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusForbidden, rec.Code)

	err = fixtures.ParseResponse(rec, &errResponse)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "invalid origin", errResponse["error"])
}

func TestSubscriptionSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}
