package v1

// Define your routes and route groups here

import (
	"ape-go-services/analytics-service/internal/api/v1/handlers"

	"github.com/gin-gonic/gin"
)

// ConfigurarRutas define todas las rutas para la v1 de la API de analítica.
func ConfigurarRutas(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	manejadorAnalitica := handlers.NuevoManejadorAnalitica()

	// Endpoint para obtener el conteo total
	v1.GET("/dato-del-dia/conteo", manejadorAnalitica.ObtenerConteoDatosDelDia)

	// Endpoint para obtener el último dato del día
	v1.GET("/dato-del-dia", manejadorAnalitica.ObtenerDatoMasReciente)
}
