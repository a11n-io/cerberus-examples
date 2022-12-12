package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/prinick96/elog"
)

// Default PORT to set if .env PORT var is empty or can't load
const DEFAULT_PORT_IF_EMPTY = "80"

// Env config struct
type EnvApp struct {
	// Server Envs
	APP_PORT    string
	JWT_SECRET  string
	SALT_ROUNDS int
}

// Get the env configuration
func GetEnv(env_file string) EnvApp {
	err := godotenv.Load(env_file)
	elog.New(elog.PANIC, "Error loading "+env_file+" file", err)

	// Heroku smell
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = DEFAULT_PORT_IF_EMPTY
	}

	saltRounds, err := strconv.ParseInt(os.Getenv("SALT_ROUNDS"), 10, 32)
	elog.New(elog.PANIC, "Error converting SALT_ROUNDS ["+os.Getenv("SALT_ROUNDS")+"] to int", err)

	return EnvApp{
		APP_PORT:    port,
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
		SALT_ROUNDS: int(saltRounds),
	}
}
