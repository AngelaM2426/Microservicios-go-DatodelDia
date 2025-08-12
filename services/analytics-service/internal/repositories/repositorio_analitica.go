package repositories

import (
	db_interna "ape-go-services/analytics-service/internal/db"
	"ape-go-services/analytics-service/internal/models"
	"database/sql"
	"fmt"
)

// RepositorioAnalitica maneja el acceso a los datos de analítica.
type RepositorioAnalitica struct {
	db *sql.DB
}

// NuevoRepositorioAnalitica crea una nueva instancia del repositorio.
func NuevoRepositorioAnalitica() *RepositorioAnalitica {
	return &RepositorioAnalitica{
		db: db_interna.BaseDeDatos,
	}
}

// ContarDatosDelDia devuelve el número total de registros en la tabla 'datos_del_dia'.
func (r *RepositorioAnalitica) ContarDatosDelDia() (int, error) {
	var conteo int
	err := r.db.QueryRow("SELECT COUNT(*) FROM datos_del_dia").Scan(&conteo)
	if err != nil {
		return 0, fmt.Errorf("error al contar los datos del día: %w", err)
	}
	return conteo, nil
}

// ObtenerMasReciente devuelve el último dato del día registrado.
func (r *RepositorioAnalitica) ObtenerMasReciente() (*models.DatoDelDia, error) {
	query := "SELECT id, tipo, valor, origen, timestamp FROM datos_del_dia ORDER BY timestamp DESC LIMIT 1"

	var dato models.DatoDelDia
	err := r.db.QueryRow(query).Scan(&dato.ID, &dato.Tipo, &dato.Valor, &dato.Origen, &dato.Timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			// Si no hay filas, no es un error. Devolvemos nil para que el servicio lo maneje.
			return nil, nil
		}
		// Para cualquier otro error, lo devolvemos.
		return nil, fmt.Errorf("error al obtener el dato más reciente: %w", err)
	}

	return &dato, nil
}
