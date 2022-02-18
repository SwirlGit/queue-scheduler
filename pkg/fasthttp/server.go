package fasthttp

import "github.com/gofiber/fiber/v2"

type RouteProvider interface {
	RegisterFastHTTPRouters(a fiber.Router)
}

func NewServer(routeProviders []RouteProvider) *fiber.App {
	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	for i := range routeProviders {
		routeProviders[i].RegisterFastHTTPRouters(server)
	}

	return server
}
