package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
)

// Relacion representa la estructura de la relación principal.
type RelacionList struct {
	IDRelacion           string
	Patente              string
	FechaSalida          string
	FechaRetorno         string
	ProcedenciaT         string
	DestinoA             string
	ProcedenciaA         string
	DestinoT             string
	ConductorJSON        json.RawMessage
	DetallesRelacionJSON json.RawMessage
	EmpresaJSON          json.RawMessage
}

type Relacion struct {
	IdRelacion  string
	FkPatente   string
	FkDocumento string
	Fesal       string
	Fere        string
	Prot        string
	Desa        string
	Proa        string
	Dest        string
}

// AgregarRelacion inserta una nueva relación en la base de datos
func (r *Repository) AgregarRelacion(rel Relacion) error {
	// query := `
	// 	INSERT INTO tbrelacion (Idrelacion, fkpatente, fkdocumento, fesal, fere, prot, desa, proa, dest)
	// 	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	// `
	// _, err := r.db.Exec(query, rel.IdRelacion, rel.FkPatente, rel.FkDocumento, rel.Fesal, rel.Fere, rel.Prot, rel.Desa, rel.Proa, rel.Dest)
	// if err != nil {
	// 	return fmt.Errorf("error al insertar relación: %v", err)
	// }
	// return nil

	query := `
		INSERT INTO tbrelacion (Idrelacion, fkpatente, fkdocumento, fesal, fere, prot, desa, proa, dest)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			fkpatente = VALUES(fkpatente), 
			fkdocumento = VALUES(fkdocumento),
			fesal = VALUES(fesal), 
			fere = VALUES(fere),
			prot = VALUES(prot),
			desa = VALUES(desa),
			proa = VALUES(proa),
			dest = VALUES(dest)
	`
	_, err := r.db.Exec(query, rel.IdRelacion, rel.FkPatente, rel.FkDocumento, rel.Fesal, rel.Fere, rel.Prot, rel.Desa, rel.Proa, rel.Dest)
	if err != nil {
		return fmt.Errorf("error al insertar o actualizar relación: %v", err)
	}
	return nil
}

// ListarUltimaRelacion obtiene la última relación en la tabla tbrelacion
func (r *Repository) ListarUltimaRelacion() (Relacion, error) {
	query := `
		SELECT Idrelacion, fkpatente, fkdocumento, fesal, fere, prot, desa, proa, dest 
		FROM tbrelacion 
		ORDER BY CAST(Idrelacion AS UNSIGNED) DESC 
		LIMIT 1
	`
	row := r.db.QueryRow(query)

	// Variables intermedias con sql.NullString para manejar valores nulos
	var idRelacionSQL, fkPatenteSQL, fkDocumentoSQL, fesalSQL, fereSQL, protSQL, desaSQL, proaSQL, destSQL sql.NullString

	err := row.Scan(
		&idRelacionSQL,
		&fkPatenteSQL,
		&fkDocumentoSQL,
		&fesalSQL,
		&fereSQL,
		&protSQL,
		&desaSQL,
		&proaSQL,
		&destSQL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Relacion{}, fmt.Errorf("no se encontraron registros en la tabla tbrelacion")
		}
		return Relacion{}, fmt.Errorf("error al escanear fila: %v", err)
	}

	// Convertir sql.NullString a string, asignando un valor vacío si es NULL
	relacion := Relacion{
		IdRelacion:  idRelacionSQL.String,
		FkPatente:   fkPatenteSQL.String,
		FkDocumento: fkDocumentoSQL.String,
		Fesal:       fesalSQL.String,
		Fere:        fereSQL.String,
		Prot:        protSQL.String,
		Desa:        desaSQL.String,
		Proa:        proaSQL.String,
		Dest:        destSQL.String,
	}

	return relacion, nil
}

func (r *Repository) ObtenerRelaciones(page, pageSize int) ([]RelacionList, int, int, error) {
	// Validar que page y pageSize sean mayores a cero
	if page <= 0 || pageSize <= 0 {
		return nil, 0, 0, fmt.Errorf("los valores de página y tamaño de página deben ser mayores a cero")
	}

	// Calcular el offset
	offset := (page - 1) * pageSize

	// Consulta principal para obtener las relaciones con paginación
	query := `
    SELECT tr.Idrelacion,
       tr.fkpatente,
       tr.fesal,
       tr.fere,
       tr.prot,
       tr.desa,
       tr.proa,
       tr.dest,
       IFNULL(JSON_OBJECT(
                      'Documento', c.documento,
                      'Licencia', c.licencia,
                      'Nombre', c.nombre,
                      'Apellidopa', c.apellidopa,
                      'Apellidoma', c.apellidoma,
                      'Direccion', c.direccion,
                      'Nacionalidad', c.nacionalidad,
                      'Residencia', c.residencia,
                      'Profesion', c.profesion,
                      'Fechanac', c.fechanac,
                      'Estcivil', c.estcivil,
                      'Sexo', c.sexo,
                      'Mva1', c.mva1,
                      'Mva2', c.mva2,
                      'Tipoa', c.tipoa,
                      'Nauto', c.nauto
              ), '{}') AS conductor_json,
       IFNULL(JSON_ARRAYAGG(
                      JSON_OBJECT(
                              'Documento', tdr.fkdocumento,
                              'Nombre', tdr.fknombre,
                              'ApellidoPa', tdr.fkapellidopa,
                              'ApellidoMa', tdr.fkapellidoma,
                              'Direccion', tdr.fkdireccion,
                              'Nacionalidad', tdr.fknacionalidad,
                              'Residencia', tdr.fkresidencia,
                              'Profesion', tdr.fkprofesion,
                              'Fechanac', tdr.fkfechanac,
                              'Estcivil', tdr.fkestcivil,
                              'Sexo', tdr.fksexo,
                              'Mva1', tdr.fkmva1,
                              'Mva2', tdr.fkmva2,
                              'TipoA', tdr.fktipoa
                      )
              ), '[]') AS detalles_relacion_json,
       IFNULL(JSON_OBJECT(
                      'Documento', e.Documento,
                      'Nombre', e.Nombre,
                      'Nauto', e.nauto,
                      'Resolucion', e.Resolucion,
                      'Permiso', e.Permiso
              ), '{}') AS empresa_json
		FROM tbrelacion tr
				INNER JOIN
			tbconductor c ON c.documento = tr.fkdocumento
				INNER JOIN
			tbempresa e on c.nauto = e.nauto
				LEFT JOIN
			tbdetrelacion tdr ON tr.Idrelacion = tdr.fkrelacion
		GROUP BY tr.Idrelacion
		ORDER BY (tr.Idrelacion + 0) DESC
		LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error ejecutando la consulta: %v", err)
	}
	defer rows.Close()

	var relaciones []RelacionList
	for rows.Next() {
		// Variables que manejan valores NULL
		var relacionID sql.NullString
		var patente sql.NullString
		var fechaSalida sql.NullString
		var fechaRetorno sql.NullString
		var procedenciaT sql.NullString
		var destinoA sql.NullString
		var procedenciaA sql.NullString
		var destinoT sql.NullString
		var conductorJSON sql.NullString
		var detallesJSON sql.NullString
		var detalleEmpresa sql.NullString
		err := rows.Scan(
			&relacionID,
			&patente,
			&fechaSalida,
			&fechaRetorno,
			&procedenciaT,
			&destinoA,
			&procedenciaA,
			&destinoT,
			&conductorJSON,
			&detallesJSON,
			&detalleEmpresa,
		)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("error escaneando fila: %v", err)
		}

		// Asignar valores seguros para evitar NULL
		relacion := RelacionList{
			IDRelacion:           safeString(relacionID),
			Patente:              safeString(patente),
			FechaSalida:          safeString(fechaSalida),
			FechaRetorno:         safeString(fechaRetorno),
			ProcedenciaT:         safeString(procedenciaT),
			DestinoA:             safeString(destinoA),
			ProcedenciaA:         safeString(procedenciaA),
			DestinoT:             safeString(destinoT),
			ConductorJSON:        safeJSON(conductorJSON, "{}"),
			DetallesRelacionJSON: safeJSON(detallesJSON, "[]"),
			EmpresaJSON:          safeJSON(detalleEmpresa, "{}"),
		}
		relaciones = append(relaciones, relacion)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("error en las filas de resultados: %v", err)
	}

	// Consulta adicional para contar el total de registros
	countQuery := `
    SELECT COUNT(*)
    FROM tbrelacion tr
    INNER JOIN tbconductor c ON c.documento = tr.fkdocumento
	INNER JOIN tbempresa e ON c.nauto = e.nauto
    LEFT JOIN tbdetrelacion tdr ON tr.Idrelacion = tdr.fkrelacion
    `

	var totalRecords int
	err = r.db.QueryRow(countQuery).Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error al contar registros: %v", err)
	}

	// Calcular el número total de páginas
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	// Retornar los datos, el total de registros y el número total de páginas
	return relaciones, totalRecords, totalPages, nil
}

// safeString devuelve un string seguro si el valor es NULL
func safeString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// EliminarRelacionesPorFecha elimina las relaciones y sus detalles en un intervalo de fechas
func (r *Repository) EliminarRelacionesPorFecha(fechaInicio, fechaFin string) error {
	// Iniciar una transacción
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %v", err)
	}

	// Primero, eliminar los detalles de la relación en el intervalo de fechas
	queryDetalle := `
		DELETE FROM tbdetrelacion
		WHERE fkrelacion IN (
			SELECT Idrelacion FROM tbrelacion
			WHERE created_at BETWEEN ? AND ?
		)`
	_, err = tx.Exec(queryDetalle, fechaInicio, fechaFin)
	if err != nil {
		tx.Rollback() // Revertir la transacción en caso de error
		return fmt.Errorf("error al eliminar detalles de la relación: %v", err)
	}

	// Luego, eliminar las relaciones principales en el intervalo de fechas
	queryRelacion := `
		DELETE FROM tbrelacion
		WHERE created_at BETWEEN ? AND ?`
	_, err = tx.Exec(queryRelacion, fechaInicio, fechaFin)
	if err != nil {
		tx.Rollback() // Revertir la transacción en caso de error
		return fmt.Errorf("error al eliminar relaciones principales: %v", err)
	}

	// Confirmar la transacción
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error al confirmar la transacción: %v", err)
	}

	return nil
}

// safeJSON devuelve un json.RawMessage seguro si el valor es NULL
func safeJSON(ns sql.NullString, defaultValue string) json.RawMessage {
	if ns.Valid {
		return json.RawMessage(ns.String)
	}
	return json.RawMessage(defaultValue)
}

func (r *Repository) ObtenerRelacionPorID(idRelacion string) (*RelacionList, error) {
	query := `
    SELECT tr.Idrelacion,
       tr.fkpatente,
       tr.fesal,
       tr.fere,
       tr.prot,
       tr.desa,
       tr.proa,
       tr.dest,
       CONCAT(
           '{"Documento":"', IFNULL(c.documento, ''), '",',
           '"Licencia":"', IFNULL(c.licencia, ''), '",',
           '"Nombre":"', IFNULL(c.nombre, ''), '",',
           '"Apellidopa":"', IFNULL(c.apellidopa, ''), '",',
           '"Apellidoma":"', IFNULL(c.apellidoma, ''), '",',
           '"Direccion":"', IFNULL(c.direccion, ''), '",',
           '"Nacionalidad":"', IFNULL(c.nacionalidad, ''), '",',
           '"Residencia":"', IFNULL(c.residencia, ''), '",',
           '"Profesion":"', IFNULL(c.profesion, ''), '",',
           '"Fechanac":"', IFNULL(c.fechanac, ''), '",',
           '"Estcivil":"', IFNULL(c.estcivil, ''), '",',
           '"Sexo":"', IFNULL(c.sexo, ''), '",',
           '"Mva1":"', IFNULL(c.mva1, ''), '",',
           '"Mva2":"', IFNULL(c.mva2, ''), '",',
           '"Tipoa":"', IFNULL(c.tipoa, ''), '",',
           '"Nauto":"', IFNULL(c.nauto, ''), '"}'
       ) AS conductor_json,
       CONCAT(
           '[',
           GROUP_CONCAT(
               CONCAT(
                   '{"Documento":"', IFNULL(tdr.fkdocumento, ''), '",',
                   '"Nombre":"', IFNULL(tdr.fknombre, ''), '",',
                   '"ApellidoPa":"', IFNULL(tdr.fkapellidopa, ''), '",',
                   '"ApellidoMa":"', IFNULL(tdr.fkapellidoma, ''), '",',
                   '"Direccion":"', IFNULL(tdr.fkdireccion, ''), '",',
                   '"Nacionalidad":"', IFNULL(tdr.fknacionalidad, ''), '",',
                   '"Residencia":"', IFNULL(tdr.fkresidencia, ''), '",',
                   '"Profesion":"', IFNULL(tdr.fkprofesion, ''), '",',
                   '"Fechanac":"', IFNULL(tdr.fkfechanac, ''), '",',
                   '"Estcivil":"', IFNULL(tdr.fkestcivil, ''), '",',
                   '"Sexo":"', IFNULL(tdr.fksexo, ''), '",',
                   '"Mva1":"', IFNULL(tdr.fkmva1, ''), '",',
                   '"Mva2":"', IFNULL(tdr.fkmva2, ''), '",',
                   '"TipoA":"', IFNULL(tdr.fktipoa, ''), '"}'
               )
               SEPARATOR ','
           ),
           ']'
       ) AS detalles_relacion_json,
       CONCAT(
           '{"Documento":"', IFNULL(e.Documento, ''), '",',
           '"Nombre":"', IFNULL(e.Nombre, ''), '",',
           '"Nauto":"', IFNULL(e.nauto, ''), '",',
           '"Resolucion":"', IFNULL(e.Resolucion, ''), '",',
           '"Permiso":"', IFNULL(e.Permiso, ''), '"}'
       ) AS empresa_json
FROM tbrelacion tr
        LEFT JOIN
    tbconductor c ON c.documento = tr.fkdocumento
        LEFT JOIN
    tbempresa e on c.nauto = e.nauto
        LEFT JOIN
    tbdetrelacion tdr ON tr.Idrelacion = tdr.fkrelacion
WHERE tr.Idrelacion = ?
GROUP BY tr.Idrelacion
    `

	row := r.db.QueryRow(query, idRelacion)

	var relacionID, patente, fechaSalida, fechaRetorno sql.NullString
	var procedenciaT, destinoA, procedenciaA, destinoT sql.NullString
	var conductorJSON, detallesJSON, detalleEmpresa sql.NullString

	err := row.Scan(
		&relacionID, &patente, &fechaSalida, &fechaRetorno,
		&procedenciaT, &destinoA, &procedenciaA, &destinoT,
		&conductorJSON, &detallesJSON, &detalleEmpresa,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontró relación con ID: %s", idRelacion)
		}
		return nil, fmt.Errorf("error escaneando relación: %v", err)
	}

	relacion := &RelacionList{
		IDRelacion:           safeString(relacionID),
		Patente:              safeString(patente),
		FechaSalida:          safeString(fechaSalida),
		FechaRetorno:         safeString(fechaRetorno),
		ProcedenciaT:         safeString(procedenciaT),
		DestinoA:             safeString(destinoA),
		ProcedenciaA:         safeString(procedenciaA),
		DestinoT:             safeString(destinoT),
		ConductorJSON:        safeJSON(conductorJSON, "{}"),
		DetallesRelacionJSON: safeJSON(detallesJSON, "[]"),
		EmpresaJSON:          safeJSON(detalleEmpresa, "{}"),
	}

	return relacion, nil
}
