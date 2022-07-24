package api

import (
	"github.com/PPSKSY-Cluster/backend/auth"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/PPSKSY-Cluster/backend/mail"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initReservationHandlers(reservationRouter fiber.Router) {
	reservationRouter.Use(auth.CheckToken(db.UserUT))
	reservationRouter.Get("/", reservationListHandler())
	reservationRouter.Get("/:id", reservationDetailHandler())
	reservationRouter.Get("/users/:uId", reservationUserHandler())
	reservationRouter.Get("/clusters/:cId", reservationClusterHandler())
	reservationRouter.Post("/", reservationCreateHandler())
	reservationRouter.Put("/:id", reservationUpdateHandler())
	reservationRouter.Delete("/:id", reservationDeleteHandler())

	return
}

// @Description	Get all Reservations
// @Tags		reservations
// @Produce		json
// @Success		200 {array} db.Reservation
// @Failure		404 {object} string
// @Router		/api/reservations/ [get]
func reservationListHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		reservations, err := db.GetAllReservations()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(reservations)
		return c.SendStatus(200)
	}
}

// @Description	Get reservation by ID
// @Tags		reservations
// @Produce		json
// @Param 		id path string true "Reservation ID"
// @Success		200 {object} db.Reservation
// @Failure		400 {object} string
// @Failure		404 {object} string
// @Router		/api/reservations/{id} [get]
func reservationDetailHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		reservation, err := db.GetReservationById(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(reservation)
		return c.SendStatus(200)
	}
}

// @Description	Get reservation by userId
// @Tags		reservations
// @Produce		json
// @Param 		uId path string true "Reservation ID"
// @Success		200 {array} db.Reservation
// @Failure		400 {object} string
// @Failure		404 {object} string
// @Router		/api/reservations/{uId} [get]
func reservationUserHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("uId")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		reservations, err := db.GetReservationsByUserId(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(reservations)
		return c.SendStatus(200)
	}
}

// @Description	Get reservation by clusterID
// @Tags		reservations
// @Produce		json
// @Param 		cId path string true "Reservation ID"
// @Success		200 {array} db.Reservation
// @Failure		400 {object} string
// @Failure		404 {object} string
// @Router		/api/reservations/{cId} [get]
func reservationClusterHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("cId")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		reservations, err := db.GetReservationsByClusterId(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		c.JSON(reservations)
		return c.SendStatus(200)
	}
}

// @Description  Create reservations
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Success      201  {object}  db.Reservation
// @Failure      404  {object}  string
// @Failure		 400  {object}  string
// @Failure      500  {object}  string
// @Router       /api/reservations/ [post]
func reservationCreateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		r := new(db.Reservation)
		if err := c.BodyParser(r); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if r.StartTime > r.EndTime {
			return fiber.NewError(fiber.StatusBadRequest, "Start-date cannot be later than end-date")
		}

		b, err := isAvailable(*r)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if !b {
			return fiber.NewError(fiber.StatusBadRequest, "Not enough nodes available to make reservation")
		}

		reservation, err := db.AddReservation(*r)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = mail.ScheduleMail(reservation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		c.JSON(reservation)
		return c.SendStatus(201)
	}
}

// @Description  Update Reservation
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Reservation ID"
// @Success      200  {object}  db.Reservation
// @Failure		 400  {object}  string
// @Failure      500  {object}  string
// @Router       /api/reservations/{id} [put]
func reservationUpdateHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		r := new(db.Reservation)
		if err := c.BodyParser(r); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if r.StartTime > r.EndTime {
			return fiber.NewError(fiber.StatusBadRequest, "Start-date cannot be later than end-date")
		}

		b, err := isAvailable(*r)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if !b {
			return fiber.NewError(fiber.StatusBadRequest, "Not enough nodes available to make reservation")
		}

		reservation, err := db.EditReservation(id, *r)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		reservation.ID = id
		err = mail.ScheduleMail(reservation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		c.JSON(reservation)
		return c.SendStatus(200)
	}
}

// @Description  Delete reservations
// @Tags         reservations
// @Param        id   path      string  true  "Reservation ID"
// @Success      204
// @Failure      400  {object}  string
// @Failure		 404  {object}  string
// @Router       /api/reservations/{id} [delete]
func reservationDeleteHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Could not convert given object ID, did you use the right ID-format?")
		}

		err = db.DeleteReservation(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		mail.RemoveIfExists(id)

		return c.SendStatus(204)
	}
}

func isAvailable(reservation db.Reservation) (bool, error) {

	cluster, err := db.GetCResourceById(reservation.ClusterID)
	if err != nil {
		return false, err
	}

	reservations, err := db.GetReservationsByClusterId(reservation.ClusterID)
	if err != nil {
		return false, err
	}

	clusterReservations := make(map[int64]int64) //Maps startTime to nodes used

	for _, r := range reservations {
		if r.ID == reservation.ID { //If reservation already exists, skip this reservation
			continue
		}
		for start := reservation.StartTime; start <= reservation.EndTime; start += 86400 {
			if start < r.EndTime {
				clusterReservations[start] += r.Nodes
			} else {
				clusterReservations[start] += 0
			}
		}
	}

	for _, n := range clusterReservations {
		if cluster.Nodes-n < reservation.Nodes {
			return false, nil
		}
	}

	return true, nil
}
