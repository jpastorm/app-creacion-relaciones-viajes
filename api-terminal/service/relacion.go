package service

import (
	"api-terminal/repository"
	"errors"
	"fmt"
	"strings"
)

// RelacionService es la estructura que contiene el repositorio
type RelacionService struct {
	repo *repository.Repository
}

// NewRelacionService crea una nueva instancia del servicio
func NewRelacionService(repo *repository.Repository) *RelacionService {
	return &RelacionService{repo: repo}
}

// AgregarRelacion agrega una nueva relación
func (s *RelacionService) AgregarRelacion(rel repository.Relacion) error {
	// Validaciones básicas
	if rel.IdRelacion == "" {
		return errors.New("el campo ID Relacion es obligatorio")
	}

	// Llamada al repositorio para agregar la relación
	err := s.repo.AgregarRelacion(rel)
	if err != nil {
		return fmt.Errorf("error al agregar relación: %v", err)
	}

	return nil
}

// ObtenerUltimaRelacion obtiene la última relación registrada
func (s *RelacionService) ObtenerUltimaRelacion() (repository.Relacion, error) {
	relacion, err := s.repo.ListarUltimaRelacion()
	if err != nil {
		return repository.Relacion{}, fmt.Errorf("error al obtener la última relación: %v", err)
	}

	return relacion, nil
}

func (s *RelacionService) ObtenerRelaciones(page, pageSize int) ([]repository.RelacionList, int, int, error) {
	relacion, totalRecords, totalPages, err := s.repo.ObtenerRelaciones(page, pageSize)
	if err != nil {
		return []repository.RelacionList{}, 0, 0, fmt.Errorf("error al obtener la relacion paginada: %v", err)
	}

	return relacion, totalRecords, totalPages, nil
}

func (s *RelacionService) EliminarRelacionesPorFecha(fechaInicio, fechaFin string) error {
	if strings.TrimSpace(fechaInicio) == "" || strings.TrimSpace(fechaFin) == "" {
		return errors.New("fecha inicio o fecha fin estan vacios")
	}

	return s.repo.EliminarRelacionesPorFecha(fechaInicio, fechaFin)
}

func (s *RelacionService) ObtenerRelacionPorID(idRelacion string) (*repository.RelacionList, error) {
	if strings.TrimSpace(idRelacion) == "" {
		return nil, errors.New("El id de relacion esta vacio")
	}

	return s.repo.ObtenerRelacionPorID(idRelacion)
}
