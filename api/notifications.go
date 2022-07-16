package api

import (
	"github.com/PPSKSY-Cluster/backend/auth"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/PPSKSY-Cluster/backend/mail"
	"github.com/gofiber/fiber/v2"
)

func initNotificationHandlers(notificationRouter fiber.Router) {
	notificationRouter.Use(auth.CheckToken())
	notificationRouter.Post("/", notificationCreateHandler())
}

// @Description  Let server notify for reservations
// @Tags         reservationNotifications
// @Accept       json
// @Success      201
// @Failure      400  {object}  string
// @Failure		 500  {object}  string
// @Router       / [post]
func notificationCreateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		n := new(db.ReservationNotification)
		if err := c.BodyParser(n); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		notification, err := db.AddNotification(*n)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err := mail.InitSchedule(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		c.JSON(notification)
		return c.SendStatus(201)
	}
}
