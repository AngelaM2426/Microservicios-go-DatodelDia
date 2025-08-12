// FILE: pkg/db/duckdb.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/marcboeker/go-duckdb"
)

var PatoDB *sql.DB

func ConectarPatoDB(rutaDB string) (*sql.DB, error) {
	if rutaDB == "" {
		return nil, fmt.Errorf("la ruta de duckdb no puede estar vacía")
	}

	directorio := filepath.Dir(rutaDB)
	if err := os.MkdirAll(directorio, 0755); err != nil {
		return nil, fmt.Errorf("fallo al crear el directorio de duckdb %s: %w", directorio, err)
	}

	conector, err := sql.Open("duckdb", rutaDB)
	if err != nil {
		return nil, fmt.Errorf("fallo al abrir la conexión con duckdb: %w", err)
	}

	if err := conector.Ping(); err != nil {
		return nil, fmt.Errorf("fallo al hacer ping a duckdb: %w", err)
	}

	PatoDB = conector
	log.Printf("Conectado exitosamente a DuckDB en: %s", rutaDB)
	return PatoDB, nil
}

func DesconectarPatoDB() error {
	if PatoDB == nil {
		return nil
	}

	log.Println("Cerrando la conexión con DuckDB...")
	if err := PatoDB.Close(); err != nil {
		return fmt.Errorf("fallo al cerrar la conexión con duckdb: %w", err)
	}

	log.Println("Conexión con DuckDB cerrada correctamente.")
	return nil
}

func ObtenerPatoDB() *sql.DB {
	if PatoDB == nil {
		log.Fatal("La conexión con DuckDB no ha sido inicializada. Llama a db.ConectarPatoDB() primero.")
	}
	return PatoDB
}
