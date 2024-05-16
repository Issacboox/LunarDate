package route

import (
	h "bam/internal/adapter/handler"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	*fiber.App
}

func NewRouter(ordianHandler h.OrdianHandler) (*Router, error) {
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
			form.Post("/ordian-regis", ordianHandler.OrdianRegister)
			form.Get("/ordian-list", ordianHandler.ListOrdian)
			form.Get("/ordian-info/:id", ordianHandler.OrdianIdEndpoint)
		}
	}

	return &Router{App: app}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Listen(listenAddr)
}
