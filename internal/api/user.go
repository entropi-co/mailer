package api

import (
	"github.com/gofiber/fiber/v2"
	"mailer/internal/snowflakes"
)

type (
	UserCreatePayload struct {
		ID string `json:"id"`
	}
)

func (a *API) setupUserRoutes() {
	a.App.Use(AuthorizeService).Post("/user", a.handleUserCreate)
}

func (a *API) handleUserCreate(ctx *fiber.Ctx) error {
	var payload UserCreatePayload
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	id, err := snowflakes.ValueFromString(payload.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := a.Instance.Storage.CreateUser(id)
	if err != nil {
		return err
	}

	return ctx.JSON(user)
}
