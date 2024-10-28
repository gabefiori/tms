package config

import (
	"os"

	"github.com/gabefiori/tms/internal/targets"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	File
	Cli
}

type Cli struct {
	List         bool
	OutputTarget bool
	Filter       string
	Path         string
}

type File struct {
	Selector []string              `json:"selector"`
	Targets  []targets.InputTarget `json:"targets"`
}

func Load(cli Cli) (*Config, error) {
	path, err := homedir.Expand(cli.Path)

	if err != nil {
		return nil, err
	}

	fb, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	var file File

	decoder := jsoniter.ConfigFastest.NewDecoder(fb)
	if err := decoder.Decode(&file); err != nil {
		return nil, err
	}

	return &Config{file, cli}, nil
}
