package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/prinick96/elog"
)

// DefaultPortIfEmpty Default PORT to set if .env PORT var is empty or can't load
const DefaultPortIfEmpty = "80"

// App config struct
type App struct {
	// Server Envs
	AppPort    string
	JwtSecret  string
	SaltRounds int
}

// GetEnv Get the env configuration
func GetEnv(envFile string) App {
	err := godotenv.Load(envFile)
	elog.New(elog.PANIC, "Error loading "+envFile+" file", err)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = DefaultPortIfEmpty
	}

	saltRounds, err := strconv.ParseInt(os.Getenv("SALT_ROUNDS"), 10, 32)
	elog.New(elog.PANIC, "Error converting SALT_ROUNDS ["+os.Getenv("SALT_ROUNDS")+"] to int", err)

	return App{
		AppPort:    port,
		JwtSecret:  os.Getenv("JWT_SECRET"),
		SaltRounds: int(saltRounds),
	}
}
