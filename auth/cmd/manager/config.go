package main

// Config represents application configuration
type Config struct {
	Log struct {
		Level string `yml:"level" envconfig:"LOG_LEVEL"`
	} `yml:"log"`
	Database struct {
		File string `yml:"file" envconfig:"DATABASE_FILE"`
	} `yml:"database"`
}
