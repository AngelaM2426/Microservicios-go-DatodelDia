package handlers

import (
	"ape-go-services/analytics-service/internal/services"
	"net/http"

	// La corrección está aquí: se ha eliminado el ".com" extra.
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

// ObtenerConteoDatosDelDia maneja la petición GET /api/v1/dato-del-dia/conteo
func (m *ManejadorAnalitica) ObtenerConteoDatosDelDia(c *gin.Context) {
	conteo, err := m.servicio.ObtenerConteoDatosDelDia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo obtener el conteo de datos del día",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conteo_datos_del_dia": conteo,
	})
}

// ObtenerDatoMasReciente maneja la petición GET /api/v1/dato-del-dia
func (m *ManejadorAnalitica) ObtenerDatoMasReciente(c *gin.Context) {
	dato, err := m.servicio.ObtenerDatoMasReciente()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ocurrió un error al obtener el dato más reciente",
		})
		return
	}

	if dato == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"mensaje": "No se encontraron datos del día",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dato_del_dia": dato,
	})
}
