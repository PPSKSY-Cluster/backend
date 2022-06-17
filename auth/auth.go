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
func CheckCredentials() func(username, password string) (string, error) {
	privKeyPem, err := readKeyFromFile(os.Getenv("PRIV_KEY_PATH"))
	if err != nil {
		panic(err.Error())
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(privKeyPem, os.Getenv("PRIVATE_PEM_PW"))
	if err != nil {
		panic(err.Error())
	}

	return func(username, password string) (string, error) {
		user, err := db.GetUserCredentials(username)
		if err != nil {
			return "", err
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return "", err
		}

		token := jwt.New(jwt.SigningMethodRS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // expires in 24 hours
		t, err := token.SignedString(privKey)
		if err != nil {
			return "", err
		}

		return t, nil
	}
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
	pubKeyPem, err := readKeyFromFile(os.Getenv("PUB_KEY_PATH"))
	if err != nil {
		panic(err.Error())
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyPem)
	if err != nil {
		panic(err.Error())
	}

	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			return c.SendStatus(401)
		}

		_, err = jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		if err != nil {
			return c.SendStatus(401)
		}

		return c.Next()
	}
}

func readKeyFromFile(path string) ([]byte, error) {
	finfo, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	key := make([]byte, finfo.Size())
	keyF, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer keyF.Close()

	keyF.Read(key)
	return key, nil
}
