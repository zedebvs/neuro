package config

import (
	"encoding/json"
	"neuro/utils"
	"os"
)

const FILENAME = "config.json"

var Conf Config

type Config struct {
	Data          string            `json:"Data"`
	LogsDir       string            `json:"LogsDir"`
	Storage       string            `json:"Storage"`
	LogFiles      map[string]string `json:"LogFiles"`
	PromptsDir    string            `json:"Prompts"`
	SystemPrompts map[string]string `json:"SystemPrompts"`
}

func (C *Config) DataPath(base string, filename string) (string, error) {
	path, _ := utils.GetPath([]string{C.Data, base, filename})
	return path, nil
}

func init() {
	config, _ := utils.GetPath([]string{FILENAME})
	data, err := os.ReadFile(config)
	if err != nil {
		utils.RedPanic("Не удалось прочитать конфигурационный файл [по умолчанию - config.json]")
	}

	err = json.Unmarshal(data, &Conf)
	if err != nil {
		utils.RedPanic("Не удалось загрузить конфиг в память")
	}
}
