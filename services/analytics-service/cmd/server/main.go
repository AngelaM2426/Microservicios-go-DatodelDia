package main

import (
	v1 "ape-go-services/analytics-service/internal/api/v1"
	db_interna "ape-go-services/analytics-service/internal/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, se continuará con las variables de entorno del sistema")
	}

	db_interna.Conectar()

	configurarCierreControlado()

	puerto := os.Getenv("PORT")
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
		// The Ping method is on the underlying sql.DB object.
		if err := db_interna.BaseDeDatos.Ping(); err != nil {
			estadoDB = "error"
		}
		c.JSON(http.StatusOK, gin.H{
			"estado":               "ok",
			"estado_base_de_datos": estadoDB,
		})
	})

	// Configure the API routes.
	v1.ConfigurarRutas(enrutador)

	log.Printf("Iniciando analytics-service en el puerto %s...\n", puerto)
	err = enrutador.Run(":" + puerto)
	if err != nil {
		log.Fatal("Fallo al iniciar el servidor:", err)
	}
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
