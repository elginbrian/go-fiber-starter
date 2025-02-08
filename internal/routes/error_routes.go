package routes

import "github.com/gofiber/fiber/v2"

func setupErrorRoutes(app *fiber.App) {
	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"status":  "error",
			"message": "The route you requested does not exist. Please check the URL and try again.",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet && c.Method() != fiber.MethodPost {
			return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
				"code":    fiber.StatusMethodNotAllowed,
				"status":  "error",
				"message": "Method Not Allowed.",
			})
		}
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"status":  "error",
				"message": "Internal Server Error. Please try again later.",
			})
		}
		return nil
	})
}
