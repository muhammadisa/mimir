package mimir

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/option"
	"os"
	"time"
)

type JWTToken struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

func GenJSONWebToken(id int64) (*JWTToken, error) {
	// AccessToken
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return nil, err
	}
	// RefreshToken
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = id
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refresh, err := refreshToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return nil, err
	}
	return &JWTToken{
		Token:   token,
		Refresh: refresh,
	}, nil
}

func InitFirebaseApp(ctx context.Context, servicePath string) (*firebase.App, error) {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(servicePath))
	if err != nil {
		return nil, err
	}
	return app, nil
}

func GetPhoneFromToken(token *auth.Token) string {
	return TruncateCountryCode(token.Firebase.Identities["phone"].([]interface{})[0].(string))
}

func ReadGoogleIDToken(ctx context.Context, idToken string, app *firebase.App) (*auth.Token, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}
