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
// @Success      200
// @Failure      400  {object}  string
// @Failure		 500  {object}  string
// @Router       / [post]
func notificationCreateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		notification := new(db.ReservationNotification)
		if err := c.BodyParser(notification); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err := db.AddNotification(*notification); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err := mail.InitSchedule(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(200)
	}
}
