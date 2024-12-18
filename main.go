package main

import (
	"PortalCRG/internal"
	"PortalCRG/internal/repository"
	"PortalCRG/internal/util"
	"PortalCRG/server"
	"cmp"
	"fmt"
	"log"
	"os"
)

func main() {

	util.DeleteURLDuplicada()

	// Inicializar repositorios
	userRepository := repository.NewUserRepositoryMongo()
	portalRepository := repository.NewPortalRepositoryMongo()

	// Inicializar casos de uso
	userService := internal.NewUserService(*userRepository, *portalRepository)
	driveService := internal.NewDriveService()
	retroEmailService := internal.NewRetroEmailService()

	retroEmailService.Init()

	errDrive := driveService.Connect()
	if errDrive == nil {
		driveService.CreateTable()
	} else {
		log.Println("ATENCION REPOSITORIO DE DESCARGA APAGADO!!")
	}
	userRepository.Init()

	//userService.PortalRepository.DeleteTipsFromDate("11-09-2024")

	// Inicializar servidor HTTP
	httpServer := server.NewHTTPServer(userService, *driveService, *retroEmailService)

	port := cmp.Or(os.Getenv("PORT"), "80")

	// Iniciar el servidor HTTP
	err := httpServer.Start(port)
	if err != nil {
		fmt.Printf("Error iniciando servidor HTTP: %v\n", err)
	}
}
