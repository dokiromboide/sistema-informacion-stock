package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080"
	}

	log.Printf("Servidor iniciando en el puerto %s", puerto)
}
