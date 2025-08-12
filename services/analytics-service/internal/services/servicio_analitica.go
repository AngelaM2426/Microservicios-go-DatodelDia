package services

import (
	"ape-go-services/analytics-service/internal/models"
	"ape-go-services/analytics-service/internal/repositories"
	"fmt"
)

// ServicioAnalitica encapsula la lógica de negocio para la analítica.
type ServicioAnalitica struct {
	repositorio *repositories.RepositorioAnalitica
}

// NuevoServicioAnalitica crea una nueva instancia del servicio.
func NuevoServicioAnalitica() *ServicioAnalitica {
	return &ServicioAnalitica{
		repositorio: repositories.NuevoRepositorioAnalitica(),
	}
}

// ObtenerConteoDatosDelDia obtiene el número total de datos.
func (s *ServicioAnalitica) ObtenerConteoDatosDelDia() (int, error) {
	conteo, err := s.repositorio.ContarDatosDelDia()
	if err != nil {
		return 0, fmt.Errorf("servicio: %w", err)
	}
	return conteo, nil
}

// ObtenerDatoMasReciente devuelve el último dato del día registrado.
func (s *ServicioAnalitica) ObtenerDatoMasReciente() (*models.DatoDelDia, error) {
	dato, err := s.repositorio.ObtenerMasReciente()
	if err != nil {
		return nil, fmt.Errorf("servicio: %w", err)
	}
	return dato, nil
}
