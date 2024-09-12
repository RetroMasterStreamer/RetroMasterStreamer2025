package main

import (
	"PortalCRG/internal"
	"PortalCRG/internal/repository"
	"PortalCRG/internal/util"
	"PortalCRG/server"
	"cmp"
	"fmt"
	"os"
)

func main() {

	util.DeleteURLDuplicada()

	// Inicializar repositorios
	userRepository := repository.NewUserRepositoryMongo()
	portalRepository := repository.NewPortalRepositoryMongo()

	// Inicializar casos de uso
	userService := internal.NewUserService(*userRepository, *portalRepository)

	userRepository.Init()

	//userService.PortalRepository.DeleteTipsFromDate("11-09-2024")

	// Inicializar servidor HTTP
	httpServer := server.NewHTTPServer(userService)

	port := cmp.Or(os.Getenv("PORT"), "80")

	// Iniciar el servidor HTTP
	err := httpServer.Start(port)
	if err != nil {
		fmt.Printf("Error iniciando servidor HTTP: %v\n", err)
	}
}
