package mimir

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

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
