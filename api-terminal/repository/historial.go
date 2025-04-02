package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Historial representa un registro en la base de datos
type Historial struct {
	ID          int       `json:"id"`
	CreatedTime time.Time `json:"created_time"`
	Datos       string    `json:"datos"`
}

// RespuestaPaginada estructura la salida en formato JSON
type RespuestaPaginada struct {
	Data         []interface{} `json:"data"` // Ahora contiene objetos JSON
	Error        string        `json:"error"`
	Message      string        `json:"message"`
	TotalPages   int           `json:"totalPages"`
	TotalRecords int           `json:"totalRecords"`
}

// initHistorialDB inicializa la base de datos SQLite y crea la tabla si no existe
func initHistorialDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	query := `
    CREATE TABLE IF NOT EXISTS historial (
        created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
        datos TEXT NOT NULL
    );`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla: %w", err)
	}

	return db, nil
}

// GuardarHistorial guarda un nuevo registro en la base de datos
func (Repository) GuardarHistorial(dbPath string, datos string) error {
	jsonDatos, err := json.Marshal(datos)
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %w", err)
	}

	db, err := initHistorialDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO historial (datos) VALUES (?)`
	_, err = db.Exec(query, string(jsonDatos))
	if err != nil {
		return fmt.Errorf("error al guardar historial: %w", err)
	}

	log.Println("Historial guardado correctamente en SQLite")
	return nil
}

// EliminarHistorialPorRango elimina registros entre dos fechas
func (Repository) EliminarHistorialPorRango(dbPath string, fechaInicio, fechaFin string) error {
	db, err := initHistorialDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM historial WHERE created_time BETWEEN ? AND ?`
	result, err := db.Exec(query, fechaInicio, fechaFin)
	if err != nil {
		return fmt.Errorf("error al eliminar historial: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Se eliminaron %d registros en el rango de fechas", rowsAffected)
	return nil
}

// ListarHistorial obtiene registros con paginación
func (Repository) ListarHistorial(dbPath string, page, pageSize int) (RespuestaPaginada, error) {
	var respuesta RespuestaPaginada
	if page <= 0 || pageSize <= 0 {
		return RespuestaPaginada{Error: "Parámetros de paginación inválidos"}, fmt.Errorf("los valores de página y tamaño deben ser mayores a cero")
	}

	offset := (page - 1) * pageSize
	db, err := initHistorialDB(dbPath)
	if err != nil {
		return respuesta, err
	}
	defer db.Close()

	// Consulta con paginación
	query := `SELECT created_time, datos FROM historial ORDER BY created_time DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return respuesta, fmt.Errorf("error al consultar historial: %w", err)
	}
	defer rows.Close()

	var historial []interface{}
	for rows.Next() {
		var h Historial
		var createdTimeStr string
		err := rows.Scan(&createdTimeStr, &h.Datos)
		if err != nil {
			return respuesta, fmt.Errorf("error escaneando fila: %w", err)
		}

		historial = append(historial, h.Datos)
	}

	// Obtener total de registros
	countQuery := `SELECT COUNT(*) FROM historial`
	var totalRecords int
	err = db.QueryRow(countQuery).Scan(&totalRecords)
	if err != nil {
		return respuesta, fmt.Errorf("error al contar registros: %w", err)
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	// Formar la respuesta
	respuesta = RespuestaPaginada{
		Data:         historial,
		Error:        "",
		Message:      "ok",
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	return respuesta, nil
}
