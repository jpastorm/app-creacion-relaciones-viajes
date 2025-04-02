package service

import (
	"api-terminal/repository"
	"errors"
	"fmt"
)

// Service es la estructura que contiene el repositorio
type ConductorService struct {
	repo *repository.Repository
}

// NewConductorService crea una nueva instancia del servicio
func NewConductorService(repo *repository.Repository) *ConductorService {
	return &ConductorService{repo: repo}
}

// ListarConductoresPaginados lista los conductores con paginación
func (s *ConductorService) ListarConductoresPaginados(page int, pageSize int, isConductor bool, documento, nombre string) ([]repository.Conductor, int, int, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, 0, errors.New("los valores de página y tamaño de página deben ser mayores a cero")
	}

	conductores, totalRecords, totalPages, err := s.repo.ListarConductores(page, pageSize, isConductor, documento, nombre)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error al listar conductores: %v", err)
	}

	if len(conductores) == 0 {
		return nil, 0, 0, nil
	}

	return conductores, totalRecords, totalPages, nil
}

// ObtenerConductorPorDocumento obtiene un conductor por su documento
func (s *ConductorService) ObtenerConductorPorDocumento(documento string) (repository.Conductor, error) {
	if documento == "" {
		return repository.Conductor{}, errors.New("el documento no puede estar vacío")
	}

	conductor, err := s.repo.ListarConductorPorDocumento(documento)
	if err != nil {
		return repository.Conductor{}, fmt.Errorf("error al obtener conductor: %v", err)
	}

	return conductor, nil
}

// CrearConductor crea un nuevo conductor
func (s *ConductorService) CrearConductor(c repository.Conductor) error {
	if c.Documento == "" {
		return errors.New("El campo documento es obligatorio")
	}

	err := s.repo.AgregarConductor(c)
	if err != nil {
		return fmt.Errorf("error al crear conductor: %v", err)
	}

	return nil
}

// ActualizarConductor actualiza un conductor existente
func (s *ConductorService) ActualizarConductor(c repository.Conductor) error {
	if c.Documento == "" {
		return errors.New("el documento no puede estar vacío")
	}

	err := s.repo.ModificarConductor(c)
	if err != nil {
		return fmt.Errorf("error al actualizar conductor: %v", err)
	}

	return nil
}

// CrearOActualizarConductor crea o actualiza un conductor según su existencia
func (s *ConductorService) CrearOActualizarConductor(c repository.Conductor) error {
	if c.Documento == "" {
		return errors.New("el documento no puede estar vacío")
	}

	err := s.repo.CrearOActualizarConductor(c)
	if err != nil {
		return fmt.Errorf("error al crear o actualizar conductor: %v", err)
	}

	return nil
}

func (s *ConductorService) EliminarConductor(docu string) error {
	if docu == "" {
		return errors.New("el documento no puede estar vacio")
	}

	return s.repo.EliminarConductor(docu)
}
