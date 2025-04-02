package service

import (
	"api-terminal/repository"
	"errors"
	"fmt"
)

// DetalleRelacionService es la estructura que contiene el repositorio
type DetalleRelacionService struct {
	repo *repository.Repository
}

// NewDetalleRelacionService crea una nueva instancia del servicio
func NewDetalleRelacionService(repo *repository.Repository) *DetalleRelacionService {
	return &DetalleRelacionService{repo: repo}
}

// AgregarDetalleRelacion agrega un nuevo detalle de relación
func (s *DetalleRelacionService) AgregarDetalleRelacion(d repository.DetalleRelacion) error {
	// Validaciones básicas
	if d.FkRelacion == "" || d.FkDocumento == "" {
		return errors.New("los campos documento y FkRelacion son obligatorio")
	}

	// Llamada al repositorio para agregar el detalle de relación
	err := s.repo.AgregarDetalleRelacion(d)
	if err != nil {
		return fmt.Errorf("error al agregar detalle de relación: %v", err)
	}

	return nil
}
