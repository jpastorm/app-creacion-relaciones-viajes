package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2/log"
)

type Conductor struct {
	Documento     string
	Licencia      string
	Nombre        string
	ApellidoPa    string
	ApellidoMa    string
	Direccion     string
	Nacionalidad  string
	Residencia    string
	Profesion     string
	FechaNac      string
	EstCivil      string
	Sexo          string
	Mva1          string
	Mva2          string
	Tipoa         string
	Nauto         string
	TipoDocumento string
}

func (r *Repository) ListarConductores(page int, pageSize int, isConductor bool, documento string, nombre string) ([]Conductor, int, int, error) {
	if page <= 0 || pageSize <= 0 {
		log.Errorf("los valores de página y tamaño de página deben ser mayores a cero")
		return nil, 0, 0, fmt.Errorf("los valores de página y tamaño de página deben ser mayores a cero")
	}

	offset := (page - 1) * pageSize

	// Construcción de la consulta dinámica
	baseQuery := `
        SELECT c.documento, c.licencia, c.nombre, c.apellidopa, c.apellidoma, c.direccion, c.nacionalidad,
               c.residencia, c.profesion, c.fechanac, c.estcivil, c.sexo, c.mva1, c.mva2, c.tipoa, c.nauto, c.tipo_documento
        FROM tbconductor c 
        WHERE 1=1`

	// Filtrar conductores según la condición
	if isConductor {
		baseQuery += " AND (c.nauto IS NOT NULL OR c.licencia IS NOT NULL)"
	} else {
		baseQuery += " AND (c.nauto IS NULL AND c.licencia IS NULL)"
	}

	// Filtrar por documento y/o nombre
	var filterParams []interface{}
	if documento != "" {
		baseQuery += " AND LOWER(c.documento) LIKE LOWER(?)" // Búsqueda insensible a mayúsculas
		filterParams = append(filterParams, "%"+documento+"%")
	}
	if nombre != "" {
		baseQuery += " AND LOWER(c.nombre) LIKE LOWER(?)" // Búsqueda insensible a mayúsculas
		filterParams = append(filterParams, "%"+nombre+"%")
	}
	// Agregar paginación
	baseQuery += " LIMIT ? OFFSET ?"
	queryParams := append(filterParams, pageSize, offset)

	// Ejecutar consulta
	rows, err := r.db.Query(baseQuery, queryParams...)
	if err != nil {
		log.Errorf("error al ejecutar la consulta: %v", err)
		return nil, 0, 0, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer rows.Close()

	var conductores []Conductor
	for rows.Next() {
		var c Conductor
		var documento, licencia, nombre, apellidoPa, apellidoMa,
			direccion, nacionalidad, residencia,
			profesion, fechanac, estcivil,
			sexo, mva1, mva2, tipoa, nauto, tipoDocumento sql.NullString

		if err := rows.Scan(
			&documento, &licencia, &nombre, &apellidoPa, &apellidoMa, &direccion, &nacionalidad,
			&residencia, &profesion, &fechanac, &estcivil, &sexo, &mva1, &mva2, &tipoa, &nauto, &tipoDocumento,
		); err != nil {
			log.Errorf("error al escanear fila: %v", err)
			return nil, 0, 0, fmt.Errorf("error al escanear fila: %v", err)
		}

		c.Documento = getStringOrEmpty(documento)
		c.Licencia = getStringOrEmpty(licencia)
		c.Nombre = getStringOrEmpty(nombre)
		c.ApellidoPa = getStringOrEmpty(apellidoPa)
		c.ApellidoMa = getStringOrEmpty(apellidoMa)
		c.Direccion = getStringOrEmpty(direccion)
		c.Nacionalidad = getStringOrEmpty(nacionalidad)
		c.Residencia = getStringOrEmpty(residencia)
		c.Profesion = getStringOrEmpty(profesion)
		c.FechaNac = getStringOrEmpty(fechanac)
		c.EstCivil = getStringOrEmpty(estcivil)
		c.Sexo = getStringOrEmpty(sexo)
		c.Mva1 = getStringOrEmpty(mva1)
		c.Mva2 = getStringOrEmpty(mva2)
		c.Tipoa = getStringOrEmpty(tipoa)
		c.Nauto = getStringOrEmpty(nauto)
		c.TipoDocumento = getStringOrEmpty(tipoDocumento)
		conductores = append(conductores, c)
	}

	// Contar total de registros
	countQuery := `SELECT COUNT(*) FROM tbconductor c WHERE 1=1`
	if isConductor {
		countQuery += " AND c.nauto IS NOT NULL AND c.licencia IS NOT NULL"
	} else {
		countQuery += " AND (c.nauto IS NULL OR c.licencia IS NULL)"
	}
	if documento != "" {
		countQuery += " AND LOWER(c.documento) LIKE LOWER(?)"
	}
	if nombre != "" {
		countQuery += " AND LOWER(c.nombre) LIKE LOWER(?)"
	}

	var totalRecords int
	err = r.db.QueryRow(countQuery, filterParams...).Scan(&totalRecords)
	if err != nil {
		log.Errorf("error al contar registros: %v", err)
		return nil, 0, 0, fmt.Errorf("error al contar registros: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return conductores, totalRecords, totalPages, nil
}

// Función auxiliar para manejar valores nulos
func getStringOrEmpty(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// Listar conductores por documento
func (r *Repository) ListarConductorPorDocumento(documento string) (Conductor, error) {
	query := `
		SELECT c.documento, c.licencia, c.nombre, c.apellidopa, c.apellidoma, c.direccion, c.nacionalidad,
			   c.residencia, c.profesion, c.fechanac, c.estcivil, c.sexo, c.mva1, c.mva2, c.tipoa, c.nauto, c.tipo_documento
		FROM tbconductor c
		WHERE c.documento = ?
	`
	row := r.db.QueryRow(query, documento)

	// Variables intermedias con tipos Null para manejar valores nulos
	var documentoSQL, licenciaSQL, nombreSQL, apellidoPaSQL, apellidoMaSQL,
		direccionSQL, nacionalidadSQL, residenciaSQL, profesionSQL, fechaNacSQL,
		estCivilSQL, sexoSQL, mva1SQL, mva2SQL, tipoaSQL, nautoSQL, tipoDocumentoSQL sql.NullString

	// Escaneamos los resultados en las variables intermedias
	err := row.Scan(
		&documentoSQL, &licenciaSQL, &nombreSQL, &apellidoPaSQL, &apellidoMaSQL, &direccionSQL,
		&nacionalidadSQL, &residenciaSQL, &profesionSQL, &fechaNacSQL, &estCivilSQL, &sexoSQL,
		&mva1SQL, &mva2SQL, &tipoaSQL, &nautoSQL, &tipoDocumentoSQL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Conductor{}, fmt.Errorf("no se encontró ningún conductor con el documento %s", documento)
		}
		return Conductor{}, fmt.Errorf("error al escanear fila: %v", err)
	}

	// Asignamos los valores al objeto Conductor
	c := Conductor{
		Documento:     documentoSQL.String,
		Licencia:      licenciaSQL.String,
		Nombre:        nombreSQL.String,
		ApellidoPa:    apellidoPaSQL.String,
		ApellidoMa:    apellidoMaSQL.String,
		Direccion:     direccionSQL.String,
		Nacionalidad:  nacionalidadSQL.String,
		Residencia:    residenciaSQL.String,
		Profesion:     profesionSQL.String,
		FechaNac:      fechaNacSQL.String,
		EstCivil:      estCivilSQL.String,
		Sexo:          sexoSQL.String,
		Mva1:          mva1SQL.String,
		Mva2:          mva2SQL.String,
		Tipoa:         tipoaSQL.String,
		Nauto:         nautoSQL.String,
		TipoDocumento: tipoDocumentoSQL.String,
	}

	return c, nil
}

// Agregar un conductor
func (r *Repository) AgregarConductor(c Conductor) error {
	convertToNullString := func(value string) sql.NullString {
		if value == "" {
			return sql.NullString{String: "", Valid: false} // NULL
		}
		return sql.NullString{String: value, Valid: true} // Valor no NULL
	}

	// Convertir los campos Licencia y Nauto
	licencia := convertToNullString(c.Licencia)
	nauto := convertToNullString(c.Nauto)
	tipoDocumento := convertToNullString(c.TipoDocumento)
	query := `
		INSERT INTO tbconductor (documento, licencia, nombre, apellidopa, apellidoma, direccion, nacionalidad, residencia, profesion, fechanac, estcivil, sexo, mva1, mva2, tipoa, nauto, tipo_documento)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, c.Documento, licencia, c.Nombre, c.ApellidoPa, c.ApellidoMa, c.Direccion, c.Nacionalidad, c.Residencia, c.Profesion, c.FechaNac, c.EstCivil, c.Sexo, c.Mva1, c.Mva2, c.Tipoa, nauto, tipoDocumento)
	if err != nil {
		return fmt.Errorf("error al insertar conductor: %v", err)
	}
	return nil
}

// Modificar un conductor
func (r *Repository) ModificarConductor(c Conductor) error {
	convertToNullString := func(value string) sql.NullString {
		if value == "" {
			return sql.NullString{String: "", Valid: false} // NULL
		}
		return sql.NullString{String: value, Valid: true} // Valor no NULL
	}

	// Convertir los campos Licencia y Nauto
	licencia := convertToNullString(c.Licencia)
	nauto := convertToNullString(c.Nauto)
	query := `
		UPDATE tbconductor
		SET licencia = ?, nombre = ?, apellidopa = ?, apellidoma = ?,
		 	direccion = ?, nacionalidad = ?, residencia = ?, profesion = ?,
		  	fechanac = ?, estcivil = ?, sexo = ?, mva1 = ?, mva2 = ?, tipoa = ?,
		   	nauto = ?, tipo_documento = ?
		WHERE documento = ?	
	`
	_, err := r.db.Exec(query, licencia, c.Nombre, c.ApellidoPa, c.ApellidoMa, c.Direccion, c.Nacionalidad, c.Residencia, c.Profesion, c.FechaNac, c.EstCivil, c.Sexo, c.Mva1, c.Mva2, c.Tipoa, nauto, c.TipoDocumento, c.Documento)
	if err != nil {
		fmt.Println("err", err.Error())
		return fmt.Errorf("error al modificar conductor: %v", err)
	}
	return nil
}

// CrearOActualizarConductor
// func (r *Repository) CrearOActualizarConductor(c Conductor) error {
// 	// Primero, verificamos si el conductor ya existe en la base de datos
// 	var licencia sql.NullString
// 	var nauto sql.NullString
// 	checkQuery := `SELECT licencia, nauto FROM tbconductor WHERE documento = ?;`

// 	// Verificar si el conductor ya existe por documento
// 	err := r.db.QueryRow(checkQuery, c.Documento).Scan(&licencia, &nauto)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// No existe, realizar un INSERT
// 			return r.AgregarConductor(c)
// 		}
// 		// Error inesperado
// 		return fmt.Errorf("error al verificar existencia de conductor: %v", err)
// 	}

// 	// Si existe, actualizar datos y realizar un UPDATE
// 	c.Licencia = licencia.String
// 	c.Nauto = nauto.String
// 	return r.ModificarConductor(c)
// }

func (r *Repository) EliminarConductor(docu string) error {
	query := `
        DELETE FROM tbconductor
        WHERE documento = ?
    `

	_, err := r.db.Exec(query, docu)
	if err != nil {
		fmt.Printf("error al eliminar conductor: %v", err)
		return fmt.Errorf("error al eliminar conductor: %v", err)
	}

	return nil
}

// CrearOActualizarConductor
func (r *Repository) CrearOActualizarConductor(c Conductor) error {
	// Verificar si el conductor ya existe en la base de datos
	var licencia sql.NullString
	var nauto sql.NullString
	checkQuery := `SELECT licencia, nauto FROM tbconductor WHERE documento = ?;`

	err := r.db.QueryRow(checkQuery, c.Documento).Scan(&licencia, &nauto)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existe, verificar si el nauto ya pertenece a otro conductor
			if c.Nauto != "" {
				existeOtro, err := r.existeOtroConductorConNauto(c.Nauto, c.Documento)
				if err != nil {
					return fmt.Errorf("error al verificar nauto: %v", err)
				}
				if existeOtro {
					err := r.liberarNauto(c.Nauto)
					if err != nil {
						return fmt.Errorf("error al liberar nauto: %v", err)
					}
				}
			}
			return r.AgregarConductor(c)
		}
		return fmt.Errorf("error al verificar existencia de conductor: %v", err)
	}

	// Si existe y se va a actualizar el nauto, verificar si ya pertenece a otro conductor
	if c.Nauto != "" && c.Nauto != nauto.String {
		existeOtro, err := r.existeOtroConductorConNauto(c.Nauto, c.Documento)
		if err != nil {
			return fmt.Errorf("error al verificar nauto: %v", err)
		}
		if existeOtro {
			err := r.liberarNauto(c.Nauto)
			if err != nil {
				return fmt.Errorf("error al liberar nauto: %v", err)
			}
		}
	}

	// Actualizar conductor
	return r.ModificarConductor(c)
}

// existeOtroConductorConNauto verifica si el `nauto` ya está en uso por otro conductor distinto al actual
func (r *Repository) existeOtroConductorConNauto(nauto, documento string) (bool, error) {
	var otroDocumento string
	query := `SELECT documento FROM tbconductor WHERE nauto = ? AND documento != ?;`
	err := r.db.QueryRow(query, nauto, documento).Scan(&otroDocumento)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No hay otro conductor con ese nauto
		}
		return false, fmt.Errorf("error al verificar otro conductor con nauto: %v", err)
	}
	return true, nil // Sí hay otro conductor con ese nauto
}

// liberarNauto actualiza el conductor que tenga ese nauto y lo establece en NULL
func (r *Repository) liberarNauto(nauto string) error {
	updateQuery := `UPDATE tbconductor SET nauto = NULL WHERE nauto = ?;`
	_, err := r.db.Exec(updateQuery, nauto)
	if err != nil {
		return fmt.Errorf("error al liberar nauto de otro conductor: %v", err)
	}
	return nil
}
