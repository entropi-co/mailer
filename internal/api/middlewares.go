package api

import (
	"github.com/gofiber/fiber/v2"
	"mailer/internal"
)

func AuthorizeService(ctx *fiber.Ctx) error {
	serviceKey, ok := ctx.GetReqHeaders()["ServiceKey"]
	if !ok {
		return fiber.ErrUnauthorized
	}

	if len(serviceKey) != 1 {
		return fiber.ErrBadRequest
	}

	if serviceKey[0] != internal.Config.GlobalServiceKey {
		return fiber.ErrUnauthorized
	}

	return ctx.Next()
}
