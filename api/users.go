package api

import (
	"github.com/PPSKSY-Cluster/backend/auth"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initUserHandlers(userRouter fiber.Router) {
	userRouter.Post("/", userCreateHandler())

	userRouter.Use(auth.CheckToken())
	userRouter.Get("/", userListHandler())
	userRouter.Get("/:id", userDetailHandler())
	userRouter.Put("/:id", userUpdateHandler())
	userRouter.Delete("/:id", userDeleteHandler())

	return
}

// @Description  Get all users
// @Tags         users
// @Produce      json
// @Success      200  {array}  db.User
// @Failure	     404  {object}  string
// @Router       /api/users/ [get]
func userListHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		users, err := db.GetAllUsers()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(users)
		return c.SendStatus(200)
	}
}

// @Description  Get user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  db.User
// @Failure	     404  {object}  string
// @Failure	     400  {object}  string
// @Router       /api/users/{id} [get]
func userDetailHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		user, err := db.GetUserById(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(user)
		return c.SendStatus(200)
	}
}

// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201  {object}  db.User
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /api/users/ [post]
func userCreateHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		hashedPW, err := auth.HashPW(u.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		u.Password = hashedPW

		user, err := db.AddUser(*u)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
// @Success      200  {object}  db.User
// @Failure 	 400  {object}  string
// @Failure	     404  {object}  string
// @Router       /api/users/{id} [put]
func userUpdateHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		user, err := db.EditUser(id, *u)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		user.ID = id
		c.JSON(user)
		return c.SendStatus(200)
	}
}

// @Description  Delete user
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Success      204
// @Failure	     400  {object}  string
// @Failure      404  {object}  string
// @Router       /api/users/{id} [delete]
func userDeleteHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		err = db.DeleteUser(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return c.SendStatus(204)
	}
}
