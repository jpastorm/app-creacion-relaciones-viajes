package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2/log"
)

type Vehiculo struct {
	Patente string
	Tipo    string
	Modelo  string
	Motor   string
	Marca   string
	Anio    string
	Pais    string
	Empresa string
	Chasis  string
	NroAuto string
	Costo   string
}

type VehiculoConConductorEmpresa struct {
	Vehiculo  `json:"vehiculo"`
	Conductor `json:"conductor"`
	Empresa   `json:"empresa"`
}

func (r *Repository) ListarVehiculos(page int, pageSize int, patente string) ([]Vehiculo, int, int, error) {
	// Validar que page y pageSize sean mayores a cero
	if page <= 0 || pageSize <= 0 {
		return nil, 0, 0, fmt.Errorf("los valores de página y tamaño de página deben ser mayores a cero")
	}

	// Calcular el offset
	offset := (page - 1) * pageSize

	// Consulta base para obtener los vehículos
	query := `
        SELECT v.patente, v.tipo, v.modelo, v.motor, v.marca, v.anio, v.pais, v.empresa, v.chasis, v.nroauto, v.costo 
        FROM tbvehiculo v
    `
	// Consulta base para contar el total de registros
	countQuery := `
        SELECT COUNT(*) 
        FROM tbvehiculo
    `

	// Variables para almacenar los argumentos de la consulta
	var args []interface{}
	var countArgs []interface{}

	// Agregar el filtro de patente si se proporciona
	if patente != "" {
		query += " WHERE v.patente LIKE ?"
		countQuery += " WHERE patente LIKE ?"
		patenteFilter := "%" + patente + "%"
		args = append(args, patenteFilter)
		countArgs = append(countArgs, patenteFilter)
	}

	// Agregar paginación a la consulta principal
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	// Ejecutar la consulta para obtener los vehículos
	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Errorf("error al ejecutar la consulta: %v", err)
		return nil, 0, 0, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer rows.Close()

	var vehiculos []Vehiculo
	for rows.Next() {
		var v Vehiculo
		if err := rows.Scan(&v.Patente, &v.Tipo, &v.Modelo, &v.Motor, &v.Marca, &v.Anio, &v.Pais, &v.Empresa, &v.Chasis, &v.NroAuto, &v.Costo); err != nil {
			log.Errorf("error al escanear fila: %v", err)
			return nil, 0, 0, fmt.Errorf("error al escanear fila: %v", err)
		}
		vehiculos = append(vehiculos, v)
	}

	// Ejecutar la consulta para contar el total de registros
	var totalRecords int
	err = r.db.QueryRow(countQuery, countArgs...).Scan(&totalRecords)
	if err != nil {
		log.Errorf("error al contar registros: %v", err)
		return nil, 0, 0, fmt.Errorf("error al contar registros: %v", err)
	}

	// Calcular el número total de páginas
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	// Retornar los vehículos, el total de registros y el número total de páginas
	return vehiculos, totalRecords, totalPages, nil
}

// Listar vehículos por patente
func (r *Repository) ListarVehiculoPorPatente(patente string) (Vehiculo, error) {
	query := `
		SELECT v.patente, v.tipo, v.modelo, v.motor, v.marca, v.anio, v.pais, v.empresa, v.chasis, v.nroauto, v.costo 
		FROM tbvehiculo v
		WHERE v.patente = ?
	`
	row := r.db.QueryRow(query, patente)

	// Variables intermedias con sql.NullString para manejar valores nulos
	var patenteSQL, tipoSQL, modeloSQL, motorSQL, marcaSQL,
		anioSQL, paisSQL, empresaSQL, chasisSQL, nroAutoSQL, costoSQL sql.NullString

	// Escaneamos los valores en las variables intermedias
	err := row.Scan(
		&patenteSQL, &tipoSQL, &modeloSQL, &motorSQL, &marcaSQL,
		&anioSQL, &paisSQL, &empresaSQL, &chasisSQL, &nroAutoSQL, &costoSQL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Vehiculo{}, fmt.Errorf("no se encontró ningún vehículo con la patente %s", patente)
		}
		return Vehiculo{}, fmt.Errorf("error al escanear fila: %v", err)
	}

	// Asignamos los valores al objeto Vehiculo
	v := Vehiculo{
		Patente: patenteSQL.String,
		Tipo:    tipoSQL.String,
		Modelo:  modeloSQL.String,
		Motor:   motorSQL.String,
		Marca:   marcaSQL.String,
		Anio:    anioSQL.String,
		Pais:    paisSQL.String,
		Empresa: empresaSQL.String,
		Chasis:  chasisSQL.String,
		NroAuto: nroAutoSQL.String,
		Costo:   costoSQL.String,
	}

	return v, nil
}

// Listar vehículos con conductor por número de auto
func (r *Repository) ListarVehiculoConConductorPorNroAuto(nroAuto string) (VehiculoConConductorEmpresa, error) {
	query := `
		SELECT v.patente, v.tipo, v.modelo, v.motor, v.marca, v.anio, v.pais, v.empresa, v.chasis, v.nroauto, v.costo, 
			   c.documento, c.licencia, c.nombre, c.apellidopa, c.apellidoma, c.direccion, c.nacionalidad,
			   c.residencia, c.profesion, c.fechanac, c.estcivil, c.sexo, c.mva1, c.mva2, c.tipoa, c.nauto,
			   e.Documento as documentoEmpresa, e.Nombre as nombreEmpresa, e.Permiso as permisoEmpresa, e.Resolucion as resolucionEmpresa

		FROM tbvehiculo v 
		LEFT JOIN tbconductor c ON v.nroauto = c.nauto
		LEFT JOIN tbempresa e ON v.nroauto = e.nauto
		WHERE v.nroauto = ?
	`
	row := r.db.QueryRow(query, nroAuto)

	// Variables intermedias con sql.NullString para manejar valores nulos
	var patenteSQL, tipoSQL, modeloSQL, motorSQL, marcaSQL sql.NullString
	var anioSQL, paisSQL, empresaSQL, chasisSQL, nroAutoSQL, costoSQL sql.NullString
	var documentoSQL, licenciaSQL, nombreSQL, apellidoPaSQL, apellidoMaSQL sql.NullString
	var direccionSQL, nacionalidadSQL, residenciaSQL, profesionSQL, fechaNacSQL sql.NullString
	var estCivilSQL, sexoSQL, mva1SQL, mva2SQL, tipoaSQL, nautoSQL sql.NullString
	var documentoEmpresa, nombreEmpresa, permisoEmpresa, resolucionEmpresa sql.NullString
	// Escaneamos los valores en las variables intermedias
	err := row.Scan(
		&patenteSQL, &tipoSQL, &modeloSQL, &motorSQL, &marcaSQL, &anioSQL, &paisSQL, &empresaSQL, &chasisSQL, &nroAutoSQL, &costoSQL,
		&documentoSQL, &licenciaSQL, &nombreSQL, &apellidoPaSQL, &apellidoMaSQL, &direccionSQL, &nacionalidadSQL,
		&residenciaSQL, &profesionSQL, &fechaNacSQL, &estCivilSQL, &sexoSQL, &mva1SQL, &mva2SQL, &tipoaSQL, &nautoSQL,
		&documentoEmpresa, &nombreEmpresa, &permisoEmpresa, &resolucionEmpresa,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return VehiculoConConductorEmpresa{}, fmt.Errorf("no se encontró ningún vehículo con conductor para el nroAuto %s", nroAuto)
		}
		return VehiculoConConductorEmpresa{}, fmt.Errorf("error al escanear fila: %v", err)
	}

	// Asignamos los valores a las estructuras anidadas
	vcc := VehiculoConConductorEmpresa{
		Vehiculo: Vehiculo{
			Patente: patenteSQL.String,
			Tipo:    tipoSQL.String,
			Modelo:  modeloSQL.String,
			Motor:   motorSQL.String,
			Marca:   marcaSQL.String,
			Anio:    anioSQL.String,
			Pais:    paisSQL.String,
			Empresa: empresaSQL.String,
			Chasis:  chasisSQL.String,
			NroAuto: nroAutoSQL.String,
			Costo:   costoSQL.String,
		},
		Conductor: Conductor{
			Documento:    documentoSQL.String,
			Licencia:     licenciaSQL.String,
			Nombre:       nombreSQL.String,
			ApellidoPa:   apellidoPaSQL.String,
			ApellidoMa:   apellidoMaSQL.String,
			Direccion:    direccionSQL.String,
			Nacionalidad: nacionalidadSQL.String,
			Residencia:   residenciaSQL.String,
			Profesion:    profesionSQL.String,
			FechaNac:     fechaNacSQL.String,
			EstCivil:     estCivilSQL.String,
			Sexo:         sexoSQL.String,
			Mva1:         mva1SQL.String,
			Mva2:         mva2SQL.String,
			Tipoa:        tipoaSQL.String,
			Nauto:        nautoSQL.String,
		},
		Empresa: Empresa{
			Documento:  documentoEmpresa.String,
			Permiso:    permisoEmpresa.String,
			Resolucion: residenciaSQL.String,
			Nombre:     nombreEmpresa.String,
		},
	}

	return vcc, nil
}

// Agregar un vehículo
func (r *Repository) AgregarVehiculo(v Vehiculo) error {
	query := `
		INSERT INTO tbvehiculo (patente, tipo, modelo, motor, marca, anio, pais, empresa, chasis, nroauto, costo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, v.Patente, v.Tipo, v.Modelo, v.Motor, v.Marca, v.Anio, v.Pais, v.Empresa, v.Chasis, v.NroAuto, v.Costo)
	if err != nil {
		return fmt.Errorf("error al insertar vehículo: %v", err)
	}
	return nil
}

// Modificar un vehículo
func (r *Repository) ModificarVehiculo(v Vehiculo) error {
	query := `
		UPDATE tbvehiculo
		SET tipo = ?, modelo = ?, motor = ?, marca = ?, anio = ?, pais = ?, empresa = ?, chasis = ?, nroauto = ?, costo = ?
		WHERE patente = ?
	`
	_, err := r.db.Exec(query, v.Tipo, v.Modelo, v.Motor, v.Marca, v.Anio, v.Pais, v.Empresa, v.Chasis, v.NroAuto, v.Costo, v.Patente)
	if err != nil {
		return fmt.Errorf("error al modificar vehículo: %v", err)
	}
	return nil
}

func (r *Repository) EliminarVehiculo(docu string) error {
	query := `
        DELETE FROM tbvehiculo
        WHERE patente = ?
    `

	_, err := r.db.Exec(query, docu)
	if err != nil {
		return fmt.Errorf("error al eliminar vehiculo: %v", err)
	}

	return nil
}

// vehiculoExiste verifica si un vehículo con la patente dada ya existe en la base de datos.
func (r *Repository) vehiculoExiste(patente string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM tbvehiculo 
		WHERE patente = ?
	`
	var count int
	err := r.db.QueryRow(query, patente).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error al verificar si el vehículo existe: %v", err)
	}
	return count > 0, nil
}

// CreateOrUpdateVehiculo verifica si un vehículo existe por su patente.
// Si existe, lo actualiza; si no existe, lo crea.
func (r *Repository) CreateOrUpdateVehiculo(v Vehiculo) error {
	// Verificar si el vehículo ya existe
	exists, err := r.vehiculoExiste(v.Patente)
	if err != nil {
		return fmt.Errorf("error al verificar si el vehículo existe: %v", err)
	}

	// Si el vehículo existe, actualizarlo
	if exists {
		return r.ModificarVehiculo(v)
	}

	// Si el vehículo no existe, crearlo
	return r.AgregarVehiculo(v)
}
