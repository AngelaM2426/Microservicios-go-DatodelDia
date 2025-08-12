package main

import (
	v1 "ape-go-services/analytics-service/internal/api/v1"
	db_interna "ape-go-services/analytics-service/internal/db"
	"ape-go-services/analytics-service/internal/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, se continuará con las variables de entorno del sistema")
	}

	db_interna.Conectar()
	prepararBaseDeDatos()

	configurarCierreControlado()

	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8081"
	}

	enrutador := gin.New()
	enrutador.Use(gin.Logger())
	enrutador.Use(gin.Recovery())

	err := enrutador.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Error al configurar los proxies de confianza: %v", err)
	}

	enrutador.GET("/health", func(c *gin.Context) {
		estadoDB := "ok"
		if err := db_interna.BaseDeDatos.Ping(); err != nil {
			estadoDB = "error"
		}
		c.JSON(http.StatusOK, gin.H{
			"estado":               "ok",
			"estado_base_de_datos": estadoDB,
		})
	})

	// Configurar las rutas de la API, incluyendo nuestro nuevo endpoint
	v1.ConfigurarRutas(enrutador)

	log.Printf("Iniciando analytics-service en el puerto %s...\n", puerto)
	err = enrutador.Run(":" + puerto)
	if err != nil {
		log.Fatal("Fallo al iniciar el servidor:", err)
	}
}

// prepararBaseDeDatos crea la tabla 'datos_del_dia' si no existe y la puebla con datos.
func prepararBaseDeDatos() {
	declaracionSQL := `
	CREATE TABLE IF NOT EXISTS datos_del_dia (
		id INTEGER PRIMARY KEY,
		tipo VARCHAR,
		valor VARCHAR,
		origen VARCHAR,
		timestamp TIMESTAMP
	);`

	_, err := db_interna.BaseDeDatos.Exec(declaracionSQL)
	if err != nil {
		log.Fatalf("Error al crear la tabla 'datos_del_dia': %v", err)
	}

	var conteo int
	db_interna.BaseDeDatos.QueryRow("SELECT COUNT(*) FROM datos_del_dia").Scan(&conteo)
	if conteo > 0 {
		log.Println("La tabla 'datos_del_dia' ya contiene datos.")
		return
	}

	datos := []models.DatoDelDia{
		{Tipo: "ofertas de empleo", Valor: "Nueva oferta para 'Desarrollador Go' publicada", Origen: "job-posting-service"},
		{Tipo: "solicitudes", Valor: "Usuario 'Ana' aplicó a 'Desarrollador Go'", Origen: "job-offers-service"},
		{Tipo: "solicitudes", Valor: "Usuario 'Juan' aplicó a 'Diseñador UX'", Origen: "job-offers-service"},
		{Tipo: "ofertas de empleo", Valor: "Se actualizó la oferta 'Diseñador UX'", Origen: "job-posting-service"},
		{Tipo: "solicitudes", Valor: "Usuario 'Pedro' aplicó a 'Analista de Datos'", Origen: "job-offers-service"},
	}

	for i, dato := range datos {
		// Retrasamos ligeramente la inserción para garantizar un orden de timestamp único
		timestamp := time.Now().Add(time.Duration(i) * time.Second)
		_, err := db_interna.BaseDeDatos.Exec(
			"INSERT INTO datos_del_dia (tipo, valor, origen, timestamp) VALUES (?, ?, ?, ?)",
			dato.Tipo,
			dato.Valor,
			dato.Origen,
			timestamp,
		)
		if err != nil {
			log.Printf("No se pudo insertar el dato '%s': %v\n", dato.Valor, err)
		}
	}
	log.Println("Tabla 'datos_del_dia' creada y poblada exitosamente.")
}

func configurarCierreControlado() {
	canal := make(chan os.Signal, 1)
	signal.Notify(canal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-canal
		log.Println("Señal de apagado recibida. Cerrando de forma controlada...")
		if err := db_interna.Desconectar(); err != nil {
			log.Printf("Error al desconectar de DuckDB: %v", err)
		}
		log.Println("Todas las conexiones han sido cerradas. Saliendo de la aplicación.")
		os.Exit(0)
	}()
}
