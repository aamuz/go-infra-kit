package config

import (
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

// Read function reads configuration from either a yaml file
// or from the environment variable
//
// filename: path to your yaml file
// config: Pointer to a struct type
func Read(filename string, config interface{}) error {
	configValue := reflect.ValueOf(config)
	if typ := configValue.Type(); typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config should be a pointer to a struct type")
	}

	if err := readFile(filename, config); err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	if err := readEnv(config); err != nil {
		return errors.New(fmt.Sprintf("Error reading ENV variables: %v", err))
	}

	return nil
}

func readFile(filename string, config interface{}) error {
	if len(filename) == 0 {
		return nil
	}

	f, err := os.OpenFile("config.yml", os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return err
	}

	return nil
}

func readEnv(config interface{}) error {
	err := envconfig.Process("", config)
	if err != nil {
		return err
	}
	return nil
}
