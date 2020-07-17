package env

import (
	"os"
	"strings"
)

var (
	debug bool
	env   string
)

func init() {
	env = strings.ToLower(os.Getenv("MOT_ENV"))
	if env == "" {
		env = "development"
	}
	switch env {
	case "development":
		debug = true
	case "production":
		debug = false
	}
}

func Debug() bool {
	return debug
}

func Environment() string {
	return env
}
