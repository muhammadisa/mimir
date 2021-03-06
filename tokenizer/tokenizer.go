package tokenizer

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/muhammadisa/mimir/strutil"
	"google.golang.org/api/option"
	"os"
	"time"
)

type JWTToken struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

func ExtractIDJSONWebToken(ctx context.Context) (interface{}, string, error) {
	bearer, err := strutil.Bearer(ctx)
	if err != nil {
		return nil, "", err
	}
	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["timezone"] == nil {
			return claims["user_id"], "", nil
		} else {
			return claims["user_id"], claims["timezone"].(string), nil
		}
	}
	return nil, "", errors.New("token invalid")
}

func GenJSONWebToken(id int64, tz string, accessExpr time.Duration) (*JWTToken, error) {
	// AccessToken
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["timezone"] = tz
	claims["exp"] = time.Now().Add(accessExpr).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return nil, err
	}
	// RefreshToken
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = id
	rtClaims["timezone"] = tz
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
