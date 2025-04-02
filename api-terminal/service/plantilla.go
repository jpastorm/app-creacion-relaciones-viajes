package service

import (
	"api-terminal/repository"
	"errors"
	"fmt"
)

type PlantillaService struct{}

func NewPlantillaService() *PlantillaService {
	return &PlantillaService{}
}

func (s *PlantillaService) GuardarPlantilla(titulo string, fuente string, data []map[string]interface{}) error {
	if titulo == "" {
		return errors.New("el campo titulo es obligatorio")
	}

	err := repository.GuardarPlantilla("./plantillas.db", titulo, fuente, data)
	if err != nil {
		return fmt.Errorf("error al guardar plantilla: %v", err)
	}

	return nil
}

func (s *PlantillaService) BuscarPlantilla(id int) (repository.Impresion, error) {
	if id == 0 {
		return repository.Impresion{}, errors.New("el campo titulo es obligatorio")
	}

	plantillas, err := repository.BuscarPlantilla("./plantillas.db", id)
	if err != nil {
		return repository.Impresion{}, fmt.Errorf("error al guardar plantilla: %v", err)
	}

	return *plantillas, nil
}

func (s *PlantillaService) ObtenerTodasLasPlantillas() ([]repository.Impresion, error) {
	plantillas, err := repository.ObtenerTodasLasPlantillas("./plantillas.db")
	if err != nil {
		return make([]repository.Impresion, 0), fmt.Errorf("error al guardar plantilla: %v", err)
	}

	return plantillas, nil
}

func (s *PlantillaService) ActualizarPlantilla(id int, titulo string, fuente string, datos []map[string]interface{}) error {
	if id == 0 {
		return fmt.Errorf("el campo ID es obligatorio")
	}
	err := repository.ActualizarPlantilla("./plantillas.db", id, titulo, fuente, datos)
	if err != nil {
		return fmt.Errorf("error al actualizar la plantilla %s", err.Error())
	}

	return nil
}

func (s *PlantillaService) EliminarPlantilla(id int) error {
	if id == 0 {
		return fmt.Errorf("el campo ID es obligatorio")
	}

	err := repository.EliminarPlantilla("./plantillas.db", id)
	if err != nil {
		return fmt.Errorf("error al eliminar la plantilla %s", err.Error())
	}

	return nil
}
