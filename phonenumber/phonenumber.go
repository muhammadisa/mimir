package phonenumber

import (
	"firebase.google.com/go/v4/auth"
	"github.com/muhammadisa/mimir/strutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetPhoneFromToken(token *auth.Token) string {
	return strutil.TruncateCountryCode(token.Firebase.Identities["phone"].([]interface{})[0].(string))
}

func ValidatePhoneNumberLength(phoneNumber string, errMsg string) error {
	if len(phoneNumber) == 0 || len(phoneNumber) <= 10 {
		return status.Error(codes.InvalidArgument, errMsg)
	}
	return nil
}

func ValidatePhoneNumberFromGoogleIDToken(
	phoneNumber string,
	verifiedPhoneNumber *string,
	token *auth.Token,
	errMsg string,
) error {
	if phoneNumber == GetPhoneFromToken(token) {
		*verifiedPhoneNumber = phoneNumber
	} else {
		return status.Error(codes.Unauthenticated, errMsg)
	}
	return nil
}
