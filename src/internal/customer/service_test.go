package customer

import (
	"testing"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetCustomerById(id int) (models.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(models.Customer), args.Error(1)
}

func (m *mockRepository) GetAllCustomers() ([]models.Customer, error) {
	args := m.Called()
	return args.Get(0).([]models.Customer), args.Error(1)
}

func TestGetCustomerById_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepository)
	svc := NewCustomerService(mockRepo)

	expectedCustomer := models.Customer{
		Dni:          "1",
		CompleteName: "John Doe",
		// Add other fields as necessary
	}

	mockRepo.On("GetCustomerById", 1).Return(expectedCustomer, nil)

	// Act
	customer, err := svc.GetCustomerById(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, customer)
	mockRepo.AssertExpectations(t)
}
