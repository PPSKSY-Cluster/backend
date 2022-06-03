package handlers

import "github.com/gofiber/fiber/v2"

type User struct {
	ID string `json:"_id"`
}

func InitUserHandlers(userRouter fiber.Router) {
	userRouter.Get("/", userListHandler())
	userRouter.Get("/:id", userDetailHandler())
	userRouter.Post("/", userCreateHandler())
	userRouter.Put("/:id", userUpdateHandler())
	userRouter.Delete("/:id", userDeleteHandler())

	return
}

func userListHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON([]User{})
		return c.SendStatus(200)
	}
}

func userDetailHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(User{})
		return c.SendStatus(200)
	}
}

func userCreateHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(User{})
		return c.SendStatus(201)
	}
}

func userUpdateHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(User{})
		return c.SendStatus(200)
	}
}

func userDeleteHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	}
}
