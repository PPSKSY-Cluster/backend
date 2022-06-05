package handlers

import (
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		users, err := db.GetAllUsers()
		if err != nil {
			return c.SendStatus(500)
		}

		c.JSON(users)
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
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return c.SendStatus(500)
		}

		user, err := db.GetUserById(id)
		if err != nil {
			return c.SendStatus(500)
		}

		c.JSON(user)
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
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return c.SendStatus(500)
		}

		user, err := db.AddUser(*u)
		if err != nil {
			return c.SendStatus(500)
		}

		c.JSON(user)
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
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return c.SendStatus(500)
		}

		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)

		user, err := db.EditUser(id, *u)
		if err != nil {
			return c.SendStatus(500)
		}

		c.JSON(user)
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
