package auth

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Helper to generate a JWT token on login
func CheckCredentials(username, password string) (string, error) {
	user, err := db.GetUserCredentials(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // expires in 24 hours

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}

// Helper to generate the hashed password on user creation
func HashPW(password string) (string, error) {
	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		return "", err
	}

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}

// Middleware that checks if the user is authenticated
func CheckToken() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			return c.SendStatus(401)
		}

		_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}

			return []byte("secret"), nil
		})

		if err != nil {
			return c.SendStatus(401)
		}

		return c.Next()
	}
}
