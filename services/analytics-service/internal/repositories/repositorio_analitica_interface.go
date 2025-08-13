package repositories

import "ape-go-services/analytics-service/internal/models"

type IRepositorioAnalitica interface {
	ObtenerMasReciente() (*models.DatoDelDia, error)
}
