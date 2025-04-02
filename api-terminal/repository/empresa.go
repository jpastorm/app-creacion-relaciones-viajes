package repository

import (
	"database/sql"
	"fmt"
)

// Estructura que representa una empresa
type Empresa struct {
	Nauto      string `json:"nauto"`
	Nombre     string `json:"nombre"`
	Resolucion string `json:"resolucion"`
	Documento  string `json:"documento"`
	Permiso    string `json:"permiso"`
}

// CrearEmpresa inserta una nueva empresa en la base de datos
func (r *Repository) CrearEmpresa(e Empresa) error {
	query := `
        INSERT INTO tbempresa (nauto, Nombre, Resolucion, Documento, Permiso)
        VALUES (?, ?, ?, ?, ?)
    `
	_, err := r.db.Exec(query, e.Nauto, e.Nombre, e.Resolucion, e.Documento, e.Permiso)
	if err != nil {
		return fmt.Errorf("error al crear empresa: %v", err)
	}
	return nil
}

// ObtenerEmpresas devuelve todas las empresas de la base de datos
func (r *Repository) ObtenerEmpresas() ([]Empresa, error) {
	query := `
        SELECT nauto, Nombre, Resolucion, Documento, Permiso
        FROM tbempresa
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresas: %v", err)
	}
	defer rows.Close()

	var empresas []Empresa
	for rows.Next() {
		var e Empresa
		err := rows.Scan(&e.Nauto, &e.Nombre, &e.Resolucion, &e.Documento, &e.Permiso)
		if err != nil {
			return nil, fmt.Errorf("error al escanear empresa: %v", err)
		}
		empresas = append(empresas, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre las filas: %v", err)
	}

	return empresas, nil
}

// ObtenerEmpresaPorID devuelve una empresa por su ID (nauto)
func (r *Repository) ObtenerEmpresaPorID(id string) (*Empresa, error) {
	query := `
        SELECT nauto, Nombre, Resolucion, Documento, Permiso
        FROM tbempresa
        WHERE nauto = ?
    `
	row := r.db.QueryRow(query, id)

	var e Empresa
	err := row.Scan(&e.Nauto, &e.Nombre, &e.Resolucion, &e.Documento, &e.Permiso)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("empresa con ID %s no encontrada", id)
	} else if err != nil {
		return nil, fmt.Errorf("error al obtener empresa: %v", err)
	}

	return &e, nil
}

// ActualizarEmpresa actualiza los datos de una empresa en la base de datos
func (r *Repository) ActualizarEmpresa(id string, e Empresa) error {
	query := `
        UPDATE tbempresa
        SET Nombre = ?, Resolucion = ?, Documento = ?, Permiso = ?
        WHERE nauto = ?
    `
	_, err := r.db.Exec(query, e.Nombre, e.Resolucion, e.Documento, e.Permiso, id)
	if err != nil {
		return fmt.Errorf("error al actualizar empresa: %v", err)
	}
	return nil
}

// EliminarEmpresa elimina una empresa de la base de datos
func (r *Repository) EliminarEmpresa(id string) error {
	query := `
        DELETE FROM tbempresa
        WHERE nauto = ?
    `
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar empresa: %v", err)
	}
	return nil
}

func (r *Repository) CrearOActualizarEmpresa(e Empresa) error {
	// Primero, verificamos si la empresa ya existe en la base de datos
	var nauto sql.NullString
	checkQuery := `SELECT nauto FROM tbempresa WHERE nauto = ?;`

	// Verificar si la empresa ya existe por nauto
	err := r.db.QueryRow(checkQuery, e.Nauto).Scan(&nauto)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existe, realizar un INSERT
			return r.CrearEmpresa(e)
		}
		// Error inesperado
		return fmt.Errorf("error al verificar existencia de empresa: %v", err)
	}

	// Si existe, realizar un UPDATE
	return r.ActualizarEmpresa(e.Nauto, e)
}
