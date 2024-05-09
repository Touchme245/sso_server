package jwt

import (
	"fmt"
	"time"

	"github.com/Touchme245/sso_server/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration)
	claims["app_id"] = app.Id

	stringToken, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", "jwt sign failed", err)
	}
	return stringToken, nil

}
