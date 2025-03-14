package settingsManager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	PathToSettings string `json:"path_to_settings"`
	DownloadPath   string `json:"download_path"`
	UseRealDebrid bool `json:"use_real_debrid"`
	DebridToken string `json:"debrid_token"`
}

func (settings Settings) SaveSettings() error {
	// Idk random perms idk
	file, err := os.OpenFile(settings.PathToSettings, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not save settings: %w", err)
	}

	jsonData, err := json.MarshalIndent(settings, "", "    ")
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("could not write json data to settings: %w", err)
	}

	return nil
}

func LoadSettings(path string) (*Settings, error) {
	settings := &Settings{}
	settings.PathToSettings = path

	file, err := os.Open(path)
	if err != nil {
		settings.DownloadPath = filepath.Join("downloads")
		err := settings.GenerateSettings()
		if err != nil {
			return nil, err
		}
		file.Close()

		// Reopen file when done generating settings
		file, err = os.Open(path)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(settings)
	if err != nil {
		return nil, fmt.Errorf("Could not decode %v: %w", settings.PathToSettings, err.Error())
	}
	return settings, nil
}

func (settings *Settings) GenerateSettings() error {
	settingsFile, err := os.Create(settings.PathToSettings)
	if err != nil {
		return err
	}
	defer settingsFile.Close()

	err = settings.SaveSettings()
	if err != nil {
		return err
	}
	
	settings, err = LoadSettings(settings.PathToSettings)
	if err != nil {
		return fmt.Errorf("could not load settings: %w", err)
	}
	return nil
}
