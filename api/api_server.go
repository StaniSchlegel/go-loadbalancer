package api

import (
	"errors"
	"loadbalancer/api/routes"

	"github.com/gofiber/fiber/v2"
)

type ApiServer struct {
	ListeningAddress string
	fiberApp         *fiber.App
}

func NewApiServer(listeningAddress string) *ApiServer {
	a := ApiServer{
		ListeningAddress: listeningAddress,
		fiberApp:         fiber.New(),
	}

	routes.SetupNodesApiRoutes(a.fiberApp)

	return &a
}

func (a *ApiServer) StartListening() error {

	if a.fiberApp == nil {
		return errors.New("Api server was not initialized")
	}

	err := a.fiberApp.Listen(a.ListeningAddress)

	if err != nil {
		return err
	}

	return nil
}
