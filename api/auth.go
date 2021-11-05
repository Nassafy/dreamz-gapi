package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"dreamz.com/api/db"
	"dreamz.com/api/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string
	Password string
}

type UserPayload struct {
	UserId   string
	Username string
	jwt.StandardClaims
}

type RefreshPayload struct {
	UserId string
	jwt.StandardClaims
}

type AuthPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshDto struct {
	RefreshToken string `json:"refreshToken"`
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
				log.Panic("Error parsing public rsa key: ", err)
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

func (server *Server) Login(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting body"})
	}
	var jsonBody user
	json.Unmarshal(body, &jsonBody)
	user := db.GetUserByUsername(server.store, jsonBody.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(jsonBody.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	} else {
		accessToken := getAccessToken(user)
		refreshToken := createRefreshToken(server, user.Id)
		// db.SaveRefreshToken(model.RefreshToken#)
		c.JSON(
			http.StatusOK, AuthPayload{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			})
	}
}

func (server *Server) Refresh(c *gin.Context) {
	var tokenDto RefreshDto
	err := c.BindJSON(&tokenDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting body"})
	}
	claim, err := decodeRefreshToken(server, tokenDto.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}
	user := db.GetUserById(server.store, claim.UserId)
	accessToken := getAccessToken(user)
	refreshToken := createRefreshToken(server, user.Id)
	// db.SaveRefreshToken(model.RefreshToken#)
	c.JSON(
		http.StatusOK, AuthPayload{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
}

func getAccessToken(user model.User) string {
	expireAt := time.Now().Add(time.Minute * 30)
	claim := &UserPayload{
		UserId:   user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
			NotBefore: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	brsa := []byte(strings.ReplaceAll(os.Getenv("RSA_PRIVATE_KEY"), "\\n", "\n"))
	rsa, err := jwt.ParseRSAPrivateKeyFromPEM(brsa)
	if err != nil {
		log.Panic("Error parsing rsa key: ", err)
	}
	signed, err := token.SignedString(rsa)
	if err != nil {
		log.Panic("Error sigin token: ", err)

	}
	return signed
}

func createRefreshToken(server *Server, userId string) string {
	expireAt := time.Now().Add(time.Hour * 24 * 100)
	claim := RefreshPayload{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			ExpiresAt: expireAt.Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	brsa := []byte(strings.ReplaceAll(os.Getenv("RSA_PRIVATE_KEY"), "\\n", "\n"))
	rsa, err := jwt.ParseRSAPrivateKeyFromPEM(brsa)
	if err != nil {
		log.Panic("Error parsing rsa key: ", err)
	}
	signed, err := token.SignedString(rsa)
	if err != nil {
		log.Panic("Error sigin token: ", err)

	}
	dbRefreshToken := model.RefreshToken{
		Id:       claim.Id,
		Token:    signed,
		Valid:    true,
		ExpireAt: claim.ExpiresAt,
		UserId:   claim.UserId,
	}
	db.CreateRefreshToken(server.store, dbRefreshToken)
	return signed
}

func decodeRefreshToken(server *Server, refresh string) (*RefreshPayload, error) {
	token, err := jwt.ParseWithClaims(refresh, &RefreshPayload{}, func(t *jwt.Token) (interface{}, error) {
		brsa := []byte(strings.ReplaceAll(os.Getenv("RSA_PUBLIC_KEY"), "\\n", "\n"))
		rsa, err := jwt.ParseRSAPublicKeyFromPEM(brsa)
		if err != nil {
			log.Panic("Error parsing public rsa key: ", err)
		}
		return rsa, nil
	})
	if err != nil {
		log.Panic("Error decoding refresh token: ", err)
	}
	claim, ok := token.Claims.(*RefreshPayload)
	if ok && token.Valid {
		dbToken := db.GetRefreshToken(server.store, claim.Id, claim.UserId)
		if !dbToken.Valid {
			db.InvalidateUserRefreshToken(server.store, claim.UserId)
			return nil, errors.New("Invalid refresh token")
		} else {
			db.InvalidateRefreshToken(server.store, claim.Id)
		}
	}
	return claim, nil
}
