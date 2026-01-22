package api

import (
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "api ", log.LstdFlags|log.LUTC)
}
