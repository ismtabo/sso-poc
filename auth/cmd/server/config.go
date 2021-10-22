package main

// Config represents application configuration
type Config struct {
	Server struct {
		Host string `yml:"host" envconfig:"SERVER_HOST"`
		Port string `yml:"port" envconfig:"SERVER_PORT"`
	} `yml:"server"`
	Log struct {
		Level string `yml:"level" envconfig:"LOG_LEVEL"`
	} `yml:"log"`
	Web struct {
		Static  string `yml:"static" envconfig:"WEB_STATIC"`
		Pages   string `yml:"pages" envconfig:"WEB_PAGES"`
		Session struct {
			Key    string `yml:"key" envconfig:"WEB_SESSION_KEY"`
			Cookie string `yml:"cookie" envconfig:"WEB_SESSION_COOKIE"`
		} `yml:"session"`
	} `yml:"web"`
	Database struct {
		File string `yml:"file" envconfig:"DATABASE_FILE"`
	} `yml:"database"`
}
