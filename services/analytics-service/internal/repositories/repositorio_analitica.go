package repositories

import (
	db_interna "ape-go-services/analytics-service/internal/db"
	"ape-go-services/analytics-service/internal/models"
	"database/sql"
	"fmt"
)

// RepositorioAnalitica maneja el acceso a los datos de analítica.
// It now implicitly implements the IRepositorioAnalitica interface.
type RepositorioAnalitica struct {
	db *sql.DB
}

// NuevoRepositorioAnalitica crea una nueva instancia del repositorio.
func NuevoRepositorioAnalitica() *RepositorioAnalitica {
	return &RepositorioAnalitica{
		db: db_interna.BaseDeDatos,
	}
}

// ObtenerMasReciente devuelve el último dato del día registrado.
func (r *RepositorioAnalitica) ObtenerMasReciente() (*models.DatoDelDia, error) {
	query := "SELECT id, nombre, valor, timestamp FROM datos_del_dia ORDER BY timestamp DESC LIMIT 1"

	var dato models.DatoDelDia
	err := r.db.QueryRow(query).Scan(&dato.ID, &dato.Nombre, &dato.Valor, &dato.Timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener el dato más reciente: %w", err)
	}

	return &dato, nil
}
