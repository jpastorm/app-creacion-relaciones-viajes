package repository

import (
	"fmt"
)

type DetalleRelacion struct {
	FkRelacion     string
	FkDocumento    string
	FkNombre       string
	FkApellidoPa   string
	FkApellidoMa   string
	FkDireccion    string
	FkNacionalidad string
	FkResidencia   string
	FkProfesion    string
	FkFechanac     string
	FkEstcivil     string
	FkSexo         string
	FkMva1         string
	FkMva2         string
	FkTipoa        string
}

// Agregar un detalle de relación
func (r *Repository) AgregarDetalleRelacion(d DetalleRelacion) error {
	query := `
		INSERT INTO tbdetrelacion (fkrelacion, fkdocumento, fknombre, fkapellidopa, fkapellidoma, fkdireccion, fknacionalidad, fkresidencia, fkprofesion, fkfechanac, fkestcivil, fksexo, fkmva1, fkmva2, fktipoa)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, d.FkRelacion, d.FkDocumento, d.FkNombre, d.FkApellidoPa, d.FkApellidoMa, d.FkDireccion, d.FkNacionalidad, d.FkResidencia, d.FkProfesion, d.FkFechanac, d.FkEstcivil, d.FkSexo, d.FkMva1, d.FkMva2, d.FkTipoa)
	if err != nil {
		return fmt.Errorf("error al insertar detalle de relación: %v", err)
	}
	return nil
}
