package admin

import (
	"easyblog/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func auth(r *http.Request, jwtSecret string) error {
	tokenCookie, err := r.Cookie("auth_token")
	if err != nil {
		return err
	}
	token, err := utils.ParseJWT([]byte(jwtSecret), tokenCookie.Value)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid claims type, expected *jwt.RegisteredClaims")
	}
	sub, ok := claims["sub"].(string)
	if !ok || sub != "admin" {
		return fmt.Errorf("invalid token")
	}
	return nil
}
