package services

import (
	"ape-go-services/analytics-service/internal/models"
	"ape-go-services/analytics-service/internal/repositories"
	"fmt"
)

// ServicioAnalitica encapsula la lógica de negocio para la analítica.
type ServicioAnalitica struct {
	// The service now depends on the interface, not the concrete repository.
	repositorio repositories.IRepositorioAnalitica
}

func NuevoServicioAnalitica() *ServicioAnalitica {
	return &ServicioAnalitica{
		repositorio: repositories.NuevoRepositorioAnalitica(),
	}
}

func NewServicioAnaliticaConRepositorio(repo repositories.IRepositorioAnalitica) *ServicioAnalitica {
	return &ServicioAnalitica{
		repositorio: repo,
	}
}

// ObtenerDatoMasReciente devuelve el último dato del día registrado.
func (s *ServicioAnalitica) ObtenerDatoMasReciente() (*models.DatoDelDia, error) {
	dato, err := s.repositorio.ObtenerMasReciente()
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}
	return dato, nil
}
