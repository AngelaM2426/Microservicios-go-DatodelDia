package handlers

import (
	"ape-go-services/analytics-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ManejadorAnalitica gestiona las peticiones HTTP para analítica.
type ManejadorAnalitica struct {
	servicio *services.ServicioAnalitica
}

// NuevoManejadorAnalitica crea una nueva instancia del manejador.
func NuevoManejadorAnalitica() *ManejadorAnalitica {
	return &ManejadorAnalitica{
		servicio: services.NuevoServicioAnalitica(),
	}
}

// ObtenerDatoMasReciente maneja la petición GET /api/v1/dato-del-dia
func (m *ManejadorAnalitica) ObtenerDatoMasReciente(c *gin.Context) {
	dato, err := m.servicio.ObtenerDatoMasReciente()
	if err != nil {
		// If the service returns any error, respond with a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ocurrió un error al obtener el dato más reciente",
		})
		return
	}

	// The service returns `nil` (without an error) if no record was found.
	// This should be handled as a 404 Not Found.
	if dato == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"mensaje": "No se encontraron datos del día",
		})
		return
	}

	// On success, return a 200 OK with the data object.
	c.JSON(http.StatusOK, gin.H{
		"dato_del_dia": dato,
	})
}
