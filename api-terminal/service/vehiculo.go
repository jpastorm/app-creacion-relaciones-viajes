package service

import (
	"api-terminal/repository"
	"errors"
	"fmt"
)

// VehiculoService es la estructura que contiene el repositorio
type VehiculoService struct {
	repo *repository.Repository
}

// NewVehiculoService crea una nueva instancia del servicio
func NewVehiculoService(repo *repository.Repository) *VehiculoService {
	return &VehiculoService{repo: repo}
}

// ListarVehiculosPaginados lista los vehículos con paginación
func (s *VehiculoService) ListarVehiculosPaginados(page int, pageSize int, patente string) ([]repository.Vehiculo, int, int, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, 0, errors.New("los valores de página y tamaño de página deben ser mayores a cero")
	}

	vehiculos, totalRecords, totalPages, err := s.repo.ListarVehiculos(page, pageSize, patente)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error al listar vehículos: %v", err)
	}

	if len(vehiculos) == 0 {
		return nil, 0, 0, nil
	}

	return vehiculos, totalRecords, totalPages, nil
}

// ObtenerVehiculoPorPatente obtiene un vehículo por su patente
func (s *VehiculoService) ObtenerVehiculoPorPatente(patente string) (repository.Vehiculo, error) {
	if patente == "" {
		return repository.Vehiculo{}, errors.New("la patente no puede estar vacía")
	}

	vehiculo, err := s.repo.ListarVehiculoPorPatente(patente)
	if err != nil {
		return repository.Vehiculo{}, fmt.Errorf("error al obtener vehículo: %v", err)
	}

	return vehiculo, nil
}

// ObtenerVehiculoConConductorPorNroAuto obtiene un vehículo con conductor por número de auto
func (s *VehiculoService) ObtenerVehiculoConConductorPorNroAuto(nroAuto string) (repository.VehiculoConConductorEmpresa, error) {
	if nroAuto == "" {
		return repository.VehiculoConConductorEmpresa{}, errors.New("el número de auto no puede estar vacío")
	}

	vcc, err := s.repo.ListarVehiculoConConductorPorNroAuto(nroAuto)
	if err != nil {
		return repository.VehiculoConConductorEmpresa{}, fmt.Errorf("error al obtener vehículo con conductor: %v", err)
	}

	return vcc, nil
}

// AgregarVehiculo agrega un nuevo vehículo
func (s *VehiculoService) AgregarVehiculo(v repository.Vehiculo) error {
	if v.Patente == "" || v.NroAuto == "" {
		return errors.New("los campos patente y NroAuto son obligatorios")
	}

	err := s.repo.AgregarVehiculo(v)
	if err != nil {
		return fmt.Errorf("error al agregar vehículo: %v", err)
	}

	return nil
}

// ActualizarVehiculo actualiza un vehículo existente
func (s *VehiculoService) ActualizarVehiculo(v repository.Vehiculo) error {
	if v.Patente == "" || v.NroAuto == "" {
		return errors.New("los campos patente y NroAuto son obligatorios")
	}

	err := s.repo.ModificarVehiculo(v)
	if err != nil {
		return fmt.Errorf("error al actualizar vehículo: %v", err)
	}

	return nil
}

func (s *VehiculoService) EliminarVehiculo(docu string) error {
	if docu == "" {
		return errors.New("el documento no puede estar vacio")
	}

	return s.repo.EliminarVehiculo(docu)
}

// CrearOActualizarVehiculo crea o actualiza un conductor según su existencia
func (s *VehiculoService) CrearOActualizarVehiculo(c repository.Vehiculo) error {
	if c.Patente == "" || c.NroAuto == "" {
		return errors.New("los campos patente y NroAuto son obligatorios")
	}

	err := s.repo.CreateOrUpdateVehiculo(c)
	if err != nil {
		return fmt.Errorf("error al crear o actualizar vehiculo: %v", err)
	}

	return nil
}
