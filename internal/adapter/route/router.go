package route

import (
	h "bam/internal/adapter/handler"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	*fiber.App
}

func NewRouter(orbianHandler h.OrbianHandler) (*Router, error) {
	app := fiber.New()

	api := app.Group("/api")
	v1 := api.Group("/v1")
	{
		// lunar := v1.Group("/lunar")
		// {
		// 	lunar.Post("/info")
		// 	lunar.Get("/info")
		// 	lunar.Post("/check")
		// }
		form := v1.Group("/form")
		{
			form.Post("/ordian-regis", orbianHandler.OrbianRegister)
		}
	}

	return &Router{App: app}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Listen(listenAddr)
}
