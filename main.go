package main

import (
	"PortalCRG/internal"
	"PortalCRG/internal/repository"
	"PortalCRG/server"
	"cmp"
	"fmt"
	"os"
)

func main() {
	// Inicializar repositorios
	userRepository := repository.NewUserRepositoryMongo()

	// Inicializar casos de uso
	userService := internal.NewUserService(*userRepository)

	userRepository.Init()

	// Inicializar servidor HTTP
	httpServer := server.NewHTTPServer(userService)

	port := cmp.Or(os.Getenv("PORT"), "80")

	// Iniciar el servidor HTTP
	err := httpServer.Start(port)
	if err != nil {
		fmt.Printf("Error iniciando servidor HTTP: %v\n", err)
	}
}
