package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"dreamz.com/api/common"
	"dreamz.com/api/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct {
	store *common.Store
}

type User struct {
	Username string
	Password string
}

type UserPayload struct {
	UserId   string
	Username string
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.GetHeader("Authorization")
		splitToken := strings.Split(reqToken, "Bearer")
		if len(splitToken) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing bearer token"})
		}

		reqToken = strings.TrimSpace(splitToken[1])
		token, err := jwt.ParseWithClaims(reqToken, &UserPayload{}, func(t *jwt.Token) (interface{}, error) {
			brsa := []byte(strings.ReplaceAll(os.Getenv("RSA_PUBLIC_KEY"), "\\n", "\n"))
			rsa, err := jwt.ParseRSAPublicKeyFromPEM(brsa)
			if err != nil {
				log.Fatal("Error parsing public rsa key: ", err)
			}
			return rsa, nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Error parsing bearer token"})
		}
		if claims, ok := token.Claims.(*UserPayload); ok && token.Valid {
			if c.Keys == nil {
				c.Keys = make(map[string]interface{})
			}
			c.Keys["userId"] = claims.UserId
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid jwt token"})
		}
	}
}

func AddAuthRoute(r *gin.Engine, s *common.Store) {
	server := authServer{store: s}
	r.POST("auth/login", server.login)
}

func (server *authServer) login(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting body"})
	}
	var jsonBody User
	json.Unmarshal(body, &jsonBody)
	user := user.DbGetUser(server.store, jsonBody.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(jsonBody.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	} else {
		claim := &UserPayload{UserId: user.ID.Hex(), Username: user.Username}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

		brsa := []byte(strings.ReplaceAll(os.Getenv("RSA_PRIVATE_KEY"), "\\n", "\n"))
		rsa, err := jwt.ParseRSAPrivateKeyFromPEM(brsa)
		if err != nil {
			log.Fatal("Error parsing rsa key: ", err)
		}
		signed, err := token.SignedString(rsa)
		if err != nil {
			log.Fatal("Error sigin token: ", err)

		} else {
			c.String(http.StatusOK, signed)
		}
	}

}
