package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/utils/common"
	"yes4all/ads-noti-api/services/ads-noti/model/request"
	"yes4all/ads-noti-api/services/ads-noti/model/response"
)

type Auth interface {
	GenerateJWTHandler(ctx context.Context, req request.JWTRequest) (resp response.JWTResp, err error)
	Authentication(c *gin.Context)
}
type auth struct {
	websocketToken string
}

func NewAuth() Auth {
	websocketCfg := config.WebsocketConfig()
	return &auth{
		websocketToken: websocketCfg.WebsocketSecretKey,
	}
}

func (a *auth) GenerateJWTHandler(ctx context.Context, req request.JWTRequest) (resp response.JWTResp, err error) {
	// Định nghĩa các claims cho JWT
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_name"] = req.Username
	claims["exp"] = time.Now().Add(time.Second * 30).Unix() // Token hết hạn sau 30 giây

	// Tạo token với signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token và trả về chuỗi token
	tokenString, err := token.SignedString([]byte(a.websocketToken))
	if err != nil {
		return response.JWTResp{}, err
	}

	return response.JWTResp{
		Token: tokenString,
	}, nil
}

func (a *auth) Authentication(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}
	entry := logger.NewLogger().
		WithKeyword(c.Request.Context(), "auth.VerifyAuth").
		WithField("tokenString", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra signing method của token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.websocketToken), nil
	})

	if err != nil {
		common.HandleError(c, common.NewUnAuthorizedErrorWithMessage("Bạn không có quyền thực hiện theo user này"))
		entry.WithError(err).Error()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn không có quyền thực hiện theo user này"})
		return
	}

	// Kiểm tra tính hợp lệ của token
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		logger.NewLogger().WithKeyword(c.Request.Context(), "auth.Authentication").WithField("user_name", claims["user_name"]).Info()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		entry.WithError(err).Error()
		common.HandleError(c, common.NewUnAuthorizedErrorWithMessage("Invalid token"))
		return
	}
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username is required"})
		return
	}
	if username != claims["user_name"] {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user name not match"})
		entry.WithError(err).Error()
		common.HandleError(c, common.NewUnAuthorizedErrorWithMessage("user name not match"))
	}
	entry.Info()
	c.Next()

}
