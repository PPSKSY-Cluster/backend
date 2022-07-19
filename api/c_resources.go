package api

import (
	"github.com/PPSKSY-Cluster/backend/auth"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initCResourceHandlers(cresourceRouter fiber.Router) {
	cresourceRouter.Get("/", cResourceListHandler())
	cresourceRouter.Get("/:id", cResourceDetailHandler())

	cresourceRouter.Use(auth.CheckToken(db.AdminUT))
	cresourceRouter.Post("/", cResourceCreateHandler())
	cresourceRouter.Put("/:id", cResourceUpdateHandler())
	cresourceRouter.Delete("/:id", cResourceDeleteHandler())

	return
}

// @Description  Get all cluster resources
// @Tags         cluster-resources
// @Produce      json
// @Success      200  {array}  db.CResource
// @Failure		 404  {object} string
// @Router       /api/cresources/ [get]
func cResourceListHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cResources, err := db.GetAllCResources()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(cResources)
		return c.SendStatus(200)
	}
}

// @Description  Get cluster resource by ID
// @Tags         cluster-resources
// @Produce      json
// @Param        id   path      string  true  "CResource ID"
// @Success      200  {object}  db.CResource
// @Failure      400  {object}  string
// @Failure 	 404  {object}  string
// @Router       /api/cresources/{id} [get]
func cResourceDetailHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		cResource, err := db.GetCResourceById(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(cResource)
		return c.SendStatus(200)
	}
}

// @Description  Create cluster resource
// @Tags         cluster-resources
// @Accept       json
// @Produce      json
// @Success      201  {object}  db.CResource
// @Failure      400  {object}  string
// @Failure 	 404  {object}  string
// @Router       /api/cresources/ [post]
func cResourceCreateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cr := new(db.CResource)
		if err := c.BodyParser(cr); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		cResource, err := db.AddCResource(*cr)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(cResource)
		return c.SendStatus(201)
	}
}

// @Description  Update cluster resource
// @Tags         cluster-resources
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "CResource ID"
// @Success      200  {object}  db.CResource
// @Failure      400  {object}  string
// @Failure 	 404  {object}  string
// @Router       /api/cresources/{id} [put]
func cResourceUpdateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cr := new(db.CResource)
		if err := c.BodyParser(cr); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		cResource, err := db.EditCResource(id, *cr)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		cResource.ID = id
		c.JSON(cResource)
		return c.SendStatus(200)
	}
}

// @Description  Delete cluster resource
// @Tags         cluster-resources
// @Param        id   path      string  true  "CResource ID"
// @Success      204
// @Failure      500
// @Router       /api/cresources/{id} [delete]
func cResourceDeleteHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		err = db.DeleteCResource(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return c.SendStatus(204)
	}
}
