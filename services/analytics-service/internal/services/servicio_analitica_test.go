package services

import (
	"ape-go-services/analytics-service/internal/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositorioAnalitica is a mock implementation of the IRepositorioAnalitica interface.
// It allows us to control the behavior of the repository during tests.
type MockRepositorioAnalitica struct {
	mock.Mock
}

// ObtenerMasReciente is the mock's method that fulfills the interface.
// It records that it was called and returns whatever we tell it to.
func (m *MockRepositorioAnalitica) ObtenerMasReciente() (*models.DatoDelDia, error) {
	args := m.Called()
	// Handle the case where the returned object is nil, to avoid a panic.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DatoDelDia), args.Error(1)
}

func TestServicioAnalitica_ObtenerDatoMasReciente(t *testing.T) {

	t.Run("Success case - data is found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockRepositorioAnalitica)
		servicio := NewServicioAnaliticaConRepositorio(mockRepo)
		expectedDato := &models.DatoDelDia{
			ID:        1,
			Nombre:    "ofertas de empleo",
			Valor:     "Nueva oferta para 'Desarrollador Go' publicada",
			Timestamp: time.Now(),
		}

		// Configure the mock: when ObtenerMasReciente is called, return our expected data.
		mockRepo.On("ObtenerMasReciente").Return(expectedDato, nil)

		// Act
		dato, err := servicio.ObtenerDatoMasReciente()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, dato)
		assert.Equal(t, expectedDato.Valor, dato.Valor)
		mockRepo.AssertExpectations(t) // Verify that the mocked method was called as expected.
	})

	t.Run("Not found case - repository returns nil", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockRepositorioAnalitica)
		servicio := NewServicioAnaliticaConRepositorio(mockRepo)

		// Configure the mock to simulate "not found": return (nil, nil).
		mockRepo.On("ObtenerMasReciente").Return(nil, nil)

		// Act
		dato, err := servicio.ObtenerDatoMasReciente()

		// Assert
		assert.NoError(t, err) // The service should not treat "not found" as an error.
		assert.Nil(t, dato)    // The data object returned should be nil.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure case - repository returns a database error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockRepositorioAnalitica)
		servicio := NewServicioAnaliticaConRepositorio(mockRepo)
		expectedError := errors.New("database connection lost")

		// Configure the mock to return an error.
		mockRepo.On("ObtenerMasReciente").Return(nil, expectedError)

		// Act
		dato, err := servicio.ObtenerDatoMasReciente()

		// Assert
		assert.Error(t, err)                                   // An error should be returned.
		assert.Nil(t, dato)                                    // No data should be returned.
		assert.Contains(t, err.Error(), expectedError.Error()) // The original error should be wrapped.
		mockRepo.AssertExpectations(t)
	})
}
