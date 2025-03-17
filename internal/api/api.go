package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"mailer/internal/instance"
)

type keyType struct{}

var contextApiKey keyType

type API struct {
	App      *fiber.App
	Instance *instance.Instance
}

// ServeAPI creates API server listening on the address set in environment variable
// This function blocks current routine
func ServeAPI(instance *instance.Instance) {
	app := fiber.New(fiber.Config{
		Immutable: false,
	})

	api := API{
		App:      app,
		Instance: instance,
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(contextApiKey, api)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	logrus.Infoln("API Listening on port 3000")
	err := app.Listen(":3000")
	if err != nil {
		logrus.Fatalf("Failed to start api server: %+v", err)
		return
	}
}
