package api

import "github.com/gofiber/fiber/v2"

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

	user, err := a.Instance.Storage.CreateUser(payload.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(user)
}
