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

// @Description  Get all users
// @Tags         users
// @Produce      json
// @Success      200  {array}  User
// @Router       /api/users/ [get]
func userListHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON([]User{})
		return c.SendStatus(200)
	}
}

// @Description  Get user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  User
// @Router       /api/users/{id} [get]
func userDetailHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(User{})
		return c.SendStatus(200)
	}
}

// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201  {object}  User
// @Router       /api/users/ [post]
func userCreateHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(User{})
		return c.SendStatus(201)
	}
}

// @Description  Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  User
// @Router       /api/users/{id} [put]
func userUpdateHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(User{})
		return c.SendStatus(200)
	}
}

// @Description  Delete user
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Success      204
// @Router       /api/users/{id} [delete]
func userDeleteHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	}
}
