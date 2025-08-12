package db

import (
	"ape-go-services/pkg/db"
	"database/sql"
	"log"
	"os"
)

var BaseDeDatos *sql.DB

func Conectar() {
	rutaDB := os.Getenv("RUTA_PATO_DB")
	if rutaDB == "" {
		log.Fatal("La variable de entorno RUTA_PATO_DB no está definida")
	}

	conexion, err := db.ConectarPatoDB(rutaDB)
	if err != nil {
		log.Fatalf("Fallo al conectar a DuckDB: %v", err)
	}
	BaseDeDatos = conexion
	log.Println("Conexión con DuckDB exitosa")
}

func Desconectar() error {
	return db.DesconectarPatoDB()
}
