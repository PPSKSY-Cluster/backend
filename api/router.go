package api

import (
	"os"

	"github.com/PPSKSY-Cluster/backend/auth"
	_ "github.com/PPSKSY-Cluster/backend/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/bson"
)

func InitRouter() (*fiber.App, error) {
	router := fiber.New()

	// define CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CLIENT_URL"),
		AllowHeaders: "Content-Type, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// setup logging
	router.Use(logger.New())

	// define api routes
	api := router.Group("/api")
	api.Get("/ping", pingHandler())
	api.Get("/docs/*", docsHandler())
	api.Post("/login", loginHandler())

	tokenRoutes := api.Group("/token-check")
	tokenRoutes.Use(auth.CheckToken())
	tokenRoutes.Post("/", tokenCheckHandler())

	userRoutes := api.Group("/users")
	initUserHandlers(userRoutes)

	cResourceRoutes := api.Group("/cresources")
	initCResourceHandlers(cResourceRoutes)

	reservationRoutes := api.Group("/reservations")
	initReservationHandlers(reservationRoutes)

	return router, nil
}

// @Description  Route for testing jwt token validity
// @Tags         general
// @Success      200
// @Failure		 401
// @Router       /api/login [post]
func tokenCheckHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}
}

// @Description  Route for login
// @Tags         general
// @Accept       json
// @Produce      json
// @Param        username   body   string  true  "Username"
// @Param        password   body   string  true  "Password"
// @Success      200  {string} string
// @Failure		 401
// @Router       /api/login [post]
func loginHandler() func(c *fiber.Ctx) error {
	type LoginPair struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	checkCredentialsF := auth.CheckCredentials()

	return func(c *fiber.Ctx) error {
		var login LoginPair

		if err := c.BodyParser(&login); err != nil {
			return c.SendStatus(500)
		}

		user, token, err := checkCredentialsF(login.Username, login.Password)
		if err != nil {
			return c.SendStatus(401)
		}

		c.JSON(bson.M{"token": token, "user": user})
		return c.SendStatus(200)
	}
}

// @Description  The route that serves the swagger documentation
// @Tags         general
// @Produce      html
// @Success      200  {html} html
// @Router       /api/docs/ [get]
func docsHandler() func(*fiber.Ctx) error {
	return swagger.HandlerDefault
}

// @Description  Ping route to act as healthcheck
// @Tags         general
// @Produce      json
// @Success      200  {string} string
// @Router       /api/ping [get]
func pingHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
}
