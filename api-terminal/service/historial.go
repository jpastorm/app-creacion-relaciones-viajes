package service

import (
	"api-terminal/repository"
	"errors"
	"fmt"
	"strings"
)

// RelacionService es la estructura que contiene el repositorio
type HistorialService struct {
	repo *repository.Repository
}

// NewHistorialService crea una nueva instancia del servicio
func NewHistorialService(repo *repository.Repository) *HistorialService {
	return &HistorialService{repo: repo}
}

// AgregarRelacion agrega una nueva relación
func (s *HistorialService) AgregarHistorial(jsonData string) error {
	// Validaciones básicas
	if jsonData == "" {
		return errors.New("no se enviaron datos")
	}

	// Llamada al repositorio para agregar la relación
	err := s.repo.GuardarHistorial("./historial.db", jsonData)
	if err != nil {
		return fmt.Errorf("error al crear historial: %v", err)
	}

	return nil
}

func (s *HistorialService) ObtenerHistorialPaginado(page, pageSize int) (repository.RespuestaPaginada, error) {
	datos, err := s.repo.ListarHistorial("./historial.db", page, pageSize)
	if err != nil {
		return repository.RespuestaPaginada{}, fmt.Errorf("error al obtener la respuesta paginada: %v", err)
	}

	return datos, nil
}

func (s *HistorialService) EliminarPorfechas(fechaInicio, fechaFin string) error {
	if strings.TrimSpace(fechaInicio) == "" || strings.TrimSpace(fechaFin) == "" {
		return errors.New("fecha inicio o fecha fin estan vacios")
	}

	return s.repo.EliminarHistorialPorRango("./historial.db", fechaInicio, fechaFin)
}
