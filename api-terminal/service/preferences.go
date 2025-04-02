package service

import (
	"api-terminal/repository"
	"fmt"
)

type PreferenceService struct{}

func NewPreferenceService() *PreferenceService {
	return &PreferenceService{}
}

func (s *PreferenceService) SavePreferencias(prefs map[string]interface{}) error {
	err := repository.SavePreferencias(prefs)
	if err != nil {
		return fmt.Errorf("error al guardar plantilla: %v", err)
	}

	return nil
}

func (s *PreferenceService) GetPreferencias() (map[string]interface{}, error) {
	return repository.GetPreferencias()
}
