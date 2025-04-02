package service

import (
	"errors"
	"fmt"

	"api-terminal/repository"
)

// EmpresaService es la estructura que contiene el repositorio
type EmpresaService struct {
	repo *repository.Repository
}

// NewEmpresaService crea una nueva instancia del servicio
func NewEmpresaService(repo *repository.Repository) *EmpresaService {
	return &EmpresaService{repo: repo}
}

// CrearOActualizarEmpresa crea o actualiza una empresa
func (s *EmpresaService) CrearOActualizarEmpresa(e repository.Empresa) error {
	// Validaciones básicas
	if e.Nauto == "" {
		return errors.New("El campo nauto es obligatorio")
	}

	// Llamada al repositorio para crear o actualizar la empresa
	err := s.repo.CrearOActualizarEmpresa(e)
	if err != nil {
		return fmt.Errorf("error al crear o actualizar empresa: %v", err)
	}

	return nil
}

// ObtenerEmpresas obtiene todas las empresas
func (s *EmpresaService) ObtenerEmpresas() ([]repository.Empresa, error) {
	// Llamada al repositorio para obtener todas las empresas
	empresas, err := s.repo.ObtenerEmpresas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresas: %v", err)
	}

	return empresas, nil
}

// ObtenerEmpresaPorID obtiene una empresa por su ID (nauto)
func (s *EmpresaService) ObtenerEmpresaPorID(id string) (*repository.Empresa, error) {
	// Validación básica
	if id == "" {
		return nil, errors.New("el ID de la empresa no puede estar vacío")
	}

	// Llamada al repositorio para obtener la empresa por ID
	empresa, err := s.repo.ObtenerEmpresaPorID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresa por ID: %v", err)
	}

	return empresa, nil
}

// EliminarEmpresa elimina una empresa por su ID (nauto)
func (s *EmpresaService) EliminarEmpresa(id string) error {
	// Validación básica
	if id == "" {
		return errors.New("el ID de la empresa no puede estar vacío")
	}

	// Llamada al repositorio para eliminar la empresa
	err := s.repo.EliminarEmpresa(id)
	if err != nil {
		return fmt.Errorf("error al eliminar empresa: %v", err)
	}

	return nil
}
