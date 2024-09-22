package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

type AppConfig struct {
	DBConnectionString string `json:"DB_CONNECTION_STRING"`
	SignedSecretKey    string `json:"SIGNED_SECRET_KEY"`
}

var globalAppConfig AppConfig

func GetAppConfigs() AppConfig {
	return globalAppConfig
}

func LoadConfig(filePath string, config interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, config)
}

func SetupAppConfigs() error {
	baseConfig := &AppConfig{}
	err := LoadConfig("./config.json", baseConfig)
	if err != nil {
		log.Fatalf("Error loading base config: %v", err)
	}

	// Determine the environment (default to "development")
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Load the environment-specific config
	envConfig := &AppConfig{}
	envConfigFile := fmt.Sprintf("./config.%s.json", env)
	err = LoadConfig(envConfigFile, envConfig)
	if err != nil {
		log.Printf("Warning: Could not load environment config file %s, using base config: %v", envConfigFile, err)
	}

	// Merge base config with environment-specific config
	globalAppConfig = *MergeConfigs(baseConfig, envConfig).(*AppConfig)
	return nil
}

func MergeConfigs(base, override interface{}) interface{} {
	// Handle env override
	baseValue := reflect.ValueOf(base).Elem()
	overrideValue := reflect.ValueOf(override).Elem()

	// Ensure both base and override are structs
	if baseValue.Kind() != reflect.Struct || overrideValue.Kind() != reflect.Struct {
		return base
	}

	// Loop through the fields of the override struct
	for i := 0; i < overrideValue.NumField(); i++ {
		overrideField := overrideValue.Field(i)
		baseField := baseValue.Field(i)

		// Only override if the override field is set (not zero value)
		if !isZeroValue(overrideField) {
			baseField.Set(overrideField)
		}
	}
	return base
}

// Check if a reflect.Value is the zero value for its type
func isZeroValue(val reflect.Value) bool {
	return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
}
