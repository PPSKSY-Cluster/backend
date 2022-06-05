package handlers

import "github.com/gofiber/fiber/v2"

type CResource struct {
	ID string `json:"_id"`
}

func InitCResourceHandlers(cresourceRouter fiber.Router) {
	cresourceRouter.Get("/", cresourceListHandler())
	cresourceRouter.Get("/:id", cresourceDetailHandler())
	cresourceRouter.Post("/", cresourceCreateHandler())
	cresourceRouter.Put("/:id", cresourceUpdateHandler())
	cresourceRouter.Delete("/:id", cresourceDeleteHandler())

	return
}

// @Description  Get all cluster resources
// @Tags         cluster-resources
// @Produce      json
// @Success      200  {array}  CResource
// @Router       /api/cresources/ [get]
func cresourceListHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON([]CResource{})
		return c.SendStatus(200)
	}
}

// @Description  Get cluster resource by ID
// @Tags         cluster-resources
// @Produce      json
// @Param        id   path      string  true  "CResource ID"
// @Success      200  {object}  CResource
// @Router       /api/cresources/{id} [get]
func cresourceDetailHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(CResource{})
		return c.SendStatus(200)
	}
}

// @Description  Create cluster resource
// @Tags         cluster-resources
// @Accept       json
// @Produce      json
// @Success      201  {object}  CResource
// @Router       /api/cresources/ [post]
func cresourceCreateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(CResource{})
		return c.SendStatus(201)
	}
}

// @Description  Update cluster resource
// @Tags         cluster-resources
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "CResource ID"
// @Success      200  {object}  CResource
// @Router       /api/cresources/{id} [put]
func cresourceUpdateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(CResource{})
		return c.SendStatus(200)
	}
}

// @Description  Delete cluster resource
// @Tags         cluster-resources
// @Param        id   path      string  true  "CResource ID"
// @Success      204
// @Router       /api/cresources/{id} [delete]
func cresourceDeleteHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	}
}
