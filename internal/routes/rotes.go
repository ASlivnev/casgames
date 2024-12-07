package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"casualgames/internal/handler"
)

func NewRoutes(h *handler.Handler) *fiber.App {
	app := fiber.New()
	app.Get("/ping", func(c *fiber.Ctx) error {
		log.Info().Msg("Healthcheck is working!")
		return c.SendString("pong")
	})

	app.Get("/deallocateAll", h.DeallocateAll)
	app.Get("/startGamesParser", h.GamesParser)
	app.Get("/api/getGames/:page", h.GetGames)
	app.Get("/api/incrementGameRang/:gameId", h.IncrementGameRang)
	//app.Get("/keywords", h.GetKeyWordsList)

	app.Static("/uploads", "./static")
	return app
}
