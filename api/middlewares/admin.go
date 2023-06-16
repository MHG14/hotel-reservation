package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	if !user.IsAdmin {
		return fmt.Errorf("unauthorized")

	}
	return c.Next()
}
