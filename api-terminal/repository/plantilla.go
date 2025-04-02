package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite
)

// Impresion representa el formato de una plantilla de impresión
type Impresion struct {
	ID     int                      `json:"id"`     // ID único de la plantilla
	Fuente string                   `json:"fuente"` //tamanio de letra
	Titulo string                   `json:"titulo"` // Título o nombre de la plantilla
	Datos  []map[string]interface{} `json:"datos"`  // Datos de la plantilla en formato JSON
}

// initDB inicializa la base de datos SQLite
func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// Crear tabla si no existe
	query := `
    CREATE TABLE IF NOT EXISTS plantillas (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		fuente TEXT NOT NULL,
        titulo TEXT NOT NULL,
        datos TEXT NOT NULL
    );`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla: %w", err)
	}

	return db, nil
}

// Guardar guarda una nueva plantilla de impresión en la base de datos
func GuardarPlantilla(dbPath string, titulo string, fuente string, datos []map[string]interface{}) error {
	// Convertir datos a JSON
	jsonDatos, err := json.Marshal(datos)
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %w", err)
	}

	// Inicializar la base de datos
	db, err := initDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insertar la plantilla en la base de datos
	query := `INSERT INTO plantillas (titulo, fuente, datos) VALUES (?,?,?)`
	_, err = db.Exec(query, titulo, fuente, string(jsonDatos))
	if err != nil {
		return fmt.Errorf("error al guardar la plantilla: %w", err)
	}

	log.Println("Plantilla guardada correctamente en SQLite")
	return nil
}

// Buscar busca una plantilla por su título
func BuscarPlantilla(dbPath string, idParam int) (*Impresion, error) {
	// Inicializar la base de datos
	db, err := initDB(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Consultar la plantilla por título
	query := `SELECT id, titulo, fuente, datos FROM plantillas WHERE id = ?`
	row := db.QueryRow(query, idParam)

	var id int
	var tituloGuardado string
	var fuenteGuardada string
	var datosJSON string
	err = row.Scan(&id, &tituloGuardado, &fuenteGuardada, &datosJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no se encontró ninguna plantilla con el id '%d'", id)
		}
		return nil, fmt.Errorf("error al buscar la plantilla: %w", err)
	}

	// Convertir datos JSON a mapa
	var datos []map[string]interface{}
	err = json.Unmarshal([]byte(datosJSON), &datos)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar datos JSON: %w", err)
	}

	// Devolver la plantilla encontrada
	return &Impresion{
		ID:     id,
		Fuente: fuenteGuardada,
		Titulo: tituloGuardado,
		Datos:  datos,
	}, nil
}

// ObtenerTodasLasPlantillas devuelve todas las plantillas almacenadas en la base de datos.
func ObtenerTodasLasPlantillas(dbPath string) ([]Impresion, error) {
	// Inicializar la base de datos
	db, err := initDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al inicializar la base de datos: %w", err)
	}
	defer db.Close()

	// Consultar todas las plantillas
	query := `SELECT id, titulo, fuente, datos FROM plantillas`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar las plantillas: %w", err)
	}
	defer rows.Close()

	// Slice para almacenar todas las plantillas
	var plantillas []Impresion

	// Iterar sobre los resultados
	for rows.Next() {
		var id int
		var tituloGuardado string
		var fuenteGuardada string
		var datosJSON string

		// Escanear los valores de cada fila
		err := rows.Scan(&id, &tituloGuardado, &fuenteGuardada, &datosJSON)
		if err != nil {
			return nil, fmt.Errorf("error al escanear los datos de una plantilla: %w", err)
		}

		// Convertir datos JSON a mapa
		var datos []map[string]interface{}
		err = json.Unmarshal([]byte(datosJSON), &datos)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar datos JSON: %w", err)
		}

		// Crear un objeto Impresion y añadirlo al slice
		plantilla := Impresion{
			ID:     id,
			Fuente: fuenteGuardada,
			Titulo: tituloGuardado,
			Datos:  datos,
		}
		plantillas = append(plantillas, plantilla)
	}

	// Verificar si hubo errores durante la iteración
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados: %w", err)
	}

	// Devolver el listado completo de plantillas
	return plantillas, nil
}

// ActualizarPlantilla actualiza una plantilla existente en la base de datos
func ActualizarPlantilla(dbPath string, id int, titulo string, fuente string, datos []map[string]interface{}) error {
	// Convertir datos a JSON
	jsonDatos, err := json.Marshal(datos)
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %w", err)
	}

	// Inicializar la base de datos
	db, err := initDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Actualizar la plantilla en la base de datos
	query := `UPDATE plantillas SET titulo = ?, fuente = ?, datos = ? WHERE id = ?`
	result, err := db.Exec(query, titulo, fuente, string(jsonDatos), id)
	if err != nil {
		return fmt.Errorf("error al actualizar la plantilla: %w", err)
	}

	// Verificar si se actualizó algún registro
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ninguna plantilla encontrada con ID %d", id)
	}

	log.Printf("Plantilla con ID %d actualizada correctamente en SQLite\n", id)
	return nil
}

// EliminarPlantilla elimina una plantilla existente de la base de datos
func EliminarPlantilla(dbPath string, id int) error {
	// Validar que el ID sea válido
	if id == 0 {
		return fmt.Errorf("el campo ID es obligatorio y debe ser un número entero válido")
	}

	// Inicializar la base de datos
	db, err := initDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Eliminar la plantilla de la base de datos
	query := `DELETE FROM plantillas WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la plantilla: %w", err)
	}

	// Verificar si se eliminó algún registro
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ninguna plantilla encontrada con ID %d", id)
	}

	log.Printf("Plantilla con ID %d eliminada correctamente de SQLite\n", id)
	return nil
}
