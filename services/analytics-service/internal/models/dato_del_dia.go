package models

import "time"

// DatoDelDia representa un registro de analítica en la base de datos.
type DatoDelDia struct {
	ID        int       `json:"id"`
	Tipo      string    `json:"tipo"`      // Por ej: "ofertas de empleo", "solicitudes"
	Valor     string    `json:"valor"`     // El contenido o descripción del dato.
	Origen    string    `json:"origen"`    // El servicio que originó el dato.
	Timestamp time.Time `json:"timestamp"` // El día y hora en que se generó.
}
